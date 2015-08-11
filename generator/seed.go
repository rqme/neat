/*
Copyright (c) 2015, Brian Hummer (brian@redq.me)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package generator

import (
	"math/rand"

	"github.com/rqme/neat"
)

// Returns a genome build from the parameters
func createSeed(ctx neat.Context, cfg ClassicSettings) (adam neat.Genome) {
	// Create the genome
	inputs := cfg.NumInputs()
	outputs := cfg.NumOutputs()
	adam = neat.Genome{
		Nodes: make(map[int]neat.Node, 1+inputs+outputs),
	}
	nodes := make([]neat.Node, len(adam.Nodes))
	node := neat.Node{NeuronType: neat.Bias, ActivationType: neat.Direct, X: 0, Y: 0}
	node.Innovation = ctx.Innovation(neat.NodeInnovation, node.Key())
	adam.Nodes[node.Innovation] = node
	nodes = append(nodes, node)
	for i := 0; i < inputs; i++ {
		node = neat.Node{NeuronType: neat.Input, ActivationType: neat.Direct, X: float64(i+1) / float64(inputs), Y: 0}
		node.Innovation = ctx.Innovation(neat.NodeInnovation, node.Key())
		adam.Nodes[node.Innovation] = node
		nodes = append(nodes, node)
	}
	x := 0.5
	for i := 0; i < outputs; i++ {
		if outputs > 1 {
			x = float64(i) / float64(outputs-1)
		}
		node = neat.Node{NeuronType: neat.Output, ActivationType: cfg.OutputActivation(), X: x, Y: 1}
		node.Innovation = ctx.Innovation(neat.NodeInnovation, node.Key())
		adam.Nodes[node.Innovation] = node
		nodes = append(nodes, node)
	}

	rng := rand.New(rand.NewSource(rand.Int63()))
	adam.Conns = make(map[int]neat.Connection, (1+inputs)*outputs)
	for i := 0; i < 1+inputs; i++ {
		for j := 0; j < outputs; j++ {
			w := (rng.Float64()*2.0 - 1.0) * cfg.WeightRange()
			conn := neat.Connection{Source: nodes[i].Innovation, Target: nodes[j+1+inputs].Innovation, Enabled: true, Weight: w}
			conn.Innovation = ctx.Innovation(neat.ConnInnovation, conn.Key())
			adam.Conns[conn.Innovation] = conn
		}
	}

	ts := cfg.Traits()
	adam.Traits = make([]float64, len(ts))
	for i, trait := range ts {
		adam.Traits[i] = rng.Float64()*(trait.Max-trait.Min) + trait.Min // TODO: Get setting values from configuration
	}
	return adam
}
