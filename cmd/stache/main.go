package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/tetsuo/stache"
	"github.com/tetsuo/stache/jsx"
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
			fmt.Fprintf(os.Stderr, "invalid option: %s (did you mean '--%s'?)\n", arg, arg[1:])
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
			fmt.Fprintln(os.Stderr, "WARNING: --outdir option is ignored (no input files)")
		}
		node, err := stache.Parse(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		src := jsx.NewReader(node)
		if _, err = io.Copy(os.Stdout, src); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dirIncludes, err := resolveGlobs(baseDir, patterns)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(dirIncludes) == 0 {
		fmt.Fprintln(os.Stderr, "no files matched the given patterns")
		os.Exit(1)
	}

	for dir, includes := range dirIncludes {
		if len(includes) == 1 {
			_, err := stache.ParseFile(filepath.Join(dir, includes[0]))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error parsing file %q: %v\n", includes[0], err)
				os.Exit(1)
			}
			continue
		}
		_, err := stache.ParseDir(dir, includes, stache.WithParallelism(parallelism))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing dir %q: %v\n", dir, err)
			os.Exit(1)
		}
	}
}

func resolveGlobs(baseDir string, patterns []string) (dirs map[string][]string, err error) {
	var (
		matches []string
		info    os.FileInfo
		dir     string
		p       string
	)
	dirs = make(map[string][]string, 0)
	for _, p = range patterns {
		matches, err = filepath.Glob(filepath.Join(baseDir, p))
		if err != nil {
			err = fmt.Errorf("failed to match files with %q: %w", p, err)
			return
		}
		for _, p = range matches {
			info, err = os.Lstat(p) // will ignore symlinks
			if err != nil {
				err = fmt.Errorf("could not lstat file %q: %w", p, err)
				return
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

}
