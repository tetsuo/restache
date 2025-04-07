package main

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/tetsuo/restache"
	"github.com/tetsuo/restache/jsx"
	"golang.org/x/sync/errgroup"
)

const PROGRAM_NAME = "restache"

var VERSION = "development"

func usage() {
	fmt.Printf("Usage: %s [OPTION] [GLOB ...]\n", PROGRAM_NAME)
	fmt.Println("Transpile Restache templates into React/JSX components.")
	fmt.Println("With no GLOB, or when GLOB is -, read standard input.")
	fmt.Println("  -h, --help            display this help and exit")
	fmt.Println("  -v, --version         output version information and exit")
	fmt.Println("  -p, --parallelism N   number of files to process in parallel (default: number of CPUs)")
	fmt.Println("  -o, --outdir DIR      write output files to DIR (default: same as input files)")
	fmt.Println("      --redux           generate components pre-wired for Redux")
}

func main() {
	// Reject single-dash long options like -parallelism
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 && !strings.Contains(arg, "=") {
			fmt.Fprintf(os.Stderr, "%s: invalid option: %s (did you mean --%s?)\n", PROGRAM_NAME, arg, arg[1:])
			fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", PROGRAM_NAME)
			os.Exit(1)
		}
	}

	var (
		help        bool
		version     bool
		parallelism int
		outdir      string
	)

	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&help, "h", false, "")

	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&version, "v", false, "")

	numCPUS := runtime.NumCPU()
	flag.IntVar(&parallelism, "parallelism", numCPUS, "")
	flag.IntVar(&parallelism, "p", numCPUS, "")

	flag.StringVar(&outdir, "outdir", "", "output dir")
	flag.StringVar(&outdir, "o", "", "output dir")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}
	if version {
		fmt.Printf("%s version %s\n", PROGRAM_NAME, VERSION)
		os.Exit(0)
	}

	patterns := flag.Args()
	n := len(patterns)

	if n == 0 || patterns[0] == "-" {
		if outdir != "" {
			fmt.Fprintf(os.Stderr, "%s: ignoring --outdir (no input files)\n", PROGRAM_NAME)
		}
		processStdin()
		os.Exit(0)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		fatalf("could not get current directory: %v", err)
	}

	filesByDir := resolveGlobs(baseDir, patterns)

	if len(filesByDir) == 0 {
		fatalf("no files matched the provided pattern")
	}

	if outdir != "" {
		if !filepath.IsAbs(outdir) {
			outdir = filepath.Join(baseDir, outdir)
		}
	}

	commonDir := commonBaseDir(slices.Collect(maps.Keys(filesByDir)))

	if parallelism < 1 {
		parallelism = runtime.NumCPU()
	}

	parallelism = min(parallelism, 32)

	for dir, files := range filesByDir {
		if len(files) == 1 {
			emitFile(dir, outdir, files[0])
		} else {
			var dst string
			if outdir == "" {
				dst = dir
			} else {
				dst, err = filepath.Rel(commonDir, dir)
				if err != nil {
					fatalf("could not determine relative path from %q to %q: %v", commonDir, dir, err)
				}
				dst = filepath.Join(outdir, dst)
			}
			emitModule(dir, dst, files, parallelism)
		}
	}
}

func processStdin() {
	node, err := restache.Parse(os.Stdin)
	if err != nil {
		fatalf("failed to parse stdin: %v", err)
	}
	if _, err = io.Copy(os.Stdout, jsx.NewReader(node)); err != nil {
		fatalf("failed to write to stdout: %v", err)
	}
}

func emitFile(dir, outdir, file string) {
	node, err := restache.ParseFile(filepath.Join(dir, file))
	if err != nil {
		fatalf("failed to parse file %q: %v", file, err)
	}
	ext := filepath.Ext(file)
	if ext != "" {
		file = file[:len(file)-len(ext)]
	}
	node.Data = file
	file += ".jsx"
	if outdir != "" {
		file = filepath.Join(outdir, file)
		if err := os.MkdirAll(outdir, 0755); err != nil {
			fatalf("could not create output directory %q: %v", outdir, err)
		}
	} else {
		file = filepath.Join(dir, file)
	}
	dst, err := os.Create(file)
	if err != nil {
		fatalf("could not create file %q: %v", file, err)
	}
	defer dst.Close()
	if _, err = io.Copy(dst, jsx.NewReader(node)); err != nil {
		dst.Close()
		fatalf("failed to write file %q: %v", file, err)
	}
}

func emitModule(dir, outdir string, files []string, parallelism int) {
	nodes, err := restache.ParseDir(dir, files, restache.WithParallelism(parallelism))
	if err != nil {
		fatalf("failed to parse directory %q: %v", dir, err)
	}
	if err := os.MkdirAll(outdir, 0755); err != nil {
		fatalf("could not create output directory %q: %v", outdir, err)
	}
	var g errgroup.Group
	g.SetLimit(min(len(nodes), parallelism))
	for _, node := range nodes {
		node := node
		g.Go(func() error {
			outfile := filepath.Join(outdir, node.Path[:len(node.Path)-1][0].Key) + ".jsx"
			dst, err := os.Create(outfile)
			if err != nil {
				return fmt.Errorf("could not create file %q: %v", outfile, err)
			}
			defer dst.Close()
			if _, err = io.Copy(dst, jsx.NewReader(node)); err != nil {
				return fmt.Errorf("failed to write file %q: %v", outfile, err)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		fatalf("%s", err.Error())
	}
}

func resolveGlobs(baseDir string, patterns []string) (dirs map[string][]string) {
	var (
		info os.FileInfo
		dir  string
		p    string
	)
	dirs = make(map[string][]string)
	for _, p = range patterns {
		matches, err := filepath.Glob(filepath.Join(baseDir, p))
		if err != nil {
			fatalf("invalid glob %q: %v", p, err)
		}
		for _, p = range matches {
			info, err = os.Lstat(p) // will ignore symlinks
			if err != nil {
				fatalf("could not access file %q: %v", p, err)
			}
			if info.IsDir() {
				continue
			}
			dir = filepath.Dir(p)
			dirs[dir] = append(dirs[dir], info.Name())
		}
	}
	return
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", PROGRAM_NAME, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func commonBaseDir(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	sep := string(filepath.Separator)
	segments := strings.Split(filepath.Clean(paths[0]), sep)
	for _, path := range paths[1:] {
		curr := strings.Split(filepath.Clean(path), sep)
		n := min(len(segments), len(curr))
		i := 0
		for i < n && segments[i] == curr[i] {
			i++
		}
		segments = segments[:i]
	}
	return sep + filepath.Join(segments...)
}
