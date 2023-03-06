package main

import (
	"fmt"
	"grpc-research/bff-rest/forum/auth"
	"grpc-research/bff-rest/forum/user"
	"log"
	"net/http"
	"os"
	"time"

	handler "grpc-research/bff-rest/internal/feature/delivery/http/v1"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	userClient, err := grpc.Dial(os.Getenv("USER_SERVICE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial user service client: %v", err)
	}
	defer userClient.Close()

	userServiceClient := user.NewUserClient(userClient)

	authClient, err := grpc.Dial(os.Getenv("AUTH_SERVICE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial user service client: %v", err)
	}
	defer authClient.Close()

	authServiceClient := auth.NewAuthClient(authClient)

	userH := handler.NewUserHandler(userServiceClient)
	authH := handler.NewAuthHandler(authServiceClient)

	handlers := handler.Routes(handler.Route{User: userH, Auth: authH})
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler:      handlers,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
