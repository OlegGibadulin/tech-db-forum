package delivery

import (
	"net/http"

	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ServiceHandler struct {
	serviceUcase service.ServiceUsecase
}

func NewServiceHandler(serviceUcase service.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUcase: serviceUcase,
	}
}

func (sh *ServiceHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/service/clear", sh.ClearServiceHandler())
	e.GET("/api/service/status", sh.GetServiceStatusHandler())
}

func (sh *ServiceHandler) ClearServiceHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		if err := sh.serviceUcase.Clear(); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, "Success")
	}
}

func (sh *ServiceHandler) GetServiceStatusHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		status, err := sh.serviceUcase.GetStatus()
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, status)
	}
}
