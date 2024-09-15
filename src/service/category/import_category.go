package service_category

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"hshelby-tkcled-product/src/database/collection"
	"hshelby-tkcled-product/src/database/model"
	"hshelby-tkcled-product/src/utilities"
	"io"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ImportPCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	// get file
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("[ImportCategory] r.FormFile()", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := context.Background()
	if strings.HasSuffix(handler.Filename, ".csv") {
		// get current categories
		categories := []model.Category{}

		cur, err := collection.Category().Collection().Find(ctx, bson.M{})
		if err != nil {
			log.Println("ImportCategory get parentCategories", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		}

		err = cur.All(ctx, &categories)
		if err != nil {
			log.Println("ImportCategory cur.All(ctx, &categories)", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		}

		parentCategories := map[string]string{}
		for _, ele := range categories {
			parentCategories[ele.Name] = ele.ID
		}

		// handle import
		err = processCSVFileCategory(file, parentCategories)
		if err != nil {
			log.Println("[ImportCategory] err", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
			return
		}

		jsonData, err := json.Marshal("success")
		if err != nil {
			log.Println("[ImportCategory] JSON marshal error", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	} else {
		utilities.WriteJSON(w, http.StatusBadRequest, nil)
	}
}

func processCSVFileCategory(file io.Reader, parentCategories map[string]string) (err error) {
	reader := csv.NewReader(file)

	firstRow := false

	var operations []mongo.WriteModel

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !firstRow { // skip title row
			firstRow = true
			continue
		}

		category := model.Category{
			ID: primitive.NewObjectID().Hex(),

			Name:        record[0],
			Description: record[1],
		}

		operation := mongo.NewInsertOneModel().SetDocument(category)
		operations = append(operations, operation)
	}

	bulkOptions := options.BulkWrite().SetOrdered(false)
	res, err := collection.Category().Collection().BulkWrite(context.TODO(), operations, bulkOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted documents, count:", res.InsertedCount)
	return nil
}
