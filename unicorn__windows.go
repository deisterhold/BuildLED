package main

import "fmt"

type Unicorn struct {
	Width  int
	Height int
}

func (u *Unicorn) Init() error {
	return nil
}

func (u *Unicorn) SetPixel(x, y int, red, green, blue uint8) {
	fmt.Printf("SetPixel(X:%d, Y:%d, R:%d, G:%d, B:%d)\n", x, y, red, green, blue)
}

func (u *Unicorn) Show() error {
	fmt.Println("Show")
	return nil
}

func (u *Unicorn) Clear() {
	fmt.Println("Clear")
}

func (u *Unicorn) CleanUp() {
	fmt.Println("CleanUp")
}
