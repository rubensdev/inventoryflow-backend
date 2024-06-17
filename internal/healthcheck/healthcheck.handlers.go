package healthcheck

import (
	"log"
	"net/http"

	"github.com/rubensdev/inventoryflow-backend/internal/jsonutil"
)

type SystemInfo struct {
	Env     string // env can be "development", "staging" or "production"
	Version string
}

type HealthCheckHandler struct {
	Logger     *log.Logger
	systemInfo map[string]string
}

func NewHealthCheckHandler(log *log.Logger, systemInfo SystemInfo) *HealthCheckHandler {
	return &HealthCheckHandler{
		Logger: log,
		systemInfo: map[string]string{
			"environment": systemInfo.Env,
			"version":     systemInfo.Version,
		},
	}
}

func (h *HealthCheckHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := jsonutil.H{
		"status":      "available",
		"system_info": h.systemInfo,
	}

	err := jsonutil.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		jsonRes := jsonutil.NewJSONResponse(h.Logger)
		jsonRes.ServerErrorResponse(w, r, err)
	}
}
