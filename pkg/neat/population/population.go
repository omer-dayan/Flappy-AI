package population

import (
	"fmt"
	"github.com/omer-dayan/flappyai/pkg/camera"
	"github.com/omer-dayan/flappyai/pkg/common"
	"github.com/omer-dayan/flappyai/pkg/neat/brain"
	"github.com/omer-dayan/flappyai/pkg/neat/fitnessArranger"
	"github.com/omer-dayan/flappyai/pkg/player"
	"sort"
)

const (
	populationSize = 100
	bestPopulationPlayersToNextGeneration = 2
)

var MainPopulation *population

type population struct {
	generation int
	currentPopulation []*player.Player
}

func (p *population) GetGeneration() int {
	return p.generation
}

func Populate(camera camera.Camera, oldGenerationPopulation []common.Object) []*player.Player {
	if MainPopulation == nil {
		return generatePopulationGeneration1(camera)
	}
	MainPopulation.generation++
	return generatePopulationNextGeneration(camera, oldGenerationPopulation)
}

func generatePopulationGeneration1(camera camera.Camera) []*player.Player {
	newPopulation := []*player.Player{}
	for i := 0; i < populationSize; i++ {
		newPopulation = append(newPopulation, player.New(camera, brain.NewBrain()))
	}
	MainPopulation = &population{
		generation: 1,
		currentPopulation: newPopulation,
	}
	return newPopulation
}

func generatePopulationNextGeneration(camera camera.Camera, oldGenerationPopulation []common.Object) []*player.Player {
	oldGenerationArrangedPopulation := []*player.Player{}

	for _, oldGenerationPlayerObject := range oldGenerationPopulation {
		oldGenerationArrangedPopulation = append(oldGenerationArrangedPopulation, oldGenerationPlayerObject.(*player.Player))
	}

	sort.Sort(fitnessArranger.ByFitness(oldGenerationArrangedPopulation))

	brainsToKeep := []*brain.Brain{}
	for i := 0; i < bestPopulationPlayersToNextGeneration; i++ {
		player := oldGenerationArrangedPopulation[i]
		brainsToKeep = append(brainsToKeep, player.Brain)
	}
	fmt.Printf("Best fitness for generation: %v\n", oldGenerationArrangedPopulation[0].Fitness)

	nextGenerationPopulation := []*player.Player{}
	for i := 0; i < populationSize; i++ {
		brainToPopulate := brainsToKeep[i % bestPopulationPlayersToNextGeneration]
		nextGenerationPopulation = append(nextGenerationPopulation, player.New(camera, brainToPopulate.CloneBrain()))
	}
	MainPopulation.currentPopulation = nextGenerationPopulation
	return nextGenerationPopulation
}
