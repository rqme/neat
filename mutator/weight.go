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

type Weight struct {

	// The mutation range of the weight. If x, range is [-x,x]
	WeightRange float64 "neat:config"

	// Probability that the weight will be mutated
	MutateWeightProbability float64 "neat:config"

	// Probability that the weight will be replaced
	ReplaceWeightProbability float64 "neat:config"
}

// Configures the helper from a JSON string
func (m *Weight) Configure(cfg string) error {
	return neat.Configure(cfg, m)
}

// Mutates a genome's weights
func (m Weight) Mutate(g *neat.Genome) error {
	rng := rand.New(rand.NewSource(rand.Int63()))
	for k, conn := range g.Conns {
		if rng.Float64() < m.MutateWeightProbability {
			if rng.Float64() < m.ReplaceWeightProbability {
				conn.Weight = m.replaceWeight(rng)
			} else {
				conn.Weight = m.mutateWeight(rng, conn.Weight)
			}
			g.Conns[k] = conn
		}
	}
	return nil
}

// Returns a modified weight
func (m Weight) mutateWeight(rng *rand.Rand, w float64) float64 {
	w = w + rng.NormFloat64()
	/*
		if w < -m.WeightRange*2 {
			w = -m.WeightRange * 2
		} else if w > m.WeightRange*2 {
			w = m.WeightRange * 2
		}
	*/
	return w
}

// Returns a new weight
func (m Weight) replaceWeight(rng *rand.Rand) float64 {
	return (rng.Float64()*2.0 - 1.0) * m.WeightRange
}
