package store

import "github.com/vector-10/url-shortner/internal/models"


type Store interface {
 Save(record *models.URLRecord) error
 GetBySlug(slug string) (*models.URLRecord, error)
 IncrementClicks(slug string) error
 Delete(slug string) error
 ListByUser(userID string) ([]*models.URLRecord, error)

}



