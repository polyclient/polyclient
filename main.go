package main

import (
	"embed"
	"log"
	"os"
	"os/exec"

	"github.com/polyclient/polyclient/pkg/app/services"
	"github.com/polyclient/polyclient/pkg/sysinfo"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "polyclient",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(services.NewGreetService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Linux: application.LinuxWindow{
			WebviewGpuPolicy: func() application.WebviewGpuPolicy {
				if sysinfo.HasNvidiaGPU(exec.Command) {
					return application.WebviewGpuPolicyNever // Workaround for https://github.com/wailsapp/wails/issues/2977
				}

				return application.WebviewGpuPolicyAlways
			}(),
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows:          application.WindowsWindow{},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start polyclient: %v", err)
		os.Exit(1)
	}
}
