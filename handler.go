package main

import (
	"context"
	"net/http"

	"github.com/Tomoka64/stock-notifier/mongo"
)

type Handler struct {
	mongoClient *mongo.Client
}

func newHandler() *Handler {
	return &Handler{
		mongoClient: mongo.New(),
	}
}

func (*Handler) index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	email := r.FormValue("email")
	err := h.mongoClient.Insert(ctx, email)
	if err != nil {
		tpl.ExecuteTemplate(w, "error.html", err)
		return
	}
	tpl.ExecuteTemplate(w, "result.html", nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	list, err := h.mongoClient.List(ctx)
	if err != nil {
		tpl.ExecuteTemplate(w, "error.html", err)
		return
	}
	tpl.ExecuteTemplate(w, "list.html", list)
}
