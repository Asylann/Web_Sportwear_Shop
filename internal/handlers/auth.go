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
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
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

	signed, err := middleware.Generate(userInDB.ID, userInDB.Email, userInDB.RoleId)
	if err != nil {
		httpresponse.WriteJSON(w, http.StatusInternalServerError, nil, "During creation of token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signed,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("User by email: %v loged in!", userInDB.Email)
	httpresponse.WriteJSON(w, http.StatusOK, "Logged in!", "")
}

func ProviderLoginHandle(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

func ProviderLoggedInHandle(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	provider := vars["provider"]

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	userEmail := user.Email

	if provider == "google" {
		verified, ok := user.RawData["email_verified"].(bool)
		if !ok {
			log.Println("Error during converting email_verified")
			http.Error(res, "Error during converting email", http.StatusInternalServerError)
			return
		}

		if verified != true {
			log.Println("Not verified email")
			http.Error(res, "Not verified email", http.StatusBadRequest)
			return
		}

		userEmail, ok = user.RawData["email"].(string)
		if !ok {
			log.Println("Error during converting email")
			http.Error(res, "Error during converting email", http.StatusInternalServerError)
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
		err = db.CreateUser(ctx, &u)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, "Error during creation of User", http.StatusInternalServerError)
		}
		log.Printf("User was created %v \n", u)
	}

	userInDB, _ := db.GetUserByEmail(context.Background(), userEmail)
	tokenStr, err := middleware.Generate(userInDB.ID, userInDB.Email, userInDB.RoleId)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Error during creation of token", http.StatusInternalServerError)
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenStr,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("User by email %v logged in!", userInDB.Email)
	http.Redirect(res, req, "http://localhost:8081/pages/dashboard.html", http.StatusSeeOther)
}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("Cookie not found")
		httpresponse.WriteJSON(w, http.StatusBadRequest, "", "Cookie not found")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	tokenStr := cookie.Value

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid Claims", http.StatusBadRequest)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}

	log.Printf("User by email= %v logged out!", email)
	httpresponse.WriteJSON(w, http.StatusOK, "Logged out!", "")
}
