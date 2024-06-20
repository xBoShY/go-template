package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Version struct {
	Name        string
	BuildNumber int
	CommitHash  string
	Branch      string
	License     string
	Repo        string
}

func (v Version) String() string {
	return fmt.Sprintf("%d-[%s]#%s", v.BuildNumber, v.Branch, v.CommitHash)
}

func convertToInt(val string) int {
	if val == "" {
		return 0
	}
	value, _ := strconv.ParseInt(val, 10, 0)
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
		BuildNumber: convertToInt(BuildNumber), // set using -ldflags
		CommitHash:  CommitHash,
		Branch:      Branch,
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
