package main

import (
	"math"

	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
)

type NEATEval struct{}

func (e *NEATEval) Evaluate(p neat.Phenome) (r neat.Result) {
	inputs := [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	expected := []float64{0, 1, 1, 0}
	actual := make([]float64, 4)

	// Run experiment
	var err error
	var sum float64
	stop := true
	for i, in := range inputs {
		outputs, err := p.Activate(in)
		if err != nil {
			break
		}
		actual[i] = outputs[0]
		sum += math.Abs(outputs[0] - expected[i])
		if expected[i] == 0 {
			stop = stop && outputs[0] < 0.5
		} else {
			stop = stop && outputs[0] > 0.5
		}
	}

	// Calculate the result
	r = result.New(p.ID(), math.Pow(4.0-sum, 2.0), err, stop)
	return
}

func neatContext() *starter.Context {
	cfg := initSettings()
	cfg.ExperimentName = "NEAT"
	cfg.ArchivePath = "./proof-out/neat"
	cfg.ArchiveName = "neat"
	cfg.WebPath = cfg.ArchivePath

	ctx := starter.NewContext(&NEATEval{})
	ctx.Settings = cfg
	return ctx
}
