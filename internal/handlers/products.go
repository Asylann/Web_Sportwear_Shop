package handlers

import (
	"WebSportwareShop/internal/cache"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/go-redis/redis/v8"
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

	cache.Rdc.Del(ctx, "products")
	log.Println("Cached products were deleted!")

	categoryIdStr := strconv.Itoa(p.CategoryID)
	cache.Rdc.Del(ctx, "products:category/"+categoryIdStr)
	log.Printf("Cached products:category/%v were deleted!", categoryIdStr)

	sellerIdStr := strconv.Itoa(p.SellerID)
	cache.Rdc.Del(ctx, "products:seller/"+sellerIdStr)
	log.Printf("Cached products:seller/%v were deleted!", sellerIdStr)

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

	p, err := db.GetProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, "", "No found product!!!")
		return
	}

	err = db.DeleteProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	res1, err := c.ChangeEtagVersionOfCartsByProductId(ctx, &pb.ChangeEtagVersionOfCartsByProductIdRequest{ProductId: int64(id)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if res1.IsChanged {
		log.Printf("Etag versions of carts that contained product by id = %v were changed", id)
	}

	res2, err := c.DeleteProductOfCarts(ctx, &pb.DeleteProductOfCartsRequest{ProductId: int64(id)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	log.Printf("Product by id = %v was deleted from all carts", res2.ProductId)

	cache.Rdc.Del(ctx, "products")
	log.Println("Cached products were deleted!")

	categoryIdStr := strconv.Itoa(p.CategoryID)
	cache.Rdc.Del(ctx, "products:category/"+categoryIdStr)
	log.Printf("Cached products:category/%v were deleted!", categoryIdStr)

	sellerIdStr := strconv.Itoa(p.SellerID)
	cache.Rdc.Del(ctx, "products:seller/"+sellerIdStr)
	log.Printf("Cached products:seller/%v were deleted!", sellerIdStr)

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Product by id=%v was deleted! \n", id)
}

func ListOfProductsHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	cacheKey := "products"

	if jsonBytes, err := cache.Rdc.Get(ctx, cacheKey).Bytes(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		log.Printf("All Products list were recieved!!! By caching")
		return
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	products, err := db.ListOfProducts(ctx)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	payload, err := httpresponse.MarshalResponse(products, "")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", err.Error())
		return
	}

	cache.Rdc.Set(ctx, cacheKey, payload, 5*time.Minute)

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
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	oldProduct, err := db.GetProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, "", "No such product")
	}
	var p models.Product
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	p.ID = id
	err = db.UpdateProduct(ctx, &p)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	cache.Rdc.Del(ctx, "products")
	log.Println("Cached products were deleted!")

	categoryIdStr := strconv.Itoa(p.CategoryID)
	cache.Rdc.Del(ctx, "products:category/"+categoryIdStr)
	log.Printf("Cached products:category/%v were deleted!", categoryIdStr)

	sellerIdStr := strconv.Itoa(p.SellerID)
	cache.Rdc.Del(ctx, "products:seller/"+sellerIdStr)
	log.Printf("Cached products:seller/%v were deleted!", sellerIdStr)

	if p.CategoryID != oldProduct.CategoryID {
		oldCategoryIdStr := strconv.Itoa(oldProduct.CategoryID)
		cache.Rdc.Del(ctx, "products:category/"+oldCategoryIdStr)
		log.Printf("Cached products:category/%v were deleted!", oldCategoryIdStr)
	}
	if p.SellerID != oldProduct.SellerID {
		oldSellerIdStr := strconv.Itoa(oldProduct.SellerID)
		cache.Rdc.Del(ctx, "products:seller/"+oldSellerIdStr)
		log.Printf("Cached products:category/%v were deleted!", oldSellerIdStr)
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

	cacheKey := "products:category/" + vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if jsonBytes, err := cache.Rdc.Get(ctx, cacheKey).Bytes(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		log.Printf("Products bu categoryId= %v were received By Caching!!", id)
		return
	}

	var products []models.Product
	products, err = db.ListOfProductsByCategory(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}

	payload, err := httpresponse.MarshalResponse(products, "")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", err.Error())
		return
	}

	cache.Rdc.Set(ctx, cacheKey, payload, 5*time.Minute)

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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	cacheKey := "products:seller/" + vars["id"]

	if jsonBytes, err := cache.Rdc.Get(ctx, cacheKey).Bytes(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		log.Printf("Products by seller id = %v were received By caching !!", vars["id"])
		return
	}

	var products []models.Product

	products, err = db.ListOfProductsBySellerID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}

	payload, err := httpresponse.MarshalResponse(products, "")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", err.Error())
		return
	}

	cache.Rdc.Set(ctx, cacheKey, payload, 5*time.Minute)

	httpresponse.WriteJSON(w, http.StatusOK, products, "")
	log.Printf("Products by userId= %v were received", id)
}
