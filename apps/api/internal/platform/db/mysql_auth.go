package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLAuthStore struct {
	db *sql.DB
}

type AuthUser struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
}

func NewMySQLAuthStore(dsn string) (*MySQLAuthStore, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	store := &MySQLAuthStore{db: conn}
	if err := store.ensureSchema(); err != nil {
		return nil, err
	}
	return store, nil
}

func randomUserID() (string, error) {
	var buf [10]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return "u_" + hex.EncodeToString(buf[:]), nil
}

func (s *MySQLAuthStore) ensureSchema() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(64) PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  keywords TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`)
	if err != nil {
		return err
	}
	if err := s.ensureProfileColumns(); err != nil {
		return err
	}
	return s.ensureJobSchema()
}

func (s *MySQLAuthStore) CreateUser(email, username, passwordHash string, keywords []string) (AuthUser, error) {
	id, err := randomUserID()
	if err != nil {
		return AuthUser{}, err
	}
	keywordCSV := strings.Join(keywords, ",")
	_, err = s.db.Exec(
		"INSERT INTO users (id, email, username, password_hash, keywords) VALUES (?, ?, ?, ?, ?)",
		id,
		email,
		username,
		passwordHash,
		keywordCSV,
	)
	if err != nil {
		return AuthUser{}, err
	}
	return AuthUser{ID: id, Email: email, Username: username, PasswordHash: passwordHash}, nil
}

func (s *MySQLAuthStore) FindUserByEmail(email string) (AuthUser, bool, error) {
	row := s.db.QueryRow(
		"SELECT id, email, username, password_hash FROM users WHERE LOWER(TRIM(email)) = ? LIMIT 1",
		email,
	)
	var user AuthUser
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthUser{}, false, nil
		}
		return AuthUser{}, false, err
	}
	return user, true, nil
}
