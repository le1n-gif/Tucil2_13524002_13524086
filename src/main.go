package main

import (
	"fmt"
)

func main() {
	var pathFile string
	fmt.Scanln(&pathFile)
	var alltriangle []Triangle
	var allVertices []Vector3
	var err error
	alltriangle, allVertices, err = ReadOBJ(pathFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// -=-=-=---=-=-test/debug=-=-=-=-=-=--=-=-=-=
	for i := range alltriangle {
		fmt.Printf("triangle[%d] vertex[1]nya : %f-%f-%f\n", i, alltriangle[i].V1.X, alltriangle[i].V1.Y, alltriangle[i].V1.Z)
		fmt.Printf("triangle[%d] vertex[2]nya : %f-%f-%f\n", i, alltriangle[i].V2.X, alltriangle[i].V2.Y, alltriangle[i].V2.Z)
		fmt.Printf("triangle[%d] vertex[3]nya : %f-%f-%f\n", i, alltriangle[i].V3.X, alltriangle[i].V3.Y, alltriangle[i].V3.Z)
	}
	for i := range allVertices {
		fmt.Printf("vertex[%d] koornya: %f-%f-%f\n", i, allVertices[i].X, allVertices[i].Y, allVertices[i].Z)
	}

}
