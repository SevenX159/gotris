package main

import (
	col "image/color" // Needed to create the Color Background (Light gray)
	"math/rand"
	"strconv" // Needed to convert Score var from INT to STR
	"time"

	"github.com/algosup/game"
	"github.com/algosup/game/color"
	"github.com/algosup/game/key"
)

type shape [4][4]uint8

var x = columns / 2
var y = 0
var screen [rows][columns]uint8

const columns = 12
const rows = 24

var colorBackground = col.RGBA{15, 15, 15, 255} // Create a New Color to display Background Lines (Light gray)

var score = 0
var frame = 0
var ignoreKey = 0
var currentShape shape
var drop = false
var isOver = false
var ignoreSpaceKey = false

var shapes = []shape{
	{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{0, 2, 2, 0},
		{0, 2, 2, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 3, 0},
		{0, 3, 3, 0},
		{0, 3, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 4, 0, 0},
		{0, 4, 4, 0},
		{0, 4, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 5, 0, 0},
		{0, 5, 0, 0},
		{0, 5, 5, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 6, 0},
		{0, 0, 6, 0},
		{0, 6, 6, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 7, 0, 0},
		{0, 7, 7, 0},
		{0, 0, 7, 0},
		{0, 0, 0, 0},
	},
}

var colors = []color.Color{
	color.Black,
	color.LightBlue,
	color.Yellow,
	color.Red,
	color.Purple,
	color.Orange,
	color.Blue,
	color.Green,
}

func drawShape(surface game.Surface, s shape, col int, row int) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			color := s[r][c]
			if color != 0 {
				game.DrawRect(surface, (col+c)*20, (row+r)*20, 20-1, 20-1, colors[color]) // 20-1 to let 1px space between each square of each shapes
			}
		}
	}
}

func isPositionValid(s shape, col int, row int) bool {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			if s[r][c] != 0 {
				if c+col < 0 {
					return false
				}
				if c+col >= columns {
					return false
				}
				if r+row >= rows {
					return false
				}
				if screen[row+r][col+c] != 0 {
					return false
				}
			}
		}
	}

	return true
}

func isRowFull(row int) bool {
	for c := 0; c < columns; c++ {
		if screen[row][c] == 0 {
			return false
		}
	}

	return true
}

func clearRow(row int) {
	for r := row; r > 0; r-- {
		for c := 0; c < columns; c++ {
			screen[r][c] = screen[r-1][c]
		}
	}

	for c := 0; c < columns; c++ {
		screen[0][c] = 0
	}
}

func clearFullRows() {
	for r := 0; r < rows; r++ {
		if isRowFull(r) {
			clearRow(r)
			score++
		}
	}
}

func shapeToScreen(s shape, col int, row int) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			color := s[r][c]
			if color != 0 {
				screen[row+r][col+c] = color
			}
		}
	}
}

func rotate(s shape) shape {
	var newShape shape
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			newShape[r][c] = s[3-c][r]
		}
	}

	return newShape
}

func pickNewShape() {
	frame = 0
	ignoreKey = 0
	y = 0
	x = columns / 2
	currentShape = shapes[rand.Intn(len(shapes))]
	if isPositionValid(currentShape, x, y) == false {
		isOver = true
	}
}

func update() {
	if isOver {
		return
	}

	frame++

	if frame == 30 || drop {
		frame = 0
		if isPositionValid(currentShape, x, y+1) == false {
			shapeToScreen(currentShape, x, y)
			pickNewShape()
			drop = false
			clearFullRows()
		} else {
			y++
		}
		return
	}

	ignoreKey--
	if ignoreKey <= 0 {
		if key.IsPressed(key.Right) {
			if isPositionValid(currentShape, x+1, y) {
				x++
			}
			ignoreKey = 10
		}
		if key.IsPressed(key.Left) {
			if isPositionValid(currentShape, x-1, y) {
				x--
			}
			ignoreKey = 10
		}

		if key.IsPressed(key.Down) {
			s := rotate(currentShape)

			if isPositionValid(s, x, y) {
				currentShape = s
			}
			ignoreKey = 10
		}

		if key.IsPressed(key.Space) {
			if ignoreSpaceKey == false {
				drop = true
			}
			ignoreSpaceKey = true
		} else {
			ignoreSpaceKey = false
		}
	}

}

func draw(surface game.Surface) {
	update()
	for c := 0; c < columns; c++ {
		for r := 0; r < rows; r++ {
			color := screen[r][c]
			if color != 0 {
				game.DrawRect(surface, c*20, r*20, 20-1, 20-1, colors[color]) // 20-1 to let 1px space between each square of each shapes

			}
		}
	}

	if isOver {
		game.DrawText(surface, "GAME OVER", 80, 40)
		t := strconv.Itoa(score)           // Convert Int to Str
		game.DrawText(surface, t, 105, 60) // Display of Score at the end of game           105 for Centered in X axis   60 to be displayed just under "Game Over"

		for c := 0; c < columns; c++ { // Display the Cols and Rows at the end of the Game
			for r := 0; r < rows; r++ {
				game.DrawRect(surface, c*20-1, r*20-1, rows*20, 1, colorBackground)
				game.DrawRect(surface, c*20-1, r*20-1, 1, columns*20, colorBackground)

			}
		}

	} else {
		drawShape(surface, currentShape, x, y)

		for c := 0; c < columns; c++ { // Display the Cols and Rows while in game
			for r := 0; r < rows; r++ {
				game.DrawRect(surface, c*20-1, r*20-1, rows*20, 1, colorBackground)
				game.DrawRect(surface, c*20-1, r*20-1, 1, columns*20, colorBackground)

			}
		}

		t := strconv.Itoa(score)           // Convert Int to Str
		game.DrawText(surface, t, 105, 40) // Display of Score
	}
}

func main() {
	// Seed random generator with current time so that each game is a different one
	rand.Seed(time.Now().UnixNano())
	pickNewShape()
	game.Run("GOTRIS", columns*20, rows*20, draw)
}
