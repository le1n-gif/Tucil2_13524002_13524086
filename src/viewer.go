package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Matrix4 = matriks 4x4 buat tranformasi 3D
type Matrix4 [4][4]float64

// Vector2 = koordinat 2D di layar
type Vector2 struct {
	X, Y float64
}

// Camera = kamera buat liat model 3D
type Camera struct {
	Position Vector3 // posisi mata
	Target   Vector3 // titik yang diliat
	Up       Vector3 // arah atas
	FOV      float64 // sudut pandang (derajat)
}

// Viewer = struktur utama buat render 3D interaktif
type Viewer struct {
	Voxels        []Voxel
	Camera        Camera
	ViewMatrix    Matrix4
	ProjMatrix    Matrix4
	RotX          float64 // rotasi X (euler)
	RotY          float64 // rotasi Y (euler)
	Distance      float64 // jarak mata dari objek (buat zoom)
	ScreenWidth   int
	ScreenHeight  int
	PixelBuffer   map[string]color.Color
	UseOrthogonal bool // pake proyeksi ortografis (lebih simpel)
}

// NewViewer membuat instance Viewer baru
func NewViewer(voxels []Voxel, screenWidth, screenHeight int) *Viewer {
	camera := Camera{
		Position: Vector3{X: 0, Y: 0, Z: 15},
		Target:   Vector3{X: 0, Y: 0, Z: 0},
		Up:       Vector3{X: 0, Y: 1, Z: 0},
		FOV:      45,
	}

	viewer := &Viewer{
		Voxels:        voxels,
		Camera:        camera,
		RotX:          0,
		RotY:          0,
		Distance:      15, // jarak awal dari target
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		PixelBuffer:   make(map[string]color.Color),
		UseOrthogonal: true, // pake ortografis dulu (lebih gampang)
	}

	viewer.UpdateMatrices()
	return viewer
}

// ============== MATRIX OPERATIONS (Manual Implementation) ==============

// Identity = buat matriks identitas 4x4
func Identity() Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// kaliin dua matriks 4x4
func MatrixMul(a, b Matrix4) Matrix4 {
	var result Matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// kaliin matriks 4x4 sama vektor 3D (pake w=1)
func MatrixVecMul(m Matrix4, v Vector3) Vector3 {
	x := m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]
	y := m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]
	z := m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]
	w := m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]

	// bagi sama w (normalisasi perspective)
	if w != 0 {
		return Vector3{x / w, y / w, z / w}
	}
	return Vector3{x, y, z}
}

// buat matriks rotasi sumbu X (radian)
func RotationX(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{1, 0, 0, 0},
		{0, c, -s, 0},
		{0, s, c, 0},
		{0, 0, 0, 1},
	}
}

// buat matriks rotasi sumbu Y (radian)
func RotationY(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{c, 0, s, 0},
		{0, 1, 0, 0},
		{-s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

// buat matriks rotasi sumbu Z (radian)
func RotationZ(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{c, -s, 0, 0},
		{s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// translasi
func Translation(x, y, z float64) Matrix4 {
	return Matrix4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

// buat matriks view (posisi mata)
func LookAt(position, target, up Vector3) Matrix4 {
	// hitung arah pandang dan sumbu lokal
	forward := target.Sub(position)
	forwardLen := math.Sqrt(forward.X*forward.X + forward.Y*forward.Y + forward.Z*forward.Z)
	if forwardLen > 0 {
		forward = Vector3{forward.X / forwardLen, forward.Y / forwardLen, forward.Z / forwardLen}
	}

	// hitung arah kanan
	right := Cross(forward, up)
	rightLen := math.Sqrt(right.X*right.X + right.Y*right.Y + right.Z*right.Z)
	if rightLen > 0 {
		right = Vector3{right.X / rightLen, right.Y / rightLen, right.Z / rightLen}
	}

	// hitung arah atas yang tegak lurus
	up = Cross(right, forward)

	// buat matriks view
	view := Matrix4{
		{right.X, right.Y, right.Z, -Dot(right, position)},
		{up.X, up.Y, up.Z, -Dot(up, position)},
		{-forward.X, -forward.Y, -forward.Z, Dot(forward, position)},
		{0, 0, 0, 1},
	}

	return view
}

// buat matriks proyeksi perspektif
func Perspective(fov, aspect, near, far float64) Matrix4 {
	// fov udah derajat, ubah ke radian
	fovRad := fov * math.Pi / 180.0
	f := 1.0 / math.Tan(fovRad/2.0)

	return Matrix4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) / (near - far), (2 * far * near) / (near - far)},
		{0, 0, -1, 0},
	}
}

// buat matriks proyeksi ortografis
func Orthographic(left, right, bottom, top, near, far float64) Matrix4 {
	return Matrix4{
		{2 / (right - left), 0, 0, -(right + left) / (right - left)},
		{0, 2 / (top - bottom), 0, -(top + bottom) / (top - bottom)},
		{0, 0, -2 / (far - near), -(far + near) / (far - near)},
		{0, 0, 0, 1},
	}
}

// update matriks view dan proyeksi berdasarkan posisi + rotasi
func (v *Viewer) UpdateMatrices() {
	// hitung posisi mata dari rotasi + jarak
	v.Camera.Position = Vector3{
		X: v.Distance * math.Sin(v.RotY) * math.Cos(v.RotX),
		Y: v.Distance * math.Sin(v.RotX),
		Z: v.Distance * math.Cos(v.RotY) * math.Cos(v.RotX),
	}

	// buat view matrix dari posisi mata
	v.ViewMatrix = LookAt(v.Camera.Position, v.Camera.Target, v.Camera.Up)

	// buat proyeksi matrix
	if v.UseOrthogonal {
		scale := v.Distance
		v.ProjMatrix = Orthographic(
			-float64(v.ScreenWidth)/200*scale,
			float64(v.ScreenWidth)/200*scale,
			-float64(v.ScreenHeight)/200*scale,
			float64(v.ScreenHeight)/200*scale,
			0.1, 1000)
	} else {
		aspect := float64(v.ScreenWidth) / float64(v.ScreenHeight)
		v.ProjMatrix = Perspective(v.Camera.FOV, aspect, 0.1, 1000)
	}
}

// ubah koordinat 3D jadi 2D buat layar
func (v *Viewer) Project(point Vector3) Vector2 {
	// transformasi ke view space (relative ke mata)
	viewPos := MatrixVecMul(v.ViewMatrix, point)

	// transformasi ke projection space
	projPos := MatrixVecMul(v.ProjMatrix, viewPos)

	// transformasi ke pixel di layar
	screenX := (projPos.X + 1) * float64(v.ScreenWidth) / 2
	screenY := (1 - projPos.Y) * float64(v.ScreenHeight) / 2

	return Vector2{screenX, screenY}
}

// gambar garis dari p1 ke p2
func (v *Viewer) DrawLine(p1, p2 Vector2, img *ebiten.Image) {
	ebitenutil.DrawLine(img, p1.X, p1.Y, p2.X, p2.Y, color.White)
}

// gambar satu voxel sebagai kerangka kubus
func (v *Viewer) DrawVoxel(voxel Voxel, img *ebiten.Image) {
	// hitung 8 titik sudut kubus
	h := voxel.HalfSize
	center := voxel.Center

	corners := [8]Vector3{
		{center.X - h, center.Y - h, center.Z - h},
		{center.X + h, center.Y - h, center.Z - h},
		{center.X + h, center.Y + h, center.Z - h},
		{center.X - h, center.Y + h, center.Z - h},
		{center.X - h, center.Y - h, center.Z + h},
		{center.X + h, center.Y - h, center.Z + h},
		{center.X + h, center.Y + h, center.Z + h},
		{center.X - h, center.Y + h, center.Z + h},
	}

	// proyeksi ke 2D layar
	projected := [8]Vector2{}
	for i, corner := range corners {
		projected[i] = v.Project(corner)
	}

	// gambar 12 garis tepi kubus
	edges := [12][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0}, // muka belakang
		{4, 5}, {5, 6}, {6, 7}, {7, 4}, // muka depan
		{0, 4}, {1, 5}, {2, 6}, {3, 7}, // garis vertikal
	}

	for _, edge := range edges {
		v.DrawLine(projected[edge[0]], projected[edge[1]], img)
	}
}

// handle input dari user, update posisi kamera
func (v *Viewer) Update() error {
	// rotasi pake arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		v.RotY -= 0.03
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		v.RotY += 0.03
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		v.RotX = math.Min(v.RotX+0.03, math.Pi/2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		v.RotX = math.Max(v.RotX-0.03, -math.Pi/2)
	}

	// zoom pake W dan S
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if v.Distance > 2 {
			v.Distance -= 0.5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if v.Distance < 100 {
			v.Distance += 0.5
		}
	}

	// P = ganti mode proyeksi
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		v.UseOrthogonal = !v.UseOrthogonal
	}

	// R = reset semua
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		v.RotX = 0
		v.RotY = 0
		v.Distance = 15
		v.Camera.FOV = 45
	}

	v.UpdateMatrices()
	return nil
}

// gambar semua voxel ke layar
func (v *Viewer) Draw(screen *ebiten.Image) {
	screen.Clear()

	// gambar setiap voxel
	for _, voxel := range v.Voxels {
		v.DrawVoxel(voxel, screen)
	}

	// tampilkan info kontrol
	ebitenutil.DebugPrint(screen,
		"← → ↑ ↓: putar | W/S: zoom | P: ganti mode | R: reset\n"+
			fmt.Sprintf("Rotasi: X=%.2f° Y=%.2f° | Jarak: %.1f\n",
				v.RotX*180/math.Pi, v.RotY*180/math.Pi, v.Distance)+
			fmt.Sprintf("Total voxel: %d\n", len(v.Voxels)),
	)
}

// return ukuran layar yang dipake
func (v *Viewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return v.ScreenWidth, v.ScreenHeight
}

// run viewer interaktif
func RunViewer(voxels []Voxel) error {
	screenWidth := 1024
	screenHeight := 768

	viewer := NewViewer(voxels, screenWidth, screenHeight)

	ebiten.SetWindowTitle("3D Voxel Viewer - Putar/Zoom/Ganti Mode/Reset")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(viewer)
}
