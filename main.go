package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lehoon/hook_api/v2/library/config"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/routes"
	md "github.com/lehoon/hook_api/v2/routes/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	route := chi.NewRouter()
	route.Use(middleware.RequestID)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(md.RequestLoggerFilter)
	route.Use(render.SetContentType(render.ContentTypeJSON))

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	route.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", routes.Routes())
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		routes.PushRoute(method, route)
		return nil
	}

	if err := chi.Walk(route, walkFunc); err != nil {
		fmt.Printf("Logging error: %s\n", err.Error())
	}

	logger.Log().Info("Hook api启动成功.")
	http.ListenAndServe(config.GetLocalAddress(), route)
}
