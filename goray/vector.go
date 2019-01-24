package goray

import "math"

type Vec3f struct {
	x, y, z float64
}

func vSubtract(a, b Vec3f) Vec3f {
	return Vec3f{a.x - b.x, a.y - b.y, a.z - b.z}
}

func vDot(a, b Vec3f) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func (v Vec3f) norm() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func vMult(v Vec3f, m float64) Vec3f {
	return Vec3f{v.x * m, v.y * m, v.z * m}
}

func (v Vec3f) normalize() Vec3f {
	f := 1 / v.norm()
	return vMult(v, f)
}
