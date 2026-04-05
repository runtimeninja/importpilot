package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/runtimeninja/importpilot/internal/service"
)

type ClientHandler struct {
	clientService *service.ClientService
}

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

type createClientRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	ShopURL string `json:"shop_url"`
	Plan    string `json:"plan"`
}

type createClientResponse struct {
	Success bool   `json:"success"`
	ID      int64  `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type clientItem struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	ShopURL string `json:"shop_url"`
	Status  string `json:"status"`
	Plan    string `json:"plan"`
}

type listClientsResponse struct {
	Success bool         `json:"success"`
	Data    []clientItem `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

type getClientResponse struct {
	Success bool        `json:"success"`
	Data    *clientItem `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type updateClientStatusRequest struct {
	Status string `json:"status"`
}

type updateClientStatusResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (h *ClientHandler) HandleClients(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateClient(w, r)
	case http.MethodGet:
		h.ListClients(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
	}
}

func (h *ClientHandler) HandleClientByID(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/status") {
		h.UpdateClientStatus(w, r)
		return
	}

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/admin/clients/")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "client id is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid client id",
		})
		return
	}

	client, err := h.clientService.GetClientByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, getClientResponse{
			Success: false,
			Error:   "client not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, getClientResponse{
		Success: true,
		Data: &clientItem{
			ID:      client.ID,
			Name:    client.Name,
			Email:   client.Email,
			ShopURL: client.ShopURL,
			Status:  client.Status,
			Plan:    client.Plan,
		},
	})
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req createClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, createClientResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	id, err := h.clientService.CreateClient(r.Context(), service.CreateClientInput{
		Name:    req.Name,
		Email:   req.Email,
		ShopURL: req.ShopURL,
		Plan:    req.Plan,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, createClientResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, createClientResponse{
		Success: true,
		ID:      id,
	})
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if rawPage := r.URL.Query().Get("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit

	clients, err := h.clientService.ListClients(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, listClientsResponse{
			Success: false,
			Error:   "failed to list clients",
		})
		return
	}

	items := make([]clientItem, 0, len(clients))
	for _, c := range clients {
		items = append(items, clientItem{
			ID:      c.ID,
			Name:    c.Name,
			Email:   c.Email,
			ShopURL: c.ShopURL,
			Status:  c.Status,
			Plan:    c.Plan,
		})
	}

	writeJSON(w, http.StatusOK, listClientsResponse{
		Success: true,
		Data:    items,
	})
}

func (h *ClientHandler) UpdateClientStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/admin/clients/")
	idStr = strings.TrimSuffix(idStr, "/status")

	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, updateClientStatusResponse{
			Success: false,
			Error:   "client id is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, updateClientStatusResponse{
			Success: false,
			Error:   "invalid client id",
		})
		return
	}

	var req updateClientStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, updateClientStatusResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	err = h.clientService.UpdateClientStatus(r.Context(), id, req.Status)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, updateClientStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, updateClientStatusResponse{
		Success: true,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
