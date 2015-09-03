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

package mutator

import (
	"math/rand"

	"github.com/rqme/neat"
)

// Complexifying mutation settings
type ComplexifySettings interface {
	WeightRange() float64                  // The mutation range of the weight. If x, range is [-x,x]
	AddNodeProbability() float64           // Probablity a node will be added to the genome
	AddConnProbability() float64           // Probability a connection will be added to the genome
	AllowRecurrent() bool                  // Allow recurrent connections to be added
	HiddenActivation() neat.ActivationType // Activation type to assign to new nodes
}

type Complexify struct {
	ComplexifySettings
	ctx neat.Context
}

func (m *Complexify) SetContext(x neat.Context) error {
	m.ctx = x
	return nil
}

// Mutates a genome's weights
func (m *Complexify) Mutate(g *neat.Genome) error {
	rng := rand.New(rand.NewSource(rand.Int63()))
	if rng.Float64() < m.AddNodeProbability() {
		m.addNode(rng, g)
	} else if rng.Float64() < m.AddConnProbability() {
		m.addConn(rng, g)
	}
	return nil
}

// Adds a new node to the genome
//
// In the add node mutation, an existing connection is split and the new node placed where the old
// connection used to be. The old connection is disabled and two new connections are added to the
// genome. The connection between the first node in the chain and the new node is given a weight
// of one, and the connection between the new node and the last node in the chain is given the
// same weight as the connection being split. Splitting the connection in this way introduces a
// nonlinearity (i.e. sigmoid function) where there was none before. Because the new node is
// immediately integrated into the network, its effect on fitness can be evaluated right away.
// Preexisting network structure is not destroyed and performs the same function, while the new
// structure provides an opportunity to elaborate on the original behaviors. (Stanley, 35)
func (m *Complexify) addNode(rng *rand.Rand, g *neat.Genome) {

	// Pick a connection to split
	var inno int
	var c0 neat.Connection
	found := false
	for k, conn := range g.Conns {
		c0 = conn
		inno = k

		// Ensure resultant node doesn't already exist
		found = true
		src := g.Nodes[c0.Source]
		tgt := g.Nodes[c0.Target]
		x := (src.X + tgt.X) / 2.0
		y := (src.Y + tgt.Y) / 2.0
		for _, node := range g.Nodes {
			if node.X == x && node.Y == y {
				found = false
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return
	}
	c0.Enabled = false
	g.Conns[inno] = c0

	// Add the new node
	src := g.Nodes[c0.Source]
	tgt := g.Nodes[c0.Target]
	n0 := neat.Node{NeuronType: neat.Hidden, ActivationType: m.HiddenActivation(), X: (src.X + tgt.X) / 2.0, Y: (src.Y + tgt.Y) / 2.0}
	n0.Innovation = m.ctx.Innovation(neat.NodeInnovation, n0.Key())
	g.Nodes[n0.Innovation] = n0

	// Add the new connections
	c1 := neat.Connection{Source: src.Innovation, Target: n0.Innovation, Enabled: true, Weight: 1.0}
	c1.Innovation = m.ctx.Innovation(neat.ConnInnovation, c1.Key())
	g.Conns[c1.Innovation] = c1

	c2 := neat.Connection{Source: n0.Innovation, Target: tgt.Innovation, Enabled: true, Weight: c0.Weight}
	c2.Innovation = m.ctx.Innovation(neat.ConnInnovation, c2.Key())
	g.Conns[c2.Innovation] = c2
}

// Adds a new connection to the genome
//
// In the add connection mutation, a single new connection gene is added connecting two previously
// unconnected nodes. (Stanley, 35)
func (m *Complexify) addConn(rng *rand.Rand, g *neat.Genome) {

	// Identify two unconnected nodes
	conns := make(map[int]neat.Connection)
	c := 0
	for _, src := range g.Nodes {
		for _, tgt := range g.Nodes {
			if src.Innovation == tgt.Innovation {
				continue // Must be two unconnected. TODO: Allow self-connections?
			}
			if tgt.NeuronType == neat.Bias || tgt.NeuronType == neat.Input {
				continue // cannot connect back to input layer
			} else if src.NeuronType == neat.Output && tgt.NeuronType == neat.Output {
				continue
			}
			if !m.AllowRecurrent() && tgt.Y <= src.Y {
				continue
			}
			found := false
			for _, c2 := range g.Conns {
				if c2.Source == src.Innovation && c2.Target == tgt.Innovation {
					found = true
					break
				}
			}
			if !found {
				conns[c] = neat.Connection{
					Source:  src.Innovation,
					Target:  tgt.Innovation,
					Enabled: true,
					Weight:  (rng.Float64()*2.0 - 1.0) * m.WeightRange(),
				}
				c += 1
			}
		}
	}

	// Go's range over maps is random, so take the first, if any, availble connection
	for _, conn := range conns {
		conn.Innovation = m.ctx.Innovation(neat.ConnInnovation, conn.Key())
		g.Conns[conn.Innovation] = conn
		break
	}
}
