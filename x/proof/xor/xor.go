package main

import (
	"flag"
	"log"
	"math"
	"path"
	"strconv"

	"github.com/rqme/neat"
	"github.com/rqme/neat/archiver"
	"github.com/rqme/neat/visualizer"
	"github.com/rqme/neat/x/example"
	"github.com/rqme/neat/x/trials"
)

type evaluator struct{}

// Evaluate computes the error for the XOR problem with the phenome
//
// To compute fitness, the distance of the output from the correct answer was summed for all four
// input patterns. The result of this error was subtracted from 4 so that higher fitness would mean
// better networks. The resulting number was squared to give proportionally more fitness the closer
// a network was to a solution. (Stanley, 43)
func (e evaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	inputs := [][]float64{
		[]float64{0, 0},
		[]float64{0, 1},
		[]float64{1, 0},
		[]float64{1, 1},
	}

	expected := []float64{0, 1, 1, 0}

	// Run experiment
	var err error
	var sum float64
	stop := true
	for i, in := range inputs {
		outputs, err := p.Activate(in)
		if err != nil {
			break
		}
		sum += math.Abs(outputs[0] - expected[i])
		if expected[i] == 0 {
			stop = stop && outputs[0] < 0.5
		} else {
			stop = stop && outputs[0] > 0.5
		}
	}

	// Calculate the result
	r = example.NewResult(p.ID(), math.Pow(4.0-sum, 2.0), err, stop)
	return
}

func main() {
	flag.Parse()
	if err := trials.Run(10, true, func(i int) (*neat.Experiment, error) {
		e, err := example.NewExperiment(
			&evaluator{},
			func(e *neat.Experiment) error {
				a, _ := e.Archiver.(*archiver.File)
				a.ArchivePath = path.Join(a.ArchivePath, strconv.Itoa(i))
				e.Archiver = a
				return nil
			},
			func(e *neat.Experiment) error {
				v, _ := e.Visualizer.(*visualizer.Web)
				v.WebPath = path.Join(v.WebPath, strconv.Itoa(i))
				e.Visualizer = v
				return nil
			},
		)
		return e, err
	}); err != nil {
		log.Fatal("Could not run XOR: ", err)
	}

}
