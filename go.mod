module github.com/qbradq/annalon

go 1.25.5

require (
	github.com/go-gl/mathgl v1.2.0
	github.com/hajimehoshi/ebiten/v2 v2.9.7
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/qbradq/q2d v0.0.0-20251129230231-f407b354fc1d
	github.com/qbradq/q3d v0.0.0
	golang.org/x/image v0.33.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace (
	github.com/qbradq/q2d => ./lib/q2d
	github.com/qbradq/q3d => ./lib/q3d
)
