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
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/rqme/neat"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

var (
	Phased   = flag.Bool("phased", false, "Use the phased mutator during evolution")
	Duration = flag.Int("duration", 90, "Maximum duration in minutes of each trial")
)

type Evaluator struct {
	stopTime time.Time
	show     bool
	useTrial bool
	trialNum int
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

func (e *Evaluator) ShowWork(s bool) {
	e.show = s
}

func (e Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	// Iterate the inputs and ask about
	letters := []uint8{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	guesses := make([][]uint8, 26)
	values := make([]float64, 26)
	sum := 0.0
	cnt := 0.0
	stop := true
	for i, input := range inputs {
		// Query the network for this letter
		outputs, err := p.Activate(input)
		if err != nil {
			return result.New(p.ID(), 0, err, false)
		}

		// Identify the max
		max, _ := stats.Max(outputs)
		values[i] = max

		// Determine success
		s2 := 0.0
		for j := 0; j < len(outputs); j++ {
			if outputs[j] == max {
				guesses[i] = append(guesses[i], letters[j])
				if j != i {
					stop = false // picked another letter
					s2 += 1.0
				} else {
					s2 += 1.0 - max
				}
			}
			cnt += 1
		}
		sum += s2
	}

	if e.show {
		b := bytes.NewBufferString("\n")
		fmt.Println()
		if e.useTrial {
			b.WriteString(fmt.Sprintf("Trial %d ", e.trialNum))
		}
		b.WriteString(fmt.Sprintf("OCR Evaluation for genome %d. Letter->Guess(confidence)\n", p.ID()))
		b.WriteString(fmt.Sprintf("------------------------------------------\n"))
		for i := 0; i < len(letters); i++ {
			b.WriteString(fmt.Sprintf("%s (%0.2f)", string(letters[i]), values[i]))
			cl := ""
			il := ""
			for j := 0; j < len(guesses[i]); j++ {
				if guesses[i][j] == letters[i] {
					cl = "   correct"
				} else {
					if il != "" {
						il += ", "
					}
					il += string(guesses[i][j])
				}
			}
			if cl == "" {
				cl = " incorrect"
			}
			b.WriteString(cl)
			if il != "" {
				b.WriteString(" but also guessed ")
				b.WriteString(il)
			}
			b.WriteString("\n")
		}
		fmt.Println(b.String())
	}

	return result.New(p.ID(), math.Pow(cnt-sum, 2), nil, stop || !time.Now().Before(e.stopTime))
}

func main() {
	flag.Parse()
	//defer profile.Start(profile.CPUProfile).Stop()

	if *Phased {
		fmt.Println("Using phased mutator")
	} else {
		fmt.Println("Using complexifying-only mutator")
	}
	fmt.Println("Each trial will run for a maximum of", *Duration, "minutes.")

	if err := trials.Run(func(i int) (*neat.Experiment, error) {
		eval := &Evaluator{stopTime: time.Now().Add(time.Minute * time.Duration(*Duration))}
		var ctx *starter.Context
		if *Phased {
			ctx = starter.NewContext(eval, func(ctx *starter.Context) {
				ctx.SetMutator(mutator.NewComplete(ctx, ctx, ctx, ctx, ctx, ctx))
			})
		} else {
			ctx = starter.NewContext(eval)
		}
		if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
			return nil, err
		} else {
			return exp, nil
		}

	}); err != nil {
		log.Fatal("Could not run OCR: ", err)
	}
}
