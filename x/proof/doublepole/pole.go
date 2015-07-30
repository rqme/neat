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


The algorithm for this experiment is taken directly from Fernando Torres's reorganized version of
Stanley's original code: https://github.com/FernandoTorres/NEAT. The BSD 2 clause license used by
the RedQ NEAT library is compatibile with the original Apache license of Dr. Stanley's code.
*/

package main

import (
	"flag"
	"log"
	"math"
	"path"
	"strconv"

	"github.com/rqme/neat"
	"github.com/rqme/neat/archiver"
	"github.com/rqme/neat/searcher"
	"github.com/rqme/neat/visualizer"
	"github.com/rqme/neat/x/example"
	"github.com/rqme/neat/x/trials"
)

var (
	Velocity = flag.Bool("velocity", false, "true if velocity info used")
)

const (
	NumInputs        = 7
	MUP              = 0.000002
	MUC              = 0.0005
	Gravity          = -9.8
	MassCart         = 1.0
	MassPole1        = 0.1
	Length1          = 0.5
	ForceMag         = 10.0
	Tau              = 0.01
	OneDegree        = 0.0174532
	SixDegrees       = 0.1047192
	TwelveDegrees    = 0.2094384
	FifteenDegrees   = 0.2617993
	ThirtySixDegrees = 0.628329
	FiftyDegrees     = 0.87266
)

type CartPole struct {
	MaxFitness         int
	Markov             bool
	LastHundred        bool
	NMarkovLong        bool
	GeneralizationTest bool
	State              []float64
	JiggleStep         []float64

	length2   float64
	masspole2 float64
	minInc    float64
	poleInc   float64
	massInc   float64

	balancedSum int
	cartposSum  float64
	cartvSum    float64
	poleposSum  float64
	polevSum    float64
}

func NewCartPole(randomize bool, velocity bool) *CartPole {
	return &CartPole{
		MaxFitness: 100000,
		Markov:     velocity,
		minInc:     0.001,
		poleInc:    0.05,
		massInc:    0.01,
		length2:    0.05,
		masspole2:  0.01,
		State:      make([]float64, 6),
		JiggleStep: make([]float64, 1000),
	}
}

func (cp *CartPole) SimplifyTask() {
	if cp.poleInc > cp.minInc {
		cp.poleInc = cp.poleInc / 2.0
		cp.massInc = cp.massInc / 2.0
		cp.length2 -= cp.poleInc
		cp.masspole2 -= cp.massInc
	}
}

func (cp *CartPole) NextTask() {
	cp.length2 += cp.poleInc
	cp.masspole2 += cp.massInc
}

func (cp *CartPole) EvalNet(net neat.Network, thresh int) (float64, error) {

	input := make([]float64, NumInputs)
	var nmarkovmax int
	var nmarkov_fitness float64
	var jiggletotal float64
	var steps int
	if cp.NMarkovLong {
		nmarkovmax = 100000
	} else if cp.GeneralizationTest {
		nmarkovmax = 1000
	} else {
		nmarkovmax = 1000
	}

	cp.init(false)

	if cp.Markov {
		for steps < cp.MaxFitness {
			steps += 1
			input[0] = cp.State[0] / 4.8
			input[1] = cp.State[1] / 2
			input[2] = cp.State[2] / 0.52
			input[3] = cp.State[3] / 2
			input[4] = cp.State[4] / 0.52
			input[5] = cp.State[5] / 2
			input[6] = .5

			outputs, err := net.Activate(input)
			if err != nil {
				return 1.0, err
			}
			cp.performAction(outputs[0], steps)
			if cp.outsideBounds() {
				break
			}
		}
		return float64(steps), nil
	} else {
		for steps < nmarkovmax {
			steps += 1
			input[0] = cp.State[0] / 4.8
			input[1] = cp.State[2] / 0.52
			input[2] = cp.State[4] / 0.52
			input[3] = .5

			outputs, err := net.Activate(input)
			if err != nil {
				return 1.0, err
			}
			cp.performAction(outputs[0], steps)
			if cp.NMarkovLong && cp.outsideBounds() {
				break
			}
		}

		if cp.GeneralizationTest {
			return float64(cp.balancedSum), nil
		}

		if steps > 100 && !cp.NMarkovLong {
			jiggletotal = 0
			for count := steps - 99 - 2; count <= steps-2; count++ {
				jiggletotal += cp.JiggleStep[count]
			}
		}

		if !cp.NMarkovLong {
			if cp.balancedSum > 100 {
				nmarkov_fitness = 0.1*float64(cp.balancedSum)/1000.0 + 0.9*0.75/jiggletotal
			} else {
				nmarkov_fitness = 0.1 * float64(cp.balancedSum) / 1000.0
			}
			return nmarkov_fitness, nil
		} else {
			return float64(steps), nil
		}
	}
}

func (cp *CartPole) init(randomize bool) {
	if !cp.Markov {
		cp.cartposSum = 0.0
		cp.cartvSum = 0.0
		cp.poleposSum = 0.0
		cp.polevSum = 0.0
	}
	cp.balancedSum = 0.0
	cp.LastHundred = false

	if !cp.GeneralizationTest {
		cp.State[0] = 0
		cp.State[1] = 0
		cp.State[3] = 0
		cp.State[4] = 0
		cp.State[5] = 0
		cp.State[2] = 0.07 // one degree
	} else {
		cp.State[4] = 0
		cp.State[5] = 0
	}
}

func (cp *CartPole) performAction(output float64, stepnum int) {
	dydx := make([]float64, 6)
	RK4 := true
	EULER_TAU := Tau / 4.0

	if RK4 {
		for i := 0; i < 2; i++ {
			dydx[0] = cp.State[1]
			dydx[2] = cp.State[3]
			dydx[4] = cp.State[5]
			cp.step(output, cp.State, dydx)
			cp.rk4(output, cp.State, dydx, cp.State)
		}
	} else {
		for i := 0; i < 8; i++ {
			cp.step(output, cp.State, dydx)
			cp.State[0] += EULER_TAU * dydx[0]
			cp.State[1] += EULER_TAU * dydx[1]
			cp.State[2] += EULER_TAU * dydx[2]
			cp.State[3] += EULER_TAU * dydx[3]
			cp.State[4] += EULER_TAU * dydx[4]
			cp.State[5] += EULER_TAU * dydx[5]
		}
	}

	// Record this state
	cp.cartposSum += math.Abs(cp.State[0])
	cp.cartvSum += math.Abs(cp.State[1])
	cp.poleposSum += math.Abs(cp.State[2])
	cp.polevSum += math.Abs(cp.State[3])
	if stepnum <= 1000 {
		cp.JiggleStep[stepnum-1] = math.Abs(cp.State[0]) + math.Abs(cp.State[1]) + math.Abs(cp.State[2]) + math.Abs(cp.State[3])
	}

	if false {

	} else if !cp.outsideBounds() {
		cp.balancedSum += 1
	}
}

func (cp *CartPole) step(action float64, st []float64, derivs []float64) {
	var force, costheta_1, costheta_2, sintheta_1, sintheta_2, gsintheta_1, gsintheta_2, temp_1, temp_2, ml_1, ml_2, fi_1, fi_2, mi_1, mi_2 float64

	force = (action - 0.5) * ForceMag * 2
	costheta_1 = math.Cos(st[2])
	sintheta_1 = math.Sin(st[2])
	gsintheta_1 = Gravity * sintheta_1
	costheta_2 = math.Cos(st[4])
	sintheta_2 = math.Sin(st[4])
	gsintheta_2 = Gravity * sintheta_2

	ml_1 = Length1 * MassPole1
	ml_2 = cp.length2 * cp.masspole2
	temp_1 = MUP * st[3] / ml_1
	temp_2 = MUP * st[5] / ml_2
	fi_1 = (ml_1 * st[3] * st[3] * sintheta_1) +
		(0.75 * MassPole1 * costheta_1 * (temp_1 + gsintheta_1))
	fi_2 = (ml_2 * st[5] * st[5] * sintheta_2) +
		(0.75 * cp.masspole2 * costheta_2 * (temp_2 + gsintheta_2))
	mi_1 = MassPole1 * (1 - (0.75 * costheta_1 * costheta_1))
	mi_2 = cp.masspole2 * (1 - (0.75 * costheta_2 * costheta_2))

	derivs[1] = (force + fi_1 + fi_2) / (mi_1 + mi_2 + MassCart)
	derivs[3] = -0.75 * (derivs[1]*costheta_1 + gsintheta_1 + temp_1) / Length1
	derivs[5] = -0.75 * (derivs[1]*costheta_2 + gsintheta_2 + temp_2) / cp.length2
}

func (cp *CartPole) rk4(f float64, y []float64, dydx []float64, yout []float64) {
	var hh, h6 float64
	var dym, dyt, yt []float64
	dym = make([]float64, 6)
	dyt = make([]float64, 6)
	yt = make([]float64, 6)

	hh = Tau * 0.5
	h6 = Tau / 6.0
	for i := 0; i <= 5; i++ {
		yt[i] = y[i] + hh*dydx[i]
	}
	cp.step(f, yt, dyt)
	dyt[0] = yt[1]
	dyt[2] = yt[3]
	dyt[4] = yt[5]
	for i := 0; i <= 5; i++ {
		yt[i] = y[i] + hh*dyt[i]
	}
	cp.step(f, yt, dym)
	dym[0] = yt[1]
	dym[2] = yt[3]
	dym[4] = yt[5]
	for i := 0; i <= 5; i++ {
		yt[i] = y[i] + Tau*dym[i]
		dym[i] += dyt[i]
	}
	cp.step(f, yt, dyt)
	dyt[0] = yt[1]
	dyt[2] = yt[3]
	dyt[4] = yt[5]
	for i := 0; i <= 5; i++ {
		yout[i] = y[i] + h6*(dydx[i]+dyt[i]+2.0*dym[i])
	}
}

func (cp *CartPole) outsideBounds() bool {
	failureAngle := ThirtySixDegrees
	return cp.State[0] < -2.4 ||
		cp.State[0] > 2.4 ||
		cp.State[2] < -failureAngle ||
		cp.State[2] > failureAngle ||
		cp.State[4] < -failureAngle ||
		cp.State[4] > failureAngle
}

type evaluator struct {
	velocity bool
	thecart  *CartPole
}

func (e evaluator) Evaluate(p neat.Phenome) (r neat.Result) {
	thresh := 100
	fit, err := e.thecart.EvalNet(p, thresh)

	stop := false
	if e.thecart.Markov {
		if fit >= float64(e.thecart.MaxFitness-1) {
			stop = true
		}
	} else if e.thecart.NMarkovLong {
		if fit >= 99999 {
			stop = true
		}
	} else if e.thecart.GeneralizationTest {
		if fit >= 999 {
			stop = true
		}
	}

	// Calculate the result
	r = example.NewResult(p.ID(), fit, err, stop)
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
			func(e *neat.Experiment) error {
				e.Searcher = searcher.Serial{
					Evaluator: &evaluator{
						velocity: *Velocity,
						thecart:  NewCartPole(true, *Velocity),
					},
				}
				return nil
			},
		)
		return e, err
	}); err != nil {
		log.Fatal("Could not run pole experiment: ", err)
	}

}
