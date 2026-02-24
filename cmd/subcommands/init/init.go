package init

import (
	"fmt"
	"os"
	"zuno/cmd/config"
	"zuno/cmd/generators/barebone"
	"zuno/cmd/generators/fiber"
	"zuno/cmd/generators/mongodb"
	"zuno/cmd/utils"
	"zuno/pkg/logger"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	packageName     string
	databaseAdapter string
	httpProvider    string
	grpcProvider    string
	wsProvider      string
)

var directCloneDirs = []string{"locales"}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		var fields []huh.Field
		// Only ask if flag NOT provided
		if !cmd.Flags().Changed("package") {
			fields = append(fields,
				huh.NewInput().
					Title("What's the package name gonna be?").
					Placeholder("github.com/<username>/<package-name>").
					Value(&packageName),
			)
		}

		if !cmd.Flags().Changed("database-adapter") {
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Choose database adapter").
					Options(
						huh.NewOption("mongodb", "mongodb"),
						huh.NewOption("postgres", "postgres"),
					).
					Value(&databaseAdapter),
			)
		}

		if !cmd.Flags().Changed("http-provider") {
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Select HTTP provider").
					Options(
						huh.NewOption("fiber", "fiber"),
						huh.NewOption("gin", "gin"),
						huh.NewOption("none", "none"),
					).
					Value(&httpProvider),
			)
		}
		if !cmd.Flags().Changed("grpc-provider") {
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Select gRPC provider").
					Options(
						huh.NewOption("grpc", "grpc"),
						huh.NewOption("gin", "gin"),
						huh.NewOption("none", "none"),
					).
					Value(&grpcProvider),
			)
		}
		if !cmd.Flags().Changed("ws-provider") {
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Select WebSocket provider").
					Options(
						huh.NewOption("gorilla", "gorilla"),
						huh.NewOption("none", "none"),
					).
					Value(&wsProvider),
			)
		}

		// Run form only if something needs prompting
		if len(fields) > 0 {
			if err := huh.NewForm(
				huh.NewGroup(fields...),
			).Run(); err != nil {
				return err
			}
		}
		config := &config.Config{
			Package: packageName,
			Adapter: config.Adapter{
				Database: config.DatabaseAdapter{
					Enabled:  true,
					Provider: databaseAdapter,
				},
			},
			Transport: config.Transport{
				Http: config.HttpTransport{
					Enabled:  httpProvider != "none",
					Provider: httpProvider,
				},
				Grpc: config.GrpcTransport{
					Enabled:  grpcProvider != "none",
					Provider: grpcProvider,
				},
				Ws: config.WsTransport{
					Enabled:  wsProvider != "none",
					Provider: wsProvider,
				},
			},
		}
		defer func() {
			logger.Info(fmt.Sprintf(`
Re-run with : zuno init --package=%s \
--database-adapter=%s \
--http-provider=%s \
--grpc-provider=%s \
--ws-provider=%s`, packageName, databaseAdapter, httpProvider, grpcProvider, wsProvider))
		}()

		return initializeNewProject(config, verbose)
	},
}

func Cmd() *cobra.Command {
	return initCmd
}

func initializeNewProject(config *config.Config, verbose bool) error {
	// At this point ALL values are guaranteed
	// packageName, adapter, transports are set
	// 1 go mod init
	logger.Info("Initializing new go module...")
	if err := utils.RunCmd("go", verbose, "mod", "init", config.Package); err != nil {
		return err
	}

	logger.Info("Initializing barebone...")
	if err := barebone.Initialize(config); err != nil {
		return err
	}

	if config.Transport.Http.Enabled {
		switch config.Transport.Http.Provider {
		case "fiber":
			logger.Info("Initializing fiber...")
			if err := fiber.Initialize(config, verbose); err != nil {
				return err
			}
		}
	}

	if config.Adapter.Database.Enabled {
		switch config.Adapter.Database.Provider {
		case "mongodb":
			logger.Info("Initializing mongodb...")
			if err := mongodb.Initialize(config); err != nil {
				return err
			}
			// case "postgres":
			// 	if err := postgres.InitializePostgres(config); err != nil {
			// 		return err
			// 	}
		}

	}

	//  go mod tidy
	logger.Info("Running go mod tidy...")
	if err := utils.RunCmd("go", verbose, "mod", "tidy"); err != nil {
		return err
	}

	conf, err := config.ToYaml()
	if err != nil {
		return err
	}

	logger.Info("Writing zuno.yml...")
	if err := os.WriteFile("zuno.yml", conf, 0644); err != nil {
		return err
	}

	return nil
}

func init() {

	initCmd.Flags().StringVar(
		&packageName,
		"package",
		"",
		"Go module package name",
	)

	initCmd.Flags().StringVar(
		&databaseAdapter,
		"database-adapter",
		"",
		"Database adapter (mongodb|postgres)",
	)

	initCmd.Flags().StringVar(
		&httpProvider,
		"http-provider",
		"",
		"HTTP provider (fiber|gin)",
	)

	initCmd.Flags().StringVar(
		&grpcProvider,
		"grpc-provider",
		"",
		"gRPC provider (grpc|grpc-go)",
	)

	initCmd.Flags().StringVar(
		&wsProvider,
		"ws-provider",
		"",
		"WebSocket provider (gorilla|gorilla-mux)",
	)
}
