package container

import (
	"log/slog"

	"github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/features"
	"github.com/acme-corp/rs-cli/internal/network"
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/platform/aws"
	"github.com/acme-corp/rs-cli/internal/platform/uv"
	"github.com/acme-corp/rs-cli/internal/runner"
	"github.com/acme-corp/rs-cli/internal/services"
	"github.com/acme-corp/rs-cli/internal/state"
	"github.com/acme-corp/rs-cli/internal/telemetry"
	"github.com/acme-corp/rs-cli/internal/ui"
	"github.com/acme-corp/rs-cli/internal/validator"
)

// Container holds all instantiated dependencies for the CLI application.
type Container struct {
	Config           *config.Config
	Logger           *slog.Logger
	Platform         *platform.Info
	Features         *features.Registry
	State            *state.Manager
	Telemetry        telemetry.Tracker
	Runner           runner.Runner
	Network          network.Checker
	Validator        validator.Validator
	UIManager        *ui.Manager

	// Providers
	AuthProvider    platform.AuthProvider
	RepoProvider    platform.RepositoryProvider
	PackageProvider platform.PackageManagerProvider

	// Services
	AuthService    *services.AuthService
	RepoService    *services.RepositoryService
	EnvService     *services.EnvironmentService
}

// Build assembles the dependencies into a Container.
func Build(cfg *config.Config, logger *slog.Logger, isNonInteractive bool) (*Container, error) {
	c := &Container{
		Config: cfg,
		Logger: logger,
	}

	// Basic infra
	pDetector := platform.NewDetector()
	pInfo, err := pDetector.Detect()
	if err != nil {
		return nil, err
	}
	c.Platform = pInfo
	c.Features = features.NewRegistry(cfg.Features)
	
	statePath := pInfo.RSConfigDir + "/state.json"
	c.State = state.NewManager(statePath)
	
	c.Telemetry = telemetry.NewNoOpTracker()
	c.Runner = runner.NewExecRunner(logger)
	c.Network = network.NewHTTPChecker()
	c.Validator = validator.NewSystemValidator(c.Runner, c.Platform, c.Features)
	c.UIManager = ui.NewManager(isNonInteractive)

	// Providers
	c.AuthProvider = aws.NewSSOAuthProvider(cfg, pInfo, c.Runner, logger)
	c.RepoProvider = aws.NewCodeArtifactProvider(cfg)
	c.PackageProvider = uv.NewUVPackageManager(cfg, pInfo)

	// Services
	c.AuthService = services.NewAuthService(c.AuthProvider, c.State, logger)
	c.RepoService = services.NewRepositoryService(c.AuthService, c.RepoProvider, c.State, logger)
	c.EnvService = services.NewEnvironmentService(pInfo, logger)

	return c, nil
}
