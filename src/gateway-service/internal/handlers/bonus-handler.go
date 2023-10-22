package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab2/src/gateway-service/internal/service"
)

func (gs *GatewayService) GetPrivilege(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privilegeInfo, err := service.UserPrivilegeController(
		gs.Config.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(privilegeInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
