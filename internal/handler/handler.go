package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/vector-10/url-shortner/internal/models"
	"github.com/vector-10/url-shortner/internal/store"
)


type Handler struct {
	store store.Store
	cache *store.RedisCache
}

func NewHandler(s store.Store, c *store.RedisCache) *Handler {
	return &Handler{store: s, cache: c}
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

	parsed, err := url.ParseRequestURI(record.LongURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.Error(w, "long_url must be a valid http or https URL", http.StatusBadRequest)
		return
	}

	record.ID = uuid.New().String()
	record.CreatedAt = time.Now()
	record.TotalClicks = 0
	record.IsActive = true
	expiresAt := record.CreatedAt.Add(3 * time.Hour)
	record.ExpiresAt = &expiresAt
	if record.LinkType == "" {
		record.LinkType = "general"
	}

	userID, _ := r.Context().Value(UserIDKey).(string)
	record.UserID = userID

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

	record, err := h.cache.GetCachedSlug(slug)
	if err != nil {
		log.Printf("cache error for slug %s: %v", slug, err)
	}

	if record == nil {
		record, err = h.store.GetBySlug(slug)
		if err != nil {
			logClickEvent(h, slug, r, false, "slug not found")
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		h.cache.CacheSlug(record)
	}

	if !record.IsActive {
		logClickEvent(h, slug, r, false, "link inactive")
		http.Error(w, "The link has been used or is no longer active", http.StatusGone)
		return
	}

	if record.ExpiresAt != nil && time.Now().After(*record.ExpiresAt) {
		logClickEvent(h, slug, r, false, "link expired")
		http.Error(w, "The link has expired", http.StatusGone)
		return
	}

	if record.MaxClicks != nil && record.TotalClicks >= *record.MaxClicks {
        logClickEvent(h, slug, r, false, "max_clicks_reached")
        http.Error(w, "The link has reached its maximum click limit", http.StatusGone)
        return
    }

    logClickEvent(h, slug, r, true, "")

    if record.MaxClicks != nil && *record.MaxClicks == 1 {
        redeemed, err := h.store.RedeemSlug(slug)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if redeemed == nil {
			logClickEvent(h, slug, r, false, "slug already redeemed")
			http.Error(w, "The link has already been used", http.StatusGone)
			return
		}
		h.cache.InvalidateSlug(slug)
    }

    h.store.IncrementClicks(slug)
    http.Redirect(w, r, record.LongURL, http.StatusFound)
}

func logClickEvent(h *Handler, slug string, r *http.Request, wasvalid bool, reason string) {
	event := &models.ClickEvent{
		Slug:            slug,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		WasValid: wasvalid,
		RejectionReason: reason,
	}
	if err := h.store.LogClickEvent(event); err != nil {
		log.Printf("failed to log click event: %v", err)
	}
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

	shortURL := os.Getenv("BASE_URL") + "/" + slug
	png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "could not generate QR code", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func (h *Handler) ListURLs(w http.ResponseWriter, r*http.Request) {
	userID, _ := r.Context().Value(UserIDKey).(string)

	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	records, err := h.store.ListByUser(userID)

	if err != nil {
		http.Error(w, "could not fetch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func generateSlug() string {
	return uuid.New().String()[:8]
}