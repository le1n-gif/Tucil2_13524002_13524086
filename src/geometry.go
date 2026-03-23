package main

import (
	"math"
)

// type struct
type Vector3 struct {
	X, Y, Z float64
}

type Triangle struct {
	V1, V2, V3 Vector3
}

type BoundingBox struct {
	Min Vector3
	Max Vector3
}

type Voxel struct {
	Center   Vector3
	HalfSize float64
}

// Operasi Vektor
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vector3) Mul(scalar float64) Vector3 {
	return Vector3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func Cross(a, b Vector3) Vector3 {
	return Vector3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func Dot(a, b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (v Vector3) Abs() Vector3 {
	return Vector3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

// Geometri

func (b BoundingBox) GetCenter() Vector3 {
	return Vector3{
		(b.Min.X + b.Max.X) / 2,
		(b.Min.Y + b.Max.Y) / 2,
		(b.Min.Z + b.Max.Z) / 2,
	}
}

func (b BoundingBox) GetHalfSize() float64 {
	return (b.Max.X - b.Min.X) / 2
}

//  Fungsi Inisialisasi

// CalculateInitialBounds mencari batas ruang 3D yang membungkus seluruh vertex
// Fungsi ini memastikan Bounding Box berbentuk kubus agar pembagian Octree uniform
func CalculateInitialBounds(vertices []Vector3) BoundingBox {
	if len(vertices) == 0 {
		return BoundingBox{}
	}

	// 1. Cari titik min dan max murni dari data vertex
	min := Vector3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	max := Vector3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}

	for _, v := range vertices {
		min.X = math.Min(min.X, v.X)
		min.Y = math.Min(min.Y, v.Y)
		min.Z = math.Min(min.Z, v.Z)
		max.X = math.Max(max.X, v.X)
		max.Y = math.Max(max.Y, v.Y)
		max.Z = math.Max(max.Z, v.Z)
	}

	// 2. Hitung Center
	center := Vector3{
		(min.X + max.X) / 2,
		(min.Y + max.Y) / 2,
		(min.Z + max.Z) / 2,
	}

	// 3. Cari sisi terpanjang supaya kotak berbentuk kubus
	sizeX := max.X - min.X
	sizeY := max.Y - min.Y
	sizeZ := max.Z - min.Z
	maxSide := math.Max(sizeX, math.Max(sizeY, sizeZ))

	// Tambahkan sedikit padding (misal 1%) supaya tidak ada vertex yang pas di garis batas
	maxSide *= 1.01
	halfSide := maxSide / 2

	// 4. Konstruksi Bounding Box baru yang simetris (Kubus)
	return BoundingBox{
		Min: Vector3{center.X - halfSide, center.Y - halfSide, center.Z - halfSide},
		Max: Vector3{center.X + halfSide, center.Y + halfSide, center.Z + halfSide},
	}
}
