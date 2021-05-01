package camera

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/omer-dayan/flappyai/pkg/common"
)

type Camera interface {
	Step() error
	Reset()
	Draw(screen *ebiten.Image, state *common.State) error
	GetX(virtualX int) float64
}

type Side2DCamera struct {
	imagePositionX int
}

func NewSide2DCamera() *Side2DCamera {
	return &Side2DCamera{
		imagePositionX: 0,
	}
}

func (c *Side2DCamera) Step() error {
	c.imagePositionX = c.imagePositionX + common.GameVelocity
	return nil
}

func (c *Side2DCamera) Draw(screen *ebiten.Image, state *common.State) error {
	for _, sprite := range state.BackgroundSprites {
		if err := sprite.Draw(screen); err != nil {
			fmt.Printf("Could not draw background sprite due to: %v\n", err)
		}
	}

	for _, object := range state.Objects {
		err := object.Draw(screen)
		if err != nil {
			fmt.Printf("Could not draw object due to: %v\n", err)
		}
	}

	for _, object := range state.IndexedObject {
		err := object.Draw(screen)
		if err != nil {
			fmt.Printf("Could not draw indexed object due to: %v\n", err)
		}
	}
	return nil
}

func (c *Side2DCamera) GetX(virtualX int) float64 {
	return float64(virtualX - c.imagePositionX)
}

func (c *Side2DCamera) Reset() {
	c.imagePositionX = 0
}