package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Version struct {
	Name        string
	Major       int
	Minor       int
	BuildNumber int
	Branch      string
	CommitHash  string
	License     string
	Repo        string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.BuildNumber)
}

func convertToInt(val string) int {
	if val == "" {
		return 0
	}

	value, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		return 0
	}

	return int(value)
}

var currentVersion = func() Version {
	absPath, err := os.Executable()
	name := filepath.Base(absPath)
	if err != nil {
		name = ""
	}

	return Version{
		Name:        name,
		Major:       convertToInt(VersionMajor),
		Minor:       convertToInt(VersionMinor),
		BuildNumber: convertToInt(BuildNumber),
		Branch:      Branch,
		CommitHash:  CommitHash,
		License:     License,
		Repo:        Repo,
	}
}()

func GetCurrentVersion() Version {
	return currentVersion
}

func GetLicenseInfo() string {
	name := currentVersion.Name
	if name != "" {
		name += " is "
	}

	return fmt.Sprintf(
		"%s licensed with %s\nsource code available at %s",
		name,
		currentVersion.License,
		currentVersion.Repo,
	)
}
