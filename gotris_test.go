package main

import "testing"

var s1 = shape{
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}

var s1r = shape{
	{0, 0, 1, 0},
	{0, 0, 1, 0},
	{0, 0, 1, 0},
	{0, 0, 1, 0},
}

func TestRotateR(t *testing.T) {
	if rotate(s1) != s1r {
		t.Error("Failed to rotate clockwise")
	}
}

func TestRotate360(t *testing.T) {
	if rotate(rotate(rotate(rotate(s1)))) != s1 {
		t.Error("4 rotations should return to original shape")
	}
}
