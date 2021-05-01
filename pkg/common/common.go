package common

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

type FileImageSprite struct {
	x int
	y int
	width int
	height int
}

func NewFileImageSprite(x, y, width, height int, imageSrc string) *FileImageSprite {
	if backgroundImage == nil {
		img, _, err := ebitenutil.NewImageFromFile(imageSrc)
		if err != nil {
			panic(fmt.Sprintf("Could not load image of background due to: %v\n", err))
		}
		backgroundImage = img
	}
	return &FileImageSprite{
		x: x,
		y: y,
		width: width,
		height: height,
	}
}

func (c *FileImageSprite) Draw(screen *ebiten.Image) error {
	edgeOpt := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, edgeOpt)
	return nil
}
