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

package network

import (
	. "github.com/rqme/errors"
	"github.com/rqme/neat"

	"bytes"
	"fmt"
)

type Neuron struct {
	neat.NeuronType
	neat.ActivationType
	X, Y float64 // Hint at where neuron might be positioned in a 2D representation
}

type Synapse struct {
	Source, Target int
	Weight         float64
}

type Activation func(float64) float64

type Classic struct {

	// Structure
	Neurons    []Neuron
	Synapses   []Synapse
	Iterations int

	// Internal state
	biases, inputs, hiddens, outputs int
	funcs                            []Activation
}

func New(neurons []Neuron, synapses []Synapse, iterations int) (net *Classic, err error) {

	// Begin a new network
	net = &Classic{Neurons: neurons, Synapses: synapses, Iterations: iterations}

	// Create the internal state and check for errors
	oo := false    // out-of-order check
	l := neat.Bias // last neuron type created
	net.funcs = make([]Activation, len(neurons))
	for i, ng := range neurons {
		switch ng.NeuronType {
		case neat.Bias:
			net.biases += 1
			oo = oo || l > neat.Bias
		case neat.Input:
			net.inputs += 1
			oo = oo || l > neat.Input
		case neat.Hidden:
			net.hiddens += 1
			oo = oo || l > neat.Hidden
		case neat.Output:
			net.outputs += 1
			oo = oo || l > neat.Output
		}
		l = ng.NeuronType
		switch ng.ActivationType {
		case neat.Direct:
			net.funcs[i] = neat.DirectActivation
		case neat.Sigmoid:
			net.funcs[i] = neat.SigmoidActivation
		case neat.SteependSigmoid:
			net.funcs[i] = neat.SteependSigmoidActivation
		case neat.Tanh:
			net.funcs[i] = neat.TanhActivation
		case neat.InverseAbs:
			net.funcs[i] = neat.InverseAbsActivation
		default:
			err = fmt.Errorf("network.classic.New - Unknown ActivationType %v", byte(ng.ActivationType))
			break
		}
	}
	if oo {
		err = fmt.Errorf("network.classic.New - Neurons are out of order")
		return
	}

	// Ensure we have inputs and outputs
	if net.inputs == 0 {
		err = fmt.Errorf("network.classic.New - Network must have at least 1 input neuron")
		return
	}
	if net.outputs == 0 {
		err = fmt.Errorf("network.classic.New - Network must have at least 1 output neuron")
		return
	}

	// Ensure the synapses map to neurons and count the sources
	cnt := len(net.Neurons)
	for _, s := range net.Synapses {
		if s.Source > cnt || s.Target > cnt {
			err = fmt.Errorf("network.classic.New - Synapses do not map to defined neurons")
			return
		}
	}

	return
}

func (n Classic) String() string {
	b := bytes.NewBufferString("Network is \n")
	b.WriteString("\tNeurons:\n")
	for i, neuron := range n.Neurons {
		b.WriteString(fmt.Sprintf("\t [%d] Type: %v Activation: %v Position: [%f, %f]\n", i, neuron.NeuronType, neuron.ActivationType, neuron.X, neuron.Y))
	}
	b.WriteString("\tSynapses:\n")
	for i, synapse := range n.Synapses {
		b.WriteString(fmt.Sprintf("\t [%d] Source: %d Target: %d Weight: %f\n", i, synapse.Source, synapse.Target, synapse.Weight))
	}
	return b.String()
}

func (n Classic) Activate(inputs []float64) (outputs []float64, err error) {

	//neat.DBG("================= %v =================", inputs)
	// Create the data structures
	errs := new(Errors)
	val := make([]float64, len(n.Neurons))

	// Set the biases
	for i := 0; i < n.biases; i++ {
		val[i] = 1.0
	}

	// Copy inputs into the network
	if len(inputs) > n.inputs {
		errs.Add(fmt.Errorf("network.classic.Activate - There are more input values (%d) than input neurons (%d)\n%v", len(inputs), n.inputs, n))
		return
	}
	copy(val[n.biases:], inputs)

	// Iterate the network synapse by synapse
	for i := 0; i < n.Iterations; i++ {
		for _, s := range n.Synapses {
			v := n.funcs[s.Source](val[s.Source])
			val[s.Target] += v * s.Weight
			//neat.DBG("After iteration %d (%d->%d): %v", j, s.Source, s.Target, val)
		}
	}

	// Return the output values
	offset := len(val) - n.outputs
	outputs = make([]float64, n.outputs)
	for i := 0; i < len(outputs); i++ {
		v := n.funcs[i+offset](val[i+offset])
		outputs[i] = v
	}

	// Return output and any errors
	err = errs.Err()
	return
}
