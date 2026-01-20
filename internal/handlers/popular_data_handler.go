package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/service-info-aggregator/internal/models/dto"
	"github.com/service-info-aggregator/internal/services"
)

type PopularDataHandler struct {
	service *services.PopularDataService
}

func NewPopularDataHandler(service *services.PopularDataService) *PopularDataHandler {
	return &PopularDataHandler{
		service: service,
	}
}

func (h *PopularDataHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PopularDataHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PopularDataHandler) getAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetAll(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, result)
}

func (h *PopularDataHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, result)
}

func (h *PopularDataHandler) create(w http.ResponseWriter, r *http.Request) {
	var input dto.PopularDataDto

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), &input)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusCreated, result)
}

func (h *PopularDataHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var input dto.PopularDataDto

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.service.Update(r.Context(), id, &input)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, updated)
}

func (h *PopularDataHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 {
		return 0, http.ErrNoLocation
	}

	return strconv.Atoi(parts[1])
}
