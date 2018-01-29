package main

import "ws2811"

const (
	pin = 18
)

type Unicorn struct {
	Width  int
	Height int
}

func (u *Unicorn) Init() error {
	return ws2811.Init(pin, u.Height*u.Width, 255)
}

func (u *Unicorn) SetPixel(x, y int, red, green, blue uint8) {
	i := (y * u.Width) + x
	var color uint32
	color = uint32(green)<<16 | uint32(red)<<8 | uint32(blue)<<0
	ws2811.SetLed(i, color)
}

func (u *Unicorn) Show() error {
	return ws2811.Render()
}

func (u *Unicorn) Clear() {
	ws2811.Clear()
}

func (u *Unicorn) CleanUp() {
	ws2811.Fini()
}
