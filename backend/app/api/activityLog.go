package api

import (
	"net/http"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	"github.com/julienschmidt/httprouter"
)

type ActivityLogHandler struct {
	ActivityLogRepository repository.ActivityLogRepository
	DeviceRepository      repository.DeviceRepository
}

func (h *ActivityLogHandler) GetInDay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accountID := AccountID(r)
	if IsScopeAllowed(r) {
		accountID = p.ByName("accountID")
	}

	device, err := h.DeviceRepository.FindDeviceByStatus(accountID, model.StatusEnable)
	if err != nil || len(device) < 1 {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	activityLog, err := h.ActivityLogRepository.GetInDay(device[0].ID.Hex(), accountID, date)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: activityLog,
	})
}

func (h *ActivityLogHandler) CurrentLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accountID := AccountID(r)
	if IsScopeAllowed(r) {
		accountID = p.ByName("accountID")
	}

	device, err := h.DeviceRepository.FindDeviceByStatus(accountID, model.StatusEnable)
	if err != nil || len(device) < 1 {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: device[0],
	})
}
