package main

import (
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/handlers"
	"WebSportwareShop/internal/middleware"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.InitDB(cfg)
	config.InitOAuthProviders()
	defer db.CloseDB()

	authMw := middleware.JWTAuth(cfg.JWT_Secret)
	RequiredCustomer := middleware.RequireRole(1, 2, 3)
	RequiredSeller := middleware.RequireRole(2, 3)
	RequiredAdmin := middleware.RequireRole(3)

	r := mux.NewRouter()
	r.Use(middleware.Logging)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // or []string{"*"} for dev
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	r.HandleFunc("/signup", handlers.CreateUserHandle).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandle).Methods("POST")
	r.HandleFunc("/auth/{provider}/login", handlers.ProviderLoginHandle).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handlers.ProviderLoggedInHandle).Methods("GET")

	r.Handle("/logout", authMw(http.HandlerFunc(handlers.LogoutHandle))).Methods("POST")

	r.Handle("/products", authMw(RequiredCustomer(http.HandlerFunc(handlers.ListOfProductsHandle)))).Methods("GET")
	r.Handle("/products/{id}", authMw(RequiredCustomer(http.HandlerFunc(handlers.GetProductHandle)))).Methods("GET")
	r.Handle("/productsByCategory/{id}", authMw(RequiredCustomer(http.HandlerFunc(handlers.ListOfProductsByCategory)))).Methods("GET")
	r.Handle("/productsBySeller/{id}", authMw(RequiredCustomer(http.HandlerFunc(handlers.ListOfProductsBySellerID)))).Methods("GET")

	r.Handle("/products/{id}", authMw(RequiredSeller(http.HandlerFunc(handlers.DeleteProductHandle)))).Methods("DELETE")
	r.Handle("/products/{id}", authMw(RequiredSeller(http.HandlerFunc(handlers.UpdateProductHandle)))).Methods("PUT")
	r.Handle("/products", authMw(RequiredSeller(http.HandlerFunc(handlers.CreateProductHandle)))).Methods("POST")

	r.Handle("/categories", authMw(RequiredCustomer(http.HandlerFunc(handlers.ListOfCategoriesHandle)))).Methods("GET")
	r.Handle("/categories/{id}", authMw(RequiredCustomer(http.HandlerFunc(handlers.GetCategoryHandle)))).Methods("GET")

	r.Handle("/categories/{id}", authMw(RequiredSeller(http.HandlerFunc(handlers.DeleteCategoryHandle)))).Methods("DELETE")
	r.Handle("/categories/{id}", authMw(RequiredSeller(http.HandlerFunc(handlers.UpdateCategoryHandle)))).Methods("PUT")
	r.Handle("/categories", authMw(RequiredSeller(http.HandlerFunc(handlers.CreateCategoryHandle)))).Methods("POST")

	r.Handle("/users/email/{id}", authMw(RequiredCustomer(http.HandlerFunc(handlers.GetUserEmailHandle)))).Methods("GET")

	r.Handle("/users", authMw(RequiredAdmin(http.HandlerFunc(handlers.ListOfUsersHandle)))).Methods("GET")
	r.Handle("/users/{id}", authMw(RequiredAdmin(http.HandlerFunc(handlers.GetUserHandle)))).Methods("GET")
	r.Handle("/users/{id}", authMw(RequiredAdmin(http.HandlerFunc(handlers.DeleteUserHandle)))).Methods("DELETE")
	r.Handle("/users/{id}", authMw(RequiredAdmin(http.HandlerFunc(handlers.UpdateUserHandle)))).Methods("PUT")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  130 * time.Second,
	}

	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	/*quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)*/

	go func() {
		log.Printf("Server is running on :%s \n", cfg.Port)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-quit.Done()
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Println("Server shut down: %v", err.Error())
	}
	log.Println("Server exiting")
}
