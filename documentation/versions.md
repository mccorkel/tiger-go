# Project Versions

## Core Dependencies

### Go
- Version: 1.23.4
- Architecture: darwin/arm64

### Wails
- Version: v2.10.1

### Svelte
- Version: 3.49.0

## Frontend Dependencies Snapshot
### Package.json
~~~ json
{
  "name": "frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "check": "svelte-check --tsconfig ./tsconfig.json"
  },
  "devDependencies": {
    "@sveltejs/vite-plugin-svelte": "^1.0.1",
    "@tsconfig/svelte": "^5.0.0",
    "svelte": "^3.49.0",
    "svelte-check": "^2.8.0",
    "svelte-preprocess": "^4.10.7",
    "tslib": "^2.4.0",
    "typescript": "^4.6.4",
    "vite": "^3.0.7"
  }
}
~~~

## Backend Dependencies Snapshot
### go.mod
~~~ go
module tiger-go

go 1.23

require github.com/wailsapp/wails/v2 v2.10.1

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/labstack/echo/v4 v4.13.3 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.1 // indirect
	github.com/leaanthony/gosod v1.0.4 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/leaanthony/u v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/tkrajina/go-reflector v0.5.8 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/go-webview2 v1.0.19 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

## AWS SDK Dependencies
- aws-sdk-go-v2: v1.25.2
- aws-sdk-go-v2/config: v1.27.4

## Notes
- Project is running on MacOS (Apple Silicon)
- Using pre-release version of Go 1.23.4
- Using development version of Wails v2.10.1
- Frontend built with Svelte 3.49.0
- Updated @tsconfig/svelte from ^3.0.0 to ^5.0.0 to fix TypeScript configuration warnings

Last updated: March 2024 