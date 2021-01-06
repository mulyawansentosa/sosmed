package main

import (
	"os"
	"sosmed/src/config"
	"sosmed/src/modules/profile/controller"
	"sosmed/src/modules/profile/repository"
	"sosmed/src/modules/profile/usecase"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	views := iris.HTML("./web/views", ".html").Layout("layout.html").Reload(true)

	app.RegisterView(views)
	app.HandleDir("/public", "./web/public")

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Title", "Judul Website")
		ctx.ViewData("Content", "Halaman Website")
		ctx.View("index.html")
	})

	db, err := config.GetMongoDB()

	if err != nil {
		os.Exit(2)
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "profile")
	profileUsecase := usecase.NewProfileUsecase(profileRepository)

	sessionManager := sessions.New(sessions.Config{
		Cookie:  "cookiename",
		Expires: time.Minute * 10,
	})

	profile := mvc.New(app.Party("/profile"))
	profile.Register(profileUsecase, sessionManager.Start)

	profile.Handle(new(controller.ProfileController))

	app.Listen(":8080")
}
