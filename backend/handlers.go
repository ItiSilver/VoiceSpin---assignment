package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type API struct {
	store Store
}

func NewAPI(store Store) *API {
	return &API{store: store}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/conversations", a.listConversations)
	mux.HandleFunc("GET /api/conversations/{id}", a.getConversation)
	mux.HandleFunc("PATCH /api/conversations/{id}", a.patchConversation)
	return withCORS(mux)
}

func (a *API) listConversations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ListFilter{
		Status:   q.Get("status"),
		Priority: q.Get("priority"),
		Search:   q.Get("search"),
	}
	writeJSON(w, http.StatusOK, a.store.List(filter))
}

func (a *API) getConversation(w http.ResponseWriter, r *http.Request) {
	c, err := a.store.Get(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "conversation not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type patchBody struct {
	Status   *Status   `json:"status"`
	Priority *Priority `json:"priority"`
}

func (a *API) patchConversation(w http.ResponseWriter, r *http.Request) {
	var body patchBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if body.Status == nil && body.Priority == nil {
		writeError(w, http.StatusBadRequest, "provide at least one of: status, priority")
		return
	}
	if body.Status != nil && !body.Status.Valid() {
		writeError(w, http.StatusBadRequest, "invalid status value")
		return
	}
	if body.Priority != nil && !body.Priority.Valid() {
		writeError(w, http.StatusBadRequest, "invalid priority value")
		return
	}

	c, err := a.store.Update(r.PathValue("id"), UpdatePatch{
		Status:   body.Status,
		Priority: body.Priority,
	})
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "conversation not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
