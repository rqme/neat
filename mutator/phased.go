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
	"github.com/montanaflynn/stats"
	"github.com/rqme/neat"
)

// Phased mutatation settings
type PhasedSettings interface {
	PruningPhaseThreshold() float64
	MaxMPCAge() int
	MaxImprovementAge() int
	ImprovementType() neat.FitnessType
}

// Phased Searching
// As an alternative to blended searching I propose the use of 'phased' searching, so called because
// the NEAT search switches between a complexifying phase and a simplifying(or pruning) phase.
//
// Phased searching works as follows:
//
// 1) Before the search starts proper, calculate a threshold at which the prune phase will begin. This
//    threshold is the current mean population complexity (MPC)** plus a specified pruning phase
//    threshold value which typically might be between 30-100 depending on the type of experiment.
// 2) The search begins in complexifying mode and continues as traditional NEAT until the prune phase
//    threshold is reached.
// 3) The search now enters a prune phase. The prune phase is almost algorithmically identical to the
//    complexifying phase, and the normal process of selecting genomes for reproduction, generating
//    offspring, monitoring species compatibility, etc. all takes place. The difference is that the
//    additive mutations are disabled and subtractive ones enabled in their place. In addition only
//    asexual reproduction (with mutation) is allowed. Crossover is disabled because this can allow
//    genes to propagate through a population, thus increasing complexity.
// 4) During each generation of the pruning phase a reading of the MPC is taken, this will normally
//    be seen to fall as functionally redundant structures are removed from the population. As pruning
//    progresses the MPC will eventually reach a floor level when no more redundant structure remains
//    in the population to be removed. Therefore once MPC has not fallen for some number of generations
//    (this is configurable, between 10-50 works well), we can reset the next pruning phase's threshold
//    to be the current MPC floor level + the pruning phase threshold parameter and switch into a
//    complexifying phase. The whole process then begins again at (2).
//
// One small modification was made to the above process, and that is to not enter prune phase unless
// the population fitness has not risen for some specified number of generations. Therefore if the
// complexity has risen past the pruning phase threshold but the population fitness is still rising
// then we hold off the pruning phase until the fitness stops rising.
//
// from Colin Green (http://sharpneat.sourceforge.net/phasedsearch.html)
type Phased struct {
	PhasedSettings
	ctx neat.Context

	// Inner mutators
	Complexify
	Pruning

	// internal state
	isPruning               bool
	pruneThresh             float64
	targetMPC, fitness      float64
	ageMPC, ageImprovement  int
	minMPC, lastImprovement float64
}

func NewPhased(ps PhasedSettings, cs ComplexifySettings, ns PruningSettings) *Phased {
	return &Phased{
		PhasedSettings: ps,
		Complexify:     Complexify{ComplexifySettings: cs},
		Pruning:        Pruning{PruningSettings: ns},
	}
}

func (m *Phased) SetContext(x neat.Context) error {
	m.ctx = x
	return m.Complexify.SetContext(x)
}

// Mutates the Genome by through complexifiying or pruning depending on current phase
func (m *Phased) Mutate(g *neat.Genome) (err error) {
	if m.isPruning {
		err = m.Pruning.Mutate(g)
	} else {
		err = m.Complexify.Mutate(g)
	}
	return
}

// Updates the statistics of the population and determines if a phase switch is required.
func (m *Phased) SetPopulation(p neat.Population) error {

	// Calculate the tnew MPC and fitness
	var n float64
	fit := make([]float64, len(p.Genomes))
	for i, g := range p.Genomes {
		n += float64(g.Complexity())
		fit[i] = g.Improvement
	}
	mpc := n / float64(len(p.Genomes))
	//neat.DBG("mpc %f fit %f", mpc, fit)
	if mpc < m.minMPC { // Looking for a drop in MPC
		m.ageMPC = 0
		m.minMPC = mpc
	} else {
		m.ageMPC += 1
	}

	var f float64
	if m.ImprovementType() == neat.Absolute {
		f, _ = stats.Max(fit)
	} else {
		f, _ = stats.VarP(fit)
	}

	// Looking for a continued increase in fitness (Absolute) or an uptick in variance (RelativeImprovement)
	if f > m.lastImprovement {
		m.ageImprovement = 0
		m.lastImprovement = f
	} else {
		m.ageImprovement += 1
	}

	// First run, just set the initial threshold and return
	if m.targetMPC == 0 {
		m.isPruning = false
		m.targetMPC = mpc + m.PruningPhaseThreshold()
		return nil
	}

	// Check for a phase change
	if m.isPruning {
		if m.ageMPC > m.MaxMPCAge() {
			m.isPruning = false
			m.targetMPC = mpc + m.PruningPhaseThreshold()
			m.ageImprovement = 0
			m.lastImprovement = 0
		}
	} else {
		if mpc >= m.targetMPC && m.ageImprovement > m.MaxImprovementAge() {
			m.isPruning = true
			m.ageMPC = 0
			m.minMPC = mpc
		}
	}

	// Toggle crossover as necessary
	if crs, ok := m.ctx.(neat.Crossoverable); ok {
		crs.SetCrossover(!m.isPruning)
	}
	return nil
}
