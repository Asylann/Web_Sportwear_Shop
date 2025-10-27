package handlers

import (
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	opb "github.com/Asylann/OrderServiceGRPC/proto"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

var OrderClient opb.OrderServiceClient

func InitOrderServiceConn() {
	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Order service is not connected: %s", err.Error())
		return
	}
	OrderClient = opb.NewOrderServiceClient(conn)
}

type BodyOfCreate struct {
	TransportType string `json:"transport_type"`
	Address       string `json:"address"`
}

func CreateOrderHandle(res http.ResponseWriter, req *http.Request) {
	var body BodyOfCreate
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.Println("Invalid Body")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid fields")
		return
	}

	UserId, err := GetUserIdFromReq(res, req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r1, err := c.GetCartByUserId(ctx, &pb.GetCartByUserIdRequest{Id: int32(UserId)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "smt went wrong")
		return
	}

	r, err := OrderClient.CreateOrder(ctx, &opb.CreateOrderRequest{
		UserId:               int32(UserId),
		CartId:               r1.Cart.Id,
		TypeOfTransportation: body.TransportType,
		Address:              body.Address,
	})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Smt went wrong")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	r2, err := OrderClient.GetItemsOfOrderById(ctx, &opb.GetItemsOfOrderByIdRequest{OrderId: int32(r.OrderId)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Smt went wrong!")
		return
	}

	var products []models.Product
	for _, v := range r2.ListOfProductsId {
		product, err := db.GetProduct(ctx, int(v))
		if err != nil {
			continue
		}
		err = db.MakeAPayment(ctx, UserId, product.SellerID, product.Price)
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Smt went wrong!")
			return
		}
		products = append(products, product)
	}

	log.Printf("User by id %v made an order with cart is %v and transport %v to address %v", UserId, r1.Cart.Id, body.TransportType, body.Address)
	httpresponse.WriteJSON(res, http.StatusOK, r.DeliveredAt.AsTime().Format("2006-01-02 15:04:05"), "")
	return
}

type OrderDTO struct {
	Id            int32  `json:"id"`
	CartId        int32  `json:"cart_id"`
	CreateAt      string `json:"create_at"`
	DeliveredAt   string `json:"delivered_at"`
	TransportType string `json:"transport_type"`
	UserId        int32  `json:"user_id"`
	Address       string `json:"address"`
}

func GetOrdersByUserId(res http.ResponseWriter, req *http.Request) {
	UserId, err := GetUserIdFromReq(res, req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r, err := OrderClient.GetOrdersByUserId(ctx, &opb.GetOrdersByUserIdRequest{UserId: int32(UserId)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Smt went wrong")
		return
	}

	orders := make([]OrderDTO, 0, len(r.Order))

	for _, v := range r.Order {
		orders = append(orders, OrderDTO{
			Id:            v.Id,
			UserId:        v.UserId,
			CartId:        v.CartId,
			Address:       v.Address,
			TransportType: v.TransportType,
			CreateAt:      v.CreateAt.AsTime().Format("2006-01-02 15:04:05"),
			DeliveredAt:   v.DeliveredAt.AsTime().Format("2006-01-02 15:04:05"),
		})
	}

	log.Printf("User By id=%v get his orders", UserId)
	httpresponse.WriteJSON(res, http.StatusOK, orders /* it is []*opb.Orders*/, "")
	return
}

func GetItemsOfOrderById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Smt went wrong!")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	r, err := OrderClient.GetItemsOfOrderById(ctx, &opb.GetItemsOfOrderByIdRequest{OrderId: int32(id)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Smt went wrong!")
		return
	}

	var products []models.Product
	for _, v := range r.ListOfProductsId {
		product, err := db.GetProduct(ctx, int(v))
		if err != nil {
			continue
		}
		products = append(products, product)
	}

	log.Printf("All product of oorder by id=%v were received", id)
	httpresponse.WriteJSON(res, http.StatusOK, products, "")
}
