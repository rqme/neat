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

package generator

import (
	"github.com/rqme/neat"

	"math"
	"math/rand"
)

type ClassicSettings interface {
	// Size of the population
	PopulationSize() int

	// Genome used to seed the population
	SeedGenome() neat.Genome
	Traits() neat.Traits

	// Network definition used if no seed genome is provided
	NumInputs() int
	NumOutputs() int
	OutputActivation() neat.ActivationType
	WeightRange() float64

	// Percent of population to be allowed to produce offspring
	SurvivalThreshold() float64

	// Probability of producing offspring only by mutation
	MutateOnlyProbability() float64

	// Rate at which a mate is chosen from another species
	InterspeciesMatingRate() float64

	// Maximum number of generations a stagnant species may exist
	MaxStagnation() int
}

type Classic struct {
	ClassicSettings
	ctx neat.Context

	cross bool
}

func (g *Classic) SetContext(x neat.Context) error {
	g.ctx = x
	return nil
}

func (g *Classic) SetCrossover(v bool) error {
	g.cross = v
	return nil
}

func (g *Classic) Generate(curr neat.Population) (next neat.Population, err error) {
	if len(curr.Genomes) == 0 {
		return g.generateFirst()
	} else {
		return g.generateNext(curr)
	}
}

// Generates the initial population
func (g *Classic) generateFirst() (next neat.Population, err error) {
	// Create the first generation
	next = neat.Population{
		Generation: 0,
		Species:    make([]neat.Species, 1, 10),
		Genomes:    make([]neat.Genome, g.PopulationSize()),
	}

	// Ensure seed exists
	seed := createSeed(g.ctx, g.ClassicSettings)

	// Create the initial species
	next.Species[0] = neat.Species{Example: seed}

	// Create the genomes
	for i := 0; i < g.PopulationSize(); i++ {
		genome := neat.CopyGenome(seed)
		genome.ID = g.ctx.NextID()
		genome.SpeciesIdx = 0
		if err = g.ctx.Mutator().Mutate(&genome); err != nil {
			return
		}
		next.Genomes[i] = genome
	}
	return
}

// Generates a subsequent population based on the current one
//
// Every species is assigned a potentially different number of offspring in proportion to the sum
// of adjusted fitnesses fi′ of its member organisms. (Stanley, 40)
//
// The lowest performing fraction of each species is eliminated. The parents to produce the next
// generation are chosen randomly among the remaining individuals (uniform distribution with re-
// placement). The highest performing individual in each species, i.e. the species champions,
// carries over from each generation. Otherwise the next generation completely replaces the one
// before. (Stanley, 40)
func (g *Classic) generateNext(curr neat.Population) (next neat.Population, err error) {

	// Update context with current population
	for _, h := range []interface{}{g.ctx.Comparer(), g.ctx.Crosser(), g.ctx.Mutator(), g.ctx.Speciater()} {
		if ph, ok := h.(neat.Populatable); ok {
			if err = ph.SetPopulation(curr); err != nil {
				return
			}
		}
	}

	// Advance the population to the next generation
	next = neat.Population{
		Generation: curr.Generation + 1,
		Species:    make([]neat.Species, 0, len(curr.Species)),
		Genomes:    make([]neat.Genome, 0, g.PopulationSize()),
	}

	// Process existing population
	pool := createPool(curr)

	// Purge stagnant species unliess it contains the best genome
	purgeSpecies(g.ClassicSettings, curr.Species, pool)

	// Calculate offspring counts
	cnts := createCounts(g.ClassicSettings, curr.Species, pool)

	// Preserve elites
	for i, l := range pool {
		if len(l) < 5 {
			cnts[i] = cnts[i] + 1
		} else {
			next.Genomes = append(next.Genomes, l[0])
		}
	}

	// Create the offspring
	rng := rand.New(rand.NewSource(rand.Int63()))
	err = createOffspring(g.ctx, g.ClassicSettings, g.cross, rng, pool, cnts, &next)
	if err != nil {
		return
	}

	// Speciate the genomes
	next.Species, err = g.ctx.Speciater().Speciate(curr.Species, next.Genomes)
	return
}

// Every species is assigned a potentially different number of offspring in proportion to the sum
// of adjusted fitnesses fi′ of its member organisms. The net effect of fitness sharing in NEAT
// can be summarized as follows. Let Fk be the average fitness of species k and |P | be the size
// of the population. Let F tot = 􏰇k Fk be the total of all species fitness averages. The number of
// offspring nk allotted to species k is: See figure 3.3 (Stanley, 40)
func createCounts(cfg ClassicSettings, species []neat.Species, pool map[int]Improvements) (cnts map[int]int) {

	// Note the total fitness
	var tot float64
	for i, l := range pool {
		f := l.Improvement()
		if species[i].Age < 10 {
			f *= 1.2 // Youth boost
		} else if species[i].Age > 30 {
			f *= 0.2 // Old penalty
		}
		tot += f
	}

	// Calculate the target number of offspring
	avail := float64(cfg.PopulationSize() - len(pool)) // preserve room for elite
	cnt := 0
	cnts = make(map[int]int)
	for idx, l := range pool {
		f := l.Improvement()
		if species[idx].Age < 10 {
			f *= 1.2 // Youth boost
		} else if species[idx].Age > 30 {
			f *= 0.2 // Old penalty
		}

		pct := f / tot
		tgt := int(math.Ceil(pct * avail))
		cnts[idx] = tgt
		cnt += tgt
	}

	// Trim back down to overcome rounding in above calculation
	for cnt > int(avail) {
		for idx, n := range cnts { // Go's range over maps is random. Yay!
			if n > 0 {
				cnts[idx] = n - 1
				cnt -= 1
				break
			}
		}
	}
	return
}
