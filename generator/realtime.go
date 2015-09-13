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

	"github.com/rqme/neat"
)

type RealTimeSettings interface {
	ClassicSettings
	IneligiblePercent() float64 // Fraction of the population ineligible for replacement
	MinimumTimeAlive() int      // Minimum number of ticks before considered for replacement
}

type RealTime struct {
	RealTimeSettings
	ctx neat.Context

	tick    int
	replace int
	cross   bool
}

func (g *RealTime) SetContext(x neat.Context) error {
	g.ctx = x
	return nil
}

func (g *RealTime) SetCrossover(v bool) error {
	g.cross = v
	return nil
}

func (g *RealTime) Generate(curr neat.Population) (next neat.Population, err error) {
	if len(curr.Genomes) == 0 {
		next, err = generateFirst(g.ctx, g.RealTimeSettings)
	} else {
		next, err = g.generateNext(curr)
	}
	return
}

func (g *RealTime) generateNext(curr neat.Population) (next neat.Population, err error) {

	// Increment the ticks
	g.tick += 1

	// Determine ticks between replacement
	m := g.MinimumTimeAlive()
	n := int(float64(m) / (float64(g.PopulationSize()) * g.IneligiblePercent()))

	// No replacement this time
	if g.tick%n != 0 {
		next = curr
		return
	}

	// Create the pool
	pool := createPool(curr)

	// 1. Remove the agent with the worst adjusted fitness from the population assuming one has beeen
	// alive sufficiently long so it has been properley evaluated.
	if !g.removeWorst(pool, n, m) {
		next = curr
		return
	}
	// 2. Re-estimate F for all species.
	ftot := g.reestimate(pool)
	purgeSpecies(g.RealTimeSettings, curr.Species, pool)

	// 3. Choose a parent species to create the new offspring
	rng := rand.New(rand.NewSource(rand.Int63()))
	cnts := make(map[int]int, 1)
	sidx := g.pickSpecies(pool, ftot, rng)
	cnts[sidx] = 1

	next.Generation = curr.Generation + 1
	next.Genomes = make([]neat.Genome, 0, len(curr.Genomes))
	if err = createOffspring(g.ctx, g.RealTimeSettings, g.cross, rng, pool, cnts, &next); err != nil {
		return
	}

	// 4. Adjust compatibility theshold Ct dynamically and reassign all agents to species
	for _, list := range pool {
		next.Genomes = append(next.Genomes, list...)
	}
	next.Species, err = g.ctx.Speciater().Speciate(curr.Species, next.Genomes)

	// 5. Place the new agent in the world
	return
}

// Seciont 3.1.1 Step 1: Removing the worst agent (Stanley, p.3)
func (g *RealTime) removeWorst(pool map[int]Improvements, n, m int) bool {
	var worst float64 = math.Inf(1)
	var wg, ws int
	ws = -1
	for i, list := range pool {
		for j := len(list) - 1; j >= 0; j-- {
			adj := list[j].Improvement / float64(len(list))
			if list[j].Birth*n > m && adj < worst {
				worst = adj
				ws = i
				wg = j
				break
			}
		}
	}
	if ws == -1 { // Entire population is too young to be replaced
		return false
	}
	list := pool[ws]
	list = append(list[:wg], list[wg+1:]...)
	pool[ws] = list
	return true
}

// Section 3.1.2 Step 2: Re-estimagting F (Stanley, p.4)
func (g *RealTime) reestimate(pool map[int]Improvements) float64 {
	ftot := 0.0
	for _, list := range pool {
		ftot += list.Improvement()
	}
	return ftot
}

// Section 3.1.3 Step 3: Choosing the parent species (Stanley, p.4)
func (g *RealTime) pickSpecies(pool map[int]Improvements, ftot float64, rng *rand.Rand) int {
	ftgt := rng.Float64() * ftot
	fsum := 0.0
	for i, list := range pool {
		fsum += list.Improvement()
		if fsum >= ftgt {
			return i
		}
	}
	return -1 // Shouldn't get here
}
