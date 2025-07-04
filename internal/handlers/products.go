package handlers

import (
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateProductHandle(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	products, err := db.ListOfProducts(context.Background())
	for _, product := range products {
		if product.Name == p.Name {
			httpresponse.WriteJSON(w, http.StatusConflict, nil, "Already exists")
		}
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.CreateProduct(ctx, &p)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusCreated, p, "")
	log.Printf("Product was created! : %v \n", p)
}

func GetProductHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var p models.Product
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	p, err = db.GetProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, p, "")
	log.Printf("Product by id=%v was recieved! \n", id)
}

func DeleteProductHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.DeleteProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Product by id=%v was deleted! \n", id)
}

func ListOfProductsHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	products, err := db.ListOfProducts(ctx)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, products, "")
	log.Printf("All Products list were recieved!!!")
}

func UpdateProductHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var p models.Product
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	p.ID = id
	err = db.UpdateProduct(ctx, &p)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	log.Printf("Product by id = %v was updated : %v", id, p)
}

func ListOfProductsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, "can not invert str to int")
		return
	}
	var products []models.Product
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	products, err = db.ListOfProductsByCategory(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, products, "")
	log.Printf("Products bu categoryId= %v were received", id)
}

func ListOfProductsBySellerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, "can not invert str to int")
		return
	}
	var products []models.Product
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	products, err = db.ListOfProductsBySellerID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, products, "")
	log.Printf("Products by userId= %v were received", id)
}
