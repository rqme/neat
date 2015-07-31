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

package trials

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rqme/neat"
)

var (
	Trials    = flag.Int("trials", 10, "Number of trials to run. Default is 10.")
	CheckStop = flag.Bool("check-stop", false, "Consider trial a failure if stop condition not met")
)

const (
	Seconds int = iota
	Iterations
	Nodes
	Connections
	Fitness
)

const (
	Min int = iota
	Avg
	Max
	Sum
	Cnt
)

func Run(f func(int) (*neat.Experiment, error)) error {

	// Create the collection variables
	n := *Trials
	fail := make([]bool, n)
	errs := make([]error, n)
	data := make([][]float64, n)
	for i := 0; i < len(data); i++ {
		data[i] = make([]float64, 5)
	}

	// Begin the display
	fmt.Printf("Run   Iters.   Seconds    Nodes     Conns    Fitness\n")
	fmt.Printf("--- --------- --------- --------- --------- ---------\n")

	// Iterate the trials
	for i := 0; i < n; i++ {
		t0 := time.Now()
		if e, err := f(i); err != nil {
			return err
		} else {
			if err = neat.Run(e); err != nil {
				errs[i] = fmt.Errorf("Error during trial %d: %v", i, err)
				fail[i] = true
			} else {
				if *CheckStop && !e.Stopped {
					fail[i] = true
				} else {
					data[i][Seconds] = time.Now().Sub(t0).Seconds()
					b := findBest(e.Population)
					data[i][Nodes] = float64(len(b.Nodes))
					data[i][Connections] = float64(len(b.Conns))
					data[i][Fitness] = b.Fitness
					data[i][Iterations] = float64(e.Iteration)
					//fail[i] = e.Iteration == e.Iterations
				}
			}
		}

		// Update the display
		if fail[i] {
			if errs[i] != nil {
				fmt.Printf("%s failed: %s\n", padInt(i, 3), errs[i])
			} else {
				fmt.Printf("%s failed: no solution found\n", padInt(i, 3)) // Print failure message : iterations or error
			}
		} else {
			fmt.Printf("%s %s %s %s %s %s\n", padInt(i, 3), padFloat(data[i][Iterations], 0, 9), padFloat(data[i][Seconds], 3, 9), padFloat(data[i][Nodes], 0, 9), padFloat(data[i][Connections], 0, 9), padFloat(data[i][Fitness], 3, 9))
		}
	}

	// Calculate Stats
	stats := make([][]float64, 5)
	for i := 0; i < len(stats); i++ {
		stats[i] = make([]float64, 5)
		stats[i][Min] = math.Inf(1)
		stats[i][Max] = math.Inf(-1)
		for j := 0; j < n; j++ {
			if !fail[j] {
				if data[j][i] < stats[i][Min] {
					stats[i][Min] = data[j][i]
				}
				if data[j][i] > stats[i][Max] {
					stats[i][Max] = data[j][i]
				}
				stats[i][Sum] += data[j][i]
				stats[i][Cnt] += 1.0
			}
		}
		if stats[i][Cnt] > 0 {
			stats[i][Avg] = stats[i][Sum] / stats[i][Cnt]
		}
	}

	// Report the findings
	if stats[0][Cnt] > 1 {
		fmt.Println("")
		fmt.Printf("MIN %s %s %s %s %s\n", padFloat(stats[Iterations][Min], 0, 9), padFloat(stats[Seconds][Min], 3, 9), padFloat(stats[Nodes][Min], 0, 9), padFloat(stats[Connections][Min], 0, 9), padFloat(stats[Fitness][Min], 3, 9))
		fmt.Printf("AVG %s %s %s %s %s\n", padFloat(stats[Iterations][Avg], 0, 9), padFloat(stats[Seconds][Avg], 3, 9), padFloat(stats[Nodes][Avg], 0, 9), padFloat(stats[Connections][Avg], 0, 9), padFloat(stats[Fitness][Avg], 3, 9))
		fmt.Printf("MAX %s %s %s %s %s\n", padFloat(stats[Iterations][Max], 0, 9), padFloat(stats[Seconds][Max], 3, 9), padFloat(stats[Nodes][Max], 0, 9), padFloat(stats[Connections][Max], 0, 9), padFloat(stats[Fitness][Max], 3, 9))
	}
	return nil
}

func findBest(p neat.Population) (b neat.Genome) {
	for _, g := range p.Genomes {
		if g.Fitness > b.Fitness {
			b = g
		}
	}
	return
}

func padInt(i int, p int) string {
	s := strconv.Itoa(i)
	if len(s) < p {
		s = strings.Repeat(" ", p-len(s)) + s
	}
	return s
}

func padFloat(f float64, d int, p int) string {
	t := "%." + strconv.Itoa(d) + "f"
	s := fmt.Sprintf(t, f)
	if len(s) < p {
		s = strings.Repeat(" ", p-len(s)) + s
	}
	return s
}
