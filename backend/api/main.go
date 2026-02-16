package main

import (
	"log"
	"os"

	"github.com/ManoVikram/AI-Meme-Generator/backend/api/routes"
	"github.com/ManoVikram/AI-Meme-Generator/backend/api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ManoVikram/AI-Meme-Generator/backend/api/proto"
)

func main() {
	// Step 1 - Connect to the Python gRPC server
	godotenv.Load("../../.env")
	grpcAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:50051"
	}
	connection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(50*1024*1024), // 50MB
		grpc.MaxCallSendMsgSize(50*1024*1024),
	))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer connection.Close()

	// Step 2 - Initialize the AI service client
	aiServicesClient := &services.Services{
		AIClient: pb.NewAIMemeGeneratorServiceClient(connection),
	}

	// Step 2 - Initialize and set up the Gin server
	server := gin.Default()

	server.RedirectTrailingSlash = true

	// Step 3 - Register the routes to the Gin server
	routes.RegisterRoutes(server, aiServicesClient)

	// Step 4 - Start the Gin server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(server.Run(":" + port))
}
