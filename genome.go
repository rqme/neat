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

package neat

import (
	"bytes"
	"fmt"
	"sort"
)

// An encoded solution
type Genome struct {
	ID          int         // Identifier for this genome
	SpeciesIdx  int         // Index of species in the population
	Nodes       Nodes       // Neuron definitions
	Conns       Connections // Synapse definitions
	Traits      []float64   // Trait values
	Fitness     float64     // Fitness of genome as it relates to the problem itself
	Improvement float64     // Fitness of genome as it relates to the improvement of the population
	Birth       int         // Generation during which this genome was born
}

func (g Genome) Complexity() int { return len(g.Nodes) + len(g.Conns) }

func (g Genome) String() string {
	b := bytes.NewBufferString(fmt.Sprintf("Genome %d Species %d Fitness %f", g.ID, g.SpeciesIdx, g.Fitness))
	nodes, conns := g.GenesByInnovation()
	b.WriteString("\n\tNodes:")
	for i, n := range nodes {
		b.WriteString(fmt.Sprintf("\n\t\t%d %s", i, n.String()))
	}
	b.WriteString("\n\tConnections:")
	for i, c := range conns {
		b.WriteString(fmt.Sprintf("\n\t\t%d %s", i, c.String()))
	}
	b.WriteString("\n\tTraits:")
	for i, t := range g.Traits {
		b.WriteString(fmt.Sprintf("\n\t\t%d %f", i, t))
	}
	return b.String()
}

// Returns the genome's genes by their markers
func (g Genome) GenesByInnovation() ([]Node, []Connection) {

	// Sort the node genes
	n := &sortNodesByInnovation{}
	n.nodes = make([]Node, 0, len(g.Nodes))
	for _, ng := range g.Nodes {
		n.nodes = append(n.nodes, ng)
	}
	sort.Sort(n)

	// Sort the conn genes
	c := &sortConnsByInnovation{}
	c.conns = make([]Connection, 0, len(g.Conns))
	for _, cg := range g.Conns {
		c.conns = append(c.conns, cg)
	}
	sort.Sort(c)

	// Return the collections
	return n.nodes, c.conns
}

// Returns the genome's genes by their keys
func (g Genome) GenesByPosition() ([]Node, []Connection) {

	// Sort the node genes
	n := &sortNodesByKey{}
	n.nodes = make([]Node, 0, len(g.Nodes))
	for _, ng := range g.Nodes {
		n.nodes = append(n.nodes, ng)
	}
	sort.Sort(n)

	// Sort the conn genes
	c := &sortConnsByKey{}
	c.nodeMap = g.Nodes
	c.conns = make([]Connection, 0, len(g.Conns))
	for _, cg := range g.Conns {
		c.conns = append(c.conns, cg)
	}
	sort.Sort(c)

	// Return the collections
	return n.nodes, c.conns
}

type Genomes []Genome

func (gs Genomes) Len() int           { return len(gs) }
func (gs Genomes) Less(i, j int) bool { return gs[i].Fitness < gs[j].Fitness }
func (gs Genomes) Swap(i, j int)      { gs[i], gs[j] = gs[j], gs[i] }

// Returns a copy of the genome
func CopyGenome(g1 Genome) (g2 Genome) {
	g2.ID = g1.ID
	g2.SpeciesIdx = g1.SpeciesIdx
	g2.Fitness = g1.Fitness
	g2.Improvement = g1.Improvement
	g2.Conns = make(map[int]Connection, len(g1.Conns))
	for k, v := range g1.Conns {
		g2.Conns[k] = v
	}
	g2.Nodes = make(map[int]Node, len(g1.Nodes))
	for k, v := range g1.Nodes {
		g2.Nodes[k] = v
	}
	g2.Traits = make([]float64, len(g1.Traits))
	copy(g2.Traits, g1.Traits)
	return
}
