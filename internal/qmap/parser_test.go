package qmap

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	// Adjust path relative to the test execution directory (internal/qmap)
	f, err := os.Open("../../assets/maps/test.map")
	if err != nil {
		t.Fatalf("failed to open test map: %v", err)
	}
	defer f.Close()

	m, err := Parse(f)
	if err != nil {
		t.Fatalf("failed to parse map: %v", err)
	}

	if len(m.Entities) == 0 {
		t.Fatal("expected at least one entity")
	}

	// Entity 0 should be worldspawn
	ent0 := m.Entities[0]
	if val, ok := ent0.Properties["classname"]; !ok || val != "worldspawn" {
		t.Errorf("expected entity 0 classname to be worldspawn, got %s", val)
	}

	// Check Brushes count
	// Visual check of test.map: looks like 6 or 7 brushes? "brush 0" to "brush 6" -> 7 brushes.
	if len(ent0.Brushes) != 7 {
		t.Errorf("expected 7 brushes in worldspawn, got %d", len(ent0.Brushes))
	}

	// Check first brush
	b0 := ent0.Brushes[0]
	// Brush 0 has 6 planes (cube)
	if len(b0.Planes) != 6 {
		t.Errorf("expected 6 planes in brush 0, got %d", len(b0.Planes))
	}

	// Check first plane
	// ( -256 -80 -32 ) ... materials/carpet_0 0 0 0 1 1
	p0 := b0.Planes[0]
	if p0.Texture != "materials/carpet_0" {
		t.Errorf("expected texture materials/carpet_0, got %s", p0.Texture)
	}

	if p0.A[0] != -256 {
		t.Errorf("expected p0.A.x = -256, got %f", p0.A[0])
	}
	if p0.A[1] != -80 {
		t.Errorf("expected p0.A.y = -80, got %f", p0.A[1])
	}
	if p0.A[2] != -32 {
		t.Errorf("expected p0.A.z = -32, got %f", p0.A[2])
	}

	if p0.OffsetX != 0 {
		t.Errorf("expected OffsetX 0, got %f", p0.OffsetX)
	}
	if p0.ScaleX != 1 {
		t.Errorf("expected ScaleX 1, got %f", p0.ScaleX)
	}
}
