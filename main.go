package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func GetLogFilePath() string {
	tempDir := os.TempDir()

	return filepath.Join(tempDir, "ut-translation.log")
}

func main() {
	app := NewApp()

	logPath := GetLogFilePath()
	logger := NewFileLogger(logPath)
	defer logger.Close()

	sharedState := NewInstallerState()
	pickTargetService := NewPickTargetService(sharedState)
	pckExplorerService := NewPckExplorerService(sharedState)

	err := wails.Run(&options.App{
		Title:       "Until Then - Instalar Tradução",
		Logger:      logger,
		Width:       1024,
		Height:      768,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
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
