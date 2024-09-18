package service_product

import (
	"context"
	"fmt"
	"hshelby-tkcled-product/config"
	"hshelby-tkcled-product/src/utilities"
	"io"
	"log"

	"net/http"
	"net/url"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("[UploadImage] url.ParseQuery", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	code := queryValues.Get("code")

	exePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	// Lấy thư mục chứa tệp thực thi
	basedir := filepath.Dir(exePath)

	filePath := filepath.Join(basedir, "serviceAccountKey.json")
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("[UploadImage] firebase.NewApp", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		log.Println("[UploadImage] app.Storage", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	bucketName := config.Get().BucketName

	bucketHandle, err := client.Bucket(bucketName)
	if err != nil {
		log.Println("[UploadImage] client.Bucket", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("[UploadImage] r.FormFile", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}
	defer file.Close()

	writer := bucketHandle.Object(code + ".png").NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": code}
	defer writer.Close()

	_, err = io.Copy(writer, file)
	if err != nil {
		log.Println("[UploadImage] io.Copy", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	avatar := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s.png?alt=media&token=%s", config.Get().BucketName, code, code)

	utilities.WriteJSON(w, http.StatusOK, avatar)
}
