package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadOBJ membaca file .obj dan mengembalikan daftar segitiga dan semua vertex mentah
func ReadOBJ(path string) ([]Triangle, []Vector3, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var allVertices []Vector3
	var triangles []Triangle

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		opcode := parts[0]

		if opcode == "v" {
			// Parsing Vertex: v x y z
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			z, _ := strconv.ParseFloat(parts[3], 64)
			allVertices = append(allVertices, Vector3{X: x, Y: y, Z: z})

		} else if opcode == "f" {
			// Parsing Face: f i j k
			// Ingat: Indeks .obj dimulai dari 1, Go dimulai dari 0
			idx1, _ := strconv.Atoi(strings.Split(parts[1], "/")[0])
			idx2, _ := strconv.Atoi(strings.Split(parts[2], "/")[0])
			idx3, _ := strconv.Atoi(strings.Split(parts[3], "/")[0])

			triangles = append(triangles, Triangle{
				V1: allVertices[idx1-1],
				V2: allVertices[idx2-1],
				V3: allVertices[idx3-1],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return triangles, allVertices, nil
}

// WriteOBJ membuat file .obj baru berdasarkan daftar voxel yang ditemukan
func WriteOBJ(path string, voxels []Voxel) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	vertexCount := 0

	for _, voxel := range voxels {
		c := voxel.Center
		h := voxel.HalfSize

		// 1. Generate 8 titik sudut untuk satu kubus (voxel)
		vertices := []Vector3{
			{c.X - h, c.Y - h, c.Z - h}, // 1
			{c.X + h, c.Y - h, c.Z - h}, // 2
			{c.X + h, c.Y + h, c.Z - h}, // 3
			{c.X - h, c.Y + h, c.Z - h}, // 4
			{c.X - h, c.Y - h, c.Z + h}, // 5
			{c.X + h, c.Y - h, c.Z + h}, // 6
			{c.X + h, c.Y + h, c.Z + h}, // 7
			{c.X - h, c.Y + h, c.Z + h}, // 8
		}

		for _, v := range vertices {
			fmt.Fprintf(writer, "v %f %f %f\n", v.X, v.Y, v.Z)
		}

		// 2. Generate 12 wajah segitiga untuk membentuk kubus
		// Menggunakan offset vertexCount agar indeks merujuk ke vertex kubus ini
		base := vertexCount
		faces := [][]int{
			{1, 2, 3}, {1, 3, 4}, // Belakang
			{5, 6, 7}, {5, 7, 8}, // Depan
			{1, 5, 8}, {1, 8, 4}, // Kiri
			{2, 6, 7}, {2, 7, 3}, // Kanan
			{4, 3, 7}, {4, 7, 8}, // Atas
			{1, 2, 6}, {1, 6, 5}, // Bawah
		}

		for _, f := range faces {
			fmt.Fprintf(writer, "f %d %d %d\n", f[0]+base, f[1]+base, f[2]+base)
		}

		vertexCount += 8
	}

	return nil
}
