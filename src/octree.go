package main

import (
	"sync"
)

type Voxelizer struct {
	MaxDepth        int
	Triangles       []Triangle
	Voxels          []Voxel
	NodesPerLevel   []int
	SkippedPerLevel []int
	mu              sync.Mutex
}

func NewVoxelizer(maxDepth int, triangles []Triangle) *Voxelizer {
	return &Voxelizer{
		MaxDepth:        maxDepth,
		Triangles:       triangles,
		Voxels:          make([]Voxel, 0),
		NodesPerLevel:   make([]int, maxDepth+1),
		SkippedPerLevel: make([]int, maxDepth+1),
	}
}

func (v *Voxelizer) Build(rootBox BoundingBox) {
	var wg sync.WaitGroup
	v.subdivide(rootBox, 0, v.Triangles, &wg, true)
	wg.Wait()
}

func (v *Voxelizer) subdivide(box BoundingBox, depth int, triangles []Triangle, wg *sync.WaitGroup, isParallel bool) {
	intersectingTriangles := []Triangle{}
	for _, tri := range triangles {
		if IsIntersect(tri, box) {
			intersectingTriangles = append(intersectingTriangles, tri)
		}
	}

	if len(intersectingTriangles) == 0 {
		v.mu.Lock()
		v.SkippedPerLevel[depth]++
		v.mu.Unlock()
		if isParallel {
			wg.Done()
		}
		return
	}

	v.mu.Lock()
	v.NodesPerLevel[depth]++
	v.mu.Unlock()

	if depth == v.MaxDepth {
		v.mu.Lock()
		v.Voxels = append(v.Voxels, Voxel{
			Center:   box.GetCenter(),
			HalfSize: box.GetHalfSize(),
		})
		v.mu.Unlock()
		if isParallel {
			wg.Done()
		}
		return
	}

	children := v.generateChildren(box)
	useGoroutine := depth < 3

	for _, childBox := range children {
		if useGoroutine {
			wg.Add(1)
			go v.subdivide(childBox, depth+1, intersectingTriangles, wg, true)
		} else {
			v.subdivide(childBox, depth+1, intersectingTriangles, wg, false)
		}
	}

	if isParallel {
		wg.Done()
	}
}

func (v *Voxelizer) generateChildren(box BoundingBox) [8]BoundingBox {
	center := box.GetCenter()
	min := box.Min
	max := box.Max

	var children [8]BoundingBox
	// 0: Kiri - Bawah - Depan
	children[0] = BoundingBox{Min: Vector3{min.X, min.Y, min.Z}, Max: Vector3{center.X, center.Y, center.Z}}

	// 1: Kanan - Bawah - Depan
	children[1] = BoundingBox{Min: Vector3{center.X, min.Y, min.Z}, Max: Vector3{max.X, center.Y, center.Z}}

	// 2: Kiri - Atas - Depan
	children[2] = BoundingBox{Min: Vector3{min.X, center.Y, min.Z}, Max: Vector3{center.X, max.Y, center.Z}}

	// 3: Kanan - Atas - Depan
	children[3] = BoundingBox{Min: Vector3{center.X, center.Y, min.Z}, Max: Vector3{max.X, max.Y, center.Z}}

	// 4: Kiri - Bawah - Belakang
	children[4] = BoundingBox{Min: Vector3{min.X, min.Y, center.Z}, Max: Vector3{center.X, center.Y, max.Z}}

	// 5: Kanan - Bawah - Belakang
	children[5] = BoundingBox{Min: Vector3{center.X, min.Y, center.Z}, Max: Vector3{max.X, center.Y, max.Z}}

	// 6: Kiri - Atas - Belakang
	children[6] = BoundingBox{Min: Vector3{min.X, center.Y, center.Z}, Max: Vector3{center.X, max.Y, max.Z}}

	// 7: Kanan - Atas - Belakang
	children[7] = BoundingBox{Min: Vector3{center.X, center.Y, center.Z}, Max: Vector3{max.X, max.Y, max.Z}}

	return children
}
