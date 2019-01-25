package goray

import "math"

type Vec3f struct {
	x, y, z float64
}

func (v Vec3f) subtract(b Vec3f) Vec3f {
	return Vec3f{v.x - b.x, v.y - b.y, v.z - b.z}
}

func (v Vec3f) add(b Vec3f) Vec3f {
	return Vec3f{v.x + b.x, v.y + b.y, v.z + b.z}
}

func (v Vec3f) dot(b Vec3f) float64 {
	return v.x*b.x + v.y*b.y + v.z*b.z
}

func (v Vec3f) norm() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vec3f) mult(m float64) Vec3f {
	return Vec3f{v.x * m, v.y * m, v.z * m}
}

func (v Vec3f) normalize() Vec3f {
	f := 1 / v.norm()
	return v.mult(f)
}

func (v Vec3f) reflect(N Vec3f) Vec3f {
	return v.subtract(N.mult(2. * v.dot(N)))
}

func (v Vec3f) refract(N Vec3f, refractiveIndex float64) Vec3f {
	cosi := -math.Max(-1., math.Min(1., v.dot(N)))
	etai := 1.
	etat := refractiveIndex
	var n Vec3f
	if cosi < 0 {
		cosi = -cosi
		etai, etat = etat, etai
		n = N.mult(-1.)
	} else {
		n = N
	}
	eta := etai / etat
	k := 1 - eta*eta*(1.-cosi*cosi)
	var r Vec3f
	if k < 0 {
		r = Vec3f{}
	} else {
		r = v.mult(eta).add(n.mult(eta*cosi - math.Sqrt(k)))
	}
	return r
}
