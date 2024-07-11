package config

type Template struct {
	BaseLoggerDebugLevel uint32
	LogFileDir           string
	LogFileName          string
	LogArchiveDir        string
	LogArchiveName       string
	LogArchiveMaxAge     string
	LogSizeLimit         uint64

	EndpointAddress string

	RestReadTimeoutSeconds  int
	RestWriteTimeoutSeconds int
}
