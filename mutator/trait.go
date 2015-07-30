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

type Trait struct {

	// Repository of traits
	Traits neat.Traits "neat:config"

	// Probability that the trait will be mutated
	MutateTraitProbability float64 "neat:config"

	// Probability that the trait will be replaced
	ReplaceTraitProbability float64 "neat:config"

	// Probability that the setting will be mutated
	MutateSettingProbability float64 "neat:config"

	// Probability that the setting will be replaced
	ReplaceSettingProbability float64 "neat:config"
}

// Configures the helper from a JSON string
func (m *Trait) Configure(cfg string) error {
	return neat.Configure(cfg, m)
}

// Mutates a genome's traits
func (m Trait) Mutate(g *neat.Genome) error {
	rng := rand.New(rand.NewSource(rand.Int63()))
	for t, trait := range m.Traits {
		if trait.IsSetting {
			if rng.Float64() < m.MutateSettingProbability {
				if rng.Float64() < m.ReplaceSettingProbability {
					g.Traits[t] = m.replaceTrait(rng, trait)
				} else {
					g.Traits[t] = m.mutateTrait(rng, trait, g.Traits[t])
				}
			}
		} else {
			if rng.Float64() < m.MutateSettingProbability {
				if rng.Float64() < m.ReplaceSettingProbability {
					g.Traits[t] = m.replaceTrait(rng, trait)
				} else {
					g.Traits[t] = m.mutateTrait(rng, trait, g.Traits[t])
				}
			}
		}
	}
	return nil
}

// Returns a modified trait
func (m Trait) mutateTrait(rng *rand.Rand, trait neat.Trait, v float64) float64 {
	v = v + rng.NormFloat64()
	if v < trait.Min {
		v = trait.Min
	} else if v > trait.Max {
		v = trait.Max
	}
	return v
}

// Returns a new trait
func (m Trait) replaceTrait(rng *rand.Rand, trait neat.Trait) float64 {
	return (rng.Float64() * (trait.Max - trait.Min)) + trait.Min
}
