package goray

import (
	"math"
)

type Sphere struct {
	center Vec3f
	radius float64
}

func (s Sphere) rayIntersect(orig, dir Vec3f) bool {
	L := vSubtract(s.center, orig)
	tca := vDot(L, dir)
	d2 := vDot(L, L) - tca*tca
	if d2 > (s.radius * s.radius) {
		return false
	}
	thc := math.Sqrt((s.radius * s.radius) - d2)
	dist := tca - thc
	t1 := tca + thc

	if dist < 0 {
		dist = t1
	}
	if dist < 0 {
		return false
	}

	return true
}
