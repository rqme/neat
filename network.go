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

package neat

import "math"

// Represents a neural network
type Network interface {

	// Activates the neural network using the inputs. Returns the ouput values.
	Activate(inputs []float64) (outputs []float64, err error)
}

type NeuronType byte

const (
	Bias   NeuronType = iota + 1 // 1
	Input                        // 2
	Hidden                       // 3
	Output                       // 4
)

func (n NeuronType) String() string {
	switch n {
	case Bias:
		return "Bias"
	case Input:
		return "Input"
	case Hidden:
		return "Hidden"
	case Output:
		return "Output"
	default:
		return "Unknown NeuronType"
	}
}

type ActivationType byte

const (
	Direct          ActivationType = iota + 1 // 1
	SteependSigmoid                           // 2
	Sigmoid                                   // 3
	Tanh                                      // 4
	InverseAbs                                // 5
)

var (
	Activations []ActivationType = []ActivationType{SteependSigmoid, Sigmoid, Tanh, InverseAbs}
)

func (a ActivationType) String() string {
	switch a {
	case Direct:
		return "Direct"
	case SteependSigmoid:
		return "Steepend Sigmoid"
	case Sigmoid:
		return "Sigmoid"
	case Tanh:
		return "Tanh"
	case InverseAbs:
		return "Inverse ABS"
	default:
		return "Unknown ActivationType"
	}
}

func (a ActivationType) Range() (float64, float64) {
	switch a {
	case Direct:
		return math.Inf(-1), math.Inf(1)
	case SteependSigmoid:
		return 0, 1.0
	case Sigmoid:
		return 0, 1.0
	case Tanh:
		return -1.0, 1.0
	case InverseAbs:
		return -1.0, 1.0
	default:
		return math.NaN(), math.NaN()
	}
}

func DirectActivation(x float64) float64          { return x }
func SigmoidActivation(x float64) float64         { return 1.0 / (1.0 + exp1(-x)) }
func SteependSigmoidActivation(x float64) float64 { return 1.0 / (1.0 + exp1(-4.9*x)) }
func TanhActivation(x float64) float64            { return math.Tanh(0.9 * x) }
func InverseAbsActivation(x float64) float64      { return x / (1.0 + math.Abs(x)) }

// Speed up over math.Exp by using less precision
// https://codingforspeed.com/using-faster-exponential-approximation/
func exp1(x float64) float64 {
	x = 1.0 + x/256.0
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	return x
}

func exp2(x float64) float64 {
	x = 1.0 + x/1024
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	return x
}
