package fitnessArranger

import (
	"github.com/omer-dayan/flappyai/pkg/player"
)
type ByFitness []*player.Player
func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Less(i, j int) bool { return a[i].Fitness > a[j].Fitness }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }