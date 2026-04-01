package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger, err := NewFileLogger("translation_installer")
	if err != nil {
		log.Fatal("Failed to initialize logging system:", err)
	}

	sharedState := NewInstallerState()
	pickTargetService := NewPickTargetService(sharedState)
	pckExplorerService := NewPckExplorerService(sharedState)

	err = wails.Run(&options.App{
		Title:  "ut-translation-setup-v2",
		Logger: logger,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			pickTargetService.startup(ctx)
			pckExplorerService.startup(ctx)
		},
		Bind: []any{
			pickTargetService,
			pckExplorerService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
