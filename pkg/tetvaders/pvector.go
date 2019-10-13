package tetvaders

// Originally from https://github.com/zacg/boids

import (
	"math"
	"math/rand"
)

// PVector represents a vector (in 3D if needed)
type PVector struct {
	X float64
	Y float64
	Z float64
}

// NewPVector2D creates a 2D vector
func NewPVector2D(x float64, y float64) PVector {
	return PVector{x, y, 0}
}

// NewPVectorFromAngle creates a vector from an angle
func NewPVectorFromAngle(angle float64) PVector {
	result := PVector{}

	result.X = math.Cos(angle)
	result.Y = math.Sin(angle)
	result.Z = 0
	return result
}

// NewRandom2dPVector creates a random angle vector
func NewRandom2dPVector() PVector {
	return NewPVectorFromAngle(rand.Float64() * math.Pi * 2)
}

// NewRandom3dPVector creates a new random 3d PVector
func NewRandom3dPVector() PVector {
	var angle float64 = 0
	var vz float64 = 0
	var result = PVector{}

	angle = rand.Float64() * math.Pi * 2
	result.Z = rand.Float64()*2 - 1

	result.X = (math.Sqrt(1-vz*vz) * math.Cos(angle))
	result.Y = (math.Sqrt(1-vz*vz) * math.Sin(angle))

	return result
}

// Mag calculates the length of the vector
func (pvec *PVector) Mag() float64 {
	return math.Sqrt(pvec.X*pvec.X + pvec.Y*pvec.Y + pvec.Z*pvec.Z)
}

// MagSq calculates the squared magnitude of the vector
// (x*x + y*y + z*z)
func (pvec *PVector) MagSq() float64 {
	return pvec.X*pvec.X + pvec.Y*pvec.Y + pvec.Z*pvec.Z
}

// Dist returns distance to specified neighbour
func (pvec *PVector) Dist(neighbour PVector) float64 {
	dx := pvec.X - neighbour.X
	dy := pvec.Y - neighbour.Y
	dz := pvec.Z - neighbour.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

//Limit the magnitude of vector to specified max
func (pvec *PVector) Limit(max float64) {
	if pvec.MagSq() > max*max {
		pvec.Normalize()
		pvec.Mult(max)
	}
}

//Add adds 2 vectors
func (pvec *PVector) Add(pVector PVector) {
	pvec.X += pVector.X
	pvec.Y += pVector.Y
	pvec.Z += pVector.Z
}

//Div divides the vector by the specified scalar
func (pvec *PVector) Div(n float64) {
	pvec.X /= n
	pvec.Y /= n
	pvec.Z /= n
}

//Mult multiplies vector by specified scalar
func (pvec *PVector) Mult(n float64) {
	pvec.X *= n
	pvec.Y *= n
	pvec.Z *= n
}

//Sub decrements vector by 1
func (pvec *PVector) Sub() {
	pvec.X--
	pvec.Y--
	pvec.Z--
}

//Diff subtracts specified vector
//returns difference
func (pvec *PVector) Diff(n PVector) PVector {
	return PVector{
		X: pvec.X - n.X,
		Y: pvec.Y - n.Y,
		Z: pvec.Z - n.Z,
	}
}

//Inc increments vector by 1
func (pvec *PVector) Inc() {
	pvec.X++
	pvec.Y++
	pvec.Z++
}

//Normalize normalizes vector to length 1
func (pvec *PVector) Normalize() {
	m := pvec.Mag()
	if m != 0 && m != 1 {
		pvec.Div(m)
	}
}
