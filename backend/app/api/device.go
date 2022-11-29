package api

import (
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"ctigroupjsc.com/phuocnn/gps-management/model"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	serviceDevice "ctigroupjsc.com/phuocnn/gps-management/service/device"
	"github.com/julienschmidt/httprouter"
)

type DeviceHandler struct {
	DeviceRepository  repository.DeviceRepository
	AccountRepository repository.AccountRepository
}

func (h *DeviceHandler) All(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		data, err := h.DeviceRepository.FindByAccountID(AccountID(r))
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ResponseBody{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		WriteJSON(w, http.StatusOK, ResponseBody{
			Code: http.StatusOK,
			Data: struct {
				Total int
				Data  interface{}
			}{
				Total: 1,
				Data:  data,
			},
		})
		return
	}

	childAccountIDs := make([]string, 0)
	condition := make(map[string]interface{})
	if Scope(r) == model.ScopeAdmin {
		childAccounts, err := h.AccountRepository.GetChildAccounts(Email(r))
		if err != nil {
			childAccounts = make([]model.Account, 0)
		}

		for i := 0; i < len(childAccounts); i++ {
			childAccountIDs = append(childAccountIDs, childAccounts[i].ID.Hex())
		}
		childAccountIDs = append(childAccountIDs, AccountID(r))

		condition["accountID"] = bson.M{"$in": childAccountIDs}
	}

	output, ok := GetQuery(r, DATATABLE_QUERY_OUTPUT)
	if ok && output == DATATABLE_QUERY_OUTPUT_DATATABLE {
		page := 1
		pageSize := 10
		var err error

		if r.URL.Query().Get("page") != "" {
			page, err = strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ResponseBody{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}
		}

		if r.URL.Query().Get("limit") != "" {
			pageSize, err = strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ResponseBody{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}
		}

		total, data, err := h.DeviceRepository.Pagination(page, pageSize, condition)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ResponseBody{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		WriteJSON(w, http.StatusOK, ResponseBody{
			Code: http.StatusOK,
			Data: struct {
				Total int
				Data  interface{}
			}{
				Total: total,
				Data:  data,
			},
		})
		return
	}

	// normal response
	accounts, err := h.DeviceRepository.FindByAccountID(AccountID(r))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: accounts,
	})
}

func (h *DeviceHandler) GetByStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	status, ok := GetQuery(r, "status")
	if !ok {
		status = "0"
	}

	statusInt := model.StatusEnable
	var err error
	statusInt, err = strconv.Atoi(status)
	if err != nil {
		statusInt = model.StatusDisable
	}

	data, err := h.DeviceRepository.FindDeviceByStatus(AccountID(r), statusInt)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if statusInt == model.StatusEnable && len(data) > 0 {
		WriteJSON(w, http.StatusOK, ResponseBody{
			Code: http.StatusOK,
			Data: data[0],
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *DeviceHandler) Detail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	data, err := h.DeviceRepository.FindByID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *DeviceHandler) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	cmd := new(serviceDevice.AddDevice)
	if err := BindJSON(r, cmd); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	addDeviceHandler := serviceDevice.AddDeviceHandler{
		DeviceRepository:  h.DeviceRepository,
		AccountRepository: h.AccountRepository,
	}

	err := addDeviceHandler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "add device successfully",
		Code:    http.StatusOK,
	})
}

func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(serviceDevice.UpdateDevice)
	if err := BindJSON(r, cmd); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !IsScopeAllowed(r) {
		cmd.AccountID = ""
	}

	updateDeviceHandler := serviceDevice.UpdateDeviceHandler{
		DeviceRepository:  h.DeviceRepository,
		AccountRepository: h.AccountRepository,
	}

	err := updateDeviceHandler.Handle(p.ByName("id"), cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_UPDATE_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_UPDATE_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "update device successfully",
		Code:    http.StatusOK,
	})
}

func (h *DeviceHandler) Toggle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	statusStr := r.URL.Query().Get("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil && !model.IsValidStatus(status) {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Code: http.StatusBadRequest,
		})
	}

	id := p.ByName("id")
	device, err := h.DeviceRepository.FindByID(id)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
	}

	device.Status = status
	if status == model.StatusEnable {
		// disable all device of this account
		deviceIDs, err := h.AccountRepository.GetDeviceIDsByID(device.AccountID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ResponseBody{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
		}

		for i := 0; i < len(deviceIDs); i++ {
			h.DeviceRepository.UpdateStatus(deviceIDs[i], model.StatusDisable)
		}
	}

	err = h.DeviceRepository.UpdateByID(id, *device)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "update device successfully",
		Code:    http.StatusOK,
	})
}

func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	id := p.ByName("id")
	device, err := h.DeviceRepository.FindByID(id)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	err = h.DeviceRepository.RemoveByID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	oldAccount, err := h.AccountRepository.GetDeviceIDsByID(device.AccountID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	tmp := make([]string, 0)
	for i := 0; i < len(oldAccount); i++ {
		if oldAccount[i] == device.ID.Hex() {
			continue
		}

		tmp = append(tmp, oldAccount[i])
	}
	oldAccount = make([]string, 0)
	oldAccount = tmp

	err = h.AccountRepository.UpdateDeviceIDs(device.AccountID, oldAccount)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "delete device successfully",
		Code:    http.StatusOK,
	})
}
