package ui

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mitchellh/go-wordwrap"
	"github.com/qbradq/annalon/internal/command"
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

const (
	ConsoleWidth  = 640
	ConsoleHeight = 360
	MaxHistory    = 1000
	MaxCommands   = 25
)

type Visibility int

const (
	VisibilityHidden Visibility = iota
	VisibilityHalf
	VisibilityFull
)

type Console struct {
	buffer       *q2d.Image
	ebitenBuffer *ebiten.Image // To upload pixels to Ebiten

	logHistory     []string
	commandHistory []string
	historyIndex   int // -1 means new line, 0 is most recent, etc.

	input     string
	cursorPos int

	visibility  Visibility
	interpreter *command.Interpreter
}

func NewConsole(interpreter *command.Interpreter) *Console {
	c := &Console{
		buffer:       q2d.NewImage(ConsoleWidth, ConsoleHeight),
		ebitenBuffer: ebiten.NewImage(ConsoleWidth, ConsoleHeight),
		logHistory:   make([]string, 0, MaxHistory),
		historyIndex: -1,
		visibility:   VisibilityFull,
		interpreter:  interpreter,
	}

	c.Log("Welcome to Annalon Console!")
	return c
}

func (c *Console) Log(msg string) {
	// Word wrap
	wrapped := wordwrap.WrapString(msg, 80) // Approx chars for 640px width with 8px font? 640/8 = 80.
	lines := strings.Split(wrapped, "\n")

	for _, line := range lines {
		if len(c.logHistory) >= MaxHistory {
			// Remove oldest (inefficient for large history but simple for now)
			c.logHistory = c.logHistory[1:]
		}
		c.logHistory = append(c.logHistory, line)
	}
}

func (c *Console) IsVisible() bool {
	return c.visibility != VisibilityHidden
}

func (c *Console) Update() {
	// Cycle visibility
	if inpututil.IsKeyJustPressed(ebiten.KeyBackquote) { // Tilde
		switch c.visibility {
		case VisibilityFull:
			c.visibility = VisibilityHidden
		case VisibilityHidden:
			c.visibility = VisibilityHalf
		case VisibilityHalf:
			c.visibility = VisibilityFull
		}
	}

	if c.visibility == VisibilityHidden {
		return
	}

	// Input Handling
	c.handleInput()
}

func (c *Console) handleInput() {
	// Repeating key logic could be added, but JustPressed is safer for now
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		cmd := strings.TrimSpace(c.input)
		if cmd != "" {
			c.Log("> " + cmd)
			c.executeCommand(cmd)
			c.addToHistory(cmd)
		}
		c.input = ""
		c.cursorPos = 0
		c.historyIndex = -1
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if c.cursorPos > 0 {
			c.input = c.input[:c.cursorPos-1] + c.input[c.cursorPos:]
			c.cursorPos--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		if c.cursorPos < len(c.input) {
			c.input = c.input[:c.cursorPos] + c.input[c.cursorPos+1:]
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if c.cursorPos > 0 {
			c.cursorPos--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if c.cursorPos < len(c.input) {
			c.cursorPos++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if len(c.commandHistory) > 0 {
			if c.historyIndex < len(c.commandHistory)-1 {
				c.historyIndex++
				// Index 0 is most recent (end of slice).
				// Let's reverse semantics: historyIndex 0 = commandHistory[len-1]
				idx := len(c.commandHistory) - 1 - c.historyIndex
				c.input = c.commandHistory[idx]
				c.cursorPos = len(c.input)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if c.historyIndex > -1 {
			c.historyIndex--
			if c.historyIndex == -1 {
				c.input = ""
				c.cursorPos = 0
			} else {
				idx := len(c.commandHistory) - 1 - c.historyIndex
				c.input = c.commandHistory[idx]
				c.cursorPos = len(c.input)
			}
		}
	}

	// Text Input
	runes := ebiten.AppendInputChars(nil)
	if len(runes) > 0 {
		var filtered string
		for _, r := range runes {
			// ASCII 96 is backquote (`), ASCII 126 is tilde (~)
			if r != '`' && r != '~' {
				filtered += string(r)
			}
		}
		if len(filtered) > 0 {
			c.input = c.input[:c.cursorPos] + filtered + c.input[c.cursorPos:]
			c.cursorPos += len(filtered)
		}
	}
}

func (c *Console) addToHistory(cmd string) {
	// Dont add duplicates if same as last? Or just add.
	// Prompt: "Keep max 25 commands in history"
	if len(c.commandHistory) > 0 && c.commandHistory[len(c.commandHistory)-1] == cmd {
		return // Optional: Ignore duplicate sequential commands
	}

	if len(c.commandHistory) >= MaxCommands {
		c.commandHistory = c.commandHistory[1:]
	}
	c.commandHistory = append(c.commandHistory, cmd)
}

func (c *Console) executeCommand(cmd string) {
	if err := c.interpreter.Execute(cmd); err != nil {
		c.Log("Error: " + err.Error())
	}
}

func (c *Console) Draw(screen *ebiten.Image) {
	if c.visibility == VisibilityHidden {
		return
	}

	// Clear background
	// Dark blue-purple: R=20, G=20, B=40, A=255
	bgColor := q2d.Color{20, 20, 40, 255}
	c.buffer.Fill(bgColor)

	// Draw Prompt line
	// Line Height approx 10 (8px font + padding)
	lineHeight := 10
	bottomY := ConsoleHeight - lineHeight

	// Draw Input
	c.buffer.Text(q2d.Point{0, bottomY}, q2d.Color{255, 255, 255, 255}, q2d.FontNormal, false, "> %s", c.input)

	// Draw Cursor
	// Measure width of prompt up to cursor
	cursorX := font.MeasureString(q2d.FontNormal, "> "+c.input[:c.cursorPos]).Ceil()

	// Draw underline cursor
	c.buffer.HLine(bottomY+9, cursorX, cursorX+8, 2, q2d.Color{255, 255, 255, 255}) // Underline cursor

	// Draw History
	drawY := bottomY - lineHeight
	for i := len(c.logHistory) - 1; i >= 0; i-- {
		if drawY < 0 {
			break
		}
		c.buffer.Text(q2d.Point{0, drawY}, q2d.Color{200, 200, 200, 255}, q2d.FontNormal, false, "%s", c.logHistory[i])
		drawY -= lineHeight
	}

	// Upload to Ebiten
	c.ebitenBuffer.WritePixels(c.buffer.Pix)

	// Draw to Screen
	op := &ebiten.DrawImageOptions{}
	if c.visibility == VisibilityHalf {
		// Draw bottom half of console to top half of screen
		// Source: y=180 to 360
		// Dest: y=0
		halfH := ConsoleHeight / 2
		sub := c.ebitenBuffer.SubImage(image.Rect(0, halfH, ConsoleWidth, ConsoleHeight)).(*ebiten.Image)
		op.GeoM.Translate(0, 0) // At top of screen
		screen.DrawImage(sub, op)
	} else {
		// Full
		screen.DrawImage(c.ebitenBuffer, op)
	}
}
