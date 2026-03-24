package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	//baca path file
	var pathFile string
	fmt.Print("Masukkan path file .obj (contoh: ../test/cow.obj): ")
	fmt.Scanln(&pathFile)

	//baca file OBJ
	triangles, vertices, err := ReadOBJ(pathFile)
	if err != nil {
		fmt.Printf("Error baca file: %v\n", err)
		return
	}

	fmt.Printf("Berhasil load %d segitiga dan %d vertex\n", len(triangles), len(vertices))

	//hitung bounding box awal
	rootBox := CalculateInitialBounds(vertices)
	fmt.Printf("Bounding Box===  Min: (%.2f, %.2f, %.2f), Max: (%.2f, %.2f, %.2f)\n",
		rootBox.Min.X, rootBox.Min.Y, rootBox.Min.Z,
		rootBox.Max.X, rootBox.Max.Y, rootBox.Max.Z)

	//input depth
	maxDepth := getMaxDepth()
	outputPath := generateOutputPath(pathFile, maxDepth)

	//init stats
	stats := NewStats(maxDepth, pathFile, outputPath)

	//buat voxelizer dan build octree
	fmt.Printf("\nMemproses octree depth %d...\n", maxDepth)
	voxelizer := NewVoxelizer(maxDepth, triangles)

	// integration stats ke voxelizer
	voxelizer.Stats = stats

	voxelizer.Build(rootBox)

	//stop timer
	stats.StopTimer()

	//save output ke obj
	fmt.Printf("Menyimpan hasil ke %s...\n", outputPath)
	err = WriteOBJ(outputPath, voxelizer.Voxels)
	if err != nil {
		fmt.Printf("Error save file: %v\n", err)
		return
	}

	//print reports
	fmt.Println()
	stats.PrintReport()

	//tanya user mau lihat 3D ato gak
	if len(voxelizer.Voxels) > 0 && len(voxelizer.Voxels) <= 1000000 {
		var viewChoice string
		fmt.Print("\nApakah ingin melihat model dalam 3D viewer? (y/n): ")
		fmt.Scanln(&viewChoice)

		if strings.ToLower(viewChoice) == "y" || strings.ToLower(viewChoice) == "yes" {
			fmt.Println("Membuka viewer...")
			fmt.Println("Kontrol: Arrow Keys=putar, W/S=zoom, P=ganti mode, R=reset")
			err := RunViewer(voxelizer.Voxels)
			if err != nil {
				fmt.Printf("Error viewer: %v\n", err)
			}
		}
	}
}

// buat path output dari path input
func generateOutputPath(inputPath string, maxDepth int) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, name+"_voxelized_d"+fmt.Sprintf("%d", maxDepth)+".obj")
}

// minta depth dari user dengan validasi
func getMaxDepth() int {
	for {
		var depthStr string
		fmt.Print("Depth octree (default 6): ")
		fmt.Scanln(&depthStr)

		// kalo kosong, pake default 6
		if depthStr == "" {
			return 6
		}

		// coba parse
		depth, err := strconv.Atoi(depthStr)
		if err != nil {
			fmt.Printf("Input salah, harus angka.\n")
			continue
		}

		// validasi range
		if depth < 1 || depth > 15 {
			fmt.Printf("Range harus 1-15, kamu input: %d\n", depth)
			continue
		}

		return depth
	}
}
