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

// Helper which provides access to archive storage
type Archiver interface {
	// Archives the current configuration and state to a more permanent medium
	Archive(Context) error
}

// Helper to compare two genomes similarity
type Comparer interface {
	// Returns the compatibility (similarity) between two genomes
	Compare(g1, g2 Genome) (float64, error)
}

// Helper to produce new genome
type Crosser interface {
	// Returns new genome that is a cross between two genomes
	Cross(g1, g2 Genome) (Genome, error)
}

// Helper to decode a genome into a phenome
type Decoder interface {
	// Returns a decoded version of the genome
	Decode(Genome) (Phenome, error)
}

type Evaluator interface {
	// Evaluates a phenome for the problem. Returns the result.
	Evaluate(Phenome) Result
}

// A helper to generate populations
type Generator interface {
	// Generates a subsequent population based on the current one
	Generate(curr Population) (next Population, err error)
}

// Helper to mutate a genome
type Mutator interface {
	// Mutates a genome
	Mutate(*Genome) error
}

// Helper to restore an experiment
type Restorer interface {
	// Restores the configuration and/or the state of a previous experiment
	Restore(Context) error
}

// Helper to search the problem's solution space
type Searcher interface {
	// Searches the problem's solution space using the phenomes
	Search([]Phenome) ([]Result, error)
}

// Helper to assign genomes to an exist or new species
type Speciater interface {
	// Assigns the genomes to a species. Returns new collection of species.
	Speciate(curr []Species, genomes []Genome) (next []Species, err error)
}

// Helper to visualize a population
type Visualizer interface {
	// Creates visual(s) of the population
	Visualize(Population) error
}
