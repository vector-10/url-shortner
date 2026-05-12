package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/vector-10/url-shortner/internal/models"
	"github.com/vector-10/url-shortner/internal/store"
)

// this handler is the layer where HTTP requests come in and are processed
type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r*http.Request) {
	var record models.URLRecord

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if record.LongURL == "" {
		http.Error(w, "long_url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(record.LongURL); err != nil {
		http.Error(w, "long_url is not a valid URL", http.StatusBadRequest)
		return
	}

	record.ID = uuid.New().String()
	record.CreatedAt = time.Now()
	record.Clicks = 0

	if record.Slug == "" {
		record.Slug = generateSlug()
	}

	if err := h.store.Save(&record); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	record, err := h.store.GetBySlug(slug)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if err := h.store.IncrementClicks(slug); err != nil {
		http.Error(w, "could not increment clicks", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, record.LongURL, http.StatusMovedPermanently)
}

func (h *Handler) Stats(w http. ResponseWriter, r*http.Request) {
	slug := r.PathValue("slug")

	record, err := h.store.GetBySlug(slug)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func (h *Handler) QRCode(w http.ResponseWriter, r*http.Request) {
	slug := r.PathValue("slug")

	_, err := h.store.GetBySlug(slug)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	shortURL := "http://localhost:8080/" + slug
	png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "could not generate QR code", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func generateSlug() string {
	return uuid.New().String()[:8]
}