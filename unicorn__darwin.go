package main

import "fmt"

// Unicorn contains the width and height of the Unicorn pHAT
type Unicorn struct {
	Width  int
	Height int
}

// Init initializes the pHAT for writing to.
func (u *Unicorn) Init() error {
	fmt.Println("OS X")
	return nil
}

// SetPixel sets the color of the specified pixel in the display buffer.
func (u *Unicorn) SetPixel(x, y int, red, green, blue uint8) {
	fmt.Printf("SetPixel(X:%d, Y:%d, R:%d, G:%d, B:%d)\n", x, y, red, green, blue)
}

// Show writes out the contents of the display buffer to the pHAT.
func (u *Unicorn) Show() error {
	fmt.Println("Show")
	return nil
}

// Clear turns off all the pixels on the pHAT and erases the display buffer.
func (u *Unicorn) Clear() {
	fmt.Println("Clear")
}

// CleanUp clears the display and cleans up the GPIO.
func (u *Unicorn) CleanUp() {
	fmt.Println("CleanUp")
}
