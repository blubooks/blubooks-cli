package router

import (
	"net/http"
	"strings"

	"github.com/blubooks/blubooks-cli/internal/pkg/server/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func New(server *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	/*	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
	*/

	//assetsDir := http.Dir("./assets/")
	//fileServer(r, "/assets", assetsDir)
	api := http.Dir("./public/api")
	fileServer(r, "/public/api", api)

	files := http.Dir("./public/files")
	fileServer(r, "/public/files", files)
	r.Get("/healthz", server.HanlderHealth)
	r.Get("/", server.HandleIndex)

	r.Get("/*", server.HandleClient)

	return r
}
