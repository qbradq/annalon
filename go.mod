module github.com/qbradq/annalon

go 1.25.5

require (
	github.com/go-gl/mathgl v1.2.0 // indirect
	github.com/qbradq/q2d v0.0.0-20251129230231-f407b354fc1d // indirect
	github.com/qbradq/q3d v0.0.0-20251208110450-913bd2d9a85b // indirect
	golang.org/x/image v0.33.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace (
	github.com/qbradq/q2d v0.0.0-20251129230231-f407b354fc1d => lib/q2d v0.0.0-20251129230231-f407b354fc1d
	github.com/qbradq/q3d v0.0.0-20251208110450-913bd2d9a85b => lib/q3d v0.0.0-20251208110450-913bd2d9a85b
)
