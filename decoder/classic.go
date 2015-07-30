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

	"github.com/rqme/neat"
	"github.com/rqme/neat/network"
)

// Helper that decodes the genome into a neural network
type Classic struct {

	// Number of times to loop through the network
	NetworkIterations int "neat:config"
}

// Configures the helper
func (d *Classic) Configure(cfg string) error {
	return neat.Configure(cfg, d)
}

// Decodes the genome into a phenome
func (d Classic) Decode(g neat.Genome) (p neat.Phenome, err error) {
	// Return the phenome
	net, e := d.decode(g)
	if e != nil {
		err = e
	}
	p = Phenome{
		Genome:  g,
		Network: net,
	}
	return
}

func (d Classic) decode(g neat.Genome) (net *network.Classic, err error) {

	// Number of network iterations must be set
	if d.NetworkIterations == 0 {
		err = fmt.Errorf("Classic decoder must have number of network iterations greater than zero.")
		return
	}

	// Identify the genes
	nodes, conns := g.GenesByPosition()

	// Create the neurons
	nmap := make(map[int]int)
	neurons := make([]network.Neuron, len(nodes))
	for i, ng := range nodes {
		nmap[ng.Innovation] = i
		neurons[i] = network.Neuron{NeuronType: ng.NeuronType, ActivationType: ng.ActivationType, X: ng.X, Y: ng.Y}
	}

	// Create the synapses
	forward := true // Keep track of conenctions to determine if this is a feed-forward only network
	synapses := make([]network.Synapse, 0, len(conns))
	for _, cg := range conns {
		src, tgt := nodes[nmap[cg.Source]], nodes[nmap[cg.Target]]
		forward = forward && src.Y < tgt.Y
		if cg.Enabled {
			synapses = append(synapses, network.Synapse{
				Source: nmap[cg.Source],
				Target: nmap[cg.Target],
				Weight: cg.Weight,
			})
		}
	}

	iterations := d.NetworkIterations
	if forward {
		iterations = 1 // No need for iterating therough the network.
	}
	net, err = network.New(neurons, synapses, iterations)

	return
}

type sortnodes []neat.Node

func (s sortnodes) Len() int { return len(s) }
func (s sortnodes) Less(i, j int) bool {
	if s[i].Y == s[j].Y {
		return s[i].X < s[j].X
	} else {
		return s[i].Y < s[j].Y
	}
}
func (s sortnodes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type sortconns struct {
	nodes map[int]int
	conns []neat.Connection
}

func (s sortconns) Len() int { return len(s.conns) }
func (s sortconns) Less(i, j int) bool {
	si := s.nodes[s.conns[i].Source]
	ti := s.nodes[s.conns[i].Target]
	sj := s.nodes[s.conns[j].Source]
	tj := s.nodes[s.conns[j].Target]
	if ti == tj {
		return si < sj
	} else {
		return ti < tj
	}
}
func (s sortconns) Swap(i, j int) { s.conns[i], s.conns[j] = s.conns[j], s.conns[i] }
