package main

import (
	"context"
	"fmt"

	"tiger-go/internal/aws"
	"tiger-go/internal/aws/cognito"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	auth *cognito.CognitoAuth
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg := aws.NewConfig()
	a.auth = cognito.NewCognitoAuth(cfg, "us-east-1", cognito.UserPoolID, cognito.ClientID)

	// Initialize the auth context
	runtime.LogInfo(a.ctx, "App started, initializing auth")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SignIn handles the login request from the frontend
func (a *App) SignIn(username string, password string) (*cognito.CognitoAuthResponse, error) {
	return a.auth.SignIn(username, password)
}

// CompleteNewPassword handles the password change request from the frontend
func (a *App) CompleteNewPassword(newPassword string) (*cognito.CognitoAuthResponse, error) {
	return a.auth.CompleteNewPassword(newPassword)
}
