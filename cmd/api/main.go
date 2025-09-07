package main

import (
	"WebSportwareShop/internal/cache"
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"WebSportwareShop/internal/handlers"
	"WebSportwareShop/internal/middleware"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env variables are loaded")
		return
	}
	// Loading config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Initialization of Db connection
	db.InitDB(cfg)
	defer db.CloseDB()

	// Initialization of OAuth and OpenID providers (Github.com and Gooogle.com)
	config.InitOAuthProviders()

	// Initialization of Redis connection
	err = cache.InitRedisConnection()
	if err != nil {
		log.Fatalf("Error during Init Redis: %v", err.Error())
		return
	}

	// Initialization of connection to Cart Microservice
	handlers.InitCartClientConnection()

	// Creation of all middlewares
	authMw := middleware.JWTAuth(cfg.JWT_Secret)
	RequiredCustomer := middleware.RequireRole(1, 2, 3)
	RequiredSeller := middleware.RequireRole(2, 3)
	RequiredAdmin := middleware.RequireRole(3)

	// Using gorilla/mux router and putting logging middleware to every request
	r := mux.NewRouter()
	r.Use(middleware.Logging)
	/*r.Use(middleware.CSP)*/

	// Creation of Each endpoints with access control
	r.HandleFunc("/signup", handlers.CreateUserHandle).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandle).Methods("POST")
	r.HandleFunc("/auth/{provider}/login", handlers.ProviderLoginHandle).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handlers.ProviderLoggedInHandle).Methods("GET")

	r.Handle("/logout", authMw(http.HandlerFunc(handlers.LogoutHandle))).Methods("POST")

	r.Handle("/carts", authMw(http.HandlerFunc(handlers.CreateCartHandle))).Methods("POST")

	r.Handle("/addToCart/{id}", authMw(http.HandlerFunc(handlers.AddToCartHandle))).Methods("POST")
	r.Handle("/myCart", authMw(http.HandlerFunc(handlers.GetItemsOfCartByIdHandle))).Methods("GET")
	r.Handle("/myCart/{id}", authMw(http.HandlerFunc(handlers.DeleteItemFromCartHandle))).Methods("DELETE")

	r.Handle("/me", authMw(http.HandlerFunc(handlers.GetInfoAboutMe))).Methods("GET")

	r.Handle("/products", http.HandlerFunc(handlers.ListOfProductsHandle)).Methods("GET")
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

	// Managing access options with rs/cors that allowed request only set sources
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Etag"},
		AllowCredentials: true,
	})

	// Giving our main route to the cors that will give access to only set sources
	handler := c.Handler(r)

	// Creation of our server with all timeouts and settings
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  130 * time.Second,
	}

	// channel for any stops from system user with Ctrl+C for gracefully shutdown process
	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// second version but 1-st one is much better, i think
	/*quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)*/

	// running our server on goroutine for waiting our system user stops
	go func() {
		log.Printf("HTTPS Server is running on :%s \n", cfg.Port)
		if err = srv.ListenAndServeTLS("localhost+2.pem", "localhost+2-key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	// Waiting till system user stops with Ctrl+C and channel get Done()
	<-quit.Done()
	log.Println("shutting down...")

	// This time for making all request that were while shutting down moment
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Forced Shutdown if time has passed
	if err = srv.Shutdown(ctx); err != nil {
		log.Println("Server shut down: %v", err.Error())
	}
	log.Println("Server exiting")
}
