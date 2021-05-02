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
	bestPopulationPlayersToNextGeneration = 1
)

var MainPopulation *population

type population struct {
	generation int
	currentPopulation []*player.Player
}

func Populate(camera camera.Camera, oldGenerationPopulation []common.Object) []*player.Player {
	if MainPopulation == nil {
		MainPopulation = &population{}
		return MainPopulation.generatePopulationGeneration1(camera)
	}
	return MainPopulation.generatePopulationNextGeneration(camera, oldGenerationPopulation)
}

func (p *population) GetGeneration() int {
	return p.generation
}

func (p *population) generatePopulationGeneration1(camera camera.Camera) []*player.Player {
	var newPopulation []*player.Player
	for i := 0; i < populationSize; i++ {
		newPopulation = append(newPopulation, player.New(camera, brain.NewBrain()))
	}

	p.generation = 1
	p.currentPopulation = newPopulation
	return newPopulation
}

func (p *population) generatePopulationNextGeneration(camera camera.Camera, oldGenerationPopulation []common.Object) []*player.Player {
	oldGenerationArrangedPopulation := arrangeOldGenerationPopulationByFitness(oldGenerationPopulation)
	brainsToReproduce := getBrainsToReproduce(oldGenerationArrangedPopulation)
	nextGenerationPopulation := reproduceByBrains(camera, brainsToReproduce)

	p.currentPopulation = nextGenerationPopulation
	p.generation++
	return nextGenerationPopulation
}

func arrangeOldGenerationPopulationByFitness(oldGenerationPopulation []common.Object) []*player.Player {
	var oldGenerationArrangedPopulation []*player.Player

	for _, oldGenerationPlayerObject := range oldGenerationPopulation {
		oldGenerationArrangedPopulation = append(oldGenerationArrangedPopulation, oldGenerationPlayerObject.(*player.Player))
	}

	sort.Sort(fitnessArranger.ByFitness(oldGenerationArrangedPopulation))
	return oldGenerationArrangedPopulation
}

func getBrainsToReproduce(oldGenerationArrangedPopulation []*player.Player) []*brain.Brain {
	var brainsToKeep []*brain.Brain
	for i := 0; i < bestPopulationPlayersToNextGeneration; i++ {
		player := oldGenerationArrangedPopulation[i]
		brainsToKeep = append(brainsToKeep, player.Brain)
	}
	fmt.Printf("Best fitness for generation: %v\n", oldGenerationArrangedPopulation[0].Fitness)
	return brainsToKeep
}

func reproduceByBrains(camera camera.Camera, brainsToKeep []*brain.Brain) []*player.Player {
	var nextGenerationPopulation []*player.Player
	for i := 0; i < populationSize; i++ {
		brainToPopulate := brainsToKeep[i % bestPopulationPlayersToNextGeneration]
		nextGenerationPopulation = append(nextGenerationPopulation, player.New(camera, brainToPopulate.CloneBrain()))
	}
	return nextGenerationPopulation
}
