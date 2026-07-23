package settings

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func containsBytes(haystack, needle []byte) bool {
	return bytes.Contains(haystack, needle)
}

func newTestSettingsStore(t *testing.T) *SettingsStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "settings.db")
	store, err := OpenSettingsStore(dbPath, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("OpenSettingsStore() error = %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSettingsStore_SetAndGetHFToken(t *testing.T) {
	store := newTestSettingsStore(t)

	if err := store.SetHFToken("hf_abc123secrettoken"); err != nil {
		t.Fatalf("SetHFToken() error = %v", err)
	}

	got, err := store.GetHFToken()
	if err != nil {
		t.Fatalf("GetHFToken() error = %v", err)
	}
	if got != "hf_abc123secrettoken" {
		t.Errorf("GetHFToken() = %q, want %q", got, "hf_abc123secrettoken")
	}
}

func TestSettingsStore_GetHFToken_NotSet(t *testing.T) {
	store := newTestSettingsStore(t)

	got, err := store.GetHFToken()
	if err != nil {
		t.Fatalf("GetHFToken() error = %v", err)
	}
	if got != "" {
		t.Errorf("GetHFToken() = %q, want empty string when unset", got)
	}
}

func TestSettingsStore_SetHFToken_Overwrites(t *testing.T) {
	store := newTestSettingsStore(t)

	if err := store.SetHFToken("hf_first"); err != nil {
		t.Fatalf("SetHFToken(first) error = %v", err)
	}
	if err := store.SetHFToken("hf_second"); err != nil {
		t.Fatalf("SetHFToken(second) error = %v", err)
	}

	got, err := store.GetHFToken()
	if err != nil {
		t.Fatalf("GetHFToken() error = %v", err)
	}
	if got != "hf_second" {
		t.Errorf("GetHFToken() = %q, want %q (overwrite)", got, "hf_second")
	}
}

func TestSettingsStore_ClearHFToken(t *testing.T) {
	store := newTestSettingsStore(t)

	if err := store.SetHFToken("hf_toremove"); err != nil {
		t.Fatalf("SetHFToken() error = %v", err)
	}
	if err := store.ClearHFToken(); err != nil {
		t.Fatalf("ClearHFToken() error = %v", err)
	}

	got, err := store.GetHFToken()
	if err != nil {
		t.Fatalf("GetHFToken() error = %v", err)
	}
	if got != "" {
		t.Errorf("GetHFToken() after clear = %q, want empty", got)
	}
}

// TestSettingsStore_TokenStoredEncrypted is the security regression test:
// the plaintext token must never appear verbatim in the on-disk database
// file, only its ciphertext.
func TestSettingsStore_TokenStoredEncrypted(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "settings.db")
	store, err := OpenSettingsStore(dbPath, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("OpenSettingsStore() error = %v", err)
	}

	const plaintext = "hf_supersecrettoken_shouldnotappearraw"
	if err := store.SetHFToken(plaintext); err != nil {
		t.Fatalf("SetHFToken() error = %v", err)
	}
	store.Close()

	raw, err := readFileBytes(dbPath)
	if err != nil {
		t.Fatalf("reading db file: %v", err)
	}
	if containsBytes(raw, []byte(plaintext)) {
		t.Errorf("plaintext HF token found verbatim in settings.db — must be encrypted at rest")
	}
}

// TestSettingsStore_WrongMasterKey_FailsToDecrypt ensures a store opened
// with a different master key cannot recover a previously stored secret
// (proves the value is actually encrypted with that key, not just obscured).
func TestSettingsStore_WrongMasterKey_FailsToDecrypt(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "settings.db")

	store1, err := OpenSettingsStore(dbPath, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("OpenSettingsStore() error = %v", err)
	}
	if err := store1.SetHFToken("hf_original"); err != nil {
		t.Fatalf("SetHFToken() error = %v", err)
	}
	store1.Close()

	store2, err := OpenSettingsStore(dbPath, []byte("ffffffffffffffffffffffffffffffff"))
	if err != nil {
		t.Fatalf("OpenSettingsStore() with different key error = %v", err)
	}
	defer store2.Close()

	if _, err := store2.GetHFToken(); err == nil {
		t.Error("GetHFToken() with wrong master key succeeded, want decryption error")
	}
}

func TestOpenSettingsStore_RejectsShortMasterKey(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "settings.db")
	_, err := OpenSettingsStore(dbPath, []byte("tooshort"))
	if err == nil {
		t.Error("OpenSettingsStore() with a short master key succeeded, want error (AES-256 needs a 32-byte key)")
	}
}
