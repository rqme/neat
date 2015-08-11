package starter

import (
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
