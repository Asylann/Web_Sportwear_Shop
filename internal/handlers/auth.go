package handlers

import (
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/middleware"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func HashingToBytes(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func CompareHashedPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandle(res http.ResponseWriter, req *http.Request) {
	var r_user LoginReq
	if err := json.NewDecoder(req.Body).Decode(&r_user); err != nil {
		log.Println("Error during decode json req to struct")
		http.Error(res, "Error during decode json req to struct", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	userInDB, err := db.GetUserByEmail(ctx, r_user.Email)
	if err != nil {
		log.Println("Invalid email or Unauthorized")
		http.Error(res, "Invalid email or Unauthorized", http.StatusUnauthorized)
		return
	}

	enteredHashedPassword, err := HashingToBytes(r_user.Password)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", err.Error())
		return
	}

	if CompareHashedPassword(string(enteredHashedPassword), userInDB.Password) {
		log.Println("Invalid password")
		http.Error(res, "Invalid password", http.StatusUnauthorized)
		return
	}

	signedToken, err := middleware.Generate(userInDB.ID, userInDB.Email, userInDB.RoleId)
	if err != nil {
		log.Println("Error during creation of Token")
		http.Error(res, "Error during creation of Token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: true,
	})

	log.Printf("User by email = %v logged in!", userInDB.Email)
	httpresponse.WriteJSON(res, http.StatusOK, "Logged In!", "")
}

func ProviderLoginHandle(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

func ProviderLoggedInHandle(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	provider := vars["provider"]

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		log.Println("Error during get user info")
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Error during get user info")
		return
	}
	fmt.Println(user.Email)
	userEmail := user.Email
	if provider == "google" {
		verified, ok := user.RawData["email_verified"].(bool)
		if !ok {
			log.Println("Error during converting claims")
			httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Error during converting claims")
			return
		}

		if !verified {
			log.Println("Invalid Verifier")
			httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid Verifier")
			return
		}

		userEmail, ok = user.RawData["email"].(string)
		if !ok {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Error during converting email")
			return
		}
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()
	_, err = db.GetUserByEmail(ctx, userEmail)
	if err != nil {
		var u models.User
		if provider == "google" {
			u = models.User{Email: userEmail, Password: "nullByGoogle", RoleId: 1}
		} else {
			u = models.User{Email: userEmail, Password: "nullByGithub", RoleId: 1}
		}
		id, err := db.CreateUser(ctx, &u)
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Error during creation of user")
			return
		}

		_, err = c.CreateCart(ctx, &pb.CreateCartRequest{UserId: int32(id)})
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "smt went wrong")
			return
		}
		log.Printf("%v`s cart was created!!!")

		err = db.ChangeEtagVersionByName(ctx, "ListOfUsers")
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "smt went wrong")
			return
		}
		log.Println("Version of ListOfUsers was changed to +1")

		log.Printf("User was created = %v", u)
	}

	userInDB, _ := db.GetUserByEmail(ctx, userEmail)

	signedToken, err := middleware.Generate(userInDB.ID, userInDB.Email, userInDB.RoleId)
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Error during creation of Token")
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})

	log.Printf("User by email = %v logged in!", userInDB.Email)
	http.Redirect(res, req, "https://localhost:8081/pages/dashboard.html", http.StatusSeeOther)
}

func LogoutHandle(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)

	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Cannot find cookie or Unauthorized")
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	tokenStr := cookie.Value

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Error during loading config")
			return nil, err
		}
		secret := cfg.JWT_Secret
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusBadRequest, "", "Invalid or Expired token")
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println(err.Error())
		httpresponse.WriteJSON(res, http.StatusInternalServerError, "", "Error during receiving claims")
		return
	}

	userEmail := claims["email"]
	log.Printf("User by email = %v logged out!", userEmail)
	httpresponse.WriteJSON(res, http.StatusOK, "Logged out!", "")
}
