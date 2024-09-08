package service_product

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"
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

func ImportProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	// get file
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("[ImportProduct] r.FormFile()", err)
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
			log.Println("ImportProduct get parentCategories", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		}

		err = cur.All(ctx, &categories)
		if err != nil {
			log.Println("ImportProduct cur.All(ctx, &parentCategories)", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
		}

		currentCategories := map[string]model.Category{}
		for _, ele := range categories {
			currentCategories[ele.Name] = ele
		}

		err = processCSVFileProduct(file, currentCategories)
		if err != nil {
			log.Println("[ImportProduct] err", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
			return
		}

		jsonData, err := json.Marshal("success")
		if err != nil {
			log.Println("[ImportProduct] JSON marshal error", err)
			utilities.WriteJSON(w, http.StatusInternalServerError, nil)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	} else {
		utilities.WriteJSON(w, http.StatusBadRequest, nil)
	}
}

func processCSVFileProduct(file io.Reader, currentCategories map[string]model.Category) (err error) {
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

		category := currentCategories[record[6]]

		product := model.Product{
			ID: primitive.NewObjectID().Hex(),

			Name:             record[0],
			Image:            record[1],
			Description:      record[2],
			Code:             record[3],
			UnitPrice:        utilities.StringToFloat64(record[4]),
			CatalogLink:      record[5],
			CategoryID:       category.ID,
			ParentCategoryID: category.ParentID,
		}

		operation := mongo.NewInsertOneModel().SetDocument(product)
		operations = append(operations, operation)
	}

	bulkOptions := options.BulkWrite().SetOrdered(false)
	res, err := collection.Product().Collection().BulkWrite(context.TODO(), operations, bulkOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted documents, count:", res.InsertedCount)

	return nil
}
