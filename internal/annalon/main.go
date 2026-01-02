package annalon

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	game := &Game{
		Scale:   2,
		console: ui.NewConsole(),
	}

	ebiten.SetWindowSize(LogicalWidth*game.Scale, LogicalHeight*game.Scale)
	ebiten.SetWindowTitle("Annalon")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
