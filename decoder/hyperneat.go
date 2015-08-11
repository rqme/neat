/*
Copyright (c) 2015 Brian Hummer (neat@boggo.net), All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted
provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this list of conditions
and the following disclaimer. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the documentation and/or other
materials provided with the distribution. Neither the name of the nor the names of its
contributors may be used to endorse or promote products derived from this software without
specific prior written permission. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package decoder

import (
	"fmt"
	"math"

	"github.com/rqme/neat"
	"github.com/rqme/neat/network"
)

// Special case: 1 layer of nodes in this case, examine nodes to separate out "vitural layers" by neuron type
// othewise connect every neuron in one layer to the subsequent layer

type HyperNEATSettings interface {
	SubstrateLayers() []SubstrateNodes // Substrate definitions
	WeightRange() float64              // Weight range for new connections
}

type HyperNEAT struct {
	HyperNEATSettings
	CppnDecoder neat.Decoder
}

// Outputs 0..len(substrate layers) = weights. if len(outputs) = 2* that number, second set is activation function, 3rd is bias connection? Need flags for these
// n = number of layers - 1
// first n = weights
// flags for bias oututs = 1 or 2 meaning use outputs starting at 1*n or 2*n
// activation, too.

func (d *HyperNEAT) Decode(g neat.Genome) (p neat.Phenome, err error) {

	// Validate the number of inputs and outputs
	if err = d.validate(g); err != nil {
		return
	}

	// Decode the CPPN
	var cppn neat.Phenome
	cppn, err = d.CppnDecoder.Decode(g)
	if err != nil {
		return nil, err
	}

	// Set the IDs of the nodes and map them
	layers := d.SubstrateLayers()
	neurons := make([]network.Neuron, 0, len(layers)*len(layers[0]))
	cnt := 0
	i := 0
	for j, l := range layers {
		// TODO: Should I sort the nodes by position in the network?
		for k, n := range l {
			l[k].id = i
			neuron := network.Neuron{NeuronType: n.NeuronType}
			neuron.X, neuron.Y = pos2D(n.Position)
			switch n.NeuronType {
			case neat.Input, neat.Bias:
				neuron.ActivationType = neat.Direct
			default:
				neuron.ActivationType = neat.Sigmoid
			}
			neurons = append(neurons, neuron)
			i += 1
		}
		layers[j] = l
		if j > 0 {
			cnt += len(layers[j]) * len(layers[j-1])
		}
	}

	// Adjust positions for better layout
	d.adjPositions(neurons)

	// Build the network from the substrates
	var outputs []float64 // output from the Cppn
	wr := d.WeightRange()
	synapses := make([]network.Synapse, 0, cnt)
	for l := 1; l < len(layers); l++ {
		for _, src := range layers[l-1] {
			for _, tgt := range layers[l] {
				outputs, err = cppn.Activate(append(src.Position, tgt.Position...))
				if err != nil {
					return nil, err
				}
				w := math.Abs(outputs[l-1])
				if w > 0.2 {
					synapse := network.Synapse{Source: src.id, Target: tgt.id, Weight: math.Copysign((w-0.2)*wr/0.8, outputs[l-1])}
					synapses = append(synapses, synapse)
				}
			}
		}
	}

	// Return the new network
	var net *network.Classic
	net, err = network.New(neurons, synapses, len(layers))
	if err != nil {
		return nil, err
	}
	p = Phenome{g, net}
	return
}

func (d *HyperNEAT) validate(g neat.Genome) error {
	var icnt, ocnt int
	for _, n := range g.Nodes {
		if n.NeuronType == neat.Input {
			icnt += 1
		} else if n.NeuronType == neat.Output {
			min, max := n.ActivationType.Range()
			found := false
			switch {
			case math.IsNaN(min), math.IsNaN(max):
				found = true
			case min >= 0:
				found = true
			}
			if found {
				return fmt.Errorf("Invalid activation type for output: %s [%f, %f]", n.ActivationType, min, max)
			}
			ocnt += 1
		}
	}

	layers := d.SubstrateLayers()
	cnt := len(layers[0][0].Position)
	for i, l := range layers {
		for j, n := range l {
			if len(n.Position) != cnt {
				return fmt.Errorf("Inconsistent position length in substrate layer %d node %d. Expected %d but found %d.", i, j, cnt, len(n.Position))
			}
		}
	}
	if cnt*2 < icnt {
		return fmt.Errorf("Insufficient number of inputs to decode substrate. Need %d but have %d", cnt*2, icnt)
	}

	if ocnt < len(layers)-1 {
		return fmt.Errorf("Insufficient number of outputs to decode substrate. Need %d but have %d", len(layers)-1, ocnt)
	}

	return nil
}

// Adjusts neuron positions so that network is layed out in layers in 2D
func (d *HyperNEAT) adjPositions(neurons []network.Neuron) {
	layers := d.SubstrateLayers()
	x := (1.0 - float64(len(layers)-1)*0.05) / float64(len(layers))
	h := 0.0
	for _, l := range layers {
		for _, n := range l {
			neuron := neurons[n.id]
			neuron.X = (neuron.X + 1.0) / 2.0
			neuron.Y = ((neuron.Y+1.0)/2.0)*x + h
			neurons[n.id] = neuron
		}
		h += x + 0.05
	}
}

func pos2D(pos []float64) (x, y float64) {
	x = pos[0]
	y = pos[1]
	// TODO: handle more dimensions, probably just a Z dimension. Should we place it in same x/y and let visualizer move it?
	/*
		for i:=2; 2 < len(pos);i++ {
			if i%2 == 0 {
				x += 0.005
			} else {
				y += 0.005
			}
		}*/
	return
}
