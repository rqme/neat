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

package crosser

import (
	"github.com/rqme/neat"

	"math/rand"
)

type ClassicSettings interface {
	// Probability that a disabled connection will be renabled in the child
	EnableProbability() float64

	// Probability that the child's connection weight or trait is an average of the values from both parents
	MateByAveragingProbability() float64
}

type Classic struct {
	ClassicSettings
}

// Returns a new genome that is a cross between the two parent genomes
//
// When crossing over, the genes with the same innovation numbers are lined up and crossed over in
// one of two ways. In the first method, matching genes are randomly chosen for the offspring
// genome. Alternatively, the connection weights of matching genes can be averaged (Wright (1991)
// reviews both types of crossover and their merits). NEAT uses both types of crossover. Disjoint
// and excess genes are inherited from the more fit parent, or if they are equally fit, each gene
// is inherited from either parent randomly. Disabled genes have a chance of being reenabled during
// crossover, allowing networks to make use of older genes once again. (Stanley, 38)
func (c Classic) Cross(p1, p2 neat.Genome) (child neat.Genome, err error) {
	rng := rand.New(rand.NewSource(rand.Int63()))

	// Ensure the more fit parent is first
	same := (p1.Fitness == p2.Fitness)
	if p1.Fitness < p2.Fitness {
		p1, p2 = p2, p1
	}

	// Create the child
	child = neat.Genome{}

	// Add the connections
	c.addConns(rng, same, p1, p2, &child)

	// Re-enable connections
	c.enableConns(rng, &child)

	// Ensure the nodes
	c.ensureNodes(rng, p1, p2, &child)

	// Set the traits
	c.setTraits(rng, same, p1, p2, &child)
	return
}

func (c *Classic) addConns(rng *rand.Rand, same bool, p1, p2 neat.Genome, child *neat.Genome) {
	child.Conns = make(map[int]neat.Connection, len(p1.Conns))
	var i, j int
	var c1, c2 neat.Connection
	_, conns1 := p1.GenesByInnovation()
	_, conns2 := p2.GenesByInnovation()
	for i < len(conns1) && j < len(conns2) {
		c1, c2 = conns1[i], conns2[j]
		switch {
		case c1.Innovation < c2.Innovation:
			child.Conns[c1.Innovation] = c1
			i += 1

		case c1.Innovation > c2.Innovation:
			if same {
				child.Conns[c2.Innovation] = c2
			}
			j += 1

		default: // conns1[i].Innovation == conns2[j].Innovation:
			if rng.Float64() < c.MateByAveragingProbability() {
				conn := neat.Connection{
					Innovation: c1.Innovation,
					Source:     c1.Source,
					Target:     c1.Target,
					Enabled:    c1.Enabled, // From NEAT FAQ : In such a situation (which I have found to be rare) you may want to edit the mating code such that disabled genes are only disabled in the offspring if they are disabled in the more fit parent. This fix will keep the disabling of genes to a minimum.
					Weight:     (c1.Weight + c2.Weight) / 2.0,
				}

				child.Conns[conn.Innovation] = conn
			} else {
				if rng.Float64() < 0.5 {
					child.Conns[c1.Innovation] = c1
				} else {
					child.Conns[c2.Innovation] = c2
				}
			}

			i += 1
			j += 1
		}
	}
	for i < len(conns1) {
		c1 = conns1[i]
		child.Conns[c1.Innovation] = c1
		i += 1
	}
	for same && j < len(conns2) {
		c2 = conns2[j]
		child.Conns[c2.Innovation] = c2
		j += 1
	}
}

// Enables connections based on probability
func (c *Classic) enableConns(rng *rand.Rand, child *neat.Genome) {
	for k, conn := range child.Conns {
		if !conn.Enabled && rng.Float64() < c.EnableProbability() {
			conn.Enabled = true
			child.Conns[k] = conn
		}
	}
}

// Ensures that child has proper nodes for each connection
func (c *Classic) ensureNodes(rng *rand.Rand, p1, p2 neat.Genome, child *neat.Genome) {
	child.Nodes = make(map[int]neat.Node, len(p1.Nodes))
	for k, node := range p1.Nodes {
		child.Nodes[k] = node
	}
	for _, conn := range child.Conns {
		var k int
		for i := 0; i < 2; i++ {
			if i == 0 {
				k = conn.Source
			} else {
				k = conn.Target
			}
			if _, ok := child.Nodes[k]; !ok {
				node, ok := p2.Nodes[k]
				if !ok {
					panic("MISSING3")
				}
				child.Nodes[k] = node
			}
		}
	}
}

// Sets the child's traits from one or both of the parents
func (c *Classic) setTraits(rng *rand.Rand, same bool, p1, p2 neat.Genome, child *neat.Genome) {
	child.Traits = make([]float64, len(p1.Traits))
	for i := 0; i < len(child.Traits); i++ {
		if same {
			if rng.Float64() < c.MateByAveragingProbability() {
				child.Traits[i] = (p1.Traits[i] + p2.Traits[i]) / 2.0
			} else {
				if rng.Float64() < 0.5 {
					child.Traits[i] = p1.Traits[i]
				} else {
					child.Traits[i] = p2.Traits[i]
				}
			}
		} else {
			child.Traits[i] = p1.Traits[i] // Take from the more fit parent
		}
	}
}
