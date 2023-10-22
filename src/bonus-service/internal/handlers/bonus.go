package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab2/src/bonus-service/internal/models"
	"lab2/src/bonus-service/internal/repository"

	"github.com/gorilla/mux"
)

type BonusHandler struct {
	BonusRepo repository.Repository
}

func (h *BonusHandler) CreatePrivilegeHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var record models.PrivilegeHistory

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.BonusRepo.CreateHistoryRecord(&record); err != nil {
		log.Printf("Failed to create ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BonusHandler) CreatePrivilegeHandler(w http.ResponseWriter, r *http.Request) {
	var record models.Privilege

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.BonusRepo.CreatePrivilege(&record); err != nil {
		log.Printf("Failed to create ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BonusHandler) GetPrivilegeByUsernameHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	privilege, err := h.BonusRepo.GetPrvilegeByUsername(params["username"])
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(privilege); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BonusHandler) GetHistoryByIdHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	history, err := h.BonusRepo.GetHistoryById(params["privilegeId"])
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(history); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
