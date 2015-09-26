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
	"math"

	"github.com/rqme/neat"
)

type ESHyperNEATSettings interface {
	InitialDepth() int
	MaxDepth() int
	DivisionThreshold() float64
	VarianceThreshold() float64
	BandThreshold() float64
	IterationLevels() int

	HyperNEATSettings
}

type ESHyperNEAT struct {
	ESHyperNEATSettings
	CppnDecoder neat.Decoder

	dims int
	divs [][]float64
}

func NewESHyperNEAT(cfg ESHyperNEATSettings, dec neat.Decoder) *ESHyperNEAT {
	d := &ESHyperNEAT{ESHyperNEATSettings: cfg, CppnDecoder: dec}

	// Create the division factors
	d.dims = len(d.SubstrateLayers()[0][0].Position)
	k := int(math.Pow(2, float64(d.dims)))
	d.divs = make([][]float64, k)
	for i := 0; i < k; i++ {
		d.divs[i] = make([]float64, d.dims)
		for j := 0; j < d.dims; j++ {
			x := int(math.Floor(float64(i)/math.Pow(2, float64(j)))) % 2
			if x == 0 {
				d.divs[i][j] = -1
			} else {
				d.divs[i][j] = 1
			}
		}
	}
	return d
}

//
//
//
func (d *ESHyperNEAT) divAndInit(cppn neat.Network, t int, a []float64, outgoing bool) (root *espoint, err error) {
	root = &espoint{
		position: make([]float64, d.dims),
		width:    1,
		level:    1,
		children: make([]*espoint, len(d.divs)),
	}
	q := make([]*espoint, 0, 16)
	q = append(q, root)
	for len(q) > 0 {
		// Dequeue the next item
		p := q[0]
		q[0] = nil
		q = q[1:]

		// Divide into subregions
		for i := 0; i < len(p.children); i++ {
			c := &espoint{
				position: make([]float64, len(p.position)),
				width:    p.width / 2.0,
				level:    p.level + 1,
			}
			for j := 0; j < d.dims; j++ {
				c.position[j] = p.position[j] + c.width*d.divs[i][j]
			}
			p.children[i+0] = c
		}

		// Process the children
		for _, c := range p.children {
			var outputs []float64
			if outgoing {
				outputs, err = cppn.Activate(append(a, append(c.position, 1.0)...))
			} else {
				outputs, err = cppn.Activate(append(c.position, append(a, 0.0)...))
			}
			if err != nil {
				return
			}
			c.weight = outputs[t]
		}

		// Divide until intial resolution or if variance is still high
		if p.level < d.InitialDepth() || (p.level < d.MaxDepth() && variance(p) > d.DivisionThreshold()) {
			for i, _ := range p.children {
				q = append(q, p.children[i])
			}
		}
	}
	return
}

func (d *ESHyperNEAT) pruneAndExtract(cppn neat.Network, t int, a []float64, p *espoint, outgoing bool) (conns esconns, err error) {
	for _, c := range p.children {
		if variance(c) > d.VarianceThreshold() {
			var con2 esconns
			con2, err = d.pruneAndExtract(cppn, t, a, c, outgoing)
			if err != nil {
				return
			}
			conns = append(conns, con2...)
		} else {
			max := 0.0
			// Determine if point is in a band by checking neighbor CPPN values
			var outputs []float64
			for i := 0; i < d.dims; i++ {
				min := math.Inf(1)
				for j := 0; j < 2; j++ {
					x := c.position[i]
					if j == 0 {
						c.position[i] -= p.width
					} else {
						c.position[i] += p.width
					}
					if outgoing {
						outputs, err = cppn.Activate(append(a, append(c.position, 1.0)...))
					} else {
						outputs, err = cppn.Activate(append(c.position, append(a, 0.0)...))
					}
					if err != nil {
						return
					}
					if min > outputs[t] {
						min = outputs[t]
					}
					c.position[i] = x
				}
				if max < min {
					max = min
				}
			}
			if max > d.BandThreshold() {
				// Create new connection specified by source, target, weight
				// and scale weight based on weight range
				var conn *esconn
				if outgoing {
					conn = &esconn{source: a, target: c.position}
				} else {
					conn = &esconn{source: c.position, target: a}
				}
				if !conns.contains(conn) {
					conn.weight = c.weight * d.WeightRange()
					conns = append(conns, conn)
				}
			}
		}
	}
	return
}

func (d *ESHyperNEAT) Decode(g neat.Genome) (p neat.Phenome, err error) {

	// Decode the CPPN
	var cppn neat.Phenome
	if cppn, err = d.CppnDecoder.Decode(g); err != nil {
		return
	}

	// Create a new substratre
	s := &Substrate{
		Nodes: make([]SubstrateNode, 0, 100),
		Conns: make([]SubstrateConn, 0, 100),
	}

	// Create the hidden layer(s)
	layers := d.SubstrateLayers()
	layers = append(layers[:1], append(make([]SubstrateNodes, d.IterationLevels()), layers[1:]...)...)

	// Assign ids to the inputs
	id := 0
	for i := 0; i < len(layers[0]); i++ {
		layers[0][i].id = id
		id += 1
		s.Nodes = append(s.Nodes, layers[0][i])
	}

	// Create hidden nodes and connections
	var inputs SubstrateNodes = layers[0]
	var hidden SubstrateNodes
	for t := 0; t < d.IterationLevels()+1; t++ {

		hidden = make([]SubstrateNode, 0, len(inputs))
		for i := 0; i < len(inputs); i++ {

			// Analyze the outgoing connectivity pattern form this input
			var root *espoint
			if root, err = d.divAndInit(cppn, t, inputs[i].Position, true); err != nil {
				return
			}

			// Traverse the tree and add conections to the list
			var conns esconns
			if conns, err = d.pruneAndExtract(cppn, t, inputs[i].Position, root, true); err != nil {
				return
			}
			for _, c := range conns {
				n := SubstrateNode{Position: c.target, NeuronType: neat.Hidden}
				idx := hidden.IndexOf(n)
				if idx == -1 {
					n.id = id
					id += 1
					hidden = append(hidden, n)
					s.Nodes = append(s.Nodes, n)
				} else {
					n = hidden[idx]
				}
				s.Conns = append(s.Conns, SubstrateConn{Source: inputs[i].id, Target: n.id, Weight: c.weight})
			}
		}
		layers[t+1] = hidden
		inputs = hidden
	}

	// Assign IDs to output nodes
	var outputs SubstrateNodes = layers[len(layers)-1]
	for i := 0; i < len(outputs); i++ {
		//outputs[i] = layers[len(layers)-1][i]
		outputs[i].id = id
		id += 1
		s.Nodes = append(s.Nodes, outputs[i])
	}

	// Output to hidden layer(s)
	for t := 0; t < d.IterationLevels(); t++ {
		hidden = layers[t+1]
		for i := 0; i < len(outputs); i++ {

			// Analyze the outgoing connectivity pattern form this input
			var root *espoint
			if root, err = d.divAndInit(cppn, t, outputs[i].Position, false); err != nil {
				return
			}

			// Traverse the tree and add conections to the list
			var conns esconns
			if conns, err = d.pruneAndExtract(cppn, t, outputs[i].Position, root, false); err != nil {
				return
			}
			for _, c := range conns {
				n := SubstrateNode{Position: c.target, NeuronType: neat.Hidden}
				idx := hidden.IndexOf(n)
				if idx == -1 {
					// Ignore as it would create an unconnected node
				} else {
					n = hidden[idx]
					s.Conns = append(s.Conns, SubstrateConn{Source: n.id, Target: outputs[i].id, Weight: c.weight})
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

type espoint struct {
	position []float64
	width    float64
	level    int
	weight   float64
	children []*espoint
}

func variance(p *espoint) float64 {
	var sum, sqr float64
	for i := 0; i < len(p.children); i++ {
		sum += p.children[i].weight
		sqr += p.children[i].weight * p.children[i].weight
	}
	mean := sum / float64(len(p.children))
	return sqr / (float64(len(p.children)) - mean*mean)
}

type esconn struct {
	source []float64
	target []float64
	weight float64
}

type esconns []*esconn

func (e esconns) contains(c *esconn) bool {
	for _, x := range e {
		found := true
		for i := 0; i < len(x.source); i++ {
			if x.source[i] != c.source[i] || x.target[i] == c.target[i] {
				found = false
				break
			}
		}
		if found {
			return true
		}
	}
	return false
}
