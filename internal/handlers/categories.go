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

func CreateCategoryHandle(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	categories, err := db.ListOfCategories(context.Background())
	for _, category := range categories {
		if category.Name == c.Name {
			httpresponse.WriteJSON(w, http.StatusConflict, nil, "Already exists")
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.CreateCategory(ctx, &c)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusCreated, c, "")
	log.Printf("Category was created! : %v \n", c)
}

func GetCategoryHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var c models.Category
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	c, err = db.GetCategory(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, c, "")
	log.Printf("Category by id=%v was recieved! \n", id)
}

func DeleteCategoryHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.DeleteCategory(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Category by id=%v was deleted! \n", id)
}

func ListOfCategoriesHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	categories, err := db.ListOfCategories(ctx)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, categories, "")
	log.Printf("All Categories list were get!!!")
}

func UpdateCategoryHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var c models.Category
	if err = json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	c.ID = id
	err = db.UpdateCategory(ctx, &c)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, c, "")
	log.Printf("Category by id = %v was updated : %v", id, c)
}
