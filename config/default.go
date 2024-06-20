package config

var defaultConfig = Template{
	BaseLoggerDebugLevel:    4,
	LogFileDir:              "",
	LogFileName:             "node.log",
	LogArchiveDir:           "",
	LogArchiveName:          "node.archive.log",
	LogArchiveMaxAge:        "",
	LogSizeLimit:            1073741824,
	NetAddress:              "",
	EndpointAddress:         "127.0.0.1:0",
	RestReadTimeoutSeconds:  15,
	RestWriteTimeoutSeconds: 120,
}
