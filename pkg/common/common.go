package common

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	ScreenWidth = 480
	ScreenHeight = 640
	GameVelocity = 2
)

var backgroundImage *ebiten.Image

type Object interface {
	Step(state *State) error
	Draw(screen *ebiten.Image) error
	GetPositionAndSize() (int, int, int ,int)
}

type Sprite interface {
	Draw(screen *ebiten.Image) error
}

type ColorSprite struct {
	x int
	y int
	width int
	height int

	color color.Color
}

func NewColorSprite(x, y, width, height int, color color.Color) *ColorSprite {
	if backgroundImage == nil {
		img, _, err := ebitenutil.NewImageFromFile("rsc/background.png")
		if err != nil {
			panic(fmt.Sprintf("Could not load image of background due to: %v\n", err))
		}
		backgroundImage = img
	}
	return &ColorSprite{
		x: x,
		y: y,
		width: width,
		height: height,
		color: color,
	}
}

func (c *ColorSprite) Draw(screen *ebiten.Image) error {
	ebitenutil.DrawRect(screen, float64(c.x), float64(c.y), float64(c.width), float64(c.height), c.color)
	edgeOpt := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, edgeOpt)
	return nil
}
