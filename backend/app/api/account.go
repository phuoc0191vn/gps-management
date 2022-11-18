package api

import (
	"net/http"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	serviceAccount "ctigroupjsc.com/phuocnn/gps-management/service/account"

	"github.com/julienschmidt/httprouter"
)

type AccountHandler struct {
	AccountRepository repository.AccountRepository
}

func (h *AccountHandler) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(serviceAccount.AddAccount)
	if err := BindJSON(r, cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	handler := &serviceAccount.AddAccountHandler{
		AccountRepository: h.AccountRepository,
	}
	err := handler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: "unable to add user",
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "add account successfully",
		Code:    http.StatusOK,
	})
}
