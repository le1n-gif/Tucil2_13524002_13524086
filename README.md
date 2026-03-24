# Tucil 2 - 3D Object Voxelization 


## Anggota Tim 
- Daniel Anindito Nugroho (13524002)
- Al Farabi (13524086)



## Deskripsi Project

Program ini membaca model 3D dari format OBJ dan melakukan voxelisasi menggunakan struktur data **Octree**. Hasil voxelisasi dapat dilihat dalam interactive 3D viewer dengan dukungan rotasi, zoom, dan perubahan mode proyeksi.

### Fitur Utama
- Parser OBJ untuk membaca model 3D
- Voxelisasi menggunakan recursive octree subdivision
- Deteksi interseksi triangle-AABB (Axis-Aligned Bounding Box)
- Interactive 3D viewer dengan Ebiten
- Parallel processing untuk performa optimal
- Statistik detail (nodes per level, execution time, dll)

---
## Struktur Project

```
Tucil2_13524002_13524086/
├── README.md                 # File ini
├── LICENSE                   # Lisensi project
├── src/
│   ├── main.go              # Entry point program
│   ├── parser.go            # Parser file OBJ
│   ├── geometry.go          # Struktur data geometri (Vector3, Triangle, dll)
│   ├── octree.go            # Voxelizer dengan octree
│   ├── intersection.go       # Triangle-AABB intersection test
│   ├── viewer.go            # Interactive 3D viewer dengan Ebiten
│   └── stats.go             # Statistik dan pelaporan
├── test/                     # Folder untuk file OBJ test
└── bin/                      # Output executable setelah compile
```
## Requirements

### 1. **Go Programming Language**
- **Versi minimum:** Go 1.16 atau lebih baru
- **Download:** [https://golang.org/dl/](https://golang.org/dl/)

### 2. **GCC Compiler (untuk Ebiten)**
Ebiten memerlukan compiler C/C++ yang kompatibel:

- **Windows:** MinGW-w64
- **macOS:** Xcode Command Line Tools
- **Linux:** GCC dari package manager

### 3. **Dependensi Go**
```
github.com/hajimehoshi/ebiten/v2  
```

---

## Instalasi Langkah-Demi-Langkah

### Step 1: Install Go

#### Windows
1. Download installer dari [golang.org/dl](https://golang.org/dl/)
2. Jalankan installer dan ikuti petunjuk
3. Verifikasi instalasi:
   ```bash
   go version
   ```

#### macOS
```bash
# Dengan Homebrew
brew install go

# Atau download dari https://golang.org/dl/
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install golang-go
```

---

### Step 2: Install GCC Compiler

#### Windows (MinGW-w64)
**Opsi A: Menggunakan Chocolatey (Recommended)**
```bash
# Install Chocolatey terlebih dahulu jika belum ada
choco install mingw

# Atau install langsung
choco install gcc
```

**Opsi B: Manual Download**
1. Download dari [MinGW-w64](https://www.mingw-w64.org/downloads/)
2. Extract ke folder (contoh: `C:\mingw64`)
3. Tambahkan ke PATH:
   - Buka `System Properties` → `Environment Variables`
   - Tambahkan `C:\mingw64\bin` ke PATH
4. Verifikasi:
   ```bash
   gcc --version
   ```

#### macOS
```bash
# Xcode Command Line Tools
xcode-select --install
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get install build-essential
sudo apt-get install libgl1-mesa-dev libalut-dev libxrandr-dev libxcursor-dev libxi-dev
```

---

### Step 3: Download Project & Install Dependensi

```bash
# 1. Clone/navigate ke folder project
cd path/to/Tucil2_13524002_13524086

# 2. Initialize/Update Go modules (jika belum ada go.mod)
go mod init voxelizer

# 3. Download semua dependensi
go mod download
go get github.com/hajimehoshi/ebiten/v2

# 4. Tidy up go.mod
go mod tidy
```

---

## Cara Menjalankan Program

### Build & Run dari Source

```bash
# Navigate ke folder project
cd src

# Jalankan program
go run .

# Atau compile terlebih dahulu
go build -o voxelizer.exe
./voxelizer.exe
```

### Penggunaan Program

1. **Input file OBJ:**
   ```
   Masukkan path file .obj (contoh: ../test/cow.obj): 
   ```

2. **Input kedalaman octree:**
   ```
   Depth octree (default 6): 8
   ```

3. **Pilih melihat 3D viewer:**
   ```
   Apakah ingin melihat model dalam 3D viewer? (y/n): y
   ```

### Kontrol 3D Viewer

| Kontrol | Aksi |
|---------|------|
| `←` `→` | Rotasi horizontal |
| `↑` `↓` | Rotasi vertikal |
| `W` | Zoom in |
| `S` | Zoom out |
| `P` | Toggle mode proyeksi (Orthogonal ↔ Perspective) |
| `R` | Reset ke posisi awal |

---




## Troubleshooting

### Error: `gcc not found` atau `cc not found`
**Solusi:**
- Pastikan MinGW-w64 sudah diinstall
- Verifikasi PATH dengan `gcc --version`
- Restart terminal/IDE setelah instalasi

### Error: Module not found `github.com/hajimehoshi/ebiten/v2`
**Solusi:**
```bash
go get -u github.com/hajimehoshi/ebiten/v2
go mod tidy
```

### Error: Build gagal di Linux
**Solusi:**
```bash
sudo apt-get install libgl1-mesa-dev libalut-dev libxrandr-dev libxcursor-dev libxi-dev
```

### Viewer tidak muncul / Crash
**Solusi:**
- Pastikan file OBJ valid
- Gunakan depth yang wajar (1-10)
- Coba dengan file test terlebih dahulu

---
