package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/rqme/neat"
	"github.com/rqme/neat/archiver"
	"github.com/rqme/neat/generator"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/visualizer"
	"github.com/rqme/neat/x/example"
	"github.com/rqme/neat/x/trials"
)

var (
	Phased   = flag.Bool("phased", false, "Use the phased mutator during evolution")
	Duration = flag.Int("duration", 90, "Maximum duration of each trial")
)

type evaluator struct {
	stopTime time.Time
}

func (e evaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	// Iterate the inputs and ask about
	sum := 0.0
	cnt := 0.0
	stop := true
	for i, input := range inputs {
		// Query the network for this letter
		outputs, err := p.Activate(input)
		if err != nil {
			return example.NewResult(p.ID(), 0, err, false)
		}

		// Identify the max
		max, _ := stats.Max(outputs)

		// Determine success
		for j := 0; j < len(outputs); j++ {
			if outputs[j] == max && j != i {
				stop = false // picked another letter
			}

			if j == i {
				sum += 1.0 - outputs[j]
			} else {
				sum += outputs[j]
			}
			cnt += 1.0
		}
	}

	return example.NewResult(p.ID(), math.Pow(cnt-sum, 2), nil, stop || !time.Now().Before(e.stopTime))
}

func main() {
	flag.Parse()
	if *Phased {
		fmt.Println("Using phased mutator")
	} else {
		fmt.Println("Using complexifying-only mutator")
	}
	fmt.Println("Each trial will run for a maximum of", *Duration, "minutes.")
	if err := trials.Run(func(i int) (*neat.Experiment, error) {
		e, err := example.NewExperiment(
			&evaluator{stopTime: time.Now().Add(time.Minute * time.Duration(*Duration))},
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

			func(e *neat.Experiment) error {
				if *Phased {
					g, _ := e.Generator.(*generator.Classic)
					g.Mutator = &mutator.Complete{}
				}
				return nil
			},
		)
		return e, err
	}); err != nil {
		log.Fatal("Could not run OCR: ", err)
	}

}
