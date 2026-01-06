package annalon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/qbradq/annalon/internal/command"
)

// InputMap holds the mapping from Ebiten input constants to string names.
var InputMap = map[string]interface{}{}
var keyToName = map[ebiten.Key]string{}
var nameToKey = map[string]ebiten.Key{}
var mouseToName = map[ebiten.MouseButton]string{}
var nameToMouse = map[string]ebiten.MouseButton{}

// InitInputMappings initializes the input mappings.
func InitInputMappings() {
	// Function to register a key
	regKey := func(key ebiten.Key, name string) {
		keyToName[key] = name
		nameToKey[name] = key
	}

	regMouse := func(btn ebiten.MouseButton, name string) {
		mouseToName[btn] = name
		nameToMouse[name] = btn
	}

	// Letters
	regKey(ebiten.KeyA, "a")
	regKey(ebiten.KeyB, "b")
	regKey(ebiten.KeyC, "c")
	regKey(ebiten.KeyD, "d")
	regKey(ebiten.KeyE, "e")
	regKey(ebiten.KeyF, "f")
	regKey(ebiten.KeyG, "g")
	regKey(ebiten.KeyH, "h")
	regKey(ebiten.KeyI, "i")
	regKey(ebiten.KeyJ, "j")
	regKey(ebiten.KeyK, "k")
	regKey(ebiten.KeyL, "l")
	regKey(ebiten.KeyM, "m")
	regKey(ebiten.KeyN, "n")
	regKey(ebiten.KeyO, "o")
	regKey(ebiten.KeyP, "p")
	regKey(ebiten.KeyQ, "q")
	regKey(ebiten.KeyR, "r")
	regKey(ebiten.KeyS, "s")
	regKey(ebiten.KeyT, "t")
	regKey(ebiten.KeyU, "u")
	regKey(ebiten.KeyV, "v")
	regKey(ebiten.KeyW, "w")
	regKey(ebiten.KeyX, "x")
	regKey(ebiten.KeyY, "y")
	regKey(ebiten.KeyZ, "z")

	// Numbers
	regKey(ebiten.Key0, "0")
	regKey(ebiten.Key1, "1")
	regKey(ebiten.Key2, "2")
	regKey(ebiten.Key3, "3")
	regKey(ebiten.Key4, "4")
	regKey(ebiten.Key5, "5")
	regKey(ebiten.Key6, "6")
	regKey(ebiten.Key7, "7")
	regKey(ebiten.Key8, "8")
	regKey(ebiten.Key9, "9")

	// Function Keys
	regKey(ebiten.KeyF1, "f1")
	regKey(ebiten.KeyF2, "f2")
	regKey(ebiten.KeyF3, "f3")
	regKey(ebiten.KeyF4, "f4")
	regKey(ebiten.KeyF5, "f5")
	regKey(ebiten.KeyF6, "f6")
	regKey(ebiten.KeyF7, "f7")
	regKey(ebiten.KeyF8, "f8")
	regKey(ebiten.KeyF9, "f9")
	regKey(ebiten.KeyF10, "f10")
	regKey(ebiten.KeyF11, "f11")
	regKey(ebiten.KeyF12, "f12")

	// Special Keys
	regKey(ebiten.KeySpace, "space")
	regKey(ebiten.KeyEnter, "enter")
	regKey(ebiten.KeyEscape, "escape")
	regKey(ebiten.KeyBackspace, "backspace")
	regKey(ebiten.KeyTab, "tab")
	regKey(ebiten.KeyShift, "shift")
	regKey(ebiten.KeyControlLeft, "ctrl")
	regKey(ebiten.KeyControlRight, "ctrl")
	regKey(ebiten.KeyAlt, "alt")

	// Arrows
	regKey(ebiten.KeyUp, "up")
	regKey(ebiten.KeyDown, "down")
	regKey(ebiten.KeyLeft, "left")
	regKey(ebiten.KeyRight, "right")

	// Mouse
	regMouse(ebiten.MouseButtonLeft, "mouse1")
	regMouse(ebiten.MouseButtonRight, "mouse2")
	regMouse(ebiten.MouseButtonMiddle, "mouse3")
}

// activeBindings tracks which bindings are currently "held down" so we can invert them on release.
// We map the key name to the command line that was executed.
var activeBindings = make(map[string]string)

// ProcessInput handles input events and executes bound commands.
func ProcessInput(interpreter *command.Interpreter) {
	bindings := interpreter.GetBindings()

	// Check Keys
	for key, name := range keyToName {
		if inpututil.IsKeyJustPressed(key) {
			if cmd, ok := bindings[name]; ok {
				if err := interpreter.Execute(cmd); err != nil {
					command.Printf("Input Error: %v", err)
				}
				// Store the original command to calculate inverse later?
				// Or does the inverse depend on the command string?
				// "When the input event associated with this bind is released, execute the inverse of the recorded boolean set commands."
				// So we should store the command string.
				activeBindings[name] = cmd
			}
		} else if inpututil.IsKeyJustReleased(key) {
			if cmd, ok := activeBindings[name]; ok {
				inverse := interpreter.GetInverseBooleanCommands(cmd)
				if inverse != "" {
					if err := interpreter.Execute(inverse); err != nil {
						command.Printf("Input Release Error: %v", err)
					}
				}
				delete(activeBindings, name)
			}
		}
	}

	// Check Mouse
	for btn, name := range mouseToName {
		if inpututil.IsMouseButtonJustPressed(btn) {
			if cmd, ok := bindings[name]; ok {
				if err := interpreter.Execute(cmd); err != nil {
					command.Printf("Input Error: %v", err)
				}
				activeBindings[name] = cmd
			}
		} else if inpututil.IsMouseButtonJustReleased(btn) {
			if cmd, ok := activeBindings[name]; ok {
				inverse := interpreter.GetInverseBooleanCommands(cmd)
				if inverse != "" {
					if err := interpreter.Execute(inverse); err != nil {
						command.Printf("Input Release Error: %v", err)
					}
				}
				delete(activeBindings, name)
			}
		}
	}
}
