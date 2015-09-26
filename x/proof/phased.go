package main

import (
	"github.com/rqme/neat"
	"github.com/rqme/neat/x/starter"
)

func phasedContext() *starter.Context {
	cfg := initSettings()
	cfg.ExperimentName = "Phased"
	cfg.ArchivePath = "./proof-out/phased"
	cfg.ArchiveName = "phased"
	cfg.WebPath = cfg.ArchivePath
	cfg.PruningPhaseThreshold = 10
	cfg.MaxMPCAge = 5
	cfg.MaxImprovementAge = 5
	cfg.ImprovementType = neat.Absolute
	cfg.DelNodeProbability = 0.015
	cfg.DelConnProbability = 0.025

	ctx := starter.NewContext(&NEATEval{})
	ctx.Settings = cfg
	return ctx
}
