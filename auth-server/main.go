package main

import (
	"context"
	"fmt"
	"grpc-research/auth-server/forum/auth"
	"grpc-research/auth-server/forum/user"
	authSvc "grpc-research/auth-server/internal"
	"grpc-research/auth-server/internal/entities"
	"grpc-research/auth-server/internal/jwt"
	"grpc-research/auth-server/internal/pgrepository"
	"grpc-research/auth-server/internal/usecase"
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	auth authSvc.UseCase
	auth.UnimplementedAuthServer
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if req.GetEmail() != "" {
		accessToken, refreshToken, err := s.auth.LoginWithEmail(ctx, entities.EmailLogin{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		})
		if err != nil {
			return nil, handleError(ctx, err)
		}
		return &auth.LoginResponse{AccessKey: accessToken, RefreshKey: refreshToken}, nil
	}
	accessToken, refreshToken, err := s.auth.LoginWithUsername(ctx, entities.UsernameLogin{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, handleError(ctx, err)
	}
	return &auth.LoginResponse{AccessKey: accessToken, RefreshKey: refreshToken}, nil

}
func (s *server) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	err := s.auth.Logout(ctx, entities.RefreshToken{RefreshToken: req.GetRefreshKey()})
	if err != nil {
		return nil, handleError(ctx, err)
	}
	return &auth.LogoutResponse{}, nil
}

func (s *server) RefreshAccess(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	accessToken, err := s.auth.RefreshAccess(ctx, entities.RefreshToken{
		RefreshToken: req.GetRefreshKey(),
	})
	if err != nil {
		return nil, err
	}
	return &auth.RefreshResponse{AccessKey: accessToken}, nil
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
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", os.Getenv("OTLP_HOST"), os.Getenv("OTLP_PORT"))),
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

	client, err := grpc.Dial(os.Getenv("USER_SERVICE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial client: %v", err)
	}
	defer client.Close()

	userClient := user.NewUserClient(client)

	connDb, err := pgrepository.OpenDb()
	if err != nil {
		log.Fatalf("failed to openDb: %v", err)
	}
	defer connDb.Close()

	tokenizer := jwt.NewJwt(os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"), os.Getenv("ACCESS_EXPIRES"), os.Getenv("REFRESH_EXPIRES"))
	authRepo := pgrepository.NewAuthRepo(connDb)
	authUC := usecase.NewAuthUseCase(userClient, tokenizer, authRepo)

	var server server
	server.auth = authUC

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor())),
	)
	auth.RegisterAuthServer(srv, &server)

	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", srv)
	}

}
