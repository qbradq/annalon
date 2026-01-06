package annalon

import (
	"io"
	"math"
	"path/filepath"

	"os"

	"image"
	"image/color"
	_ "image/png"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/qbradq/annalon/assets"
	"github.com/qbradq/annalon/internal/command"
	"github.com/qbradq/annalon/internal/ui"
	"github.com/qbradq/q3d"
)

const (
	LogicalWidth  = 640
	LogicalHeight = 360
)

type Game struct {
	Scale   int
	console *ui.Console

	scene         *q3d.Scene
	framebuffer   *q3d.FrameBuffer
	renderContext *q3d.RenderContext
	camera        *q3d.Camera
	cube          *q3d.Entity
}

func (g *Game) Update() error {
	// Handle Scale Down (F9)
	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		g.Scale = int(math.Max(1, float64(g.Scale)/2))
		ebiten.SetWindowSize(LogicalWidth*g.Scale, LogicalHeight*g.Scale)
	}

	// Handle Scale Up (F11)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.Scale = int(math.Min(16, float64(g.Scale)*2))
		ebiten.SetWindowSize(LogicalWidth*g.Scale, LogicalHeight*g.Scale)
	}

	// Handle Fullscreen Toggle (F10)
	if inpututil.IsKeyJustPressed(ebiten.KeyF10) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	g.console.Update()

	// Rotate the cube
	g.cube.Rotation = g.cube.Rotation.Mul(mgl32.QuatRotate(0.01, mgl32.Vec3{1, 0, 0})) // X axis
	g.cube.Rotation = g.cube.Rotation.Mul(mgl32.QuatRotate(0.02, mgl32.Vec3{0, 1, 0})) // Y axis
	g.cube.Rotation = g.cube.Rotation.Mul(mgl32.QuatRotate(0.03, mgl32.Vec3{0, 0, 1})) // Z axis
	g.cube.Rotation = g.cube.Rotation.Normalize()

	g.scene.UpdateEntity(g.cube)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw 3D Viewport
	g.framebuffer.Clear(color.RGBA{0, 0, 0, 255})
	g.camera.RenderScene(g.renderContext, g.scene, nil)

	// Blit framebuffer to screen
	screen.WritePixels(g.framebuffer.Pix)

	// Draw Console
	g.console.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LogicalWidth, LogicalHeight
}

func createCubeMesh(size float32, texture *q3d.Texture) *q3d.Mesh {
	mesh := &q3d.Mesh{}
	s := size / 2

	// Helper to create a vertex
	v := func(x, y, z, u, v float32) q3d.Vertex {
		return q3d.Vertex{
			Position: mgl32.Vec3{x, y, z},
			TexCoord: mgl32.Vec2{u, v},
			Color:    color.RGBA{255, 255, 255, 255},
		}
	}

	// Front face (Z+)
	mesh.AddConvexPolygon(texture,
		v(-s, -s, s, 0, 0),
		v(s, -s, s, 1, 0),
		v(s, s, s, 1, 1),
		v(-s, s, s, 0, 1),
	)

	// Back face (Z-)
	mesh.AddConvexPolygon(texture,
		v(s, -s, -s, 0, 0),
		v(-s, -s, -s, 1, 0),
		v(-s, s, -s, 1, 1),
		v(s, s, -s, 0, 1),
	)

	// Left face (X-)
	mesh.AddConvexPolygon(texture,
		v(-s, -s, -s, 0, 0),
		v(-s, -s, s, 1, 0),
		v(-s, s, s, 1, 1),
		v(-s, s, -s, 0, 1),
	)

	// Right face (X+)
	mesh.AddConvexPolygon(texture,
		v(s, -s, s, 0, 0),
		v(s, -s, -s, 1, 0),
		v(s, s, -s, 1, 1),
		v(s, s, s, 0, 1),
	)

	// Top face (Y+)
	mesh.AddConvexPolygon(texture,
		v(-s, s, s, 0, 0),
		v(s, s, s, 1, 0),
		v(s, s, -s, 1, 1),
		v(-s, s, -s, 0, 1),
	)

	// Bottom face (Y-)
	mesh.AddConvexPolygon(texture,
		v(-s, -s, -s, 0, 0),
		v(s, -s, -s, 1, 0),
		v(s, -s, s, 1, 1),
		v(-s, -s, s, 0, 1),
	)

	return mesh
}

func Main() {
	interpreter := command.NewInterpreter()
	interpreter.RegisterBuiltins()
	interpreter.SetQuitHandler(func() {
		os.Exit(0)
	})

	console := ui.NewConsole(interpreter)
	command.SetPrinter(console.Log)

	// Startup Logic
	dataDir := "data"
	configFile := filepath.Join(dataDir, "config.con")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Copy from assets
		src, err := assets.FS.Open("data/config.default.con")
		if err == nil {
			// We can't defer in a long running loop or function easily without closure,
			// but Main is one-shot.
			// Ideally check errors.
			os.MkdirAll(dataDir, 0755)
			dst, err := os.Create(configFile)
			if err == nil {
				io.Copy(dst, src)
				dst.Close()
			}
			src.Close()
		}
	}

	// Execute init script
	if err := interpreter.ExecFile("data/init.con"); err != nil {
		command.Printf("Failed to execute init.con: %v", err)
	}

	// 3D Initialization
	width, height := LogicalWidth, LogicalHeight
	fb := q3d.NewFrameBuffer(width, height)
	scene := q3d.NewScene()
	// Set ambient light to full as per requirements
	scene.AmbientLight = color.RGBA{255, 255, 255, 255}

	rc := q3d.NewRenderContext()

	// Load Texture
	texFile, err := assets.FS.Open("textures/grid_blue.png")
	if err != nil {
		panic(err)
	}
	defer texFile.Close()
	texImg, _, err := image.Decode(texFile)
	if err != nil {
		panic(err)
	}
	rgba, ok := texImg.(*image.RGBA)
	if !ok {
		// Convert to RGBA
		b := texImg.Bounds()
		rgba = image.NewRGBA(b)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				rgba.Set(x, y, texImg.At(x, y))
			}
		}
	}
	texture := &q3d.Texture{RGBA: rgba}

	// Create Cube
	cubeMesh := createCubeMesh(64, texture)
	cube := &q3d.Entity{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(),
		Mesh:     cubeMesh,
	}
	scene.AddEntity(cube)

	// Setup Camera
	camera := &q3d.Camera{
		Entity: q3d.Entity{
			Position: mgl32.Vec3{0, 32, 512},
			Rotation: mgl32.QuatIdent(),
		},
		FOV:          60,
		Near:         0.1,
		Far:          1000,
		RenderTarget: fb,
	}
	// Look at origin
	// Camera LookAt is implicit in UpdateViewMatrix if we set rotation correctly?
	// Wait, q3d.Camera.UpdateViewMatrix uses Rotation to determine forward vector.
	// So we need to set the camera's Rotation to look at the origin.
	// q3d.Camera.UpdateViewMatrix implementation:
	// forward := c.Rotation.Rotate(mgl32.Vec3{0, 0, -1})
	// target := c.Position.Add(forward)
	// c.viewMatrix = mgl32.LookAtV(c.Position, target, up)
	// Actually LookAtV computes the view matrix based on Pos, Target, Up.
	// But UpdateViewMatrix GENERATES the target from the Rotation.
	// So we must set the Rotation such that (0,0,-1) rotated points to the target.
	// Target is Origin (0,0,0). Camera is at (0, 32, 32).
	// Direction vector = Origin - Camera = (0, -32, -32). Normalized: (0, -0.707, -0.707).
	// We need a quaternion that rotates (0,0,-1) to (0, -0.707, -0.707).
	// Use mgl32.QuatBetweenVectors.
	direction := mgl32.Vec3{0, 0, 0}.Sub(camera.Position).Normalize()
	camera.Rotation = mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, -1}, direction)

	camera.UpdateMatrices() // Initial update

	game := &Game{
		Scale:         2,
		console:       console,
		scene:         scene,
		framebuffer:   fb,
		renderContext: rc,
		camera:        camera,
		cube:          cube,
	}

	ebiten.SetWindowSize(LogicalWidth*game.Scale, LogicalHeight*game.Scale)
	ebiten.SetWindowTitle("Annalon")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
