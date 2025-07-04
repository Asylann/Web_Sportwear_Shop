package handlers

import (
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	userInDB, err := db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	if userInDB.Password != req.Password {
		httpresponse.WriteJSON(w, http.StatusUnauthorized, nil, "invalid password")
		return
	}

	claims := jwt.MapClaims{
		"sub":     userInDB.ID,
		"role_id": userInDB.RoleId,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	log.Printf("User by email: %v loged in!", userInDB.Email)
	httpresponse.WriteJSON(w, http.StatusOK, map[string]string{"token": signed}, "")
}

func CreateUserHandle(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	users, err := db.ListOfUsers(context.Background())
	for _, user := range users {
		if user.Email == u.Email {
			httpresponse.WriteJSON(w, http.StatusConflict, nil, "Already exists")
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.CreateUser(ctx, &u)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusCreated, u.Email, "")
	log.Printf("User was created! : %v \n", u)
}

func GetUserHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var u models.User
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	u, err = db.GetUser(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusUnauthorized, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, u, "")
	log.Printf("User by id=%v was recieved! \n", id)
}

func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err = db.DeleteUser(ctx, id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Printf("User by id=%v was deleted! \n", id)
}

func ListOfUsersHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	users, err := db.ListOfUsers(ctx)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, users, "")
	log.Printf("All Users list were received!!!")
}

func UpdateUserHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var u models.User
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	u.ID = id
	err = db.UpdateUser(ctx, &u)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, u, "")
	log.Printf("User by id = %v was updated : %v", id, u)
}
