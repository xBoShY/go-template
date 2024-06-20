package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xboshy/go-deadlock"
	"github.com/xboshy/go-template/config"
	"github.com/xboshy/go-template/daemon/rest-server/api"
	"github.com/xboshy/go-template/logging"
	"github.com/xboshy/go-template/node"
)

var server http.Server

const maxHeaderBytes = 4096

type ServerNode interface {
	api.APINodeInterface
	Start()
	Stop()
}

type Server struct {
	RootPath string
	log      logging.Logger
	node     ServerNode
	stopping chan struct{}
}

func (s *Server) Initialize(cfg config.Template) error {
	s.log = logging.Base()
	liveLog, archive := cfg.ResolveLogPaths(s.RootPath)

	var maxLogAge time.Duration
	var err error
	if cfg.LogArchiveMaxAge != "" {
		maxLogAge, err = time.ParseDuration(cfg.LogArchiveMaxAge)
		if err != nil {
			s.log.Fatalf("invalid config LogArchiveMaxAge: %s", err)
			maxLogAge = 0
		}
	}

	var logWriter io.Writer
	if cfg.LogSizeLimit > 0 {
		fmt.Println("Logging to: ", liveLog)
		logWriter = logging.MakeCyclicFileWriter(liveLog, archive, cfg.LogSizeLimit, maxLogAge)
	} else {
		fmt.Println("Logging to: stdout")
		logWriter = os.Stdout
	}
	s.log.SetOutput(logWriter)
	s.log.SetJSONFormatter()
	s.log.SetLevel(logging.Level(cfg.BaseLoggerDebugLevel))
	s.log.Infof("LogLevel: %v", s.log.GetLevel())

	logging.SetupDeadlockLogger(s.log)
	deadlockState := "enabled"
	if deadlock.Opts.Disable {
		deadlockState = "disabled"
	}
	s.log.Infof("Deadlock detection is set to: %s (Default state is '%s')", deadlockState, config.DefaultDeadlock)

	var serverNode ServerNode
	thisNode, err := node.MakeNode(s.log, s.RootPath, cfg)
	serverNode = api.APINode{Node: thisNode}
	if os.IsNotExist(err) {
		return fmt.Errorf("node has not been installed: %s", err)
	}
	if err != nil {
		return fmt.Errorf("couldn't initialize the node: %s", err)
	}
	s.node = serverNode

	logging.RegisterExitHandler(s.Stop)

	return nil
}

func makeListener(addr string) (net.Listener, error) {
	var listener net.Listener
	var err error
	if (addr == "127.0.0.1:0") || (addr == ":0") {
		// if port 0 is provided, prefer port 8080 first, then fall back to port 0
		preferredAddr := strings.Replace(addr, ":0", ":8080", -1)
		listener, err = net.Listen("tcp", preferredAddr)
		if err == nil {
			return listener, err
		}
	}
	// err was not nil or :0 was not provided, fall back to originally passed addr
	return net.Listen("tcp", addr)
}

func (s *Server) Start() {
	version := config.GetCurrentVersion()

	s.log.Infof("Trying to start a %s node", version.Name)
	s.node.Start()
	s.log.Infof("Successfully started a %s node.", version.Name)
	cfg := s.node.Config()

	s.stopping = make(chan struct{})

	addr := cfg.EndpointAddress
	if addr == "" {
		addr = ":http"
	}

	listener, err := makeListener(addr)
	if err != nil {
		fmt.Printf("Could not create the listener: %v\n", err)
		os.Exit(1)
	}

	addr = listener.Addr().String()
	server = http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(cfg.RestReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(cfg.RestWriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	e := api.NewRouter(
		s.log, s.node, s.stopping, listener,
	)
	errChan := make(chan error, 1)

	e.Logger = s.log.MakeEchoLogger()

	go func() {
		err := e.StartServer(&server)
		errChan <- err
	}()

	// Handle signals cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGHUP)

	fmt.Printf("Node running and accepting RPC requests over HTTP on port %v. Press Ctrl-C to exit\n", addr)
	select {
	case err := <-errChan:
		if err != nil {
			s.log.Warn(err)
		} else {
			s.log.Info("Node exited successfully")
		}
		s.Stop()
	case sig := <-c:
		s.log.Infof("Exiting on %v", sig)
		s.Stop()
		os.Exit(0)
	}
}

func (s *Server) Stop() {
}
