package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetsuo/restache"
	"github.com/tetsuo/restache/jsx"
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
		compileStdin()
		os.Exit(0)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		fatalf("could not get current directory: %v", err)
	}

	dirIncludes := resolveGlobs(baseDir, patterns)

	if len(dirIncludes) == 0 {
		fatalf("no files matched the provided pattern")
	}

	if outdir != "" {
		if !filepath.IsAbs(outdir) {
			outdir = filepath.Join(baseDir, outdir)
		}
	}

	for dir, includes := range dirIncludes {
		if len(includes) == 1 {
			os.Exit(compileSingle(dir, outdir, includes[0]))
		}
		parallelism = min(parallelism, 32)
		nodes, err := restache.ParseDir(dir, includes, restache.WithParallelism(parallelism))
		if err != nil {
			fatalf("failed to parse directory %q: %v", dir, err)
		}
		fmt.Println(nodes)
	}
}

func compileStdin() {
	node, err := restache.Parse(os.Stdin)
	if err != nil {
		fatalf("failed to parse stdin: %v", err)
	}
	src := jsx.NewReader(node)
	if _, err = io.Copy(os.Stdout, src); err != nil {
		fatalf("failed to write to stdout: %v", err)
	}
}

func compileSingle(dir, outdir, filename string) int {
	node, err := restache.ParseFile(filepath.Join(dir, filename))
	if err != nil {
		fatalf("failed to parse file %q: %v", filename, err)
	}
	ext := filepath.Ext(filename)
	if ext != "" {
		filename = filename[:len(filename)-len(ext)]
	}
	node.Data = []byte(filename)
	src := jsx.NewReader(node)
	filename += ".jsx"
	if outdir != "" {
		if err := os.MkdirAll(outdir, 0755); err != nil {
			fatalf("could not create output directory %q: %v", outdir, err)
		}
		filename = filepath.Join(outdir, filename)
	} else {
		filename = filepath.Join(dir, filename)
	}
	dst, err := os.Create(filename)
	if err != nil {
		fatalf("could not create file %q: %v", filename, err)
	}
	code := 0
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to write output: %v\n", PROGRAM_NAME, err)
		code = 1
	}
	dst.Close()
	return code
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

var sep = string(filepath.Separator)

func commonBaseDir(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

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
