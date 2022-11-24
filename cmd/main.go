package main

import (
	"log"
	"strconv"

	"github.com/go-openapi/loads"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/DrownSelf/OrderService/pkg/grpc"

	"github.com/DrownSelf/DriverService/internal/auth"
	"github.com/DrownSelf/DriverService/internal/configs"
	"github.com/DrownSelf/DriverService/internal/handlers"
	"github.com/DrownSelf/DriverService/internal/repositories"
	"github.com/DrownSelf/DriverService/internal/services"
	"github.com/DrownSelf/DriverService/restapi"
	"github.com/DrownSelf/DriverService/restapi/operations"
)

func main() {
	log.Printf("Server started")
	config, err := configs.LoadConnectionConfig()
	if err != nil {
		log.Fatalf("error during loading config: %v", err)
	}

	forger := auth.NewJwt(config.SecretKey)

	repository, err := repositories.NewDriverRepository(*config)
	if err != nil {
		log.Fatalf("error during connecting db: %v", err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	api := operations.NewDriverServiceAPI(swaggerSpec)
	server := restapi.NewServer(api)

	server.Port, err = strconv.Atoi(config.AppPort)
	if err != nil {
		log.Fatalf("error during connecting to server: %v", err)
	}

	conn, err := grpc.Dial(config.GrpcClient, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can't connect to grpc server: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)
	service := services.NewDriverService(repository, &auth.Hasher{})
	handler := handlers.NewDriverHandler(handlers.HandlerDependencies{
		DriverService: service,
		OrderClient:   client,
		Forger:        forger,
		Config:        config,
	})
	handlers.ConfigureHandlers(api, handler)
	if err = server.Serve(); err != nil {
		log.Fatalln(err)
	}
	_ = server.Shutdown()
}
