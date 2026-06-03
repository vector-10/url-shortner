package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/vector-10/url-shortner/internal/models"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (p *PostgresUserStore) CreateUser(userID string) error {
	_, err := p.db.Exec(`
	INSERT INTO users (id, email, password_hash, created_at)
	VALUES($1, $2, $3, $4)`,
		user.ID,
		user.Email,
		user.Passwordhash,
		user.createdAt,
	)
	return err
}

func (p *PostgresUserStore) GetUserByEmail(email string) (*models.User, error) {
	row := p.db.QueryRow(`
	SELECT id, email, password_hash, created_at FROm users WHERE email = $1`, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Passwordhash,
		&user.createdAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (p *PostgresUserStore) GetUserByID(userID string) (*models.User, error) {
	row := p.db.QueryRow(`
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE id = $1`, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Passwordhash,
		&user.createdAt,
	)
	if err != sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

