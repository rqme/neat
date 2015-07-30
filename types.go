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
	"fmt"
)

// Structure to hold the configuration and state of the experiment
type Archive struct {
	Config string
	State  string
}

// Provides next integer in a sequence
type IDSequence interface {

	// Returns the next int in the sequence
	Next() int
}

// Tracks innovations
type Marker interface {

	// Marks the node with its innovation number
	MarkNode(*Node)

	// Marks the connection with its innovation number
	MarkConn(*Connection)

	// Resets the innovation history
	Reset()
}

// A decoded solution
type Phenome interface {

	// The identify of the underlying genome
	ID() int

	// Trait values which could be used during evaluation
	Traits() []float64

	// The decoded genome
	Network
}

type Phenomes []Phenome

// Population of genomes for a given generation
type Population struct {
	Generation int
	Species    []Species
	Genomes    []Genome
}

// The result of an evaluation
type Result interface {

	// Returns the ID of the phenome
	ID() int

	// Returns the fitness of the phenome for the problem
	Fitness() float64

	// Returns the error, if any, occuring while evaluating the phenome.
	Err() error

	// Returns true if the stop condition was met
	Stop() bool
}

type Species struct {
	Age        int // Age in terms of generations
	Stagnation int // Number of generations since an improvement in fitness
	Fitness    float64
	Example    Genome
}

type Trait struct {
	Name      string
	Min, Max  float64
	IsSetting bool
}

type Traits []Trait

func (t Traits) IndexOf(name string) int {
	for i, trait := range t {
		if trait.Name == name {
			return i
		}
	}
	return -1
}

type FitnessType byte

const (
	AbsoluteFitness FitnessType = iota
	RelativeFitness
)

func (f FitnessType) String() string {
	switch f {
	case AbsoluteFitness:
		return "Absolute Fitness"
	case RelativeFitness:
		return "Relative Fitness"
	default:
		return fmt.Sprintf("Unknown FitnessType: %d", f)
	}
}
