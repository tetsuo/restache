package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tetsuo/commonpath"
	"github.com/tetsuo/restache"
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
}

func main() {
	// Reject single-dash long options like -parallelism
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 && !strings.Contains(arg, "=") {
			logf("invalid option: %s (did you mean --%s?)", arg, arg[1:])
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
			logf("ignoring --outdir (no input files)")
		}
		// Process stdin:
		node, err := restache.Parse(os.Stdin)
		if err != nil {
			fatalf("failed to parse stdin: %v", err)
		}
		if _, err := restache.Render(os.Stdout, node); err != nil {
			fatalf("failed to write to stdout: %v", err)
		}
		os.Exit(0)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		fatalf("could not get current directory: %v", err)
	}

	filesByDir := resolveGlobs(baseDir, patterns)

	if len(filesByDir) == 0 {
		os.Exit(1)
	}

	const maxParallelism = 32

	if parallelism < 1 {
		parallelism = numCPUS
	} else if parallelism > maxParallelism {
		parallelism = maxParallelism
	}

	if outdir != "" {
		if !filepath.IsAbs(outdir) {
			outdir = filepath.Join(baseDir, outdir)
		}
		if err := os.MkdirAll(outdir, 0755); err != nil {
			fatalf("could not create output directory %q: %v", outdir, err)
		}
		common, err := commonpath.CommonPath(collectKeys(filesByDir))
		if err != nil {
			fatalf("could not determine the common sub-path: %v", err)
		}
		for dir, files := range filesByDir {
			actualOutDir, err := filepath.Rel(common, dir)
			if err != nil {
				fatalf("could not determine relative path from %q to %q: %v", common, dir, err)
			}
			actualOutDir = filepath.Join(outdir, actualOutDir)
			if err := os.MkdirAll(actualOutDir, 0755); err != nil {
				fatalf("could not create output directory %q: %v", actualOutDir, err)
			}
			if len(files) == 1 {
				file := files[0]
				ext := filepath.Ext(file)
				var stem string
				if ext != "" {
					stem = file[:len(file)-len(ext)]
				} else {
					stem = file
				}
				inputFile := filepath.Join(dir, file)
				outputFile := filepath.Join(actualOutDir, stem+".jsx")
				if art, err := restache.TranspileFile(inputFile, outputFile); err != nil {
					fatal(err.Error())
				} else {
					logArtifact(art)
				}
				continue
			}
			if _, err := restache.TranspileModule(dir, actualOutDir,
				restache.WithIncludes(files),
				restache.WithParallelism(parallelism),
				restache.WithCallback(logArtifact),
			); err != nil {
				fatal(err.Error())
			}
		}
	} else {
		for dir, files := range filesByDir {
			if len(files) == 1 {
				file := files[0]
				if art, err := restache.TranspileFile(filepath.Join(dir, file), ""); err != nil {
					fatal(err.Error())
				} else {
					logArtifact(art)
				}
				continue
			}
			if _, err := restache.TranspileModule(dir, "",
				restache.WithIncludes(files),
				restache.WithParallelism(parallelism),
				restache.WithCallback(logArtifact),
			); err != nil {
				fatal(err.Error())
			}
		}
	}
}

func collectKeys[K comparable, V any](m map[K]V) (keys []K) {
	for key := range m {
		keys = append(keys, key)
	}
	return
}

func resolveGlobs(baseDir string, patterns []string) (dirs map[string][]string) {
	var (
		info os.FileInfo
		dir  string
	)
	dirs = make(map[string][]string)
	for _, pat := range patterns {
		var ppat string
		if filepath.IsAbs(pat) {
			ppat = pat
		} else {
			ppat = filepath.Join(baseDir, pat)
		}
		matches, err := filepath.Glob(ppat)
		if err != nil {
			fatalf("invalid glob %q: %v", pat, err)
		}
		if len(matches) == 0 {
			logf("no files matched the pattern: %s", pat)
			continue
		}
		for _, match := range matches {
			info, err = os.Lstat(match) // will ignore symlinks
			if err != nil {
				fatalf("could not access file %q: %v", match, err)
			}
			if info.IsDir() {
				continue
			}
			dir = filepath.Dir(match)
			dirs[dir] = append(dirs[dir], info.Name())
		}
	}
	return
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", PROGRAM_NAME, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", PROGRAM_NAME, fmt.Sprintf(format, args...))
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func logArtifact(art restache.Artifact) {
	logf("%d bytes written to %s (%.2fms) at %s",
		art.Bytes, art.Path, float64(art.Elapsed.Microseconds())/1000, time.Now().Format("15:04:05"))
}
