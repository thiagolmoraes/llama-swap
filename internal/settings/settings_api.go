package settings

import (
	"encoding/json"
	"net/http"

	"github.com/mostlygeek/llama-swap/internal/chain"
)

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// SettingsHandler exposes HTTP endpoints for user-configurable integration
// settings (currently the HuggingFace API token).
type SettingsHandler struct {
	store *SettingsStore
}

// NewSettingsHandler creates a settings API handler backed by store.
func NewSettingsHandler(store *SettingsStore) *SettingsHandler {
	return &SettingsHandler{store: store}
}

// RegisterRoutes adds the settings endpoints to mux, wrapped by authChain so
// they require the same API key as the rest of the app's custom endpoints.
func (h *SettingsHandler) RegisterRoutes(mux *http.ServeMux, authChain chain.Chain) {
	mux.Handle("PUT /api/settings/hf-token", authChain.ThenFunc(h.handlePutHFToken))
	mux.Handle("GET /api/settings/hf-token", authChain.ThenFunc(h.handleGetHFToken))
	mux.Handle("DELETE /api/settings/hf-token", authChain.ThenFunc(h.handleDeleteHFToken))
}

type putHFTokenRequest struct {
	Token string `json:"token"`
}

func (h *SettingsHandler) handlePutHFToken(w http.ResponseWriter, r *http.Request) {
	var req putHFTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Token == "" {
		jsonError(w, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.store.SetHFToken(req.Token); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"msg": "saved"})
}

// handleGetHFToken reports only whether a token is configured — the
// plaintext value is never returned once stored.
func (h *SettingsHandler) handleGetHFToken(w http.ResponseWriter, r *http.Request) {
	token, err := h.store.GetHFToken()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]bool{"configured": token != ""})
}

func (h *SettingsHandler) handleDeleteHFToken(w http.ResponseWriter, r *http.Request) {
	if err := h.store.ClearHFToken(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"msg": "cleared"})
}
