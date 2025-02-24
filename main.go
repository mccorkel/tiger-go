package main

import (
	"context"
	"embed"
	"tiger-go/internal/aws/cognito"
	"tiger-go/internal/aws/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		println("Error loading AWS config:", err.Error())
		return
	}

	cognitoAuth := cognito.NewCognitoAuth(
		cfg,
		"us-east-1",
		"us-east-1_qoSTRtzGT",
		"1df6gkm0eaiq4a33i6mu5gp1ji",
	)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "tiger-go",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			runtime.LogInfo(ctx, "Application started")
		},
		Bind: []interface{}{
			app,
			cognitoAuth,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
