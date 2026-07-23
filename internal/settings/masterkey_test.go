package settings

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadOrCreateMasterKey_CreatesNewKeyWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "master.key")

	key, err := LoadOrCreateMasterKey(path)
	if err != nil {
		t.Fatalf("LoadOrCreateMasterKey() error = %v", err)
	}
	if len(key) != 32 {
		t.Errorf("LoadOrCreateMasterKey() returned %d bytes, want 32", len(key))
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected key file to be created at %s: %v", path, err)
	}
}

func TestLoadOrCreateMasterKey_PersistsAcrossReloads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "master.key")

	key1, err := LoadOrCreateMasterKey(path)
	if err != nil {
		t.Fatalf("first LoadOrCreateMasterKey() error = %v", err)
	}

	key2, err := LoadOrCreateMasterKey(path)
	if err != nil {
		t.Fatalf("second LoadOrCreateMasterKey() error = %v", err)
	}

	if string(key1) != string(key2) {
		t.Error("LoadOrCreateMasterKey() returned a different key on second call, want the same persisted key")
	}
}

func TestLoadOrCreateMasterKey_GeneratesUniqueKeysAcrossPaths(t *testing.T) {
	path1 := filepath.Join(t.TempDir(), "master.key")
	path2 := filepath.Join(t.TempDir(), "master.key")

	key1, err := LoadOrCreateMasterKey(path1)
	if err != nil {
		t.Fatalf("LoadOrCreateMasterKey(path1) error = %v", err)
	}
	key2, err := LoadOrCreateMasterKey(path2)
	if err != nil {
		t.Fatalf("LoadOrCreateMasterKey(path2) error = %v", err)
	}

	if string(key1) == string(key2) {
		t.Error("LoadOrCreateMasterKey() generated the same key for two independent paths, want randomness")
	}
}

func TestLoadOrCreateMasterKey_FilePermissionsAreOwnerOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix file permissions not applicable on windows")
	}
	path := filepath.Join(t.TempDir(), "master.key")

	if _, err := LoadOrCreateMasterKey(path); err != nil {
		t.Fatalf("LoadOrCreateMasterKey() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat key file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("key file permissions = %o, want 0600 (owner read/write only)", perm)
	}
}

func TestLoadOrCreateMasterKey_RejectsCorruptExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "master.key")
	if err := os.WriteFile(path, []byte("too-short-to-be-a-valid-key"), 0o600); err != nil {
		t.Fatalf("seeding corrupt key file: %v", err)
	}

	_, err := LoadOrCreateMasterKey(path)
	if err == nil {
		t.Error("LoadOrCreateMasterKey() with a corrupt/wrong-length key file succeeded, want error")
	}
}
