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
	"github.com/rqme/neat"

	"math/rand"
)

type Pruning struct {

	// Probablity a node will be removed to the genome
	DelNodeProbability float64 "neat:config"

	// Probability a connection will be removed to the genome
	DelConnProbability float64 "neat:config"
}

// Configures the helper from a JSON string
func (m *Pruning) Configure(cfg string) error {
	return neat.Configure(cfg, m)
}

// Mutates a genome's weights
func (m Pruning) Mutate(g *neat.Genome) error {
	rng := rand.New(rand.NewSource(rand.Int63()))
	if rng.Float64() < m.DelNodeProbability {
		m.delNode(rng, g)
	} else if rng.Float64() < m.DelConnProbability {
		m.delConn(rng, g)
	}
	return nil
}

// Removes a hidden node from the genome
//
// Neuron deletion is slightly more complex. The deletion algorithm attempts to replace neurons with
// connections to maintain any circuits a neuron may have participated in, in further generations
// those connections themselves will be open to deletion. This approach provides NEAT with the ability
// to delete whole structures, not just connections.
//
// Because we replace connected neurons with connections we must be careful which neurons we delete.
// Any neuron with only incoming or only outgoing connections is at a dead-end of a circuit and can
// therefore be safely deleted with all of it's connections. However, a neuron with multiple incoming
// and multiple outgoing connections will require a large number of connections to substitute for the
// loss of the neuron - we must fully connect all of the original neuron's source neurons with its
// target neurons, this could be done but may actually be detrimental since the functionality
// represented by the neuron is now distributed over a number of connections, and this cannot easily
// be reversed. Because of this, such neurons are omitted from the process of selecting neurons for
// deletion.
//
// Neurons with only one incoming or one outgoing connection can be replaced with however many
// connections were on the other side of the neuron, therefore these are candidates for deletion.
func (m *Pruning) delNode(rng *rand.Rand, g *neat.Genome) {

	type check struct {
		AsSource []neat.Connection
		AsTarget []neat.Connection
	}

	// Build a map of available nodes to delete
	var chk check
	var node neat.Node
	avail := make(map[neat.Node]check)
	for _, node = range g.Nodes {
		if node.NeuronType == neat.Hidden {
			chk = check{make([]neat.Connection, 0, 5), make([]neat.Connection, 0, 5)}
			for _, conn := range g.Conns {
				if conn.Source == node.Innovation {
					chk.AsSource = append(chk.AsSource, conn)
				} else if conn.Target == node.Innovation {
					chk.AsTarget = append(chk.AsTarget, conn)
				}
			}
			if len(chk.AsSource) <= 1 || len(chk.AsTarget) <= 1 {
				avail[node] = chk
			}
		}
	}
	if len(avail) == 0 {
		return // there are nodes available to delete
	}

	// Pick a node to delete
	for node, chk = range avail {
		break
	}

	// Remove dead-end connections
	if len(chk.AsSource) == 0 {
		for _, conn := range chk.AsTarget {
			delete(g.Conns, conn.Innovation)
		}
	} else if len(chk.AsTarget) == 0 {
		for _, conn := range chk.AsSource {
			delete(g.Conns, conn.Innovation)
		}
	} else {
		// Bypass this node in all the connections. Only one of the slices will have more than 1 connection
		for _, sc := range chk.AsSource {
			for _, tc := range chk.AsTarget {
				sc.Target = tc.Target
				tc.Source = sc.Source
				g.Conns[sc.Innovation] = sc
				g.Conns[tc.Innovation] = tc
			}
		}
	}
}

// Removes a connection from the genome
//
// Connection deletion is very simply the deletion of a randomly selected connection, all connections
// are considered to be available for deletion. When a connection is deleted the neurons that were at
// each end of the connection are tested to check if they are no longer connected to by other
// connections, if this is the case then the stranded neuron is also deleted. Note that a more thorough
// cleanup routine could be invoked at this point that cleans up any dead-end structures that could not
// possibly be functional, but this can become complex and so we leave NEAT to eliminate such structures
// naturally.
func (m *Pruning) delConn(rng *rand.Rand, g *neat.Genome) {

	// Pick a connection at random
	var conn neat.Connection
	for _, conn = range g.Conns {
		break
	}

	// Remove the connection from the genome
	delete(g.Conns, conn.Innovation)

	// Look for orphaned nodes
	var node neat.Node
	var found bool
	for i := 0; i < 2; i++ {
		found = false
		if i == 0 {
			node = g.Nodes[conn.Source]
		} else {
			node = g.Nodes[conn.Target]
		}
		if node.NeuronType != neat.Hidden {
			continue
		}

		// Check for another connection using this node gene
		for _, conn2 := range g.Conns {
			if conn2.Source == node.Innovation || conn2.Target == node.Innovation {
				found = true
				break
			}
		}
		if !found {
			delete(g.Nodes, node.Innovation)
		}
	}

}
