/*
Copyright (c) 2015 Brian Hummer (brian@redq.me), All rights reserved.

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

	// Create a new Substrate
	layers := d.SubstrateLayers()
	ncnt := len(layers[0])
	ccnt := 0
	for i := 1; i < len(layers); i++ {
		ncnt += len(layers[i])
		ccnt += len(layers[i]) * len(layers[i-1])
	}
	s := &Substrate{
		Nodes: make([]SubstrateNode, 0, ncnt),
		Conns: make([]SubstrateConn, 0, ccnt),
	}

	// Add the nodes to the substrate
	i := 0
	for _, l := range layers {
		// TODO: Should I sort the nodes by position in the network?
		for j, n := range l {
			l[j].id = i
			s.Nodes = append(s.Nodes, n)
			i += 1
		}
	}

	// Create connections
	var outputs []float64 // output from the Cppn
	wr := d.WeightRange()
	for l := 1; l < len(layers); l++ {
		for _, src := range layers[l-1] {
			for _, tgt := range layers[l] {
				outputs, err = cppn.Activate(append(src.Position, tgt.Position...))
				if err != nil {
					return nil, err
				}
				w := math.Abs(outputs[l-1])
				if w > 0.2 {
					s.Conns = append(s.Conns, SubstrateConn{
						Source: src.id,
						Target: tgt.id,
						Weight: math.Copysign((w-0.2)*wr/0.8, outputs[l-1]),
					})
				}
			}
		}
	}

	// Return the new network
	var net neat.Network
	net, err = s.Decode()
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
