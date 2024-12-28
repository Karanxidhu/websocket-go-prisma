package controller

import (
	"fmt"
	"net/http"

	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/data/response"
	"github.com/karanxidhu/go-websocket/helper"
	"github.com/karanxidhu/go-websocket/service"
)

type FileController struct {
	FileService service.FileService
}

func NewFileController(fileService service.FileService) *FileController {
	return &FileController{
		FileService: fileService,
	}
}

func (fc *FileController) SaveFile(w http.ResponseWriter, r *http.Request) {
	fileReq := request.FileResponse{}
	helper.ReadRequest(r, &fileReq)
	fmt.Println(fileReq)

	fc.FileService.Save(r.Context(), fileReq)

	webresponse := response.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    fileReq,
	}
	helper.WriteResponse(w, webresponse)
}
