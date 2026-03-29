package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	pickTargetService := NewPickTargetService()

	err := wails.Run(&options.App{
		Title:  "ut-translation-setup-v2",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			pickTargetService.startup(ctx)
		},
		Bind: []any{
			pickTargetService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
 