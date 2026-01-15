package main

import (
	"log"
	"mock-jack/internal/app"

	"github.com/hajimehoshi/ebiten/v2"
)

func main () {
	// Window basic settings
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("MockJack")

	if err := ebiten.RunGame(app.New()); err != nil {
		log.Fatal(err)
	}
}