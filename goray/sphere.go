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
	diffuseColor Vec3f
}

func (s Sphere) rayIntersect(orig, dir Vec3f) (bool, float64) {
	var dist float64
	L := vSubtract(s.center, orig)
	tca := vDot(L, dir)
	d2 := vDot(L, L) - tca*tca
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
