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
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/rqme/neat"
)

// Generates the initial population
func generateFirst(ctx neat.Context, cfg ClassicSettings) (next neat.Population, err error) {
	// Create the first generation
	next = neat.Population{
		Generation: 0,
		Species:    make([]neat.Species, 1, 10),
		Genomes:    make([]neat.Genome, cfg.PopulationSize()),
	}

	// Create the genomes
	wg := new(sync.WaitGroup)
	for i := 0; i < len(next.Genomes); i++ {
		wg.Add(1)
		go func(i int) {
			genome := createSeed(ctx, cfg)
			genome.ID = ctx.NextID()
			genome.SpeciesIdx = 0
			next.Genomes[i] = genome
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Create the initial species
	next.Species[0] = neat.Species{Example: next.Genomes[0]}

	return
}

// Creates the pool of potential parents, grouped by species' index. The list of genomes is also
// sorted by fitness in decending order for future operations
func createPool(curr neat.Population) (pool map[int]Improvements) {
	pool = make(map[int]Improvements, len(curr.Species))
	for _, genome := range curr.Genomes {
		pool[genome.SpeciesIdx] = append(pool[genome.SpeciesIdx], neat.CopyGenome(genome))
	}
	for idx, list := range pool {
		sort.Sort(sort.Reverse(list))
		pool[idx] = list
	}
	return
}

// Removes stagnant species from the pool of possible parents. Allow the species with the most fit
// genome to continue past stagnation.
//
// TODO: Add setting so that user can control whether species with best is removed if stagnant for too long
func purgeSpecies(cfg ClassicSettings, species []neat.Species, pool map[int]Improvements) {

	// Update the species' adjusted fitness and stagnation level and note the best
	max := -1.0
	best := -1
	remove := make([]int, 0, len(species))
	for i, s := range species {

		// Update stagnation and fitness
		l := pool[i]
		f := l.Improvement()
		if f <= s.Improvement {
			species[i].Stagnation += 1
		} else {
			species[i].Stagnation = 0
			species[i].Improvement = f
		}

		// Plan to remove stagnant species
		if species[i].Stagnation > cfg.MaxStagnation() {
			remove = append(remove, i)
		} else {
			// Trim species to just most fit members
			cnt := int(math.Max(1.0, float64(len(l))*cfg.SurvivalThreshold()))
			pool[i] = l[:cnt]
		}

		// Use the same loop to note the species with the best
		// Should this be Improvement instead of Fitness?
		if l[0].Fitness > max {
			max = l[0].Fitness
			best = i
		}
	}

	// Remove any stagnant species
	for _, idx := range remove {
		if idx != best {
			delete(pool, idx)
		}
	}
}

func createOffspring(ctx neat.Context, cfg ClassicSettings, cross bool, rng *rand.Rand, pool map[int]Improvements, cnts map[int]int, next *neat.Population) (err error) {
	var child neat.Genome
	for idx, cnt := range cnts {
		l := pool[idx]
		for i := 0; i < cnt; i++ {
			p1, p2 := pickParents(cfg, cross, rng, l, pool)
			if p1.ID == p2.ID {
				child = neat.CopyGenome(p1)
			} else {
				child, err = ctx.Crosser().Cross(p1, p2)
				if err != nil {
					return
				}
			}
			child.ID = ctx.NextID()
			child.Birth = next.Generation
			err = ctx.Mutator().Mutate(&child)
			next.Genomes = append(next.Genomes, child)
		}
	}
	return
}

func pickParents(cfg ClassicSettings, cross bool, rng *rand.Rand, species Improvements, pool map[int]Improvements) (p1, p2 neat.Genome) {

	// Parent 1 comes from the species
	i := rng.Intn(len(species))
	p1 = species[i]

	if !cross || rng.Float64() < cfg.MutateOnlyProbability() { // Offspring is mutate only -- comes from one parent
		p2 = p1
	} else {
		if rng.Float64() < cfg.InterspeciesMatingRate() { // Offspring could come from any species
			for _, l := range pool {
				species = l
			}
		}
		i = rng.Intn(len(species))
		p2 = species[i]
	}

	return
}

type Improvements []neat.Genome

func (f Improvements) Len() int { return len(f) }
func (f Improvements) Less(i, j int) bool {
	if f[i].Improvement == f[j].Improvement {
		return f[i].Complexity() < f[j].Complexity()
	} else {
		return f[i].Improvement < f[j].Improvement
	}
}
func (f Improvements) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// Returns the average fitness for the list
//
// As the reproduction mechanism, NEAT uses explicit fitness sharing (Goldberg and Richard- son
// 1987), where organisms in the same species must share the fitness of their niche. Thus, a
// species cannot afford to become too big even if many of its organisms perform well. Therefore,
// any one species is unlikely to take over the entire population, which is crucial for speciated
// evolution to support a variety of topologies. The adjusted fitness fi′ for organism i is
// calculated according to its distance δ from every other organism j in the population:
//   See figure 3.2
// The sharing function sh is set to 0 when distance δ(i,j) is above the threshold δt; otherwise,
// sh(δ(i, j)) is set to 1 (Spears 1995). Thus, 􏰇nj=1 sh(δ(i, j)) reduces to the number of
// organisms in the same species as organism i. This reduction is natural since species are already
// clustered by compatibility using the threshold δt. (Stanley, 39-40)
//
// NOTE: This works out to just taking the average of the improvements
func (f Improvements) Improvement() float64 {
	if len(f) == 0 {
		return 0
	} else {
		var sum float64
		for _, genome := range f {
			sum += genome.Improvement
		}
		return sum / float64(len(f))
	}
}

type SpeciesList []*neat.Species

func (f SpeciesList) Len() int           { return len(f) }
func (f SpeciesList) Less(i, j int) bool { return f[i].Improvement < f[j].Improvement }
func (f SpeciesList) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
