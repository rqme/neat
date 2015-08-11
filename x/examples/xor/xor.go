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

package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

type Evaluator struct {
	show     bool
	useTrial bool
	trialNum int
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

// Evaluate computes the error for the XOR problem with the phenome
//
// To compute fitness, the distance of the output from the correct answer was summed for all four
// input patterns. The result of this error was subtracted from 4 so that higher fitness would mean
// better networks. The resulting number was squared to give proportionally more fitness the closer
// a network was to a solution. (Stanley, 43)
func (e Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {
	inputs := [][]float64{
		[]float64{0, 0},
		[]float64{0, 1},
		[]float64{1, 0},
		[]float64{1, 1},
	}

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

	// Display the work
	if e.show {
		fmt.Println()
		if e.useTrial {
			fmt.Printf("Trial %d ", e.trialNum)
		}
		fmt.Printf("XOR Evaluation for genome %d\n", p.ID())
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("For {0,0}, expected 0. output was %f\n", actual[0])
		fmt.Printf("For {0,1}, expected 1. output was %f\n", actual[1])
		fmt.Printf("For {1,0}, expected 1. output was %f\n", actual[2])
		fmt.Printf("For {1,1}, expected 0. output was %f\n", actual[3])
	}

	// Calculate the result
	r = result.New(p.ID(), math.Pow(4.0-sum, 2.0), err, stop)
	return
}

func (e *Evaluator) ShowWork(s bool) {
	e.show = s
}

func main() {
	flag.Parse()
	//defer profile.Start(profile.CPUProfile).Stop()
	if err := trials.Run(func(i int) (*neat.Experiment, error) {
		ctx := starter.NewClassicContext(&Evaluator{})
		if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
			return nil, err
		} else {
			return exp, nil
		}

	}); err != nil {
		log.Fatal("Could not run XOR: ", err)
	}

}
