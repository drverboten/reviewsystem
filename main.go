package main

import (
	"fmt"

	"github.com/drverboten/reviewsystem/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()

	/*Adding two handlers to recover from any http-relative panics
	and log the requests to the terminal*/
	app.Use(recover.New())
	app.Use(logger.New())

	app.RegisterView(iris.HTML("./public", ".html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Post("/send", func(ctx iris.Context) {
		alumno := Alumno{}
		err := ctx.ReadForm(&alumno)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
		alumno.Time = models.GetSubmitTime()
		alumno.Id = alumno.Alumno
		fmt.Println(alumno)

		flag, reason := AddAlumno(alumno)
		if flag {
			ctx.JSON(iris.Map{"message": "Entrega exitosa para el alumno " + alumno.Alumno + " a las: " + alumno.Time + "horas"})
		} else {
			ctx.JSON(iris.Map{"message": reason})
		}
	})

	app.Get("/getall", func(ctx iris.Context) {
		alumnos := GetAll()
		var alum []string
		for _, element := range alumnos {
			alum = append(alum, element.Alumno)
		}
		if len(alumnos) > 0 {
			ctx.JSON(alum)
		} else {
			ctx.JSON(iris.Map{"message": "No hay alumnos registrados"})
		}
	})

	assetHandler := app.StaticHandler("./public", false, false)
	app.SPA(assetHandler)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
