package qmap

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

// Parse reads a Quake map from the reader.
func Parse(r io.Reader) (*Map, error) {
	s := new(scanner.Scanner)
	s.Init(r)
	// Disable ScanInts/ScanFloats so numbers are scanned as identifiers (strings).
	// This allows us to handle negative numbers like "-256" as a single token if we include '-' in IsIdentRune,
	// and textures like "materials/carpet" as a single token if we include '/'.
	// Disable builtin skipping to handle texture paths with special chars and correct comment skipping
	s.Mode = scanner.ScanIdents | scanner.ScanStrings
	s.Whitespace = 0 // We handle whitespace manually
	s.IsIdentRune = func(ch rune, i int) bool {
		// Standard idents + special chars often found in texture names or values (except / which is comment start)
		return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.' || ch == '*' || ch == '-' || ch == '+'
	}

	p := &parser{s: s}
	return p.parseMap()
}

type parser struct {
	s *scanner.Scanner
}

// scan returns the next non-whitespace, non-comment token.
func (p *parser) scan() rune {
	for {
		tok := p.s.Scan()
		if tok == scanner.EOF {
			return tok
		}

		// Handle whitespace
		if tok > 0 && unicode.IsSpace(tok) {
			continue
		}

		// Handle comments
		if tok == '/' {
			next := p.s.Peek()
			if next == '/' {
				// Line comment
				for {
					ch := p.s.Next()
					if ch == '\n' || ch == '\r' || ch == scanner.EOF {
						break
					}
				}
				continue
			}
			// Block comments (if any? Quake maps usually //)
		}

		return tok
	}
}

// scanRaw returns next token, skipping comments but returning whitespace.
func (p *parser) scanRaw() rune {
	for {
		tok := p.s.Scan()
		if tok == scanner.EOF {
			return tok
		}

		// Handle comments
		if tok == '/' {
			next := p.s.Peek()
			if next == '/' {
				// Line comment
				for {
					ch := p.s.Next()
					if ch == '\n' || ch == '\r' || ch == scanner.EOF {
						break
					}
				}
				continue
			}
		}
		return tok
	}
}

func (p *parser) nextToken() string {
	tok := p.s.TokenText()
	// Handle quoting if Scanner didn't automatically strip it for raw strings?
	// text/scanner strips quotes for ScanStrings mode.
	if strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\"") {
		return strings.Trim(tok, "\"")
	}
	return tok
}

func (p *parser) parseMap() (*Map, error) {
	m := &Map{}
	for {
		tok := p.scan()
		if tok == scanner.EOF {
			break
		}
		if tok != '{' {
			// Skip garbage or error?
			// Some maps have comments at top level, Scanner skips them.
			// Try to find start of entity.
			// But scan() returns start char.
			// Valid map starts with entity `{`.
			// If we hit something else that isn't EOF, it's weird.
			// But the scanner returns the char itself for symbols.
			continue
		}

		entity, err := p.parseEntity()
		if err != nil {
			return nil, err
		}
		m.Entities = append(m.Entities, *entity)
	}
	return m, nil
}

func (p *parser) parseEntity() (*Entity, error) {
	e := &Entity{
		Properties: make(map[string]string),
	}

	for {
		tok := p.scan()
		if tok == scanner.EOF {
			return nil, fmt.Errorf("unexpected EOF inside entity")
		}

		text := p.nextToken()

		if text == "}" {
			return e, nil
		}

		if text == "{" {
			// Brush
			brush, err := p.parseBrush()
			if err != nil {
				return nil, err
			}
			e.Brushes = append(e.Brushes, *brush)
			continue
		}

		// Property: "key" "value"
		// text is key
		key := text

		// scanner returns next token
		tok2 := p.scan()
		if tok2 == scanner.EOF {
			return nil, fmt.Errorf("unexpected EOF looking for property value of %s", key)
		}
		val := p.nextToken()
		e.Properties[key] = val
	}
}

func (p *parser) parseBrush() (*Brush, error) {
	b := &Brush{}
	for {
		// Expect `(` or `}`
		tok := p.scan()
		if tok == scanner.EOF {
			return nil, fmt.Errorf("unexpected EOF inside brush")
		}
		text := p.nextToken()

		if text == "}" {
			return b, nil
		}

		if text == "(" {
			// Start of plane
			// We already consumed `(`
			plane, err := p.parsePlane()
			if err != nil {
				return nil, err
			}
			b.Planes = append(b.Planes, *plane)
		} else {
			// Maybe garbage or parsing error
			// Ignore or unexpected?
			// Brushes only contain planes.
		}
	}
}

func (p *parser) parsePlane() (*Plane, error) {
	// Format: ( x1 y1 z1 ) ( x2 y2 z2 ) ( x3 y3 z3 ) texture xA xB xC xD yA yB yC yD (Valve 220)
	// OR: ( x1 y1 z1 ) ( x2 y2 z2 ) ( x3 y3 z3 ) texture offX offY rot scaleX scaleY (Standard)

	// We already consumed the first `(`

	pl := &Plane{}

	// Point A
	if err := p.parsePoint(&pl.A); err != nil {
		return nil, fmt.Errorf("failed to parse point A: %w", err)
	}
	// Expect `)`
	p.scan()
	if p.nextToken() != ")" {
		return nil, fmt.Errorf("expected ')' after point A, got %s", p.nextToken())
	}

	// Expect `(`
	p.scan()
	if p.nextToken() != "(" {
		return nil, fmt.Errorf("expected '(' before point B, got %s", p.nextToken())
	}

	// Point B
	if err := p.parsePoint(&pl.B); err != nil {
		return nil, fmt.Errorf("failed to parse point B: %w", err)
	}
	// Expect `)`
	p.scan()
	if p.nextToken() != ")" {
		return nil, fmt.Errorf("expected ')' after point B, got %s", p.nextToken())
	}

	// Expect `(`
	p.scan()
	if p.nextToken() != "(" {
		return nil, fmt.Errorf("expected '(' before point C, got %s", p.nextToken())
	}

	// Point C
	if err := p.parsePoint(&pl.C); err != nil {
		return nil, fmt.Errorf("failed to parse point C: %w", err)
	}
	// Expect `)`
	p.scan()
	if p.nextToken() != ")" {
		return nil, fmt.Errorf("expected ')' after point C, got %s", p.nextToken())
	}

	// Texture
	// Texture name can contain special chars and is terminated by whitespace.
	// We scanraw to detect space.
	var textureParts []string

	// First part (already skipped preceeding whitespace if we called scan() before?
	// The caller called scan() implies we are at ')' previously.
	// We need to advance to start of texture.

	// Scan until we find start of texture (skipping space)
	p.scan()
	// scan() sets up current token.
	if p.s.TokenText() == "" {
		// EOF or error
	}
	textureParts = append(textureParts, p.nextToken())

	// Now continue scanning RAW to check for space separation
	for {
		// Peek or ScanNext?
		// check if next char is space using Peek is hard because comments?
		// Use scanRaw.
		tok := p.scanRaw()
		if tok == scanner.EOF {
			break
		}
		if tok > 0 && unicode.IsSpace(tok) {
			// Space delimiter found. Texture parsing done.
			// But wait, the space is consumed. The next field starts after space.
			// But parsePlane expects 'text' variable to be set for the next logic if we are not careful?
			// parsePlane logic below calls p.scan().
			// If we consumed space, the NEXT token is the start of next field.
			// BUT the code below expects `p.scan()` to be called to get `text`.
			// So we should put back the token?
			// We can't put back.

			// Actually, look at logic below:
			// p.scan()
			// text := p.nextToken()

			// So we just need to STOP here. The next p.scan() call (below) will skip this space (if it's scan())
			// OR if we consumed space, we need to ensure next p.scan() gets the next non-space token.
			// Yes, p.scan() skips whitespace.
			break
		}

		// Append to texture
		textureParts = append(textureParts, p.nextToken())
	}
	pl.Texture = strings.Join(textureParts, "")

	// Parameters
	// Standard: [offX] [offY] [rot] [scaleX] [scaleY]
	// Valve220: [ [ xA xB xC xD ] [ yA yB yC yD ] rot scaleX scaleY ]

	// Check for Valve 220 `[`
	// Peek isn't easy with text/scanner unless we used Peek() before Scan() but we already Scanned.
	// So scan next token.
	p.scan()
	text := p.nextToken()

	if text == "[" {
		// Valve 220 format not fully supported yet based on simple struct
		// But let's consume it to avoid breaking
		// [ 1 0 0 0 ] [ 0 -1 0 0 ] 0 1 1

		// Consume up to closing `]`
		for {
			if p.nextToken() == "]" {
				break
			}
			if p.scan() == scanner.EOF {
				return nil, fmt.Errorf("EOF inside valve 220 texture axes")
			}
		}

		// Next `[`
		p.scan() // expect `[`
		if p.nextToken() != "[" {
			// Maybe space detected? loop?
			// Standard Scanner skips space.
			// If not [, error?
		}

		for {
			if p.nextToken() == "]" {
				break
			}
			if p.scan() == scanner.EOF {
				return nil, fmt.Errorf("EOF inside valve 220 texture axes")
			}
		}

		// Now rot scaleX scaleY
		v, err := p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.Rotation = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.ScaleX = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.ScaleY = float32(v)

	} else {
		// Standard
		// text is offX
		// We already scanned it into `text` variable above
		v, err := p.parseNumberWithCurrent(text)
		if err != nil {
			return nil, err
		}
		pl.OffsetX = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.OffsetY = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.Rotation = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.ScaleX = float32(v)

		v, err = p.parseFloat()
		if err != nil {
			return nil, err
		}
		pl.ScaleY = float32(v)
	}

	return pl, nil
}

func (p *parser) parsePoint(pt *Point) error {
	// x y z
	var err error
	v, err := p.parseFloat()
	pt[0] = float32(v)
	if err != nil {
		return err
	}

	v, err = p.parseFloat()
	pt[1] = float32(v)
	if err != nil {
		return err
	}

	v, err = p.parseFloat()
	pt[2] = float32(v)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) parseFloat() (float64, error) {
	p.scan()
	text := p.nextToken()
	return p.parseNumberWithCurrent(text)
}

func (p *parser) parseNumberWithCurrent(text string) (float64, error) {
	if text == "-" {
		p.scan()
		valText := p.nextToken()
		v, err := strconv.ParseFloat(valText, 64)
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	return strconv.ParseFloat(text, 64)
}
