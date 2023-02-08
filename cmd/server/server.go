package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var HTTPServer *echo.Echo

func GetClientIP(c echo.Context) error {
	logrus.Infof("test")
	c.JSON(http.StatusOK, c.RealIP())
	return nil
}

func setRoute(e *echo.Echo) {
	ddns := e.Group("")
	ddns.GET("/client_ip", GetClientIP)
}

func main() {
	server := echo.New()
	setRoute(server)

	httpServerAddress := fmt.Sprintf("%s:%s", "0.0.0.0", "9999")
	err := server.Start(httpServerAddress)
	if err != nil {
		logrus.Errorf("start http server in child process[pid:%d] with err: %+v\n", os.Getpid(), err)
		os.Exit(2)
	}
}
