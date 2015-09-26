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
	"os"

	"github.com/rqme/neat"
	"github.com/rqme/neat/archiver"
	"github.com/rqme/neat/comparer"
	"github.com/rqme/neat/crosser"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/generator"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/searcher"
	"github.com/rqme/neat/speciater"
	"github.com/rqme/neat/visualizer"
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

func NewContext(evl neat.Evaluator, options ...func(*Context)) *Context {

	// Create the context
	ctx := &Context{
		state: make(map[string]interface{}),
		identify: identify{
			innos: make(map[innovation]int, 100),
		},
	}

	// Set the default helpers
	if ctx.Settings.ArchivePath == "" && *ConfigName == "" {
		if len(os.Args) > 0 {
			ctx.Settings.ArchiveName = os.Args[0]
		} else {
			ctx.Settings.ArchiveName = "experiment"
		}
	} else {
		ctx.Settings.ArchiveName = *ConfigName // Can be overriden in settings file
	}
	ctx.arc = &archiver.File{FileSettings: ctx}
	ctx.cmp = &comparer.Classic{ClassicSettings: ctx}
	ctx.crs = &crosser.Classic{ClassicSettings: ctx}
	ctx.dec = &decoder.Classic{}
	ctx.evl = evl
	ctx.gen = &generator.Classic{ClassicSettings: ctx}
	ctx.mut = mutator.New(ctx, ctx, ctx)
	ctx.src = &searcher.Concurrent{}
	ctx.spc = speciater.NewDynamic(ctx, ctx)
	//ctx.spc = &speciater.Classic{ClassicSettings: ctx}
	ctx.vis = &visualizer.Web{WebSettings: ctx}

	// Override with the options
	for _, option := range options {
		option(ctx)
	}

	// Connect everything up
	attachContext(ctx)
	return ctx
}

func attachContext(ctx *Context) {
	for _, h := range []interface{}{ctx.arc, ctx.cmp, ctx.crs, ctx.dec, ctx.evl, ctx.gen, ctx.mut, ctx.src, ctx.spc, ctx.vis} {
		if ch, ok := h.(neat.Contextable); ok {
			ch.SetContext(ctx)
		}
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

func (c *Context) SetArchiver(h neat.Archiver)     { c.arc = h }
func (c *Context) SetComparer(h neat.Comparer)     { c.cmp = h }
func (c *Context) SetCrosser(h neat.Crosser)       { c.crs = h }
func (c *Context) SetDecoder(h neat.Decoder)       { c.dec = h }
func (c *Context) SetEvaluator(h neat.Evaluator)   { c.evl = h }
func (c *Context) SetGenerator(h neat.Generator)   { c.gen = h }
func (c *Context) SetMutator(h neat.Mutator)       { c.mut = h }
func (c *Context) SetSearcher(h neat.Searcher)     { c.src = h }
func (c *Context) SetSpeciater(h neat.Speciater)   { c.spc = h }
func (c *Context) SetVisualizer(h neat.Visualizer) { c.vis = h }

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

// ESHyperNEAT decoder settings
func (c Context) InitialDepth() int          { return c.Settings.InitialDepth }
func (c Context) MaxDepth() int              { return c.Settings.MaxDepth }
func (c Context) DivisionThreshold() float64 { return c.Settings.DivisionThreshold }
func (c Context) VarianceThreshold() float64 { return c.Settings.VarianceThreshold }
func (c Context) BandThreshold() float64     { return c.Settings.BandThreshold }
func (c Context) IterationLevels() int       { return c.Settings.IterationLevels }

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

// Real-Time generator settings
func (c Context) IneligiblePercent() float64 { return c.Settings.IneligiblePercent }
func (c Context) MinimumTimeAlive() int      { return c.Settings.MinimumTimeAlive }

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
