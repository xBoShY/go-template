package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

func GetDefaultLocal() Template {
	return defaultConfig
}

func (cfg *Template) ResolveLogPaths(rootDir string) (liveLog, archive string) {
	// the default locations of log and archive are root
	liveLog = filepath.Join(rootDir, cfg.LogFileName)
	archive = filepath.Join(rootDir, cfg.LogArchiveName)

	// if LogFileDir is set, use it instead
	if cfg.LogFileDir != "" {
		liveLog = filepath.Join(cfg.LogFileDir, cfg.LogFileName)
		archive = filepath.Join(cfg.LogFileDir, cfg.LogArchiveName)
	}

	// if LogArchivePath is set, use it instead
	if cfg.LogArchiveDir != "" {
		archive = filepath.Join(cfg.LogArchiveDir, cfg.LogArchiveName)
	}

	return liveLog, archive
}

func ResolveDataDir(dataDirectory *string, flag string, envPrefix string) (string, error) {
	var dir string
	var envData = envPrefix + "_DATA"
	if dataDirectory == nil || *dataDirectory == "" {
		dir = os.Getenv(envData)
	} else {
		dir = *dataDirectory
	}

	absolutePath, absPathErr := filepath.Abs(dir)

	if len(dir) == 0 {
		//lint:ignore ST1005 Error message will be shown to the user
		return "", fmt.Errorf("Data directory not specified. Please use -%s or set $%s in your environment.\n", flag, envData)
	}

	if absPathErr != nil {
		//lint:ignore ST1005 Error message will be shown to the user
		return "", fmt.Errorf("Can't convert data directory's path to absolute, %v\n", dir)
	}

	if _, err := os.Stat(absolutePath); err != nil {
		//lint:ignore ST1005 Error message will be shown to the user
		return "", fmt.Errorf("Data directory %s does not appear to be valid\n", dir)
	}

	return absolutePath, nil
}

func LockFile(dir string, instance string) (*flock.Flock, error) {
	filename := instance + ".lock"
	lockPath := filepath.Join(dir, filename)
	fileLock := flock.New(lockPath)
	locked, err := fileLock.TryLock()
	if err != nil {
		//lint:ignore ST1005 Error message will be shown to the user
		return nil, fmt.Errorf("Unexpected failure in establishing %s: %s \n", filename, err.Error())
	}
	if !locked {
		//lint:ignore ST1005 Error message will be shown to the user
		return nil, fmt.Errorf("Failed to lock %s; is an instance of %s already running in this data directory?", filename, instance)
	}

	return fileLock, nil
}

func FormatVersionAndLicense() string {
	version := GetCurrentVersion()
	return fmt.Sprintf("%d\n (commit [%s]#%s)\n%s",
		version.BuildNumber,
		version.Branch,
		version.CommitHash,
		GetLicenseInfo(),
	)
}
