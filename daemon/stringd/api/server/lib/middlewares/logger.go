package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xboshy/go-template/daemon/stringd/api/server/lib/ctxkeys"
	"github.com/xboshy/go-template/logging"
)

// LoggerMiddleware provides some extra state to the logger middleware
type LoggerMiddleware struct {
	log logging.Logger
}

// MakeLogger initializes the logger middleware function
func MakeLogger(log logging.Logger) fiber.Handler {
	logger := LoggerMiddleware{
		log: log,
	}

	return logger.handler
}

// Logger is an echo middleware to add log to the API
func (logger *LoggerMiddleware) handler(c *fiber.Ctx) error {
	log := logger.log

	// Set elapsed start time
	//start := time.Now()

	// Handle request, store err for logging
	chainErr := c.Next()

	// Manually call error handler
	if chainErr != nil {
		if err := c.App().ErrorHandler(c, chainErr); err != nil {
			_ = c.SendStatus(fiber.StatusInternalServerError) //nolint:errcheck // TODO: Explain why we ignore the error here
		}
	}

	// Set elapsed stop time
	//stop := time.Now()

	requestid := c.Locals(ctxkeys.RequestID)
	if requestid != nil {
		log = log.With(ctxkeys.RequestID, requestid)
	}

	// Get a reference to the request and response objects.
	res := c.Response()
	//req := c.Request()

	// Propagate the error if the next middleware has a problem
	//if err := c.Next(); err != nil {
	//	c.Er
	//	ctx.Error(err)
	//}

	res.Header.ContentLength()
	/*
		log.WithFields(
			logging.Fields{
				ctxkeys.RemoteHost:       req.RemoteAddr,
				ctxkeys.RemoteUser:       "-",
				ctxkeys.RequestTimestamp: start,
				ctxkeys.RequestMethod:    req.Header.Method(),
				ctxkeys.RequestURI:       req.RequestURI,
				ctxkeys.RequestProto:     req.Header.Protocol(),
				ctxkeys.ResponseStatus:   res.StatusCode(),
				ctxkeys.ResponseSize:     strconv.FormatInt(res.Header.ContentLength(), 10),
				ctxkeys.RequestUserAgent: req.Header.UserAgent(),
				ctxkeys.Elapsed:          time.Since(start).Nanoseconds(),
			},
		).Info()*/

	return c.Next()
}
