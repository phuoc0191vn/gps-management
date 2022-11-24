package api

import (
	"net/http"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
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
