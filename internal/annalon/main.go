package annalon

import (
	"io"
	"math"
	"path/filepath"

	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/qbradq/annalon/assets"
	"github.com/qbradq/annalon/internal/command"
	"github.com/qbradq/annalon/internal/ui"
)

const (
	LogicalWidth  = 640
	LogicalHeight = 360
)

type Game struct {
	Scale   int
	console *ui.Console
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.console.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LogicalWidth, LogicalHeight
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

	game := &Game{
		Scale:   2,
		console: console,
	}

	ebiten.SetWindowSize(LogicalWidth*game.Scale, LogicalHeight*game.Scale)
	ebiten.SetWindowTitle("Annalon")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
