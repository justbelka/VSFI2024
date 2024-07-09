package main

import (
	"fmt"
	"time"
)

var letters = map[rune][]string{
	'V': {
		"V       V",
		" V     V ",
		"  V   V  ",
		"   V V   ",
		"    V    ",
	},
	'S': {
		" SSSSS ",
		"S      ",
		" SSSSS ",
		"      S",
		" SSSSS ",
	},
	'F': {
		"FFFFFFF",
		"F      ",
		"FFFFFF ",
		"F      ",
		"F      ",
	},
	'I': {
		"IIIIIII",
		"   I   ",
		"   I   ",
		"   I   ",
		"IIIIIII",
	},
}

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func colorizeArt(art []string, colorOffset int) []string {
	coloredArt := make([]string, len(art))
	for i, line := range art {
		for j, char := range line {
			if char == ' ' {
				coloredArt[i] += " "
			} else {
				colorIndex := (j - colorOffset + len(colors)) % len(colors)
				coloredArt[i] += colors[colorIndex] + string(char) + "\033[0m"
			}
		}
	}
	return coloredArt
}

func printArt(art [][]string) {
	for row := 0; row < len(art[0]); row++ {
		for _, lines := range art {
			fmt.Print(lines[row] + "  ")
		}
		fmt.Println()
	}
}

func main() {
	asciiLetters := []string{"V", "S", "F", "I"}
	asciiArt := make([][]string, len(asciiLetters))
	for i, letter := range asciiLetters {
		asciiArt[i] = letters[rune(letter[0])]
	}

	colorOffset := 0

	for {
		clearScreen()
		coloredArt := make([][]string, len(asciiArt))
		for i, art := range asciiArt {
			coloredArt[i] = colorizeArt(art, colorOffset)
		}
		printArt(coloredArt)
		colorOffset = (colorOffset + 1) % len(colors)
		time.Sleep(300 * time.Millisecond)
	}
}