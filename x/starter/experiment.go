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
	"flag"
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

const (
	NoTrials int = -1
)

var (
	ConfigPath = flag.String("config-path", "", "Path to configuration file to override archive.")
	ConfigName = flag.String("config-name", "", "Name prepended to all configuration and state files")
)

type ConfigSettings struct {
	path, name string
}

func (s ConfigSettings) ArchiveName() string { return s.name }
func (s ConfigSettings) ArchivePath() string { return s.path }

func NewExperiment(ctx neat.Context, cfg neat.ExperimentSettings, t int) (exp *neat.Experiment, err error) {

	// Create the experiment
	exp = &neat.Experiment{ExperimentSettings: cfg}
	exp.SetContext(ctx)

	// Restore the saved setting and, if available, state
	if *ConfigName == "" {
		*ConfigName = os.Args[0] // Use the executable's name
	}
	rst := &archiver.File{
		FileSettings: ConfigSettings{path: *ConfigPath, name: *ConfigName},
	}
	if err = rst.Restore(ctx); err != nil {
		return
	}

	// Update helpers with trial number
	if t > NoTrials {
		hs := []interface{}{
			ctx.Archiver(),
			ctx.Comparer(),
			ctx.Crosser(),
			ctx.Decoder(),
			ctx.Evaluator(),
			ctx.Generator(),
			ctx.Mutator(),
			ctx.Searcher(),
			ctx.Speciater(),
			ctx.Visualizer(),
		}
		for _, h := range hs {
			if th, ok := h.(neat.Trialable); ok {
				th.SetTrial(t)
			}
		}
	}

	// Load ids and innovations
	if ph, ok := ctx.(neat.Populatable); ok {
		ph.SetPopulation(exp.Population())
	}
	return
}

func NewClassicContext(evl neat.Evaluator) *Context {
	// Create the context
	ctx := NewContext()
	if *ConfigName == "" {
		ctx.Settings.ArchiveName = os.Args[0]
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
	ctx.vis = &visualizer.Web{WebSettings: ctx}

	attachContext(ctx)
	return ctx
}

func NewNoveltyContext(evl neat.Evaluator) *Context {
	// Create the context
	ctx := NewContext()
	if *ConfigName == "" {
		ctx.Settings.ArchiveName = os.Args[0]
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
	ctx.src = &searcher.Novelty{NoveltySettings: ctx, Searcher: &searcher.Concurrent{}}
	ctx.spc = speciater.NewDynamic(ctx, ctx)
	ctx.vis = &visualizer.Web{WebSettings: ctx}

	attachContext(ctx)
	return ctx
}

func NewHyperNEATContext(evl neat.Evaluator) *Context {
	// Create the context
	ctx := NewContext()
	if *ConfigName == "" {
		ctx.Settings.ArchiveName = os.Args[0]
	} else {
		ctx.Settings.ArchiveName = *ConfigName // Can be overriden in settings file
	}
	ctx.arc = &archiver.File{FileSettings: ctx}
	ctx.cmp = &comparer.Classic{ClassicSettings: ctx}
	ctx.crs = &crosser.Classic{ClassicSettings: ctx}
	ctx.dec = &decoder.HyperNEAT{CppnDecoder: decoder.Classic{}, HyperNEATSettings: ctx}
	ctx.evl = evl
	ctx.gen = &generator.Classic{ClassicSettings: ctx}
	ctx.mut = mutator.NewComplete(ctx, ctx, ctx, ctx)
	ctx.src = &searcher.Concurrent{}
	ctx.spc = speciater.NewDynamic(ctx, ctx)
	ctx.vis = &visualizer.Web{WebSettings: ctx}

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
