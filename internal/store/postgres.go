package store

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/vector-10/url-shortner/internal/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}


func (p *PostgresStore) Save(record *models.URLRecord) error {
	_, err := p.db.Exec(`
	    INSERT INTO url_records (id, slug, long_url, user_id, created_at, expires_at, is_active, max_clicks, total_clicks, link_type)
	VALUES
		 ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		record.ID,
		record.Slug,
		record.LongURL,
		record.UserID,
		record.CreatedAt,
		record.ExpiresAt,
		record.IsActive,
		record.MaxClicks,
		record.TotalClicks,
		record.LinkType,
	)
	return err
}

func (p *PostgresStore) GetBySlug(slug string) (*models.URLRecord, error) {
	row := p.db.QueryRow(`
	SELECT id, slug, long_url, user_id, created_at, expires_at, is_active, max_clicks, total_clicks,
	link_type FROM url_records WHERE slug = $1`, slug)

	var record models.URLRecord
	err := row.Scan(
		&record.ID,
		&record.Slug,
		&record.LongURL,
		&record.UserID,
		&record.CreatedAt,
		&record.ExpiresAt,
		&record.IsActive,
		&record.MaxClicks,
		&record.TotalClicks,
		&record.LinkType,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("slug not found")
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (p *PostgresStore) IncrementClicks(slug string) error {
	_, err := p.db.Exec(`UPDATE url_records SET total_clicks = total_clicks + 1 WHERE slug = $1`, slug)
	return err
}

func (p *PostgresStore) Delete(slug string) error {
	_, err := p.db.Exec(`DELETE FROM url_records WHERE slug = $1`, slug)
	return err
}

func (p *PostgresStore) ListByUser(userID string) ([]*models.URLRecord, error) {
	rows, err := p.db.Query(`SELECT id, slug, long_url, user_id, created_at, expires_at, is_active, max_clicks, total_clicks, link_type
	FROM url_records WHERE user_id = $1 ORDER BY created_at DESC`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.URLRecord
	for rows.Next() {
		var record models.URLRecord
		if err := rows.Scan(
			&record.ID,
			&record.Slug,
			&record.LongURL,
			&record.UserID,
			&record.CreatedAt,
			&record.ExpiresAt,
			&record.IsActive,
			&record.MaxClicks,
			&record.TotalClicks,
			&record.LinkType,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if records == nil {
		records = []*models.URLRecord{}
	}
	return records, nil
}

func (p *PostgresStore) LogClickEvent(event *models.ClickEvent) error {
	_, err := p.db.Exec(`
	INSERT INTO click_events (slug, clicked_at, ip_address, user_agent, was_valid, rejection_reason)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		event.Slug,
		time.Now(),
		event.IPAddress,
		event.UserAgent,
		event.WasValid,
		event.RejectionReason)
	return err
}

func (p *PostgresStore) DeactivateSlug(slug string) error {
	_, err := p.db.Exec(`UPDATE url_records SET is_active = false WHERE slug = $1`, slug)
	return err
}
