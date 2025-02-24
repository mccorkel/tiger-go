package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	cognito "tiger-go/internal/tiger-aws/tiger-cognito"
	tigerconfig "tiger-go/internal/tiger-aws/tiger-config"
	tigerwhip "tiger-go/internal/tiger-whip"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"
	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime/types"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/codec/opus"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	_ "github.com/pion/mediadevices/pkg/driver/microphone"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	auth        *cognito.CognitoAuth
	ivsClient   *ivsrealtime.Client
	stageArn    string
	whipClient  *tigerwhip.WHIPClient
	mediaStream mediadevices.MediaStream
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg, err := tigerconfig.LoadDefaultConfig()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to load AWS config: %v", err)
		return
	}
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

// GetStageInfo retrieves or creates an IVS stage
func (a *App) GetStageInfo() (map[string]interface{}, error) {
	if a.ivsClient == nil {
		cfg, err := config.LoadDefaultConfig(a.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}
		a.ivsClient = ivsrealtime.NewFromConfig(cfg)
	}

	// If we already have a stage ARN, return it
	if a.stageArn != "" {
		return map[string]interface{}{
			"StageArn": a.stageArn,
		}, nil
	}

	// Create a new stage
	input := &ivsrealtime.CreateStageInput{
		Name: aws.String("GoWailsStage-" + time.Now().Format("20060102150405")),
	}

	result, err := a.ivsClient.CreateStage(a.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create stage: %w", err)
	}

	a.stageArn = *result.Stage.Arn

	return map[string]interface{}{
		"StageArn": a.stageArn,
		"Name":     *result.Stage.Name,
	}, nil
}

// CreateParticipantToken generates a token for joining the stage
func (a *App) CreateParticipantToken(userId string) (map[string]interface{}, error) {
	if a.ivsClient == nil || a.stageArn == "" {
		return nil, fmt.Errorf("stage not initialized, call GetStageInfo first")
	}

	input := &ivsrealtime.CreateParticipantTokenInput{
		StageArn: &a.stageArn,
		UserId:   &userId,
		Capabilities: []types.ParticipantTokenCapability{
			"PUBLISH",
			"SUBSCRIBE",
		},
		Duration: aws.Int32(60 * 60), // 1 hour
	}

	result, err := a.ivsClient.CreateParticipantToken(a.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create participant token: %w", err)
	}

	return map[string]interface{}{
		"Token":          *result.ParticipantToken,
		"ParticipantId":  userId,
		"ExpirationTime": time.Now().Add(time.Hour).Unix(),
	}, nil
}

// StartWebcamCapture starts capturing from the webcam
func (a *App) StartWebcamCapture() (map[string]interface{}, error) {
	// Get available devices
	devices := mediadevices.EnumerateDevices()

	// Filter for cameras
	var cameras []mediadevices.MediaDeviceInfo
	for _, device := range devices {
		if device.Kind == mediadevices.VideoInput {
			cameras = append(cameras, device)
		}
	}

	if len(cameras) == 0 {
		return nil, fmt.Errorf("no camera found")
	}

	// Configure media constraints
	constraints := mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.DeviceID = prop.String(cameras[0].DeviceID)
			c.Width = prop.Int(1280)
			c.Height = prop.Int(720)
			c.FrameRate = prop.Float(30)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {
			// Configure audio constraints
		},
	}

	// Use H.264 codec as required by IVS
	h264Params, err := openh264.NewParams()
	if err != nil {
		return nil, fmt.Errorf("failed to create H.264 codec: %w", err)
	}
	h264Params.BitRate = 2_500_000
	h264Params.KeyFrameInterval = 60

	// Create Opus codec
	opusParams, err := opus.NewParams()
	if err != nil {
		return nil, fmt.Errorf("failed to create Opus codec: %w", err)
	}

	// Get media stream
	mediaStream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: constraints.Video,
		Audio: constraints.Audio,
		Codec: mediadevices.NewCodecSelector(
			mediadevices.WithVideoEncoders(&h264Params),
			mediadevices.WithAudioEncoders(&opusParams),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user media: %w", err)
	}

	a.mediaStream = mediaStream

	return map[string]interface{}{
		"Success": true,
		"Message": "Webcam capture started",
	}, nil
}

// StartScreenCapture starts capturing the screen
func (a *App) StartScreenCapture() (map[string]interface{}, error) {
	// Get available devices
	devices := mediadevices.EnumerateDevices()

	// Filter for screens
	var screens []mediadevices.MediaDeviceInfo
	for _, device := range devices {
		if device.Kind == mediadevices.VideoInput &&
			strings.Contains(strings.ToLower(device.Label), "screen") {
			screens = append(screens, device)
		}
	}

	if len(screens) == 0 {
		return nil, fmt.Errorf("no screen found")
	}

	// Stop existing media stream if any
	if a.mediaStream != nil {
		for _, track := range a.mediaStream.GetTracks() {
			track.Close()
		}
		a.mediaStream = nil
	}

	// Configure media constraints
	constraints := mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.DeviceID = prop.String(screens[0].DeviceID)
			// Limit to 720p as per AWS docs
			c.Width = prop.Int(1280)
			c.Height = prop.Int(720)
			c.FrameRate = prop.Float(30)
		},
	}

	// Use H.264 codec as required by IVS
	h264Params, err := openh264.NewParams()
	if err != nil {
		return nil, fmt.Errorf("failed to create H.264 codec: %w", err)
	}
	h264Params.BitRate = 2_500_000
	h264Params.KeyFrameInterval = 60

	// Create Opus codec
	opusParams, err := opus.NewParams()
	if err != nil {
		return nil, fmt.Errorf("failed to create Opus codec: %w", err)
	}

	// Get media stream
	mediaStream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: constraints.Video,
		// No audio for screen capture
		Codec: mediadevices.NewCodecSelector(
			mediadevices.WithVideoEncoders(&h264Params),
			mediadevices.WithAudioEncoders(&opusParams),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get display media: %w", err)
	}

	a.mediaStream = mediaStream

	return map[string]interface{}{
		"Success": true,
		"Message": "Screen capture started",
	}, nil
}

// StartStreaming starts streaming to IVS using WHIP
func (a *App) StartStreaming(token string) (map[string]interface{}, error) {
	if a.mediaStream == nil {
		return nil, fmt.Errorf("no media stream available, start capture first")
	}

	// Create WHIP client
	whipEndpoint := fmt.Sprintf("https://global.whip.live-video.net/v1/whip/%s", a.stageArn)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	whipClient, err := tigerwhip.NewWHIPClient(whipEndpoint, token, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create WHIP client: %w", err)
	}

	// Add tracks from media stream
	for _, track := range a.mediaStream.GetTracks() {
		_, err := whipClient.AddTrack(track)
		if err != nil {
			return nil, fmt.Errorf("failed to add track: %w", err)
		}

		// Log track info
		runtime.LogInfof(a.ctx, "Added track: %s", track.ID())
	}

	// Connect to WHIP endpoint
	if err := whipClient.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to WHIP endpoint: %w", err)
	}

	a.whipClient = whipClient

	return map[string]interface{}{
		"Success": true,
		"Message": "Streaming started",
	}, nil
}

// StopStreaming stops the streaming
func (a *App) StopStreaming() (map[string]interface{}, error) {
	if a.whipClient != nil {
		if err := a.whipClient.Close(); err != nil {
			return nil, fmt.Errorf("failed to close WHIP client: %w", err)
		}
		a.whipClient = nil
	}

	if a.mediaStream != nil {
		for _, track := range a.mediaStream.GetTracks() {
			track.Close()
		}
		a.mediaStream = nil
	}

	return map[string]interface{}{
		"Success": true,
		"Message": "Streaming stopped",
	}, nil
}

// GetPreviewFrame captures a single frame from the current media stream
func (a *App) GetPreviewFrame() (map[string]interface{}, error) {
	if a.mediaStream == nil {
		return nil, fmt.Errorf("no media stream available")
	}

	// This is a placeholder - actual implementation would depend on how
	// you want to handle the preview (e.g., save to file, encode to base64, etc.)
	return map[string]interface{}{
		"Success": true,
		"Message": "Preview functionality not yet implemented",
	}, nil
}
