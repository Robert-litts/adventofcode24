package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	var inputFile = flag.String("inputFile", "../input/day14.input", "Relative file path to use as input.")
	flag.Parse()
	fmt.Println("Running Part 1:")
	if err := Part1(*inputFile); err != nil {
		fmt.Println("Error in Part 2:", err)
		return
	}

}

func Part1(inputFile string) error {
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	//result := 0
	var vX, vY, x, y int
	robots := Robots{}
	width := 101 //set width per instructions
	height := 103
	//Mid used for tracking safe/unsafe areas
	midWidth := math.Ceil(float64(width)/2) - 1
	midHeight := math.Ceil(float64(height)/2) - 1

	contents := string(bytes)
	safety := 0
	//Track robots found in each of the 4 quadrants
	quadTL := 0
	quadTR := 0
	quadBL := 0
	quadBR := 0
	//Track how many found "in the middle", for Part 2
	nonSafe := 0

	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	lines := strings.Split(contents, "\n")

	// Iterate through lines in blocks of 3
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i]) // Remove spaces

		// Skip empty
		if line == "" {
			continue
		}

		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &x, &y, &vX, &vY)
		robot := Robot{X: x, Y: y, vX: vX, vY: vY}
		robots.Robots = append(robots.Robots, robot)

	}

	// fmt.Println("Initial Matrix")
	// matrix := makeMatrix(width, height, &robots)
	// for _, char := range matrix {
	// 	fmt.Println(char)
	// }
	for i := 0; i < 10000; i++ {
		nonSafe = 0
		quadBL = 0
		quadBR = 0
		quadTL = 0
		quadTR = 0

		//safety := 0
		for bots := range len(robots.Robots) {
			newX := robots.Robots[bots].X + robots.Robots[bots].vX
			newY := robots.Robots[bots].Y + robots.Robots[bots].vY

			//fmt.Printf("new X: %d, new Y: %d, height: %d, width: %d \n", newX, newY, height, width)

			// Wrap the X position within the width
			newX = (newX + width) % width

			// Wrap the Y position within the height
			newY = (newY + height) % height

			// Update the robot's coordinates
			robots.Robots[bots].X = newX
			robots.Robots[bots].Y = newY
			//fmt.Printf("X: %d, Y: %d, Midwidth: %f, Midheight: %f \n", robots.Robots[bots].X, robots.Robots[bots].Y, midWidth, midHeight)
			if robots.Robots[bots].X >= 0 && robots.Robots[bots].Y >= 0 {
				if robots.Robots[bots].X < int(midWidth) && robots.Robots[bots].Y < int(midHeight) {
					quadTL++
					//fmt.Printf("TL Added, Robot %d: X=%d, Y=%d\n", bots, robots.Robots[bots].X, robots.Robots[bots].Y)

				}
				if robots.Robots[bots].X > int(midWidth) && robots.Robots[bots].Y < int(midHeight) {
					quadTR++
					//fmt.Printf("TR Added, Robot %d: X=%d, Y=%d\n", bots, robots.Robots[bots].X, robots.Robots[bots].Y)

				}
				if robots.Robots[bots].X > int(midWidth) && robots.Robots[bots].Y > int(midHeight) {
					quadBR++
					//fmt.Printf("BR, Robot %d: X=%d, Y=%d\n", bots, robots.Robots[bots].X, robots.Robots[bots].Y)

				}
				if robots.Robots[bots].X < int(midWidth) && robots.Robots[bots].Y > int(midHeight) {
					quadBL++
					//fmt.Printf("BL, Robot %d: X=%d, Y=%d\n", bots, robots.Robots[bots].X, robots.Robots[bots].Y)

				}
				if robots.Robots[bots].X == int(midWidth) || robots.Robots[bots].Y == int(midHeight) {
					nonSafe++
				}

			}

		}
		// fmt.Printf("Matrix after %d seconds \n", i)
		// matrix = makeMatrix(width, height, &robots)
		// for _, char := range matrix {
		// 	fmt.Println(char)

		// }
		safety = quadTL * quadTR * quadBL * quadBR
		if i == 99 {
			fmt.Println("Part 1 Answer: ", safety)
		}

		//Part 2, find the "Christmas Tree"
		//Tuned so that if 25 or more robots are in a non-safe quadrant, print the matrix,
		//Out of 10,500 tries, result is 2 images (one is correct)
		//Correct answer will be i+1 (since using 0 index for loop)
		if nonSafe > 25 {
			fmt.Println("High NonSafe Score: ", i+1, ",Writing to file")
			matrix := makeMatrix(width, height, &robots)
			_, err := fmt.Fprintf(file, " %d Seconds: \n", i+1)
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
			for _, char := range matrix {
				//fmt.Println(char)

				_, err := fmt.Fprintln(file, char)
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
			}

		}

	}

	fmt.Println("Safety: ", safety)
	return nil

}

func makeMatrix(width, height int, robots *Robots) [][]string {
	var matrix [][]string
	//uniqueChar := make(map[string][]Coordinate)

	//contents := string(bytes)

	for i := 0; i < height; i++ {

		// Row to hold the strings
		row := make([]string, width)

		for i := 0; i < width; i++ {

			row[i] = "."

		}
		matrix = append(matrix, row)
	}

	for bots := range len(robots.Robots) {
		x := robots.Robots[bots].X
		y := robots.Robots[bots].Y

		if matrix[y][x] != "." {
			val, err := strconv.Atoi(matrix[y][x])
			if err != nil {
				fmt.Println("Error converting string to integer:", err)
			}
			val++
			matrix[y][x] = strconv.Itoa(val)
			continue
		}

		if matrix[y][x] == "." {
			matrix[y][x] = "1"
			continue
		}
	}

	return matrix
}

// Struct to represent a visited list of coordinates. Points is a slice of Coordinate structs.
type Robots struct {
	Robots []Robot
}

// Struct to represent a coordinate on the grid. X and Y are integers.
type Robot struct {
	X  int
	Y  int
	vX int
	vY int
}
