package settings

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mostlygeek/llama-swap/internal/chain"
)

// noopAuthChain simulates the server's apiChain with a single required key,
// without importing internal/server (which itself imports internal/mantle
// and would create an import cycle).
func testAuthChain(t *testing.T, requiredKey string) chain.Chain {
	t.Helper()
	return chain.New(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+requiredKey {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func newTestSettingsHandler(t *testing.T) (*SettingsHandler, *SettingsStore) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "settings.db")
	store, err := OpenSettingsStore(dbPath, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("OpenSettingsStore() error = %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return NewSettingsHandler(store), store
}

func mustRegisterSettingsRoutes(t *testing.T, h *SettingsHandler, authChain chain.Chain) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux, authChain)
	return mux
}

func TestSettingsAPI_PutHFToken_RequiresAuth(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	body := strings.NewReader(`{"token":"hf_shouldnotpersist"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/settings/hf-token", body)
	// no Authorization header
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("PUT without auth: status = %d, want 401", w.Code)
	}
}

func TestSettingsAPI_GetHFToken_RequiresAuth(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	req := httptest.NewRequest(http.MethodGet, "/api/settings/hf-token", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("GET without auth: status = %d, want 401", w.Code)
	}
}

func TestSettingsAPI_DeleteHFToken_RequiresAuth(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	req := httptest.NewRequest(http.MethodDelete, "/api/settings/hf-token", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("DELETE without auth: status = %d, want 401", w.Code)
	}
}

func TestSettingsAPI_PutThenGet_RoundTrip(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	putBody := strings.NewReader(`{"token":"hf_myrealtoken"}`)
	putReq := httptest.NewRequest(http.MethodPut, "/api/settings/hf-token", putBody)
	putReq.Header.Set("Authorization", "Bearer secret-key")
	putW := httptest.NewRecorder()
	mux.ServeHTTP(putW, putReq)

	if putW.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want 200, body = %s", putW.Code, putW.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/settings/hf-token", nil)
	getReq.Header.Set("Authorization", "Bearer secret-key")
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200, body = %s", getW.Code, getW.Body.String())
	}

	var resp struct {
		Configured bool `json:"configured"`
	}
	if err := json.Unmarshal(getW.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decoding GET response: %v", err)
	}
	if !resp.Configured {
		t.Error("GET response configured = false, want true after PUT")
	}
}

// TestSettingsAPI_GetHFToken_NeverReturnsPlaintext is a security regression
// test: the GET endpoint must report whether a token is configured without
// ever echoing the secret value back in the response body.
func TestSettingsAPI_GetHFToken_NeverReturnsPlaintext(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	const secret = "hf_supersecretvalue"
	putBody := strings.NewReader(`{"token":"` + secret + `"}`)
	putReq := httptest.NewRequest(http.MethodPut, "/api/settings/hf-token", putBody)
	putReq.Header.Set("Authorization", "Bearer secret-key")
	mux.ServeHTTP(httptest.NewRecorder(), putReq)

	getReq := httptest.NewRequest(http.MethodGet, "/api/settings/hf-token", nil)
	getReq.Header.Set("Authorization", "Bearer secret-key")
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)

	if strings.Contains(getW.Body.String(), secret) {
		t.Errorf("GET response leaked plaintext token: %s", getW.Body.String())
	}
}

func TestSettingsAPI_PutHFToken_RejectsEmptyToken(t *testing.T) {
	h, _ := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	putBody := strings.NewReader(`{"token":""}`)
	putReq := httptest.NewRequest(http.MethodPut, "/api/settings/hf-token", putBody)
	putReq.Header.Set("Authorization", "Bearer secret-key")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, putReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PUT empty token: status = %d, want 400", w.Code)
	}
}

func TestSettingsAPI_DeleteHFToken_ClearsIt(t *testing.T) {
	h, store := newTestSettingsHandler(t)
	mux := mustRegisterSettingsRoutes(t, h, testAuthChain(t, "secret-key"))

	if err := store.SetHFToken("hf_toremove"); err != nil {
		t.Fatalf("SetHFToken() setup error = %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/settings/hf-token", nil)
	req.Header.Set("Authorization", "Bearer secret-key")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("DELETE status = %d, want 200", w.Code)
	}

	got, err := store.GetHFToken()
	if err != nil {
		t.Fatalf("GetHFToken() error = %v", err)
	}
	if got != "" {
		t.Errorf("token still present after DELETE: %q", got)
	}
}
