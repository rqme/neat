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

package comparer

import (
	"github.com/rqme/neat"

	"math"
)

type ClassicSettings interface {
	DisjointCoefficient() float64
	ExcessCoefficient() float64
	WeightCoefficient() float64
}

// Helper to compare two genomes similarity
type Classic struct {
	ClassicSettings
}

// Compares two genomes and returns their compatibility distance
//
// The number of excess and disjoint genes between a pair of genomes is a natural measure of their
// compatibility distance. The more disjoint two genomes are, the less evolutionary history they
// share, and thus the less compatible they are. Therefore, we can measure the compatibility distance δ
// of different structures in NEAT as a simple linear combination of the number of excess E and
// disjoint D genes, as well as the average weight differences of matching genes W, including disabled genes:
//     see formula in paper
// The coefﬁcients c1, c2,and c3 allow us to adjust the importance of the three factors, and thefactor N,
// the number of genes in the larger genome, normalizes for genome size (N can be set to 1 if both genomes
// are small, i.e., consist of fewer than 20 genes (Stanley, 110)
func (c Classic) Compare(g1, g2 neat.Genome) (float64, error) {

	// Determine N
	n := 1.0
	if len(g1.Conns) > len(g2.Conns) && len(g1.Conns) >= 20 {
		n = float64(len(g1.Conns))
	} else if len(g2.Conns) >= 20 {
		n = float64(len(g2.Conns))
	}

	// Ensure connections are sorted
	_, conns1 := g1.GenesByInnovation()
	_, conns2 := g2.GenesByInnovation()

	// Calculate the components. This assumes both genomes' connections are sorted by their
	// innovation number (which is true if the NEAT library created them)
	var d, e, w, x float64
	i := 0
	j := 0
	for i < len(conns1) || j < len(conns2) {
		switch {
		case i == len(conns1):
			e += 1
			j += 1
		case j == len(conns2):
			e += 1
			i += 1
		default:
			c1 := conns1[i]
			c2 := conns2[j]
			switch {
			case c1.Innovation < c2.Innovation:
				d += 1
				i += 1
			case c1.Innovation > c2.Innovation:
				d += 1
				j += 1
			default: // Same innovation number
				w += math.Abs(c1.Weight - c2.Weight)
				x += 1
				i += 1
				j += 1
			}
		}
	}

	// Return the compatibility distance
	n = 1 // NOTE: The variable N mentioned in the paper does not seem to be used in any implemenation
	δ := c.ExcessCoefficient()*e/n + c.DisjointCoefficient()*d/n
	if x > 0 {
		δ += c.WeightCoefficient() * w / x
	}
	return δ, nil
}
