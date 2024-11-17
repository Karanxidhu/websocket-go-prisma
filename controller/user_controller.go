package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/data/response"
	"github.com/karanxidhu/go-websocket/helper"
	"github.com/karanxidhu/go-websocket/service"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (controller *UserController) Save(w http.ResponseWriter, r *http.Request) {
	userCreateReq := request.UserCreateReq{}
	helper.ReadRequest(r, &userCreateReq)
	fmt.Println(userCreateReq.Username)
	controller.UserService.Save(r.Context(), userCreateReq)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    userCreateReq,
	}
	helper.WriteResponse(w, webresponse)
}
func (controller *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userUpdateReq := request.UserUpdateReq{}
	helper.ReadRequest(r, &userUpdateReq)
	vars := mux.Vars(r)
	userId := vars["userId"]
	userUpdateReq.UserId = userId
	controller.UserService.Update(r.Context(), userUpdateReq)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    userUpdateReq,
	}
	helper.WriteResponse(w, webresponse)
}

func (controller *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	
	result := controller.UserService.FindAll(r.Context())

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data: result,
	}
	helper.WriteResponse(w, webresponse)
}


func (controller *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	result := controller.UserService.FindById(r.Context(), userId)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data: result,
	}
	helper.WriteResponse(w, webresponse)
}

func (controller *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	controller.UserService.Delete(r.Context(), userId)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data: "row deleted successfully",
	}
	helper.WriteResponse(w, webresponse)
}
