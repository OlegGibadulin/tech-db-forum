package mwares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type MiddlewareManager struct{}

func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{}
}

func (m *MiddlewareManager) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logrus.Warn(err)
			}
		}()
		return next(cntx)
	}
}

func (m *MiddlewareManager) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		logrus.Info(cntx.Request().RemoteAddr, " ", cntx.Request().Method, " ", cntx.Request().URL)

		start := time.Now()
		err := next(cntx)
		end := time.Now()

		logrus.Info("Status: ", cntx.Response().Status, " Work time: ", end.Sub(start))
		logrus.Println()
		return err
	}
}
