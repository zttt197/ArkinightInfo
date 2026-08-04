package main

import (
	"embed"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"arkinightinfo/data"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dataRoot := data.ResolveDataRoot()
	opService := &OperatorService{dataRoot: dataRoot}

	// Combine the embedded frontend assets with the avatar file server.
	fsHandler := application.AssetFileServerFS(assets)
	avatarDir := filepath.Join(dataRoot, "avatars")
	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/avatars/") {
			name := strings.TrimPrefix(r.URL.Path, "/avatars/")
			http.ServeFile(w, r, filepath.Join(avatarDir, name))
			return
		}
		fsHandler.ServeHTTP(w, r)
	})

	// Register download progress event so the frontend can listen.
	application.RegisterEvent[string]("dl-progress")

	app := application.New(application.Options{
		Name:        "明日方舟干员查询",
		Description: "Arknights 干员数据查询工具",
		Services: []application.Service{
			application.NewService(opService),
		},
		Assets: application.AssetOptions{
			Handler: combinedHandler,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Wire the service progress callback to emit Wails events.
	opService.SetProgressCallback(func(msg string) {
		app.Event.Emit("dl-progress", msg)
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
