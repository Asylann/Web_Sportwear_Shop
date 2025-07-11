package handlers

import (
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/httpresponse"
	"WebSportwareShop/internal/middleware"
	"WebSportwareShop/internal/models"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
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

func GoogleLoginHandle(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func GoogleLoggedInHandle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	code := r.URL.Query().Get("code")

	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	client := config.GoogleOAuthConfig.Client(ctx, token)
	resp, _ := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	defer resp.Body.Close()

	var profile struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&profile)

	ctx, cancel = context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	user, err := db.GetUserByEmail(ctx, profile.Email)
	IsCreatedUser := false
	if err != nil {
		var u = models.User{
			Email:    profile.Email,
			Password: "nullByGoogle",
			RoleId:   1,
		}
		err = db.CreateUser(ctx, &u)
		if err != nil {
			log.Println(err.Error())
			httpresponse.WriteJSON(w, http.StatusBadRequest, nil, err.Error())
			return
		}
		log.Printf("User was created! : %v \n", u)
		IsCreatedUser = true
	}

	jwtToken, err := middleware.Generate(user.ID, profile.Email, 1)
	if err != nil {
		log.Println("JWT token generation failed!")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("User by email: %v loged in!", profile.Email)
	http.Redirect(w, r, "http://localhost:8081/pages/dashboard.html", http.StatusSeeOther)

	if !IsCreatedUser {
		httpresponse.WriteJSON(w, http.StatusOK, "Logged in!", "")
	} else {
		httpresponse.WriteJSON(w, http.StatusOK, "Singed up!", "")
	}
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
