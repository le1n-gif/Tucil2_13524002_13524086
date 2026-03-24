package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type AppStats struct {
	MaxDepth     int
	NodesCreated []int // per-level: berapa node yang dibuat
	NodesSkipped []int // per-level: berapa node yang di-skip
	TotalVoxels  int64 // total voxel di max depth
	StartTime    time.Time
	EndTime      time.Time
	InputPath    string
	OutputPath   string
	mu           sync.Mutex // lock buat atomic read/write
}

// init
func NewStats(maxDepth int, inputPath, outputPath string) *AppStats {
	return &AppStats{
		MaxDepth:     maxDepth,
		NodesCreated: make([]int, maxDepth+1),
		NodesSkipped: make([]int, maxDepth+1),
		TotalVoxels:  0,
		StartTime:    time.Now(),
		InputPath:    inputPath,
		OutputPath:   outputPath,
	}
}

// catet increment node yang dibuat
func (s *AppStats) RecordNode(depth int) {
	if depth >= 0 && depth <= s.MaxDepth {
		s.mu.Lock()
		s.NodesCreated[depth]++
		s.mu.Unlock()
	}
}

// skip increment node yang di skip
func (s *AppStats) RecordSkip(depth int) {
	if depth >= 0 && depth <= s.MaxDepth {
		s.mu.Lock()
		s.NodesSkipped[depth]++
		s.mu.Unlock()
	}
}

// increment jumlah voxel
func (s *AppStats) RecordVoxel() {
	atomic.AddInt64(&s.TotalVoxels, 1)
}

// stop timer
func (s *AppStats) StopTimer() {
	s.EndTime = time.Now()
}

// print
func (s *AppStats) PrintReport() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("=[]=[]=[]= REPORTS =[]=[]=[]=")
	fmt.Printf("File input:  %s\n", s.InputPath)
	fmt.Printf("File output: %s\n", s.OutputPath)
	fmt.Println()

	// hitung jumlah vertex sama face dari voxel
	totalVoxels := atomic.LoadInt64(&s.TotalVoxels)
	vertices := totalVoxels * 8
	faces := totalVoxels * 12

	fmt.Printf("Total voxel:    %d\n", totalVoxels)
	fmt.Printf("Total vertex:   %d (voxel × 8)\n", vertices)
	fmt.Printf("Total face:    %d (voxel × 12)\n", faces)
	fmt.Println()

	// hitung durasi
	duration := s.EndTime.Sub(s.StartTime)
	fmt.Printf("Waktu proses: %v\n", duration)
	fmt.Println()

	// printlnn statistik per level
	fmt.Println("Statistik per level:")
	fmt.Println("Level | Dibuat | Di-skip")
	fmt.Println("------|--------|--------")
	for depth := 0; depth <= s.MaxDepth; depth++ {
		fmt.Printf("%5d | %6d | %7d\n", depth, s.NodesCreated[depth], s.NodesSkipped[depth])
	}
	fmt.Println()
}
