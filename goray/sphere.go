package goray

import (
	"math"
)

type Sphere struct {
	center   Vec3f
	radius   float64
	material Material
}

type Material struct {
	albedo           [3]float64
	diffuseColor     Vec3f
	specularExponent float64
}

func (s Sphere) rayIntersect(orig, dir Vec3f) (bool, float64) {
	var dist float64
	L := s.center.subtract(orig)
	tca := L.dot(dir)
	d2 := L.dot(L) - tca*tca
	if d2 > (s.radius * s.radius) {
		return false, dist
	}
	thc := math.Sqrt((s.radius * s.radius) - d2)
	dist = tca - thc
	t1 := tca + thc

	if dist < 0 {
		dist = t1
	}
	if dist < 0 {
		return false, dist
	}

	return true, dist
}
