package settings

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	_ "modernc.org/sqlite"
)

const hfTokenSettingKey = "hf_token"

// SettingsStore persists user-configurable integration settings (currently
// just the HuggingFace API token) in a local SQLite database. Secrets are
// encrypted at rest with AES-256-GCM under masterKey; the plaintext is never
// written to disk.
type SettingsStore struct {
	db        *sql.DB
	masterKey []byte
}

// OpenSettingsStore opens (creating if needed) the settings database at
// dbPath. masterKey must be exactly 32 bytes (AES-256).
func OpenSettingsStore(dbPath string, masterKey []byte) (*SettingsStore, error) {
	if len(masterKey) != 32 {
		return nil, fmt.Errorf("mantle: master key must be 32 bytes, got %d", len(masterKey))
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("mantle: opening settings db: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("mantle: creating settings table: %w", err)
	}

	return &SettingsStore{db: db, masterKey: masterKey}, nil
}

// Close closes the underlying database connection.
func (s *SettingsStore) Close() error {
	return s.db.Close()
}

// SetHFToken encrypts and persists the HuggingFace API token, overwriting
// any previously stored value.
func (s *SettingsStore) SetHFToken(token string) error {
	ciphertext, err := s.encrypt(token)
	if err != nil {
		return fmt.Errorf("mantle: encrypting HF token: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, hfTokenSettingKey, ciphertext)
	if err != nil {
		return fmt.Errorf("mantle: storing HF token: %w", err)
	}
	return nil
}

// GetHFToken returns the decrypted HuggingFace API token, or "" if none has
// been set.
func (s *SettingsStore) GetHFToken() (string, error) {
	var ciphertext string
	err := s.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, hfTokenSettingKey).Scan(&ciphertext)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("mantle: reading HF token: %w", err)
	}

	plaintext, err := s.decrypt(ciphertext)
	if err != nil {
		return "", fmt.Errorf("mantle: decrypting HF token: %w", err)
	}
	return plaintext, nil
}

// ClearHFToken removes the stored HuggingFace API token, if any.
func (s *SettingsStore) ClearHFToken() error {
	_, err := s.db.Exec(`DELETE FROM settings WHERE key = ?`, hfTokenSettingKey)
	if err != nil {
		return fmt.Errorf("mantle: clearing HF token: %w", err)
	}
	return nil
}

// encrypt returns a base64-encoded "nonce || ciphertext" blob for plaintext,
// sealed with AES-256-GCM under s.masterKey.
func (s *SettingsStore) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// decrypt reverses encrypt. It fails if encoded was not sealed with the same
// s.masterKey (e.g. wrong master key, or tampered ciphertext).
func (s *SettingsStore) decrypt(encoded string) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decoding ciphertext: %w", err)
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(sealed) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := sealed[:nonceSize], sealed[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (wrong master key or corrupted data): %w", err)
	}
	return string(plaintext), nil
}
