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
	"math/rand"
	"os"
	"path"

	svg "github.com/ajstarks/svgo"
	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

var (
	NEAT        = flag.Bool("neat", false, "User the Classic NEAT decoder")
	HyperNEAT   = flag.Bool("hyperneat", false, "Use HyperNEAT decoder")
	ESHyperNEAT = flag.Bool("eshyperneat", false, "Use ESHyperNEAT decoder")
	Resolution  = flag.Int("resolution", 11, "The resolution of the field. Default is 11.")
	Cases       = flag.Int("cases", 75, "The number of cases to evaluate per phenome") // called trials in the paper
	WorkPath    = flag.String("work-path", ".", "Output directory for maze diagrams")
)

// This experiment is based on the Visual Discrimination experiment described in the paper
// "A Hypercube-Based Indirect Encoding for Evolving Large-Scale Neural Networks" by Stanley,
// D'Ambrosio and Guaci. It is available at http://eplex.cs.ucf.edu/publications/2009/stanley-alife09

type Case struct {
	X, Y   int
	Inputs []float64
}

func (c Case) String() string {
	r := *Resolution
	b := bytes.NewBufferString(" ")
	for x := 0; x < r; x++ {
		b.WriteString(fmt.Sprintf("%d", x%10))
	}
	b.WriteString("\n")
	for y := 0; y < r; y++ {
		b.WriteString(fmt.Sprintf("%d", y%10))
		for x := 0; x < r; x++ {
			i := y*r + x
			if c.Inputs[i] == 0 {
				b.WriteString(".")
			} else {
				b.WriteString("*")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

type Evaluator struct {
	cases    []Case
	show     bool
	useTrial bool
	trialNum int
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

// Setup creates the cases for each round of evaluations.
//
// The trials are organized as follows. The small object appears at 25 uniformly distributed
// locations such that it is always completely within the visual field. For each of these 25
// locations, the larger object is placed five units to the right, down, and diagonally, once per
// trial. The large object wraps around to the other side of the field when it hits the border. If
// the larger object is not completely within the visual field, it is moved the smallest distance
// possible that places it fully in view. Because of wrapping, this method of evaluation tests
// cases where the small object is on all possible sides of the large object. Thus many relative
// positions (though not all) are tested for a total number of 75 trials on the 11 by 11 substrate
// for each evaluation during evolution. (Stanley, p.17)
func (e *Evaluator) Setup() error {
	e.cases = make([]Case, *Cases)
	r := *Resolution
	for i := 0; i < len(e.cases); i += 3 {
		// Position the small box
		sx := rand.Intn(r)
		sy := rand.Intn(r)
		// Make the big box cases
		for j := 0; j < 3; j++ {

			// Create the field
			field := make([]float64, r*r)
			field[sy*r+sx] = 1.0 // Set the small box

			// Make the big box relative to the small box
			var bx, by int
			switch j {
			case 0: // To the right
				bx, by = putBig(r, sx+3, sy)
			case 1: // Down
				bx, by = putBig(r, sx, sy+3)
			case 2: // Diagonal
				bx, by = putBig(r, sx+3, sy+3)
			}

			// Fill in the big box and store the case
			for kx := -1; kx <= 1; kx++ {
				for ky := -1; ky <= 1; ky++ {
					field[(by+ky)*r+(bx+kx)] = 1.0
				}
			}
			e.cases[i+j] = Case{X: bx, Y: by, Inputs: field}
		}

	}
	return nil
}

func putBig(r, x, y int) (int, int) {
	switch {
	case x <= 1:
		x = r - 3
	case x+2 > r-1:
		x = 1
	case x+1 > r-1:
		x = 1
	}

	switch {
	case y <= 1:
		y = r - 3
	case y+2 > r-1:
		y = 1
	case y+1 > r-1:
		y = 1
	}
	return x, y
}

// Within each trial, the substrate is activated over the entire visual field. The unit with the
// highest activation in the target field is interpreted as the substrate’s selection. Fitness is
// calculated from the sum of the squared distances between the target and the point of highest
// activation over all 75 trials. This fitness function rewards generalization and provides a
// smooth gradient for solutions that are close but not perfect. (Stanley, p.17)
func (e Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {
	sum := 0.0
	for i := 0; i < len(e.cases); i++ {
		d, err := e.findBox(p, i)
		if err != nil {
			return result.New(p.ID(), 0, err, false)
		}
		sum += d
	}

	res := *Resolution
	stop := sum == 0 // Perfect run
	max := (res*res + res*res)
	fit := float64(len(e.cases)*max) - sum
	return result.New(p.ID(), fit, nil, stop)
}

// In a single trial, two objects, represented as black squares, are situated in the visual field
// at different locations. One object is three times as wide and tall as the other (figure 9). The
// goal is to locate the center of the larger object in the visual field. The target field
// specifies this location as the node with the highest level of activation. (Stanley, p.16)
func (e Evaluator) findBox(p neat.Phenome, i int) (float64, error) {
	c := e.cases[i]
	outputs, err := p.Activate(c.Inputs)
	if err != nil {
		return 0, err
	}

	var max float64
	var idx int
	for j := 0; j < len(outputs); j++ {
		if outputs[j] > max {
			max = outputs[j]
			idx = j
		}
	}

	r := *Resolution
	x := idx % r
	y := (idx - x) / r
	d := float64((c.X-x)*(c.X-x)) + float64((c.Y-y)*(c.Y-y)) // Squared distance

	if e.show {
		if err = e.showGrid(p.ID(), e.cases[i], x, y); err != nil {
			return 0, err
		}
	}
	return d, nil
}

func (e Evaluator) showGrid(id int, c Case, px, py int) error {
	// Create the image
	var t string
	if e.useTrial {
		t = fmt.Sprintf("-%d", e.trialNum)
	}
	f, err := os.Create(path.Join(*WorkPath, fmt.Sprintf("boxes%s-%d.svg", t, id)))
	if err != nil {
		return err
	}
	defer f.Close()

	w := *Resolution + 2
	h := w
	img := svg.New(f)
	img.Start(w, h)
	defer img.End()

	// Draw the grid and boxes
	img.Path(fmt.Sprintf("M %f %f L %f %f", 0, 0, w, 0), `stroke-width="1" stroke="black" fill="none"`)
	img.Path(fmt.Sprintf("M %f %f L %f %f", 0, h, w, h), `stroke-width="1" stroke="black" fill="none"`)
	img.Path(fmt.Sprintf("M %f %f L %f %f", 0, 0, 0, h), `stroke-width="1" stroke="black" fill="none"`)
	img.Path(fmt.Sprintf("M %f %f L %f %f", w, 0, w, h), `stroke-width="1" stroke="black" fill="none"`)
	for x := 0; x < *Resolution; x++ {
		for y := 0; y < *Resolution; y++ {
			i := y**Resolution + x
			if c.Inputs[i] == 1.0 {
				img.Square(x+1, h-(y+1), 1, `fill="black"`)
			}
		}
	}
	img.Circle(px+1, h-(py+1), 1, `fill="green"`)
	return nil
}

func main() {
	flag.Parse()
	if *HyperNEAT {
		fmt.Println("Using HyperNEAT decoder")
	} else if *ESHyperNEAT {
		fmt.Println("Using ESHyperNEAT decoder")
	} else if *NEAT {
		fmt.Println("Using Classic NEAT decoder")
	} else {
		log.Fatal("Please specify a decoder. See --help.")
	}
	// Ensure Cases is a multiple of 3
	if *Cases%3 != 0 {
		*Cases += 3 - *Cases%3
	}
	fmt.Println("Each evaluation will consist of", *Cases, "cases.")
	fmt.Println("Using a resolution of", *Resolution, "x", *Resolution)

	if err := trials.Run(func(i int) (*neat.Experiment, error) {
		eval := &Evaluator{}
		var ctx *starter.Context
		if *HyperNEAT {
			ctx = starter.NewContext(eval, func(ctx *starter.Context) {
				ctx.SetMutator(mutator.NewComplete(ctx, ctx, ctx, ctx, ctx, ctx))
				ctx.SetDecoder(&decoder.HyperNEAT{CppnDecoder: decoder.Classic{}, HyperNEATSettings: newSettingsWithLayers(ctx)})
			})
		} else if *ESHyperNEAT {
			ctx = starter.NewContext(eval, func(ctx *starter.Context) {
				ctx.SetMutator(mutator.NewComplete(ctx, ctx, ctx, ctx, ctx, ctx))
				ctx.SetDecoder(decoder.NewESHyperNEAT(newSettingsWithLayers(ctx), decoder.Classic{}))
			})
		} else {
			ctx = starter.NewContext(eval) // Classic NEAT decoder
		}
		if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
			return nil, err
		} else {
			return exp, nil
		}
	}); err != nil {
		log.Fatal("Could not run boxes experiment: ", err)
	}
}
