package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"bookapp/internal/models"

	"github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxOpenConns(10)
	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) Create(u *models.User) error {
	// provider 'local' or 'google', etc.
	_, err := p.db.ExecContext(context.Background(), `
      INSERT INTO users (email, hashed_password, provider, provider_id)
      VALUES ($1, $2, $3, $4)
    `, u.Email, u.HashedPassword, u.Provider, u.ProviderID)
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" {
		return ErrUserExists
	}
	return err
}

func (p *PostgresStore) FindByEmail(email string) (*models.User, error) {
	row := p.db.QueryRowContext(context.Background(), `
      SELECT id, email, hashed_password, provider, provider_id
      FROM users WHERE email = $1
    `, email)

	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Provider, &u.ProviderID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (p *PostgresStore) Update(u *models.User) error {
	_, err := p.db.ExecContext(context.Background(), `
      UPDATE users 
      SET hashed_password = $1, provider = $2, provider_id = $3, updated_at = now()
      WHERE email = $4
    `, u.HashedPassword, u.Provider, u.ProviderID, u.Email)
	return err
}

// SaveRefreshToken stores a new refresh token for the user.
func (p *PostgresStore) SaveRefreshToken(token string, userID int, expiresAt time.Time) error {
	_, err := p.db.ExecContext(context.Background(), `
		INSERT INTO refresh_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, userID, expiresAt)
	return err
}

// DeleteRefreshToken removes a refresh token (used for rotation or logout).
func (p *PostgresStore) DeleteRefreshToken(token string) error {
	_, err := p.db.ExecContext(context.Background(), `
		DELETE FROM refresh_tokens WHERE token = $1
	`, token)
	return err
}

// FindRefreshToken checks if a token exists in the DB and isn't expired.
func (p *PostgresStore) FindRefreshToken(token string) (int, error) {
	var userID int
	err := p.db.QueryRowContext(context.Background(), `
		SELECT user_id FROM refresh_tokens
		WHERE token = $1 AND expires_at > now()
	`, token).Scan(&userID)

	if err == sql.ErrNoRows {
		return 0, errors.New("invalid or expired refresh token")
	}
	return userID, err
}

// FindByID retrieves a user by their primary key.
func (p *PostgresStore) FindByID(id int) (*models.User, error) {
	row := p.db.QueryRowContext(context.Background(), `
		SELECT id, email, hashed_password, provider, provider_id
		FROM users WHERE id = $1
	`, id)

	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Provider, &u.ProviderID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

// SaveVerification stores a pending signup.
func (p *PostgresStore) SaveVerification(email, hashedPw, code string, expires time.Time) error {
	_, err := p.db.ExecContext(context.Background(), `
    INSERT INTO email_verifications (email, hashed_password, code, expires_at)
    VALUES ($1,$2,$3,$4)
    ON CONFLICT (email) DO UPDATE
      SET hashed_password = excluded.hashed_password,
          code = excluded.code,
          expires_at = excluded.expires_at,
          created_at = now()
  `, email, hashedPw, code, expires)
	return err
}

// GetVerification fetches pending signup by email+code.
func (p *PostgresStore) GetVerification(email, code string) (hashedPw string, err error) {
	// log.Printf("DEBUG: GetVerification called with email='%s', code='%s'", email, code)

	row := p.db.QueryRowContext(context.Background(), `
    SELECT hashed_password FROM email_verifications
    WHERE email=$1 AND code=$2 AND expires_at>now()
  `, email, code)
	if err := row.Scan(&hashedPw); err != nil {
		return "", err
	}
	// log.Printf("DEBUG: GetVerification succeeded, returning hashed password")
	return hashedPw, nil
}

// DeleteVerification deletes after successful verify.
func (p *PostgresStore) DeleteVerification(email string) error {
	_, err := p.db.ExecContext(context.Background(), `
    DELETE FROM email_verifications WHERE email=$1
  `, email)
	return err
}

// func StartCleanupTasks(db *sql.DB) {
// 	go func() {
// 		for {
// 			_, _ = db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < now()`)
// 			_, _ = db.Exec(`DELETE FROM email_verification WHERE created_at < now() - interval '1 day'`)
// 			time.Sleep(24 * time.Hour)
// 		}
// 	}()
// }
