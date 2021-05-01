package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/omer-dayan/flappyai/pkg/camera"
	"github.com/omer-dayan/flappyai/pkg/common"
	"github.com/omer-dayan/flappyai/pkg/neat/population"
	"github.com/omer-dayan/flappyai/pkg/pipe"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
)

const (
	pipeSpace           = 300
	pipePlacementChange = -150
)

var (
	mainGame  *Flappy
	mPlusFont font.Face
)

type Flappy struct {
	camera camera.Camera
	state *common.State
}

func NewFlappy() *Flappy {
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("Flappy AI")
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	mPlusFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	flappy := &Flappy{}
	camera := camera.NewSide2DCamera()
	state := common.NewState(onGameEnd)
	initState(camera, state)
	flappy.camera = camera
	flappy.state = state
	mainGame = flappy
	return flappy
}

func initState(camera camera.Camera, state *common.State) {
	oldGenerationPlayers := state.GetAllPlayers()
	state.Init()
	state.BackgroundSprites = append(state.BackgroundSprites, common.NewFileImageSprite(0, 0, common.ScreenWidth, common.ScreenHeight, "rsc/background.png"))
	initPipes(camera, state)
	population := population.Populate(camera, oldGenerationPlayers)
	for _, player := range population {
		state.InsertPlayer(player)
	}
}

func initPipes(camera camera.Camera, state *common.State) {
	for i := 2; i < 200; i++ {
		newPipes := pipe.New(i*2, i * pipeSpace +pipePlacementChange, camera)
		state.IndexObject(newPipes[0])
		state.IndexObject(newPipes[1])
	}
}

func (f Flappy) Update() error {
	f.camera.Step()
	return f.state.Step()
}

func (f Flappy) Draw(screen *ebiten.Image) {
	f.camera.Draw(screen, f.state)
	text.Draw(screen, fmt.Sprintf("Generation: %v", population.MainPopulation.GetGeneration()), mPlusFont, 10, 50, color.White)
	text.Draw(screen, fmt.Sprintf("Players Alive: %v", f.state.GetPlayerAliveCount()), mPlusFont, 10, 85, color.White)
}

func (f Flappy) Layout(_, _ int) (screenWidth, screenHeight int) {
	return common.ScreenWidth, common.ScreenHeight
}

func reRunGame() error {
	initState(mainGame.camera, mainGame.state)
	mainGame.camera.Reset()
	return nil
}

func onGameEnd(_ *common.State) {
	reRunGame()
}
