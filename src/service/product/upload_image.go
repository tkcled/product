package service_product

import (
	"context"
	"fmt"
	"hshelby-tkcled-product/config"
	"hshelby-tkcled-product/src/database/collection"
	"hshelby-tkcled-product/src/database/model"
	"hshelby-tkcled-product/src/utilities"
	"io"
	"log"

	"net/http"
	"net/url"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("[UploadImage] url.ParseQuery", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	productID := queryValues.Get("product_id")
	product := &model.Product{}
	err = collection.Product().Collection().FindOne(context.Background(), bson.M{"_id": productID}).Decode(product)
	if err != nil {
		log.Println("[UploadImage] find product", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

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

	writer := bucketHandle.Object(productID + ".png").NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": productID}
	defer writer.Close()

	_, err = io.Copy(writer, file)
	if err != nil {
		log.Println("[UploadImage] io.Copy", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	avatar := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s.png?alt=media&token=%s", config.Get().BucketName, productID, productID)

	utilities.WriteJSON(w, http.StatusOK, avatar)
}
