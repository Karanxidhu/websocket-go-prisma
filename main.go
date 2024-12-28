package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karanxidhu/go-websocket/config"
	"github.com/karanxidhu/go-websocket/controller"
	"github.com/karanxidhu/go-websocket/repository"
	"github.com/karanxidhu/go-websocket/routes"
	"github.com/karanxidhu/go-websocket/service"
)


func main() {
	Db, err := config.ConnectToDB()
	userRespository := repository.NewUserRespository(Db)
	fileRespository := repository.NewFileRepository(Db)
	
	userService := service.NewUserServiceImpl(userRespository)
	fileService := service.NewFileServiceImpl(fileRespository)
	
	userController := controller.NewUserController(userService)
	fileController := controller.NewFileController(fileService)

	router := setupRoutes(userController, fileController)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Prisma.Disconnect()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes(userController *controller.UserController, fileController *controller.FileController) *mux.Router {
	router := mux.NewRouter()
	manager := NewManager()
	router.HandleFunc("/", handler)
	router.HandleFunc("/ws", manager.servesWs)
	router.HandleFunc("/api/users", userController.FindAll).Methods("GET")
	router.HandleFunc("/api/users/{userId}", userController.FindById).Methods("GET")
	router.HandleFunc("/api/users", userController.Save).Methods("POST")
	router.HandleFunc("/api/users/{userId}", userController.Update).Methods("PUT")
	router.HandleFunc("/api/users/{userId}", userController.Delete).Methods("DELETE")
	router.HandleFunc("/api/upload", routes.UploadHandler).Methods("POST")

	router.HandleFunc("/api/files/add", fileController.SaveFile).Methods("POST")

	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
