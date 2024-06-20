package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xboshy/go-template/config"
	server "github.com/xboshy/go-template/daemon/rest-server"
	"github.com/xboshy/go-template/logging"
)

const envPrefix = "REST_DATA"

var logFileName = fmt.Sprintf("%s.log", config.GetCurrentVersion().Name)
var logArchiveName = fmt.Sprintf("%s.archive.log", config.GetCurrentVersion().Name)

var dataDirectory = flag.String("d", "", "Daemon data path")
var versionCheck = flag.Bool("v", false, "Display and write current build version and exit")
var logToStdout = flag.Bool("o", false, fmt.Sprintf("Write to stdout instead of %s by overriding config.LogSizeLimit to 0", logFileName))
var listenIP = flag.String("l", "", "Override config.EndpointAddress (REST listening address) with ip:port")

func main() {
	flag.Parse()

	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	var err error
	cfg := config.GetDefaultLocal()

	cfg.LogFileName = logFileName
	cfg.LogArchiveName = logArchiveName

	if *versionCheck {
		fmt.Println(config.FormatVersionAndLicense())
		return 0
	}

	dataDir, err := config.ResolveDataDir(dataDirectory, "d", envPrefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fileLock, err := config.LockFile(dataDir, "rest-server")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer fileLock.Unlock()

	log := logging.Base()
	if logToStdout != nil && *logToStdout {
		cfg.LogSizeLimit = 0
	}

	s := server.Server{
		RootPath: dataDir,
	}

	if *listenIP != "" {
		cfg.EndpointAddress = *listenIP
	}

	err = s.Initialize(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Error(err)
		return 1
	}

	s.Start()

	return 0
}
