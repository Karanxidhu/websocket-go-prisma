package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karanxidhu/go-websocket/config"
	"github.com/karanxidhu/go-websocket/controller"
	"github.com/karanxidhu/go-websocket/repository"
	"github.com/karanxidhu/go-websocket/service"
)

func main() {
	db, err := config.ConnectToDB()
	userRespository := repository.NewUserRespository(db)
	
	userService := service.NewUserServiceImpl(userRespository)
	
	userController := controller.NewUserController(userService)
	
	router := setupRoutes(userController)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Prisma.Disconnect()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes(userController *controller.UserController) *mux.Router {
	router := mux.NewRouter()
	manager := NewManager()
	router.HandleFunc("/", handler)
	router.HandleFunc("/ws", manager.servesWs)
	router.HandleFunc("/api/users", userController.FindAll).Methods("GET")
	router.HandleFunc("/api/users/{userId}", userController.FindById).Methods("GET")
	router.HandleFunc("/api/users", userController.Save).Methods("POST")
	router.HandleFunc("/api/users/{userId}", userController.Update).Methods("PUT")
	router.HandleFunc("/api/users/{userId}", userController.Delete).Methods("DELETE")
	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
