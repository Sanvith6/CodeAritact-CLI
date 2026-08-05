package config

// CurrentConfigVersion is the version of the config schema defined in this codebase.
const CurrentConfigVersion = 1

// Migrate upgrades a configuration struct to the current schema version.
func Migrate(cfg Config) (Config, error) {
	// If the version is 0 (missing), we assume it's v1 or a new config.
	if cfg.Version == 0 {
		cfg.Version = 1
	}

	if cfg.Version == CurrentConfigVersion {
		return cfg, nil
	}

	// Future migration logic:
	// if cfg.Version == 1 {
	//     cfg = migrateV1ToV2(cfg)
	// }
	// ... continue until CurrentConfigVersion

	return cfg, nil
}
