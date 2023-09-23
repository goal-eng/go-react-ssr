package go_ssr

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/hot_reload"
	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/internal/type_converter"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

// Init starts the Go SSR plugin
func Init(optionalCfg ...config.Config) error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logger.L.Err(err).Msg("Error loading .env file")
	}

	// Initialize logger globally
	logger.Init()
	// Get config if it was passed to the function
	cfg := getConfig(optionalCfg)
	// Set default values for config
	if err := defaults.Set(cfg); err != nil {
		logger.L.Err(err).Msg("Failed to set defaults")
		return err
	}

	// Load config in to global variable
	if err := config.Load(*cfg); err != nil {
		logger.L.Err(err).Msg("Failed to load config")
		return err
	}

	// Compile the global css file if it exists
	if config.C.GlobalCSSFilePath != "" {
		if err := react_renderer.BuildGlobalCSSFile(); err != nil {
			logger.L.Err(err).Msg("Failed to build global css file")
			return err
		}
	}

	// If running in production mode, return and dont start hot reload or type converter
	if os.Getenv("APP_ENV") == "production" {
		logger.L.Info().Msg("Running in production mode")
		return nil
	}
	logger.L.Info().Msg("Running in development mode")
	logger.L.Debug().Msg("Starting type converter")

	// Start the type converter to convert Go types to Typescript types
	if err := type_converter.Init(); err != nil {
		logger.L.Err(err).Msg("Failed to init type converter")
		return err
	}

	logger.L.Debug().Msg("Starting hot reload")
	// Watches for changes in the frontend directory & starts a websocket server to send updates to the browser
	hot_reload.Init()
	return nil
}

// getConfig returns the config if it was passed to the function, otherwise it returns a default config
func getConfig(optionalCfg []config.Config) (cfg *config.Config) {
	if len(optionalCfg) > 0 {
		cfg = &optionalCfg[0]
	} else {
		logger.L.Info().Msg("No config provided, using defaults")
		cfg = &config.Config{}
	}
	return cfg
}
