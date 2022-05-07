package router

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/left/api"
	"github.com/go-chi/chi/v5"
)

// mountPresets mounts all presets from the given preset store.
func mountPresets(r chi.Router, app dto.App) {
	res, err := app.PresetList(context.Background())
	if err != nil {
		log.Fatalln("router.mountPresets:", err)
	}

	for _, p := range res.Presets {
		u, _ := url.Parse(p.URL)
		route, url := u.Path, p.URL

		if err := validRoute(u.Path); err != nil {
			log.Fatalf("router.mountPresets: URL=%s route=%s: %s", url, route, err)
		}

		r.Get(route, api.PresetGetURLNew(app, url))
		log.Println("router.mountPresets: mounting url", url, "to", route)
	}
}

// validRoute returns nil if the given chi route is valid.
func validRoute(route string) error {
	if !strings.HasPrefix(route, "/") {
		return fmt.Errorf("route must start with /")
	}

	return nil
}

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
	} else if err != fs.ErrNotExist {
		log.Fatalln("router.mountFS:", err)
	}
}

// indexGet returns index.html from the given filesystem.
func indexGet(httpFS http.FileSystem) http.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		log.Fatalln("router.indexGet:", err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatalln("router.indexGet:", err)
	}

	modtime := stat.ModTime()

	return func(rw http.ResponseWriter, r *http.Request) {
		http.ServeContent(rw, r, "index.html", modtime, index)
	}
}
