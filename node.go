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
	"encoding/json"
	"fmt"
	"sort"
)

// Definition of a neuron
type Node struct {
	Innovation int
	X, Y       float64
	NeuronType
	ActivationType
}

func (n Node) String() string {
	return fmt.Sprintf("Node %d at [%f, %f] Neuron %v Activation %v", n.Innovation, n.X, n.Y, n.NeuronType, n.ActivationType)
}

// Nodes is a map of nodes by innovation number
type Nodes map[int]Node

// nodeToSlice converts a nodes map into a slice for encoding
func (nm Nodes) nodesToSlice() []Node {

	// Maps with non-string keys cannot be encoded. Transfer to a slice to handle this
	items := &sortNodesByInnovation{make([]Node, 0, len(nm))}
	for _, s := range nm {
		items.nodes = append(items.nodes, s)
	}
	sort.Sort(items)
	return items.nodes
}

// nodesFromSlice converts a slice into a nodes map for decoding
func (nm *Nodes) nodesFromSlice(items []Node) {
	if *nm == nil {
		*nm = make(map[int]Node)
	}
	m := *nm
	for _, s := range items {
		m[s.Innovation] = s
	}
}

// MarshalJSON marshals the nodes map into JSON
func (nm Nodes) MarshalJSON() ([]byte, error) {
	items := nm.nodesToSlice()
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(items)
	return buf.Bytes(), err
}

// UnmarshalJSON unmarshals JSON into a nodes map
func (nm *Nodes) UnmarshalJSON(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	var items []Node
	err := dec.Decode(&items)
	if err == nil {
		nm.nodesFromSlice(items)
	}
	return err
}

type sortNodesByKey struct {
	nodes []Node
}

func (g *sortNodesByKey) Len() int { return len(g.nodes) }
func (g *sortNodesByKey) Less(i, j int) bool {
	if g.nodes[i].Y == g.nodes[j].Y {
		return g.nodes[i].X < g.nodes[j].X
	} else {
		return g.nodes[i].Y < g.nodes[j].Y
	}
}
func (g *sortNodesByKey) Swap(i, j int) { g.nodes[i], g.nodes[j] = g.nodes[j], g.nodes[i] }

type sortNodesByInnovation struct {
	nodes []Node
}

func (g *sortNodesByInnovation) Len() int { return len(g.nodes) }
func (g *sortNodesByInnovation) Less(i, j int) bool {
	return g.nodes[i].Innovation < g.nodes[j].Innovation
}
func (g *sortNodesByInnovation) Swap(i, j int) { g.nodes[i], g.nodes[j] = g.nodes[j], g.nodes[i] }
