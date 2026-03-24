package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// read file .obj dan return daftar segitiga dan semua vertex mentah
func ReadOBJ(path string) ([]Triangle, []Vector3, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Error: Tidak bisa buka file: %w", err)
	}
	defer file.Close()

	var allVertices []Vector3
	var triangles []Triangle
	var lineNum int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// skip baris kosong atau komentar
		if line == "" {
			continue
		}

		parts := strings.Fields(line)

		// validate punya parts minimal
		if len(parts) == 0 {
			continue
		}

		opcode := parts[0]

		// vertex parsing (v x y z)
		if opcode == "v" {
			// validasi: harus tepat 4 bagian (v + 3 koordinat)
			if len(parts) != 4 {
				return nil, nil, fmt.Errorf(" PARSE ERROR Baris %d:\n   vertex harus 3 angka (v x y z)\n   dapet: %d item\n   data: %v", lineNum, len(parts)-1, parts)
			}

			// parse 3 koordinat sebagai float64
			x, errX := strconv.ParseFloat(parts[1], 64)
			y, errY := strconv.ParseFloat(parts[2], 64)
			z, errZ := strconv.ParseFloat(parts[3], 64)

			// validasi parsing berhasil
			if errX != nil || errY != nil || errZ != nil {
				return nil, nil, fmt.Errorf("PARSE ERROR Baris %d:\n   koordinat vertex harus angka (bukan text/huruf)\n   dapet: '%s' '%s' '%s'", lineNum, parts[1], parts[2], parts[3])
			}

			allVertices = append(allVertices, Vector3{X: x, Y: y, Z: z})

		} else if opcode == "f" {
			// validasi: harus tepat 4 bagian (f + 3 indeks)
			if len(parts) != 4 {
				return nil, nil, fmt.Errorf("PARSE ERROR Baris %d:\n   face harus 3 index (f i j k)\n   dapet: %d item\n   data: %v", lineNum, len(parts)-1, parts)
			}

			// parsing 3 indeks (handle format i, i/vt, i/vt/vn)
			idx1Str := strings.Split(parts[1], "/")[0]
			idx2Str := strings.Split(parts[2], "/")[0]
			idx3Str := strings.Split(parts[3], "/")[0]

			// parse jadi integer
			idx1, err1 := strconv.Atoi(idx1Str)
			idx2, err2 := strconv.Atoi(idx2Str)
			idx3, err3 := strconv.Atoi(idx3Str)

			// validasi parsing berhasil
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, nil, fmt.Errorf("PARSE ERROR Baris %d:\n   index face harus angka (bukan text/huruf)\n   dapet: '%s' '%s' '%s'", lineNum, idx1Str, idx2Str, idx3Str)
			}

			// validasi index dalam range
			if idx1 < 1 || idx1 > len(allVertices) || idx2 < 1 || idx2 > len(allVertices) || idx3 < 1 || idx3 > len(allVertices) {
				return nil, nil, fmt.Errorf("PARSE ERROR Baris %d:\n   index out of range (vertex tidak ada)\n   max vertex: %d, tapi dapet: %d %d %d", lineNum, len(allVertices), idx1, idx2, idx3)
			}

			// semua valid, tambah triangle
			triangles = append(triangles, Triangle{
				V1: allVertices[idx1-1],
				V2: allVertices[idx2-1],
				V3: allVertices[idx3-1],
			})
		} else {
			// ERROR: typo atau format salah
			return nil, nil, fmt.Errorf("PARSE ERROR Baris %d:\n   opcode tidak valid: '%s'\n   hanya bisa 'v' (vertex) atau 'f' (face)\n   data: %v", lineNum, opcode, parts)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("ERROR baca file: %w", err)
	}

	return triangles, allVertices, nil
}

// create  file .obj baru berdasarkan daftar voxel yang ditemukan
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

		// Generate 8 titik sudut untuk satu kubus (voxel)
		vertices := []Vector3{
			{c.X - h, c.Y - h, c.Z - h},
			{c.X + h, c.Y - h, c.Z - h},
			{c.X + h, c.Y + h, c.Z - h},
			{c.X - h, c.Y + h, c.Z - h},
			{c.X - h, c.Y - h, c.Z + h},
			{c.X + h, c.Y - h, c.Z + h},
			{c.X + h, c.Y + h, c.Z + h},
			{c.X - h, c.Y + h, c.Z + h},
		}

		for _, v := range vertices {
			fmt.Fprintf(writer, "v %f %f %f\n", v.X, v.Y, v.Z)
		}

		// Generate 12 face triangle buat bentuk kubus
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
