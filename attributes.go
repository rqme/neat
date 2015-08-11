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

// Provides setup actions in the helper's lifestyle
type Setupable interface {
	// Sets up the helper
	Setup() error
}

// Provides takedown actions in the helper's lifecycle
type Takedownable interface {
	// Takes down the helper
	Takedown() error
}

type OriginalFitnessable interface {
	OriginalFitness() float64
}

// A helper that would like to see the population
type Populatable interface {
	// Provides the population to the helper
	SetPopulation(Population) error
}

// A helper that would like to see the active phenomes
type Phenomable interface {
	// Provides the phenomes to the helper
	SetPhenomes(Phenomes) error
}

// Behaviorable describes an item tracks behaviors expressed during evaluation
type Behaviorable interface {
	// Returns the expressed behaviors
	Behavior() []float64
}

// Demonstrable indicates that a helper is capable of showing (or hiding) detail of its work during
// execution. This is typically employed by an evaluator and turned on when evaluating the "best"
// genome after the evolution is completed.
type Demonstrable interface {
	// ShowWork toggles the flag informing the helper to show (or) not the details during execution
	ShowWork(bool)
}

// Contextable describes an item that can receive a context
type Contextable interface {
	SetContext(Context) error
}

type Crossoverable interface {
	SetCrossover(bool) error
}

type Trialable interface {
	SetTrial(int) error
}

type Improvable interface {
	Improvement() float64
}
