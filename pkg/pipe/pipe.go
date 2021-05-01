package pipe

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/omer-dayan/flappyai/pkg/camera"
	"github.com/omer-dayan/flappyai/pkg/common"
	"image/color"
	_ "image/png"
	"math/rand"
)

const (
	pipesSpace = 160
	pipeWidth = 80
	pipeEdgeImageHeight = 37
	pipeAndEdgeImageVisualDiff = 2
	pipeImageHeight = 640
)

var (
	pseudoColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	edgeImage *ebiten.Image
	image *ebiten.Image
)

type Pipe struct {
	id           int
	isBottomPipe bool
	X            int
	EdgeY        int
	camera       camera.Camera
}

func New(id, x int, camera camera.Camera) []*Pipe {
	edgeY := rand.Intn(common.ScreenHeight * 0.6) + common.ScreenHeight * 0.1
	if edgeImage == nil {
		img, _, err := ebitenutil.NewImageFromFile("rsc/pipeEdge.png")
		if err != nil {
			panic(fmt.Sprintf("Could not load image of pipe edge due to: %v\n", err))
		}
		edgeImage = img
	}
	if image == nil {
		img, _, err := ebitenutil.NewImageFromFile("rsc/pipe.png")
		if err != nil {
			panic(fmt.Sprintf("Could not load image of pipe due to: %v\n", err))
		}
		image = img
	}
	return []*Pipe{
		{
			id:           id,
			isBottomPipe: false,
			X:            x,
			EdgeY:        edgeY,
			camera:       camera,
		}, {
			id:           id + 1,
			isBottomPipe: true,
			X:            x,
			EdgeY:        edgeY + pipesSpace,
			camera:       camera,
		} ,
	}
}

func (p *Pipe) Step(state *common.State) error {
	if p.isOverScreen() {
		state.RemoveNextIndexObject()
	}
	return nil
}

func (p *Pipe) Draw(screen *ebiten.Image) error {
	physicalPositionX := p.camera.GetX(p.X)
	if p.isOnScreen() {
		if p.isBottomPipe {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(physicalPositionX + pipeAndEdgeImageVisualDiff, float64(p.EdgeY))
			screen.DrawImage(image, opt)
			edgeOpt := &ebiten.DrawImageOptions{}
			edgeOpt.GeoM.Translate(physicalPositionX, float64(p.EdgeY))
			screen.DrawImage(edgeImage, edgeOpt)
		} else {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(physicalPositionX + pipeAndEdgeImageVisualDiff, float64(p.EdgeY- pipeImageHeight))
			screen.DrawImage(image, opt)
			edgeOpt := &ebiten.DrawImageOptions{}
			edgeOpt.GeoM.Translate(physicalPositionX, float64(p.EdgeY- pipeEdgeImageHeight))
			screen.DrawImage(edgeImage, edgeOpt)
		}
	}
	return nil
}

func (p *Pipe) GetPositionAndSize() (int, int, int, int) {
	if p.isBottomPipe {
		return p.X, p.EdgeY, pipeWidth, common.ScreenHeight - p.EdgeY
	} else {
		return p.X, 0, pipeWidth, p.EdgeY
	}
}

func (p *Pipe) IsBottomPipe() bool {
	return p.isBottomPipe
}

func (p *Pipe) GetPipeWidth() int {
	return pipeWidth
}

func (p *Pipe) isOnScreen() bool {
	physicalPositionX := p.camera.GetX(p.X)
	return physicalPositionX > -pipeWidth && physicalPositionX < common.ScreenWidth
}
func (p *Pipe) isOverScreen() bool {
	physicalPositionX := p.camera.GetX(p.X)
	return physicalPositionX < -pipeWidth
}
