package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	addrFlag    = flag.String("addr", "localhost:7070", "host and port to bind the server to")
	staticFlag  = flag.String("static", "static", "path to the folder containing static files to serve")
	devModeFlag = flag.Bool("dev", false, "enable development mode")
	serveFlag   = flag.Bool("serve", false, "start the HTTP server")
)

func main() {
	flag.Parse()

	staticPath := *staticFlag
	devMode := *devModeFlag
	serve := *serveFlag

	staticPathFound, err := fileExists(staticPath)
	if err != nil {
		log.Fatal(err)
	}
	if !staticPathFound {
		log.Fatalf("directory %q doesn't exist", staticPath)
	}

	if !serve {
		// Build and exit
		err := Build(Config{
			SourceMap: devMode,
			Minify:    !devMode,
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Serve mode
	staticBasePrefix := "/"
	staticHandler := http.StripPrefix(staticBasePrefix, http.FileServer(http.FS(os.DirFS(staticPath))))

	server := &http.Server{
		Addr: *addrFlag,
	}

	if !devMode {
		// Build only once in production builds
		err := Build(Config{SourceMap: false, Minify: true})
		if err != nil {
			log.Fatal(err)
		}
		server.Handler = staticHandler
	} else {
		// Dev mode: build on every request to app.js
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/app.js" && r.Method == http.MethodGet {
				err := Build(Config{SourceMap: true, Minify: false})
				if err != nil {
					log.Println(err)
				}
			} else if strings.HasPrefix(r.URL.Path, "/api/") {
				response := map[string]any{"ok": true}
				b, err := json.Marshal(response)
				if err != nil {
					log.Println(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, err = w.Write(b)
				if err != nil {
					log.Println(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			staticHandler.ServeHTTP(w, r)
		})
	}

	log.Printf("listening on %s", server.Addr)
	log.Print(server.ListenAndServe())
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
