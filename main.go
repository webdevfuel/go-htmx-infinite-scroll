package main

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/webdevfuel/go-htmx-infinite-scroll/post"
	"github.com/webdevfuel/go-htmx-infinite-scroll/template"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := template.Hello("world")
		component.Render(r.Context(), w)
	})
	r.Get("/infinite-scroll", func(w http.ResponseWriter, r *http.Request) {
		var htmx = r.Header.Get("Hx-Request") == "true"
		var component templ.Component
		var err error
		if htmx {
			component, err = handleHTMXRequest(r)
		} else {
			component, err = handleInitialRequest()
		}
		if err != nil {
			return
		}
		err = component.Render(r.Context(), w)
		if err != nil {
			return
		}
	})
	http.ListenAndServe("localhost:3000", r)
}

func getCursorOrFallback(r *http.Request, key string, fallback int) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return fallback, nil
	}
	cursor, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return cursor, nil
}

func handleHTMXRequest(r *http.Request) (templ.Component, error) {
	cursor, err := getCursorOrFallback(r, "cursor", 5)
	if err != nil {
		return nil, err
	}
	posts, nextCursor, err := post.GetPostsByCursorAndLimit(cursor, 4)
	if err != nil {
		return nil, err
	}
	component := template.PostsAndButton(posts, nextCursor)
	return component, nil
}

func handleInitialRequest() (templ.Component, error) {
	cursor, err := post.GetFirstCursor()
	if err != nil {
		return nil, err
	}
	posts, _, err := post.GetPostsByCursorAndLimit(cursor, 4)
	if err != nil {
		return nil, err
	}
	component := template.InfiniteScroll(posts)
	return component, nil
}
