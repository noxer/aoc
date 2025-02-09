package utils

var Directions = []Vec{
	{X: -1}, // left / west
	{X: +1}, // right / east
	{Y: -1}, // up / north
	{Y: +1}, // down / south
}

type Vec struct {
	X, Y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y}
}

func (v Vec) Mul(s int) Vec {
	return Vec{v.X * s, v.Y * s}
}

func (v Vec) Zero() bool {
	return v.X == 0 && v.Y == 0
}
