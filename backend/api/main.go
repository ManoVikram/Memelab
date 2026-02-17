package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/ManoVikram/AI-Meme-Generator/backend/api/routes"
	"github.com/ManoVikram/AI-Meme-Generator/backend/api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcoauth "google.golang.org/grpc/credentials/oauth"

	"context"

	pb "github.com/ManoVikram/AI-Meme-Generator/backend/api/proto"
)

type cloudRunTokenSource struct {
	audience string
}

func (ts cloudRunTokenSource) Token() (*oauth2.Token, error) {
	idToken, err := metadata.GetWithContext(context.Background(), "instance/service-accounts/default/identity?audience="+ts.audience+"&format=full")
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{AccessToken: idToken}, nil
}

func main() {
	// Step 1 - Load .env file
	godotenv.Load("../../.env")

	// Step 2 - Connect to the Python gRPC server
	grpcAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:50051"
	}

	// Detect if running locally or on Cloud Run
	var dialOptions []grpc.DialOption
	dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(50*1024*1024),
		grpc.MaxCallSendMsgSize(50*1024*1024),
	))

	if os.Getenv("K_SERVICE") != "" {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))

		audience := "https://" + os.Getenv("GRPC_SERVER_ADDRESS")[:strings.Index(os.Getenv("GRPC_SERVER_ADDRESS"), ":")]

		tokenSource := cloudRunTokenSource{audience: audience}

		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(grpcoauth.TokenSource{TokenSource: tokenSource}))
	} else {
		// Running locally - use insecure
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	connection, err := grpc.NewClient(grpcAddress, dialOptions...)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer connection.Close()

	// Initialize the AI service client
	aiServicesClient := &services.Services{
		AIClient: pb.NewAIMemeGeneratorServiceClient(connection),
	}

	// Initialize and set up the Gin server
	server := gin.Default()
	server.RedirectTrailingSlash = true

	// Register routes
	routes.RegisterRoutes(server, aiServicesClient)

	// Start the Gin server with HTTP/2 cleartext support
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
