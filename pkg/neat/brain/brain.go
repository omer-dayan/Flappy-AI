package brain

import (
	"math"
	"math/rand"
	"time"
)

const (
	mutationFactor = 0.2
	mutationRate = 0.3
	randomFirstWeightRange = 0.15
)

var (
	randomSeed = rand.NewSource(time.Now().UnixNano())
	random = rand.New(randomSeed)
)

type Brain struct {
	inputNodes []*inputNode
	outputNode *outputNode
}

type inputNode struct {
	name string
	connection *connection
}

func (i *inputNode) clone(node *outputNode) *inputNode {
	connectionClone := i.connection.clone(node)
	return &inputNode{
		name:       i.name,
		connection: connectionClone,
	}
}

func (i *inputNode) activate(value int) {
	i.connection.outputNode.value += int(i.connection.weight * float64(value))
}

type connection struct {
	weight float64
	outputNode *outputNode
}

func (c *connection) clone(node *outputNode) *connection {
	oldClone := &connection{
		weight: c.weight,
		outputNode: node,
	}
	if random.Intn(100) <= 100 * mutationRate {
		w1 := (random.Float64() * (2 *mutationFactor)) + (1 - mutationFactor)
		result := oldClone.weight * w1
		oldClone.weight = result
	}
	return oldClone
}

type outputNode struct {
	value int
}

func (o *outputNode) clone() *outputNode {
	return &outputNode{
		value: o.value,
	}
}

func NewBrain() *Brain {
	output := &outputNode{
		value: 0,
	}
	distanceYNode := newInputNode("Distance From Ground", output)
	distanceXTopPipeNode := newInputNode("Distance X From Top Pipe", output)
	distanceYTopPipeNode := newInputNode("Distance Y From Top Pipe", output)
	distanceXBottomPipeNode := newInputNode("Distance X From Bottom Pipe", output)
	distanceYBottomPipeNode := newInputNode("Distance Y From Bottom Pipe", output)
	inputs := []*inputNode{distanceYNode, distanceXTopPipeNode, distanceYTopPipeNode, distanceXBottomPipeNode, distanceYBottomPipeNode}
	return &Brain{
		inputNodes: inputs,
		outputNode: output,
	}
}

func newInputNode(name string, node *outputNode) *inputNode {
	w := rand.Float64() * (2 * randomFirstWeightRange) - randomFirstWeightRange
	connection := &connection{
		weight: w,
		outputNode: node,
	}
	return &inputNode{
		name: name,
		connection: connection,
	}
}

func (b *Brain) CloneBrain() *Brain {
	outputNodeClone := b.outputNode.clone()
	inputNodeClones := []*inputNode{}
	for _, input := range b.inputNodes {
		inputNodeClones = append(inputNodeClones,  input.clone(outputNodeClone))
	}
	return &Brain{
		inputNodes: inputNodeClones,
		outputNode: outputNodeClone,
	}
}

func (b *Brain) flashOutput() {
	b.outputNode.value = 0
}

func (b *Brain) getBooleanOutput() bool {
	return math.Tanh(float64(b.outputNode.value)) < -0.5
}

func (b *Brain) ShouldIJump(distanceFromGround, distanceXFromUpperPipe, distanceYFromUpperPipe, distanceXFromBottomPipe, distanceYFromBottomPipe int) bool {
	b.flashOutput()

	b.inputNodes[0].activate(distanceFromGround)
	b.inputNodes[1].activate(distanceXFromUpperPipe)
	b.inputNodes[2].activate(distanceYFromUpperPipe)
	b.inputNodes[3].activate(distanceXFromBottomPipe)
	b.inputNodes[4].activate(distanceYFromBottomPipe)

	return b.getBooleanOutput()
}
