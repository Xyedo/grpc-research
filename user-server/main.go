package main

import (
	"context"
	"fmt"
	"grpc-research/user-server/forum/user"
	userSvc "grpc-research/user-server/internal"
	"grpc-research/user-server/internal/entites"
	"grpc-research/user-server/internal/pgrepository"
	"grpc-research/user-server/internal/usecase"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	userSvc userSvc.Usecase
	user.UnimplementedUserServer
}

func (s *server) AddUser(ctx context.Context, payload *user.AddUserRequest) (*user.AddUserResponse, error) {
	resp, err := s.userSvc.AddUser(ctx,
		entites.User{
			Username: payload.GetUsername(),
			Email:    payload.GetEmail(),
			Password: payload.GetPassword(),
		})
	if err != nil {
		return nil, handleError(ctx, err)
	}
	return &user.AddUserResponse{Id: resp}, nil
}
func (s *server) GetUser(ctx context.Context, payload *user.GetUserByNameOrEmail) (*user.GetUserResponse, error) {
	resp, err := s.userSvc.GetUser(ctx, entites.UserIdentifier{
		Username: payload.GetUsername(),
		Email:    payload.GetEmail(),
	})
	if err != nil {
		return nil, handleError(ctx, err)
	}
	return &user.GetUserResponse{
		Username:       resp.Username,
		Email:          resp.Email,
		Hashedpassword: resp.Password,
	}, nil
}
func (s *server) GetUserById(ctx context.Context, payload *user.GetUserByIdRequest) (*user.GetUserResponse, error) {
	resp, err := s.userSvc.GetUserById(ctx, payload.GetId())
	if err != nil {
		return nil, handleError(ctx, err)
	}
	return &user.GetUserResponse{
		Username:       resp.Username,
		Email:          resp.Email,
		Hashedpassword: resp.Password,
	}, nil
}

func handleError(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	log.Println(err)
	return err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", os.Getenv("OLTP_HOST"), os.Getenv("OLTP_PORT"))),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tracerProvider)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dbConn, err := pgrepository.OpenDb()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer dbConn.Close()
	userRepo := pgrepository.NewRepository(dbConn)
	userSvc := usecase.NewUserUseCase(userRepo)
	server := server{userSvc: userSvc}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	user.RegisterUserServer(srv, &server)
	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", srv)
	}
}
