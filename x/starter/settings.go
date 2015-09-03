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
package starter

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
)

type Settings struct {

	// Experiment settings
	Iterations     int
	Traits         neat.Traits
	FitnessType    neat.FitnessType
	ExperimentName string

	// File archiver settings
	ArchivePath string
	ArchiveName string

	// Classic comparer settings
	DisjointCoefficient float64
	ExcessCoefficient   float64
	WeightCoefficient   float64

	// Classic crosser settings
	EnableProbability          float64
	MateByAveragingProbability float64

	// HyperNEAT decoder settings
	SubstrateLayers []decoder.SubstrateNodes

	// Classic generator settings
	PopulationSize         int
	NumInputs              int
	NumOutputs             int
	OutputActivation       neat.ActivationType
	WeightRange            float64
	SurvivalThreshold      float64
	MutateOnlyProbability  float64
	InterspeciesMatingRate float64
	MaxStagnation          int
	SeedGenome             neat.Genome

	// Real-Time generator settings
	IneligiblePercent float64
	MinimumTimeAlive  int

	// Classic mutator settings
	MutateActivationProbability float64             // Probability that the node's activation will be mutated
	MutateWeightProbability     float64             // Probability that the weight will be mutated
	ReplaceWeightProbability    float64             // Probability that the weight will be replaced
	MutateTraitProbability      float64             // Probability that the trait will be mutated
	ReplaceTraitProbability     float64             // Probability that the trait will be replaced
	MutateSettingProbability    float64             // Probability that the setting will be mutated
	ReplaceSettingProbability   float64             // Probability that the setting will be replaced
	AddNodeProbability          float64             // Probablity a node will be added to the genome
	AddConnProbability          float64             // Probability a connection will be added to the genome
	AllowRecurrent              bool                // Allow recurrent connections to be added
	HiddenActivation            neat.ActivationType // Activation type to assign to new nodes
	DelNodeProbability          float64             // Probablity a node will be removed to the genome
	DelConnProbability          float64             // Probability a connection will be removed to the genome

	// Phased mutator settings
	PruningPhaseThreshold float64
	MaxMPCAge             int
	MaxImprovementAge     int
	ImprovementType       neat.FitnessType

	// Novelty searcher settings
	NoveltyEvalArchive      bool
	NoveltyArchiveThreshold float64
	NumNearestNeighbors     int

	// Classic speciater settings
	CompatibilityThreshold float64
	TargetNumberOfSpecies  int
	CompatibilityModifier  float64

	// Web visualizer settings
	WebPath string
}

func (s Settings) String() string {
	b := bytes.NewBufferString("Settings:\n")
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		b.WriteString(fmt.Sprintf("\t%s: %v\n", t.Field(i).Name, v.Field(i).Interface()))
	}
	return b.String()
}
