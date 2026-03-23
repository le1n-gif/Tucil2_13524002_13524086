package main

import (
	"math"
)

func IsIntersect(tri Triangle, box BoundingBox) bool {
	center := box.GetCenter()
	h := box.GetHalfSize()

	v0 := tri.V1.Sub(center)
	v1 := tri.V2.Sub(center)
	v2 := tri.V3.Sub(center)

	// Hitung tepi segitiga (edges)
	e0 := v1.Sub(v0)
	e1 := v2.Sub(v1)
	e2 := v0.Sub(v2)

	// Tes 1

	minX := math.Min(v0.X, math.Min(v1.X, v2.X))
	maxX := math.Max(v0.X, math.Max(v1.X, v2.X))
	if minX > h || maxX < -h {
		return false
	}

	minY := math.Min(v0.Y, math.Min(v1.Y, v2.Y))
	maxY := math.Max(v0.Y, math.Max(v1.Y, v2.Y))
	if minY > h || maxY < -h {
		return false
	}

	minZ := math.Min(v0.Z, math.Min(v1.Z, v2.Z))
	maxZ := math.Max(v0.Z, math.Max(v1.Z, v2.Z))
	if minZ > h || maxZ < -h {
		return false
	}

	// Tes 2
	normal := Cross(e0, e1)
	if !planeBoxOverlap(normal, v0, h) {
		return false
	}

	// Tes 3

	if !testAxis(e0.Z, -e0.Y, v0, v1, v2, h, "X") {
		return false
	}
	if !testAxis(e1.Z, -e1.Y, v0, v1, v2, h, "X") {
		return false
	}
	if !testAxis(e2.Z, -e2.Y, v0, v1, v2, h, "X") {
		return false
	}

	if !testAxis(-e0.Z, e0.X, v0, v1, v2, h, "Y") {
		return false
	}
	if !testAxis(-e1.Z, e1.X, v0, v1, v2, h, "Y") {
		return false
	}
	if !testAxis(-e2.Z, e2.X, v0, v1, v2, h, "Y") {
		return false
	}

	if !testAxis(e0.Y, -e0.X, v0, v1, v2, h, "Z") {
		return false
	}
	if !testAxis(e1.Y, -e1.X, v0, v1, v2, h, "Z") {
		return false
	}
	if !testAxis(e2.Y, -e2.X, v0, v1, v2, h, "Z") {
		return false
	}

	return true
}

func planeBoxOverlap(normal Vector3, v0 Vector3, h float64) bool {
	vMin := Vector3{}
	vMax := Vector3{}

	if normal.X > 0 {
		vMin.X = -h
		vMax.X = h
	} else {
		vMin.X = h
		vMax.X = -h
	}

	if normal.Y > 0 {
		vMin.Y = -h
		vMax.Y = h
	} else {
		vMin.Y = h
		vMax.Y = -h
	}

	if normal.Z > 0 {
		vMin.Z = -h
		vMax.Z = h
	} else {
		vMin.Z = h
		vMax.Z = -h
	}

	if Dot(normal, vMin.Sub(v0)) > 0 {
		return false
	}
	if Dot(normal, vMax.Sub(v0)) >= 0 {
		return true
	}

	return false
}

func testAxis(a, b float64, v0, v1, v2 Vector3, h float64, axisType string) bool {
	var p0, p1, p2 float64

	switch axisType {
	case "X":
		p0 = a*v0.Y + b*v0.Z
		p1 = a*v1.Y + b*v1.Z
		p2 = a*v2.Y + b*v2.Z
	case "Y":
		p0 = a*v0.X + b*v0.Z
		p1 = a*v1.X + b*v1.Z
		p2 = a*v2.X + b*v2.Z
	case "Z":
		p0 = a*v0.X + b*v0.Y
		p1 = a*v1.X + b*v1.Y
		p2 = a*v2.X + b*v2.Y
	}

	min := math.Min(p0, math.Min(p1, p2))
	max := math.Max(p0, math.Max(p1, p2))

	// Proyeksi radius kotak ke sumbu L
	rad := math.Abs(a)*h + math.Abs(b)*h

	// Cek overlap
	if min > rad || max < -rad {
		return false
	}
	return true
}
