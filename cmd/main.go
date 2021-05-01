package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/omer-dayan/flappyai/pkg/game"
)

func main() {
	flappyGame := game.NewFlappy()
	if err := ebiten.RunGame(*flappyGame); err != nil {
		fmt.Printf("Erro ocured while initializing the program: %v\n", err)
	}
}
