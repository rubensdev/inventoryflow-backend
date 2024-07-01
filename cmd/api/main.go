package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rubensdev/inventoryflow-backend/internal/healthcheck"
	userdom "github.com/rubensdev/inventoryflow-backend/internal/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr, found := os.LookupEnv("LISTEN_ADDR")
	if !found {
		addr = ":8080"
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	router := httprouter.New()

	// Healthcheck
	healthCheckHandler := healthcheck.NewHealthCheckHandler(logger, healthcheck.SystemInfo{
		Env:     os.Getenv("ENV"),
		Version: os.Getenv("VERSION"),
	})

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler.StatusHandler)

	inMemoryUserRepo := userdom.NewInMemoryUserRepository()
	usersHandler := userdom.NewUserHandler(logger, *userdom.NewUserService(inMemoryUserRepo))

	// Users
	router.HandlerFunc(http.MethodGet, "/v1/users", usersHandler.GetUsers)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", usersHandler.GetUserByID)
	router.HandlerFunc(http.MethodPost, "/v1/users", usersHandler.CreateUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", usersHandler.UpdateUserByID)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", usersHandler.DeleteUserByID)

	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
