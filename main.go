package main

import (
	"embed"
	"log"

	"arkinightinfo/data"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	opService := &OperatorService{
		dataRoot: data.ResolveDataRoot(),
	}

	app := application.New(application.Options{
		Name:        "明日方舟干员查询",
		Description: "Arknights 干员数据查询工具",
		Services: []application.Service{
			application.NewService(opService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "明日方舟干员查询",
		Width:  1200,
		Height: 780,
		MinWidth:  960,
		MinHeight: 600,
		BackgroundColour: application.NewRGB(18, 19, 28),
		URL: "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
