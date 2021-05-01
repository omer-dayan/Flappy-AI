package player

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/omer-dayan/flappyai/pkg/camera"
	"github.com/omer-dayan/flappyai/pkg/common"
	"github.com/omer-dayan/flappyai/pkg/neat/brain"
	"github.com/omer-dayan/flappyai/pkg/pipe"
	"image/color"
	"math"
)

const (
	g = 0.35
	jumpVerticalVelocity = -7
	deathVerticalVelocity = -9
	birdWidth = 60
	birdHeight = 50
)

var (
	pseudoColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	image *ebiten.Image
)

type Player struct {
	x float64
	y float64
	width float64
	height float64

	verticalVelocity float64
	horizontalVelocity float64

	isFallen bool
	isDead   bool
	Fitness  uint64

	Brain  *brain.Brain
	camera camera.Camera
}

func New(camera camera.Camera, brain *brain.Brain) *Player {
	if image == nil {
		img, _, err := ebitenutil.NewImageFromFile("rsc/bird.png")
		if err != nil {
			panic(fmt.Sprintf("Could not load image of bird due to: %v\n", err))
		}
		image = img
	}
	return &Player{
		x:                  common.ScreenWidth * 0.4,
		y:                  common.ScreenHeight / 2,
		width:              birdWidth,
		height:             birdHeight,
		verticalVelocity:   0,
		horizontalVelocity: common.GameVelocity,
		isFallen:           false,
		isDead:             false,
		camera:             camera,
		Fitness: 0,
		Brain:              brain,
	}
}

func (p *Player) Step(state *common.State) error {
	if !p.isDead {
		p.handleInput(state)
		for _, touchableObjects := range state.IndexedObject[0:8] {
			if p.isCollided(touchableObjects) {
				p.die(state)
			}
		}
		if p.isOutScreenY() {
			p.die(state)
		}
		p.Fitness++
	}
	p.verticalVelocity = p.verticalVelocity + g
	p.isFallen = p.verticalVelocity > 0
	p.changePosition()
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) error {
	edgeOpt := &ebiten.DrawImageOptions{}
	edgeOpt.GeoM.Translate(p.camera.GetX(int(p.x)), p.y)
	screen.DrawImage(image, edgeOpt)
	return nil
}

func (p *Player) GetPositionAndSize() (int, int, int, int) {
	return int(p.x), int(p.y), int(p.width), int(p.height)
}

func (p *Player) handleInput(state *common.State) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.jump()
	}

	distanceFromGround := int(p.y)
	nextBottomPipe, nextTopPipe := p.getNextPipes(state)
	if p.Brain.ShouldIJump(distanceFromGround, nextTopPipe.X - int(p.x), nextTopPipe.EdgeY - int(p.y), nextBottomPipe.X - int(p.x), nextBottomPipe.EdgeY - int(p.y)) {
		p.jump()
	}
}

func (p *Player) getNextPipes(state *common.State) (*pipe.Pipe, *pipe.Pipe) {
	var nextBottomPipe *pipe.Pipe
	var nextTopPipe *pipe.Pipe
	for i := 0; i < 8; i += 2 {
		maybeNextTopPipe := state.IndexedObject[i].(*pipe.Pipe)
		maybeNextBottomPipe := state.IndexedObject[i + 1].(*pipe.Pipe)
		tx, _, tw, _ := maybeNextTopPipe.GetPositionAndSize()
		if int(p.x) < tx + tw {
			if nextTopPipe == nil {
				nextTopPipe = maybeNextTopPipe
			}
		}

		bx, _, bw, _ := maybeNextBottomPipe.GetPositionAndSize()
		if int(p.x) < bx + bw {
			if nextBottomPipe == nil {
				nextBottomPipe = maybeNextBottomPipe
			}
		}

		if nextBottomPipe != nil && nextTopPipe != nil {
			return nextBottomPipe, nextTopPipe
		}
	}
	return state.IndexedObject[8].(*pipe.Pipe), state.IndexedObject[9].(*pipe.Pipe)
}

func (p *Player) jump() {
	if p.isFallen {
		p.verticalVelocity = jumpVerticalVelocity
	}
}

func (p *Player) changePosition() {
	p.y = p.y + p.verticalVelocity
	p.x = p.x + p.horizontalVelocity
}

func (p *Player) isCollided(object common.Object) bool {
	x, y, w, h := object.GetPositionAndSize()
	mx, my, mw, mh := p.GetPositionAndSize()
	return (my + mh >= y && my <= y+h) && (mx + mw >= x && mx < x + w)
}

func (p *Player) isOutScreenY() bool {
	return p.y + p.height < 0 || p.y > common.ScreenHeight
}

func (p *Player) die(state *common.State) {
	p.horizontalVelocity = 0
	p.verticalVelocity = deathVerticalVelocity
	p.isDead = true
	state.OnPlayerDeath()
}

func calculateDistance(x1, y1, x2, y2 int) int {
	xd := math.Pow(float64(x2 - x1), 2)
	yd := math.Pow(float64(y2 - y1), 2)
	return int(math.Sqrt(xd + yd))
}
