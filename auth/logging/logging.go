package logging

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/vespaiach/auth/service"
	"github.com/vespaiach/auth/store"
)

type loggingService struct {
	logger log.Logger
	service.Service
}

// NewLogging returns a new instance of a logging Service.
func NewLogging(logger log.Logger, s service.Service) service.Service {
	return &loggingService{logger, s}
}

func (s *loggingService) VerifyLogin(username string, password string) (user *store.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "VerifyLogin",
			"username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.VerifyLogin(username, password)
}
