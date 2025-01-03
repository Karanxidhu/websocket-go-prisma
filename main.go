package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://chaygger.karanxd.space"}), // Frontend URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),                 // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),                           // Allowed headers
	)

	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
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
	router.HandleFunc("/api/chat/{roomName}", fileController.GetChat).Methods("GET")
	router.HandleFunc("/api/users/profile", userController.Profile).Methods("POST")
	router.HandleFunc("/api/files/add", fileController.SaveFile).Methods("POST")
	router.HandleFunc("/env", Getenv).Methods("GET")
	return router
}
func Getenv(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(os.Getenv("TEST")))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
