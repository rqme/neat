package starter

import (
	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
)

type Context struct {
	// Component helpers
	arc neat.Archiver
	cmp neat.Comparer
	crs neat.Crosser
	dec neat.Decoder
	evl neat.Evaluator
	gen neat.Generator
	mut neat.Mutator
	src neat.Searcher
	spc neat.Speciater
	vis neat.Visualizer

	Settings
	state map[string]interface{}
	identify
}

func NewContext() *Context {
	return &Context{
		state:    make(map[string]interface{}),
		identify: newIdentify(),
	}
}

// Main context methods
func (c Context) Archiver() neat.Archiver       { return c.arc }
func (c Context) Comparer() neat.Comparer       { return c.cmp }
func (c Context) Crosser() neat.Crosser         { return c.crs }
func (c Context) Decoder() neat.Decoder         { return c.dec }
func (c Context) Evaluator() neat.Evaluator     { return c.evl }
func (c Context) Generator() neat.Generator     { return c.gen }
func (c Context) Mutator() neat.Mutator         { return c.mut }
func (c Context) Searcher() neat.Searcher       { return c.src }
func (c Context) Speciater() neat.Speciater     { return c.spc }
func (c Context) Visualizer() neat.Visualizer   { return c.vis }
func (c Context) State() map[string]interface{} { return c.state }

// Experiment settings
func (c Context) Iterations() int               { return c.Settings.Iterations }
func (c Context) Traits() neat.Traits           { return c.Settings.Traits }
func (c Context) FitnessType() neat.FitnessType { return c.Settings.FitnessType }
func (c Context) ExperimentName() string        { return c.Settings.ExperimentName }

// File archiver settings
func (c Context) ArchivePath() string { return c.Settings.ArchivePath }
func (c Context) ArchiveName() string { return c.Settings.ArchiveName }

// Classic comparer settings
func (c Context) DisjointCoefficient() float64 { return c.Settings.DisjointCoefficient }
func (c Context) ExcessCoefficient() float64   { return c.Settings.ExcessCoefficient }
func (c Context) WeightCoefficient() float64   { return c.Settings.WeightCoefficient }

// Classic crosser settings
func (c Context) EnableProbability() float64          { return c.Settings.EnableProbability }
func (c Context) MateByAveragingProbability() float64 { return c.Settings.MateByAveragingProbability }

// HyperNEAT decoder settings
func (c Context) SubstrateLayers() []decoder.SubstrateNodes { return c.Settings.SubstrateLayers }

// Classic generator settings
func (c Context) PopulationSize() int                   { return c.Settings.PopulationSize }
func (c Context) SeedGenome() neat.Genome               { return c.Settings.SeedGenome }
func (c Context) NumInputs() int                        { return c.Settings.NumInputs }
func (c Context) NumOutputs() int                       { return c.Settings.NumOutputs }
func (c Context) OutputActivation() neat.ActivationType { return c.Settings.OutputActivation }
func (c Context) WeightRange() float64                  { return c.Settings.WeightRange }
func (c Context) SurvivalThreshold() float64            { return c.Settings.SurvivalThreshold }
func (c Context) MutateOnlyProbability() float64        { return c.Settings.MutateOnlyProbability }
func (c Context) InterspeciesMatingRate() float64       { return c.Settings.InterspeciesMatingRate }
func (c Context) MaxStagnation() int                    { return c.Settings.MaxStagnation }

// Classic mutator settings
func (c Context) MutateActivationProbability() float64  { return c.Settings.MutateActivationProbability }
func (c Context) MutateWeightProbability() float64      { return c.Settings.MutateWeightProbability }
func (c Context) ReplaceWeightProbability() float64     { return c.Settings.ReplaceWeightProbability }
func (c Context) MutateTraitProbability() float64       { return c.Settings.MutateTraitProbability }
func (c Context) ReplaceTraitProbability() float64      { return c.Settings.ReplaceTraitProbability }
func (c Context) MutateSettingProbability() float64     { return c.Settings.MutateSettingProbability }
func (c Context) ReplaceSettingProbability() float64    { return c.Settings.ReplaceSettingProbability }
func (c Context) AddNodeProbability() float64           { return c.Settings.AddNodeProbability }
func (c Context) AddConnProbability() float64           { return c.Settings.AddConnProbability }
func (c Context) AllowRecurrent() bool                  { return c.Settings.AllowRecurrent }
func (c Context) HiddenActivation() neat.ActivationType { return c.Settings.HiddenActivation }
func (c Context) DelNodeProbability() float64           { return c.Settings.DelNodeProbability }
func (c Context) DelConnProbability() float64           { return c.Settings.DelConnProbability }

// Phased mutator settings
func (c Context) PruningPhaseThreshold() float64    { return c.Settings.PruningPhaseThreshold }
func (c Context) MaxMPCAge() int                    { return c.Settings.MaxMPCAge }
func (c Context) MaxImprovementAge() int            { return c.Settings.MaxImprovementAge }
func (c Context) ImprovementType() neat.FitnessType { return c.Settings.ImprovementType }

// Novelty searcher settings
func (c Context) NoveltyEvalArchive() bool         { return c.Settings.NoveltyEvalArchive }
func (c Context) NoveltyArchiveThreshold() float64 { return c.Settings.NoveltyArchiveThreshold }
func (c Context) NumNearestNeighbors() int         { return c.Settings.NumNearestNeighbors }

// Classic speciate settings
func (c Context) CompatibilityThreshold() float64      { return c.Settings.CompatibilityThreshold }
func (c *Context) SetCompatibilityThreshold(v float64) { c.Settings.CompatibilityThreshold = v }
func (c Context) TargetNumberOfSpecies() int           { return c.Settings.TargetNumberOfSpecies }
func (c Context) CompatibilityModifier() float64       { return c.Settings.CompatibilityModifier }

// Web visualizer settings
func (c Context) WebPath() string { return c.Settings.WebPath }
