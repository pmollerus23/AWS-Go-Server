package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/pmollerus23/go-aws-server/internal/aws"
	"github.com/pmollerus23/go-aws-server/internal/config"
	"github.com/pmollerus23/go-aws-server/internal/server"

	_ "github.com/pmollerus23/go-aws-server/docs" // Swagger docs
)

//	@title						AWS Go Server API
//	@version					1.0
//	@description				A production-grade Go web server with AWS Cognito authentication and AWS service integration.
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				https://github.com/pmollerus23/go-aws-server
//	@contact.email				support@example.com
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						localhost:8080
//	@BasePath					/
//	@schemes					http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer token authentication via AWS Cognito. Use format: "Bearer {access_token}"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Info("configuration loaded",
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
		"aws_region", cfg.AWS.Region,
	)

	// Initialize AWS clients
	awsClients, err := aws.NewClients(ctx, logger, cfg.AWS)
	if err != nil {
		return fmt.Errorf("failed to initialize AWS clients: %w", err)
	}

	// Create and run server
	srv := server.New(logger, cfg, awsClients)
	return srv.Run(ctx)
}
