package handlers

import (
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/middleware"
	"WebSportwareShop/internal/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func mdHashing(payload []byte) string {
	hasher := md5.New()
	sum := hasher.Sum([]byte(payload))
	return hex.EncodeToString(sum)
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

	hashedPassword, err := HashingToBytes(u.Password)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "Can not generate hash tp such password")
		return
	}

	u.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	id, err := db.CreateUser(ctx, &u)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	_, err = c.CreateCart(ctx, &pb.CreateCartRequest{UserId: int32(id)})
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "Can not create cart")
		return
	}
	log.Printf("%v`s cart was created!!!", u.Email)

	err = db.ChangeEtagVersionByName(ctx, "ListOfUsers")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "smt went wrong")
		return
	}
	log.Println("Version of ListOfUsers was changed to +1")

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

	err = DeleteCart(id)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "Can not delete the cart of user")
		return
	}

	err = db.ChangeEtagVersionByName(ctx, "ListOfUsers")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "smt went wrong")
		return
	}
	log.Println("Version of ListOfUsers was changed to +1")
	w.WriteHeader(http.StatusNoContent)
	log.Printf("User by id=%v was deleted! \n", id)
}

func ListOfUsersHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	version, err := db.GetEtagVersionByName(ctx, "ListOfUsers")
	if err != nil {
		log.Println(err.Error(), "here")
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	versionStr := "v" + strconv.Itoa(version)
	etag := `"` + mdHashing([]byte(versionStr)) + `"`

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "max-age=30 ,public")

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		log.Println("List of users were received By http caching")
		return
	}

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
	err = db.ChangeEtagVersionByName(ctx, "ListOfUsers")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "", "smt went wrong")
		return
	}
	log.Println("Version of ListOfUsers was changed to +1")

	httpresponse.WriteJSON(w, http.StatusOK, u, "")
	log.Printf("User by id = %v was updated : %v", id, u)
}

func GetUserEmailHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	var email string
	email, err = db.GetUserEmail(ctx, id)
	if err != nil {
		log.Fatal(err.Error())
		httpresponse.WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, email, "")
	log.Printf("User`s email by id = %v was recieved \n", id)
}

func GetInfoAboutMe(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "not found Cookie")
		return
	}

	token := cookie.Value

	idItf, err := middleware.GetClaimFromToken(token, "sub")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "not found claims")
		return
	}
	role_idItf, err := middleware.GetClaimFromToken(token, "role_id")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "not found claims")
		return
	}
	emailItf, err := middleware.GetClaimFromToken(token, "email")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "not found claims")
		return
	}

	idFloat := idItf.(float64)
	roleIdFloat := role_idItf.(float64)
	email := emailItf.(string)

	id := int(idFloat)
	roleId := int(roleIdFloat)

	userInfo := models.UserInfo{
		ID:     id,
		Email:  email,
		RoleId: roleId,
	}

	log.Printf("Info about user = %v was received!!!", email)
	httpresponse.WriteJSON(res, http.StatusOK, userInfo, "")
}
