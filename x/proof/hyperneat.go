package main

import (
	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/x/starter"
)

type SettingsWithLayers struct {
	decoder.ESHyperNEATSettings
	layers []decoder.SubstrateNodes
}

func (h SettingsWithLayers) SubstrateLayers() []decoder.SubstrateNodes { return h.layers }

func hyperneatSettings(cfg decoder.ESHyperNEATSettings) (hns SettingsWithLayers) {
	hns.ESHyperNEATSettings = cfg

	hns.layers = make([]decoder.SubstrateNodes, 3)
	hns.layers[0] = []decoder.SubstrateNode{
		{Position: []float64{-1.0, 0.0}, NeuronType: neat.Input},
		{Position: []float64{1.0, 0.0}, NeuronType: neat.Input},
	}
	hns.layers[1] = []decoder.SubstrateNode{
		{Position: []float64{-1.0, 0.0}, NeuronType: neat.Hidden},
		{Position: []float64{1.0, 0.0}, NeuronType: neat.Hidden},
	}
	hns.layers[2] = []decoder.SubstrateNode{
		{Position: []float64{0.0, 0.0}, NeuronType: neat.Output},
	}
	return
}

func hyperneatContext() *starter.Context {
	cfg := initSettings()
	cfg.ExperimentName = "HyperNEAT"
	cfg.ArchivePath = "./proof-out/hyperneat"
	cfg.ArchiveName = "hyperneat"
	cfg.WebPath = cfg.ArchivePath
	cfg.MutateActivationProbability = 0.25
	cfg.OutputActivation = neat.Tanh
	cfg.NumInputs = 4
	cfg.NumOutputs = 2

	ctx := starter.NewContext(&NEATEval{}, func(ctx *starter.Context) {
		ctx.SetMutator(mutator.NewComplete(ctx, ctx, ctx, ctx, ctx, ctx))
		ctx.SetDecoder(&decoder.HyperNEAT{CppnDecoder: decoder.Classic{}, HyperNEATSettings: hyperneatSettings(ctx)})
	})
	ctx.Settings = cfg
	return ctx
}
