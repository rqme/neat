package main

import (
	"github.com/rqme/neat"
	"github.com/rqme/neat/x/starter"
)

func initSettings() starter.Settings {
	return starter.Settings{
		// Experiment settings
		Iterations:  100,
		FitnessType: neat.Absolute,

		// Classic comparer settings
		DisjointCoefficient: 1.0,
		ExcessCoefficient:   1.0,
		WeightCoefficient:   1.0,

		// Classic crosser settings
		EnableProbability:          0.2,
		MateByAveragingProbability: 0.4,

		// Classic generator settings
		PopulationSize:         150,
		NumInputs:              2,
		NumOutputs:             1,
		OutputActivation:       neat.SteependSigmoid,
		HiddenActivation:       neat.SteependSigmoid,
		WeightRange:            2.5,
		SurvivalThreshold:      0.2,
		MutateOnlyProbability:  0.25,
		InterspeciesMatingRate: 0.001,
		MaxStagnation:          15,

		// Classic mutator settings
		MutateWeightProbability:  0.9,
		ReplaceWeightProbability: 0.2,
		AddNodeProbability:       0.125,
		AddConnProbability:       0.025,

		// Classic speciater settings
		CompatibilityThreshold: 3.0,
		TargetNumberOfSpecies:  15,
		CompatibilityModifier:  0.3,
	}
}
