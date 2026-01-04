package qmap

import "github.com/go-gl/mathgl/mgl32"

// Point represents a 3D coordinate.
type Point = mgl32.Vec3

// Plane defined by 3 points and texture info.
type Plane struct {
	A, B, C  Point
	Texture  string
	OffsetX  float32
	OffsetY  float32
	Rotation float32
	ScaleX   float32
	ScaleY   float32
}

// Brush is a convex volume defined by planes.
type Brush struct {
	Planes []Plane
}

// Entity has key-value pairs and a list of brushes.
type Entity struct {
	Properties map[string]string
	Brushes    []Brush
}

// Map holds all entities.
type Map struct {
	Entities []Entity
}
