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
	"bytes"
	"fmt"

	"github.com/rqme/neat"
	"github.com/rqme/neat/network"
)

type SubstrateNode struct {
	id       int // internal ID of node in substrate
	Position []float64
	neat.NeuronType
}

func (n SubstrateNode) String() string {
	return fmt.Sprintf("%d %v %s", n.id, n.Position, n.NeuronType)
}

type SubstrateNodes []SubstrateNode

func (s SubstrateNodes) Contains(n SubstrateNode) bool {
	return s.IndexOf(n) != -1
}

func (s SubstrateNodes) IndexOf(n SubstrateNode) int {
	for j, x := range s {
		f := true
		for i := 0; i < len(x.Position); i++ {
			if x.Position[i] != n.Position[i] {
				f = false
				break
			}
		}
		if f {
			return j
		}
	}
	return -1
}

type SubstrateConn struct {
	Source, Target int // IDs of the source and target nodes
	Weight         float64
}

func (c SubstrateConn) String() string {
	return fmt.Sprintf("%d -> %d : %f", c.Source, c.Target, c.Weight)
}

type SubstrateConns []SubstrateConn

type Substrate struct {
	Nodes SubstrateNodes
	Conns SubstrateConns
}

func (s Substrate) String() string {
	b := bytes.NewBufferString("Substrate")
	b.WriteString("\n\tNodes:")
	for i, n := range s.Nodes {
		b.WriteString(fmt.Sprintf("\n\t\t%d - %s", i, n))
	}
	b.WriteString("\n\tConns:")
	for i, c := range s.Conns {
		b.WriteString(fmt.Sprintf("\n\t\t%d - %s", i, c))
	}
	return b.String()
}
func (s Substrate) Decode() (neat.Network, error) {

	// Create neurons from the nodes
	ns := make([]network.Neuron, len(s.Nodes))
	nm := make(map[int]int, len(s.Nodes))
	for i, sn := range s.Nodes {
		nm[sn.id] = i
		ns[i] = network.Neuron{NeuronType: sn.NeuronType}
		ns[i].X, ns[i].Y = sn.Position[0], sn.Position[1] // TODO: Improve this as all layers will be collapsed
		switch sn.NeuronType {
		case neat.Input, neat.Bias:
			ns[i].ActivationType = neat.Direct
		default:
			ns[i].ActivationType = neat.Sigmoid
		}
	}

	// Create synapses from the connections
	cs := make([]network.Synapse, len(s.Conns))
	for i, sc := range s.Conns {
		cs[i] = network.Synapse{
			Source: nm[sc.Source],
			Target: nm[sc.Target],
			Weight: sc.Weight,
		}
	}

	// Return the new network
	return network.New(ns, cs)
}

// Trims the substrate of connections and hidden nodes that are not part of a valid path from
// input to outpout
func (s *Substrate) trim() {

	// Iterate output neurons
}
