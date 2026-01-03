# Annalon
Annalon is an experimental 3D RPG.

## License
Annalon is licensed under the GPL-3.0 License or later.

# Build and Run

## Platform Support

Annalon targets the following platforms:
- Desktop Browser
- Windows
- Linux
- macOS

## Dependencies

Annalon is built using Go 1.25+ and uses the following libraries:

- https://github.com/hajimehoshi/ebiten/v2 (platform abstraction)
- https://github.com/qbradq/q3d (3D rendering)
- https://github.com/qbradq/q2d (2D rendering)
- https://github.com/go-gl/mathgl/mgl32 (3D math)
- https://github.com/mitchellh/go-wordwrap (text wrapping)

## Building

Annalon can be built for the native platform using the following commands:

```bash
git clone https://github.com/qbradq/annalon.git
cd annalon
go build cmd/annalon
```
