package router

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// printAddresses prints all listening addresses.
func printAddresses(port string) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("router.PrintAddresses:", err)
		return
	}

	message := "\nNavigate to one of the following addresses:\n"
	for i := range addr {
		ip := net.ParseIP(strings.Split(addr[i].String(), "/")[0])
		if ip != nil && ip.To4() != nil {
			message = message + "\thttp://" + ip.String() + ":" + port + "\n"
		}
	}
	fmt.Println(message)
}

// mountFS adds GET handlers for all files and folders using the given filesystem.
func mountFS(r chi.Router, f fs.FS) {
	httpFS := http.FS(f)
	fsHandler := http.StripPrefix("/", http.FileServer(httpFS))

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				r.Get("/"+name+"/*", fsHandler.ServeHTTP)
			} else if name == "index.html" {
				indexHandler := indexGet(httpFS)
				r.Get("/", indexHandler)
				r.Get("/index.html", indexHandler)
			} else {
				r.Get("/"+name, fsHandler.ServeHTTP)
			}
		}
	} else {
		log.Fatal("router.mountFS:", err)
	}
}

// indexGet returns index.html from the given filesystem.
func indexGet(httpFS http.FileSystem) http.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		log.Fatal("router.indexGet:", err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatal("router.indexGet:", err)
	}

	modtime := stat.ModTime()

	return func(rw http.ResponseWriter, r *http.Request) {
		http.ServeContent(rw, r, "index.html", modtime, index)
	}
}
