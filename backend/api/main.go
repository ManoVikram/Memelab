package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ManoVikram/AI-Meme-Generator/backend/api/routes"
	"github.com/ManoVikram/AI-Meme-Generator/backend/api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ManoVikram/AI-Meme-Generator/backend/api/proto"
)

func main() {
	// Step 1 - Load .env file
	godotenv.Load("../../.env")

	// Step 2 - Connect to the Python gRPC server
	grpcAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:50051"
	}

	// Detect if running locally or on Cloud Run
	// K_SERVICE is an environment variable Cloud Run automatically injects into every container
	var creds grpc.DialOption
	if os.Getenv("K_SERVICE") != "" {
		creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	connection, err := grpc.NewClient(grpcAddress, creds, grpc.WithDefaultCallOptions(
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
	h2cServer := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(server, &http2.Server{}),
	}
	log.Fatal(h2cServer.ListenAndServe())
}
