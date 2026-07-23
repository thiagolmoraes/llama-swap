package settings

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

const masterKeySize = 32 // AES-256

// LoadOrCreateMasterKey returns the AES-256 master key used to encrypt
// settings at rest, reading it from path if present or generating and
// persisting a new random one (mode 0600, owner-only) otherwise. Calling it
// again with the same path always returns the same key, so restarts keep
// access to previously encrypted secrets.
func LoadOrCreateMasterKey(path string) ([]byte, error) {
	existing, err := os.ReadFile(path)
	if err == nil {
		if len(existing) != masterKeySize {
			return nil, fmt.Errorf("mantle: master key file %s has %d bytes, want %d (corrupt or wrong-version file)", path, len(existing), masterKeySize)
		}
		return existing, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("mantle: reading master key file %s: %w", path, err)
	}

	key := make([]byte, masterKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("mantle: generating master key: %w", err)
	}

	if err := os.WriteFile(path, key, 0o600); err != nil {
		return nil, fmt.Errorf("mantle: writing master key file %s: %w", path, err)
	}

	return key, nil
}
