package api

import (
	"net/http"
	"strconv"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	serviceContact "ctigroupjsc.com/phuocnn/gps-management/service/contact"
	"github.com/julienschmidt/httprouter"
)

type ContactHandler struct {
	ContactRepository repository.ContactRepository
	AccountRepository repository.AccountRepository
}

func (h *ContactHandler) All(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	condition := make(map[string]interface{})
	if Scope(r) == model.ScopeUser {
		condition["accountID"] = AccountID(r)
	}

	if Scope(r) != model.ScopeUser {
		condition["adminEmail"] = Email(r)
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

		total, data, err := h.ContactRepository.Pagination(page, pageSize, condition)
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
	data, err := h.ContactRepository.All()
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *ContactHandler) Done(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	err := h.ContactRepository.UpdateStatus(p.ByName("id"), model.StatusContactProcessed)
	if err != nil {
		WriteJSON(w, http.StatusOK, ResponseBody{
			Message: "unable to update status contact",
			Code:    http.StatusOK,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !IsScopeAllowed(r) {
		return
	}

	err := h.ContactRepository.RemoveByID(p.ByName("id"))
	if err != nil {
		WriteJSON(w, http.StatusOK, ResponseBody{
			Message: "unable to remove status contact",
			Code:    http.StatusOK,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
	})
}

func (h *ContactHandler) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(serviceContact.AddContact)
	if err := BindJSON(r, cmd); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	account, err := h.AccountRepository.FindByID(AccountID(r))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	cmd.AccountID = AccountID(r)
	cmd.Email = Email(r)
	cmd.AdminEmail = account.CreatedBy

	handler := &serviceContact.AddContactHandler{
		ContactRepository: h.ContactRepository,
	}

	err = handler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "add contact successfully",
		Code:    http.StatusOK,
	})
}
