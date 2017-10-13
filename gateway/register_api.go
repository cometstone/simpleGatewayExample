package gateway

import "github.com/labstack/echo"

func initUpdateApi() {

	go func() {
		e := echo.New()
		e.POST("/api/update",apiUpdate)
		e.Logger.Fatal(e.Start(":"+Conf.Api.ApiUpdatePort))
	}()

	go registerApiUpdate()
}

func apiUpdate(c echo.Context) error {
	return nil
}

func registerApiUpdate() {

}