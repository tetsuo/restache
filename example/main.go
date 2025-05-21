package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	addrFlag    = flag.String("addr", "localhost:7070", "host and port to bind the server to")
	staticFlag  = flag.String("static", "static", "path to the folder containing static files to serve")
	devModeFlag = flag.Bool("dev", false, "enable development mode")
)

func main() {
	flag.Parse()

	staticHandler := StaticHandler(*staticFlag, "/static/", *devModeFlag)

	server := &http.Server{
		Addr: *addrFlag,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path == "/" {
				printHTML(w)
			} else if strings.HasPrefix(r.URL.Path, "/static/") {
				staticHandler.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
		}),
	}

	log.Printf("listening on %s", server.Addr)
	log.Print(server.ListenAndServe())
}

func StaticHandler(staticPath, prefix string, devMode bool) http.Handler {
	staticPathFound, err := fileExists(staticPath)
	if err != nil {
		log.Fatal(err)
	}
	if !staticPathFound {
		log.Fatalf("directory %q doesn't exist", staticPath)
	}

	var staticHandler http.Handler

	if devMode {
		staticHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := Build(Config{EntryPoint: staticPath, Watch: false, Bundle: true})
			if err != nil {
				log.Fatal(err)
			}
			http.StripPrefix(prefix, http.FileServer(http.FS(os.DirFS(staticPath)))).ServeHTTP(w, r)
		})
	} else {
		staticHandler = http.StripPrefix(prefix, http.FileServer(http.FS(os.DirFS(staticPath))))
	}

	return staticHandler
}

func printHTML(w io.Writer) {
	const s = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<script type="importmap">
{
	"imports": {
		"react": "https://esm.sh/react?bundle",
		"react-dom/client": "https://esm.sh/react-dom/client?bundle"
	}
}
</script>
</head>
<body>
<div id="root"></div>
<script type="module" src="/static/index.js"></script>
</body>
</html>
`
	io.WriteString(w, s)
}

func fileExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		exists = true
		return
	}
	exists = false
	if os.IsNotExist(err) {
		err = nil
		return
	}
	return
}
