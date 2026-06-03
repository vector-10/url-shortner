package store

import (
    "database/sql"
    "errors"

    _ "github.com/lib/pq"
    "github.com/vector-10/url-shortner/internal/models"
)

type PostgresUserStore struct {
    db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
    return &PostgresUserStore{db: db}
}

func (p *PostgresUserStore) CreateUser(user *models.User) error {
    _, err := p.db.Exec(`
        INSERT INTO users (id, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)`,
        user.ID,
        user.Email,
        user.PasswordHash,
        user.CreatedAt,
    )
    return err
}

func (p *PostgresUserStore) GetUserByEmail(email string) (*models.User, error) {
    row := p.db.QueryRow(`
        SELECT id, email, password_hash, created_at
        FROM users
        WHERE email = $1`, email)

    var user models.User
    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (p *PostgresUserStore) GetUserByID(id string) (*models.User, error) {
    row := p.db.QueryRow(`
        SELECT id, email, password_hash, created_at
        FROM users
        WHERE id = $1`, id)

    var user models.User
    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}
