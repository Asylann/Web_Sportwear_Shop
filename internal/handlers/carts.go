package handlers

import (
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/middleware"
	"WebSportwareShop/internal/models"
	"context"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

var c pb.CartServiceClient

func InitCartClientConnection() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	c = pb.NewCartServiceClient(conn)
}

func CreateCartHandle(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println("Not found cookie")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "UnAuthorized")
		return
	}

	token := cookie.Value

	userIdRaw, err := middleware.GetClaimFromToken(token, "sub")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		log.Println("user_id is not a number")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid user ID")
		return
	}

	userId := int32(userIdFloat)

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	r, err := c.CreateCart(ctx, &pb.CreateCartRequest{UserId: userId})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}
	if !r.IsCreated {
		log.Println("smt goes wrong on cart service")
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}
	log.Printf("Cart of user by id = %v was created \n", userId)
	httpresponse.WriteJSON(res, http.StatusOK, "Cart is created", "")
}

func AddToCartHandle(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Cannot find id ")
		return
	}

	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println("Not found cookie")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "UnAuthorized")
		return
	}

	token := cookie.Value

	userIdRaw, err := middleware.GetClaimFromToken(token, "sub")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}
	userEmail, err := middleware.GetClaimFromToken(token, "email")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		log.Println("user_id is not a number")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid user ID")
		return
	}

	userId := int32(userIdFloat)

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	r1, err := c.GetCartByUserId(ctx, &pb.GetCartByUserIdRequest{Id: userId})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusNotFound, "", "Can not find your cart to add")
		return
	}
	cart := r1.GetCart()

	p, err := db.GetProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusNotFound, "", "Not found product")
		return
	}

	r2, err := c.AddToCart(ctx, &pb.AddToCardRequest{Item: &pb.Cart_Item{CartId: cart.Id, Product: &pb.Product{
		Id: int32(p.ID),
	}}})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusNotFound, "", "Not found such product or cart")
		return
	}
	if !r2.GetIsAdded() {
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Connect with developers to find problem")
		return
	}

	httpresponse.WriteJSON(res, http.StatusOK, "Added to your cart", "")
	log.Printf("Product by id = %v was added to cart of %v \n", p.ID, userEmail)
}

func GetItemsOfCartById(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println("Not found cookie")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "UnAuthorized")
		return
	}

	token := cookie.Value

	userIdRaw, err := middleware.GetClaimFromToken(token, "sub")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}
	userEmail, err := middleware.GetClaimFromToken(token, "email")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		log.Println("user_id is not a number")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid user ID")
		return
	}

	userId := int32(userIdFloat)

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	r1, err := c.GetCartByUserId(ctx, &pb.GetCartByUserIdRequest{Id: userId})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusNotFound, "", "Can not find your cart to add")
		return
	}

	cartId := r1.GetCart().GetId()

	r2, err := c.GetItemsOfCartById(ctx, &pb.GetItemsOfCartByIdRequest{Id: cartId})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Smt went wrong")
		return
	}

	productIds := r2.GetProduct()

	var products []models.Product
	for i := 0; i < len(productIds); i++ {
		p, err := db.GetProduct(ctx, int(productIds[i].GetId()))
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Smt went wrong")
			return
		}
		products = append(products, p)
	}

	httpresponse.WriteJSON(res, http.StatusOK, products, "")
	log.Printf("Products of %v cart were received!!!", userEmail)
}

func DeleteItemFromCart(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Cannot find id ")
		return
	}

	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println("Not found cookie")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "UnAuthorized")
		return
	}

	token := cookie.Value

	userIdRaw, err := middleware.GetClaimFromToken(token, "sub")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}
	userEmail, err := middleware.GetClaimFromToken(token, "email")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		log.Println("user_id is not a number")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid user ID")
		return
	}

	userId := int32(userIdFloat)

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	r1, err := c.GetCartByUserId(ctx, &pb.GetCartByUserIdRequest{Id: userId})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusNotFound, "", "Can not find your cart to add")
		return
	}

	cartId := r1.GetCart().GetId()

	r2, err := c.DeleteItemFromCart(ctx, &pb.DeleteItemFromCartRequest{CartId: cartId, ProductId: int32(id)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Smt went wrong")
		return
	}

	deletedProduct, err := db.GetProduct(ctx, int(r2.DeletedProduct.GetId()))
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Not found product deleted product here")
		return
	}

	httpresponse.WriteJSON(res, http.StatusOK, deletedProduct, "")
	log.Printf("Product by id = %v was deleted from %v cart", deletedProduct.ID, userEmail)
}
