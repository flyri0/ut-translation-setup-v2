package main

import (
	"context"
	"math"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	screens, err := wails.ScreenGetAll(ctx)
	if err == nil && len(screens) > 0 {
		activeScreen := screens[0]

		for _, screen := range screens {
			if screen.IsCurrent {
				activeScreen = screen
				break
			}
		}

		// Calculate optimal dimensions (Fraction: 0.5, Aspect Ratio: 4:3)
		width, height := a.resizeWithRatio(activeScreen.Size.Width, activeScreen.Size.Height, 4.0, 3.0, 0.5)

		wails.WindowSetSize(ctx, width, height)
		wails.WindowCenter(ctx)
	}

	wails.WindowShow(ctx)
}

func (a *App) resizeWithRatio(screenWidth, screenHeight int, ratioW, ratioH, fraction float64) (int, int) {
	maxW := float64(screenWidth) * fraction
	maxH := float64(screenHeight) * fraction

	aspect := ratioW / ratioH

	width := maxW
	height := maxW / aspect

	if height > maxH {
		height = maxH
		width = maxH * aspect
	}

	return int(math.Floor(width)), int(math.Floor(height))
}
