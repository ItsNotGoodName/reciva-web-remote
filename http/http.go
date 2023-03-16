package http

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

type RadioContext struct {
	echo.Context
	Radio hub.Radio
}

func Router(a API, fs fs.FS) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	swagger(e)

	{
		api := e.Group("/api")

		api.GET("/build", a.GetBuild)
		api.GET("/presets", a.ListPresets)
		api.POST("/presets", a.UpdatePreset)
		api.GET("/presets/*", a.GetPreset)

		api.POST("/radios", a.DiscoverRadios)
		api.GET("/radios", a.ListRadios)

		apiRadios := api.Group("/radios/:uuid")
		apiRadios.Use(a.RadioMiddleware)
		apiRadios.GET("", a.GetRadio)
		apiRadios.DELETE("", a.DeleteRadio)
		apiRadios.POST("/volume", a.RefreshRadioVolume)
		apiRadios.POST("/subscription", a.RefreshRadioSubscription)

		api.GET("/states", a.ListStates)

		apiStates := api.Group("/states/:uuid")
		apiStates.Use(a.RadioMiddleware)
		apiStates.GET("", a.GetState)
		apiStates.POST("", a.PostState)

		api.GET("/ws", a.WS(upgrader()))
	}

	mountFS(e, fs)
	mountPresets(e, a.Store)

	return e
}

func Start(e *echo.Echo, port int) {
	printAddresses(strconv.Itoa(port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func upgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}

// mountPresets mounts all presets from the given preset store.
func mountPresets(e *echo.Echo, store store.Store) {
	res, err := store.ListPresets(context.Background())
	if err != nil {
		log.Fatalln("http.mountPresets:", err)
	}

	for _, p := range res {
		u, _ := url.Parse(p.URL)
		route, url := u.Path, p.URL

		if err := validRoute(u.Path); err != nil {
			log.Fatalf("http.mountPresets: URL=%s route=%s: %s", url, route, err)
		}

		e.GET(route, GetPresetURLNew(store, url))
		log.Println("http.mountPresets: mounting url", url, "to", route)
	}
}

func GetPresetURLNew(store store.Store, url string) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := store.GetPreset(c.Request().Context(), url)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, res.URLNew)
	}
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

// validRoute returns nil if the given route is valid.
func validRoute(route string) error {
	if !strings.HasPrefix(route, "/") {
		return fmt.Errorf("route must start with /")
	}

	return nil
}

// mountFS adds GET handlers for all files and folders using the given filesystem.
func mountFS(e *echo.Echo, f fs.FS) {
	httpFS := http.FS(f)
	fsHandler := func(c echo.Context) error {
		http.StripPrefix("/", http.FileServer(httpFS)).ServeHTTP(c.Response(), c.Request())
		return nil
	}

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				e.GET("/"+name+"/*", fsHandler)
			} else if name == "index.html" {
				indexHandler := indexGet(httpFS)
				e.GET("/", indexHandler)
				e.GET("/index.html", indexHandler)
			} else {
				e.GET("/"+name, fsHandler)
			}
		}
	} else if err != fs.ErrNotExist {
		log.Fatalln("http.mountFS:", err)
	}
}

// indexGet returns index.html from the given filesystem.
func indexGet(httpFS http.FileSystem) echo.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		log.Fatalln("http.indexGet:", err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatalln("http.indexGet:", err)
	}

	modtime := stat.ModTime()

	return func(c echo.Context) error {
		http.ServeContent(c.Response(), c.Request(), "http.html", modtime, index)
		return nil
	}
}
