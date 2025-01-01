package controller

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/karanxidhu/go-websocket/config"
	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/data/response"
	"github.com/karanxidhu/go-websocket/helper"
	"github.com/karanxidhu/go-websocket/repository"
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
	token := controller.UserService.Save(r.Context(), userCreateReq)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    token,
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
		Data:    result,
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
		Data:    result,
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
		Data:    "row deleted successfully",
	}
	helper.WriteResponse(w, webresponse)
}

func (controller *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	req := request.ProfileRequest{}
	helper.ReadRequest(r, &req)

	AuthToken := req.AuthToken

	if AuthToken == "" {
		http.Error(w, "Auth token is required", http.StatusBadRequest)
		return
	}

	var secretKey = []byte("hjgdsajfsdakfhasdhfiao@!#!@$!$231231")

	token, err := jwt.Parse(AuthToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the AuthToken
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		webresponse := response.WebResponse{
			Code:    401,
			Message: "Token invalid",
		}
		helper.WriteResponse(w, webresponse)
		return
	}

	UserID := ""
	// Extract claims and access "userId"
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, found := claims["userID"]; found {
			fmt.Printf("User ID: %v\n", userId)
			UserID = userId.(string)
		} else {
			fmt.Println("userId not found in the token")
		}
	} else {
		fmt.Println("Invalid token")
	}

	user, err := repository.FindById(UserID, r.Context(), config.Db)

	if err != nil {
		http.Error(w, "Unable to find user", http.StatusBadRequest)
		return
	}
	helper.WriteResponse(w, response.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    user,
	})
}
