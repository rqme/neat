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

package speciater

import (
	"github.com/rqme/neat"
)

type ClassicSettings interface {
	CompatibilityThreshold() float64 // Threshold above which two genomes are not compatible
}

type Classic struct {
	ClassicSettings
	ctx neat.Context
}

func (s *Classic) SetContext(x neat.Context) error {
	s.ctx = x
	return nil
}

// Assigns the genomes to a species. Returns new collection of species.
//
// Throughout evolution, NEAT maintains a list of species numbered in the order they ap- peared. In
// the first generation, since there are no preexisting species, NEAT begins by creating species 1
// and placing the first genome into that species. All other genomes are placed into species as
// follows: A random member of each existing species is chosen as its permanent representative.
// Genomes are tested one at a time; if a genome’s distance to the representative of any existing
// species is less than δt, a compatibility threshold, it is placed into this species. Otherwise,
// if it is not compatible with any existing species, a new species is created and given a new
// number. After the first generation, genomes are first compared with species from the previous
// generation so that the same species numbers can be used to identify species throughout the run.
// (Stanley, 39)
//
// TODO: Pick a random reprentative each time instead of permanently recording it with the species
// record
func (s Classic) Speciate(curr []neat.Species, genomes []neat.Genome) (next []neat.Species, err error) {

	// Copy the species to the new set
	next = make([]neat.Species, len(curr))
	for i, s := range curr {
		next[i] = s
		next[i].Age = s.Age + 1
	}

	// Iterate the genomes, looking for target species
	// TODO: This could be made concurrent if it is slow
	var δ float64
	cnts := make([]int, len(curr))
	for i, genome := range genomes {
		found := false
		for j, species := range next {
			δ, err = s.ctx.Comparer().Compare(genome, species.Example)
			if err != nil {
				return
			}
			if δ < s.CompatibilityThreshold() {
				genomes[i].SpeciesIdx = j
				cnts[j] += 1
				found = true
				break
			}
		}
		if !found {
			genomes[i].SpeciesIdx = len(next)
			cnts = append(cnts, 1)
			species := neat.Species{
				Example: neat.CopyGenome(genomes[i]),
			}
			next = append(next, species)

		}
	}

	// Purge unused species
	i := 0
	for i < len(next) {
		if cnts[i] == 0 {
			next = append(next[:i], next[i+1:]...)
			cnts = append(cnts[:i], cnts[i+1:]...)
			for j := 0; j < len(genomes); j++ {
				if genomes[j].SpeciesIdx > i {
					genomes[j].SpeciesIdx -= 1
				}
			}
		} else {
			i += 1
		}
	}

	return
}
