package rater

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type Server struct {
	provider Provider
	echo     *echo.Echo
}

type APIError struct {
	Error string `json:"error"`
}

func NewServer(provider Provider) *Server {
	e := echo.New()
	e.HideBanner = true

	return &Server{
		provider: provider,
		echo:     e,
	}
}

func (s *Server) getAllRates(c echo.Context) error {
	rates, err := s.provider.GetAllRates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIError{"some internal error"})
	}

	if len(rates) == 0 {
		return c.JSON(http.StatusNotFound, APIError{"symbols haven't been found"})
	}

	return c.JSON(http.StatusOK, rates)
}

func (s *Server) getRate(c echo.Context) error {
	list := strings.Split(c.QueryParam("symbols"), ",")
	if len(list) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, APIError{"fill symbols param"})
	}

	rates, err := s.provider.GetRates(list...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIError{"some internal error"})
	}

	if len(rates) == 0 {
		return c.JSON(http.StatusNotFound, APIError{"symbols haven't been found"})
	}

	return c.JSON(http.StatusOK, rates)
}

func (s *Server) Run(port int) error {
	s.echo.GET("/rates", s.getAllRates)
	s.echo.GET("/rate", s.getRate)

	return s.echo.Start(fmt.Sprintf(":%d", port))

}
