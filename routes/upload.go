package routes

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-jwt/jwt"
	"github.com/karanxidhu/go-websocket/config"
	"github.com/karanxidhu/go-websocket/data/response"
	"github.com/karanxidhu/go-websocket/helper"
	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/repository"
)

const (
	S3_BUCKET_NAME = "nigga-bucket"
	S3_REGION      = "ap-south-1"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload handler")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	RoomName := r.FormValue("roomName")
	AuthToken := r.FormValue("authToken")

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
	UserID = user.Id
	if UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		fmt.Println("File is required")
		return
	}
	if file == nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		fmt.Println("File is required")
		return
	}

	defer file.Close()

	fileKey, err := uplaodToS3(file, fileHeader, UserID)

	if err != nil {
		http.Error(w, "Unable to upload file to S3"+err.Error(), http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3_BUCKET_NAME, S3_REGION, fileKey)

	fileData := model.File{
		Message:        fileURL,
		UploadedAt: time.Now(),
		RoomName:   RoomName,
		UserName:   user.Username,
	}

	repository.SaveFile(r.Context(), fileData, config.Db)

	helper.WriteResponse(w, map[string]string{
		"message":  "File uploaded successfully",
		"url":      fileURL,
		"userID":   UserID,
		"roomName": RoomName,
	})
}

func uplaodToS3(file multipart.File, fileHeader *multipart.FileHeader, userID string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		// CredentialsChainVerboseErrors: aws.Bool(true),
	})

	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	fileKey := fmt.Sprintf("uploads/%s/%s", userID, fileHeader.Filename)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(S3_BUCKET_NAME),
		Key:           aws.String(fileKey),
		Body:          file,
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", err
	}

	return fileKey, nil

}
