package api

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	serviceUser "ctigroupjsc.com/phuocnn/gps-management/service/user"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if bson.IsObjectIdHex(id) {
		user, err := h.UserRepository.FindByID(id)
		if err != nil {
			WriteJSON(w, HTTP_ERROR_CODE_READ_FAILED, ResponseBody{
				Message: "unable to find user",
				Code:    HTTP_ERROR_CODE_READ_FAILED,
			})
			return
		}

		WriteJSON(w, http.StatusOK, user)
		return
	}

	user, err := h.UserRepository.FindByIdentifyID(id)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_READ_FAILED, ResponseBody{
			Message: "unable to find user",
			Code:    HTTP_ERROR_CODE_READ_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(serviceUser.AddUser)
	if err := BindJSON(r, cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	handler := &serviceUser.AddUserHandler{
		UserRepository: h.UserRepository,
	}
	id, err := handler.Handle(cmd)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: "unable to add user",
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: fmt.Sprintf("add user successfully: %s", id),
		Code:    http.StatusOK,
	})
}
