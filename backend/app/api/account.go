package api

import (
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"ctigroupjsc.com/phuocnn/gps-management/model"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	serviceAccount "ctigroupjsc.com/phuocnn/gps-management/service/account"
	"github.com/julienschmidt/httprouter"
)

type AccountHandler struct {
	AccountRepository     repository.AccountRepository
	UserRepository        repository.UserRepository
	ActivityLogRepository repository.ActivityLogRepository
	ReportRepository      repository.ReportRepository
}

func (h *AccountHandler) All(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		data, err := h.AccountRepository.FindByEmail(Email(r))
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
				Data:  []model.Account{*data},
			},
		})
		return
	}

	condition := make(map[string]interface{})
	if Scope(r) == model.ScopeAdmin {
		condition["$or"] = []bson.M{
			{"email": Email(r)},
			{"createdBy": Email(r)},
		}
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

		total, data, err := h.AccountRepository.Pagination(page, pageSize, condition)
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
	accounts, err := h.AccountRepository.All()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "add account successfully",
		Code:    http.StatusOK,
		Data:    accounts,
	})
}

func (h *AccountHandler) AccountInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, err := h.AccountRepository.FindByEmail(Email(r))
	if err != nil {
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *AccountHandler) Detail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	detailAccountHandler := serviceAccount.DetailAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}

	data, err := detailAccountHandler.Handle(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *AccountHandler) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	cmd := new(serviceAccount.AddAccount)
	if err := BindJSON(r, cmd); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	cmd.CreatedBy = Email(r)

	handler := &serviceAccount.AddAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}
	err := handler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "add account successfully",
		Code:    http.StatusOK,
	})
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(serviceAccount.UpdateAccount)
	if err := BindJSON(r, cmd); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	cmd.UserID = p.ByName("userID")

	handler := &serviceAccount.UpdateAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}
	err := handler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "update account successfully",
		Code:    http.StatusOK,
	})
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	deleteAccountHandler := &serviceAccount.DeleteAccountHandler{
		AccountRepository:     h.AccountRepository,
		UserRepository:        h.UserRepository,
		ActivityLogRepository: h.ActivityLogRepository,
		ReportRepository:      h.ReportRepository,
	}

	err := deleteAccountHandler.Handle(p.ByName("id"))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "delete account successfully",
		Code:    http.StatusOK,
	})
}

// Reset func is delete all data, do not delete account
func (h *AccountHandler) Reset(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	id := p.ByName("id")

	err := h.ActivityLogRepository.RemoveByAccountID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	err = h.ReportRepository.RemoveByAccountID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "reset account successfully",
		Code:    http.StatusOK,
	})
}

func (h *AccountHandler) GetChildAccounts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	if Scope(r) == model.ScopeRoot {
		result, err := h.AccountRepository.All()
		if err == nil {
			WriteJSON(w, http.StatusOK, ResponseBody{
				Data: result,
				Code: http.StatusOK,
			})
			return
		}
	}

	result, err := h.AccountRepository.GetChildAccounts(Email(r))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Data: result,
		Code: http.StatusOK,
	})
}
