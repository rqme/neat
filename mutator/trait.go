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

package mutator

import (
	"github.com/rqme/neat"

	"math/rand"
)

type TraitSettings interface {
	Traits() neat.Traits                // Repository of traits
	MutateTraitProbability() float64    // Probability that the trait will be mutated
	ReplaceTraitProbability() float64   // Probability that the trait will be replaced
	MutateSettingProbability() float64  // Probability that the setting will be mutated
	ReplaceSettingProbability() float64 // Probability that the setting will be replaced
}

type Trait struct {
	TraitSettings
}

// Mutates a genome's traits
func (m Trait) Mutate(g *neat.Genome) error {
	rng := rand.New(rand.NewSource(rand.Int63()))
	ts := m.Traits()
	for i, _ := range g.Traits {
		t := ts[i]
		if t.IsSetting {
			if rng.Float64() < m.MutateTraitProbability() {
				if rng.Float64() < m.ReplaceTraitProbability() {
					m.replaceTrait(rng, t, &g.Traits[i])
				} else {
					m.mutateTrait(rng, t, &g.Traits[i])
				}
			}
		} else {
			if rng.Float64() < m.MutateTraitProbability() {
				if rng.Float64() < m.ReplaceTraitProbability() {
					m.replaceTrait(rng, t, &g.Traits[i])
				} else {
					m.mutateTrait(rng, t, &g.Traits[i])
				}
			}
		}
	}
	return nil
}

func (m *Trait) replaceTrait(rng *rand.Rand, t neat.Trait, v *float64) {
	*v = (rng.Float64() * (t.Max - t.Min)) + t.Min
}

func (m *Trait) mutateTrait(rng *rand.Rand, t neat.Trait, v *float64) {
	tv := *v
	tv += rng.NormFloat64()
	if tv < t.Min {
		tv = t.Min
	} else if tv > t.Max {
		tv = t.Max
	}
}
