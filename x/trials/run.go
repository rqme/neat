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
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/rqme/neat"
)

var (
	Trials     = flag.Int("trials", 10, "Number of trials to run. Default is 10.")
	CheckStop  = flag.Bool("check-stop", false, "Consider trial a failure if stop condition not met")
	ShowWork   = flag.Bool("show-work", false, "Evaluates the best genome separately, showing the detail ")
	SkipEvolve = flag.Bool("skip-evolve", false, "Skips evolution phase if population restored. Use with -best to display archived population.")
)

func Run(f func(int) (*neat.Experiment, error)) error {

	// Create the collection variables
	n := *Trials
	exps := make([]*neat.Experiment, n)
	best := make([]neat.Genome, n)
	noso := make([]bool, n)
	fail := make([]bool, n)
	skip := make([]bool, n)
	errs := make([]error, n)

	var secs stats.Float64Data = make([]float64, 0, n)
	var fits stats.Float64Data = make([]float64, 0, n)
	var iters stats.Float64Data = make([]float64, 0, n)
	var nodes stats.Float64Data = make([]float64, 0, n)
	var conns stats.Float64Data = make([]float64, 0, n)

	var se, fi, it, no, co float64

	// Begin the display
	fmt.Printf("Beginning trials %v\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Run   Iters.   Seconds    Nodes     Conns    Fitness   Fail   Comment \n")
	fmt.Printf("--- --------- --------- --------- --------- --------- ------ ---------\n")

	// Iterate the trials
	var showTime, showBest bool
	var err error
	for i := 0; i < n; i++ {
		var t0, t1 time.Time
		t0 = time.Now()
		if exps[i], err = f(i); err != nil {
			errs[i] = err
			fail[i] = true
			showBest = false
			showTime = false
		} else {
			if *SkipEvolve && len(exps[i].Population().Genomes) > 0 {
				skip[i] = true
			} else if err = neat.Run(exps[i]); err != nil {
				t1 = time.Now()
				errs[i] = err
				fail[i] = true
			} else {
				t1 = time.Now()
				if *CheckStop && !exps[i].Stopped() {
					errs[i] = fmt.Errorf("No solution found")
					noso[i] = true
				}
			}

			// Update the time stats
			if exps[i].Iteration() > 0 {
				showTime = true
				it = float64(exps[i].Iteration())
				se = t1.Sub(t0).Seconds()
				if !fail[i] && !skip[i] {
					iters = append(iters, it)
					secs = append(secs, se)
				}
			} else {
				showTime = false
			}

			// Find the best and update the results
			if !fail[i] && len(exps[i].Population().Genomes) > 0 {
				showBest = true
				best[i] = findBest(exps[i].Population())
				no = float64(len(best[i].Nodes))
				co = float64(len(best[i].Conns))
				fi = best[i].Fitness

				if !noso[i] {
					nodes = append(nodes, no)
					conns = append(conns, co)
					fits = append(fits, fi)
				}
			} else {
				showBest = false
			}
		}
		showTrial(padInt(i, 3), showTime, showBest, it, se, no, co, fi, fail[i] || noso[i], skip[i], errs[i])
	}

	// Display the summary
	funcs := []func(stats.Float64Data) (float64, error){stats.Mean, stats.Median, stats.StdDevP, stats.Min, stats.Max}
	labs := []string{"AVG", "MED", "SDV", "MIN", "MAX"} // 3 characters
	var itm, sem, nom, com, fim [5]float64
	if len(iters) > 0 {
		showTime = true
		for i := 0; i < len(funcs); i++ {
			itm[i], _ = funcs[i](iters)
			sem[i], _ = funcs[i](secs)
		}
	} else {
		showTime = false
	}
	if len(nodes) > 0 {
		showBest = true
		for i := 0; i < len(funcs); i++ {
			nom[i], _ = funcs[i](nodes)
			com[i], _ = funcs[i](conns)
			fim[i], _ = funcs[i](fits)
		}
	} else {
		showBest = false
	}

	fmt.Printf("\nSummary for trials excluding failures (and time for skipped)\n")
	fmt.Printf("      Iters.   Seconds    Nodes     Conns    Fitness\n")
	fmt.Printf("--- --------- --------- --------- --------- ---------\n")
	for i := 0; i < len(itm); i++ {
		showTrial(labs[i], showTime, showBest, itm[i], sem[i], nom[i], com[i], fim[i], false, false, nil)
	}

	// Show the evaluations of the best
	if *ShowWork {
		for i := 0; i < len(exps); i++ {
			if dh, ok := exps[i].Context().Evaluator().(neat.Demonstrable); ok {
				dh.ShowWork(true)
				if rst, ok := exps[i].Context().Archiver().(neat.Restorer); ok {
					ctx := exps[i].Context()
					if err := rst.Restore(ctx); err != nil {
						fmt.Printf("Error restoring the context for trial %d: %v\n", i, err)
					}
				}
				if p, err := exps[i].Context().Decoder().Decode(best[i]); err != nil {
					fmt.Printf("Error decoding phenome in show work: %v\n", err)
				} else {
					evl := exps[i].Context().Evaluator()
					if th, ok := evl.(neat.Trialable); ok {
						th.SetTrial(i)
					}
					r := evl.Evaluate(p)
					if r.Err() != nil {
						fmt.Printf("Error demonstrating phenome: %v\n", r.Err())
					}
				}
			}
		}
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

func showTrial(key string, showTime, showBest bool, it, se, no, co, fi float64, fail, skip bool, err error) {
	var ses, its, nos, cos, fis, fas, cms string
	if showTime {
		its = padFloat(it, 0, 9)
		ses = padFloat(se, 3, 9)
	} else {
		its = strings.Repeat(" ", 9)
		ses = strings.Repeat(" ", 9)
	}
	if showBest {
		nos = padFloat(no, 0, 9)
		cos = padFloat(co, 0, 9)
		fis = padFloat(fi, 3, 9)
	} else {
		nos = strings.Repeat(" ", 9)
		cos = strings.Repeat(" ", 9)
		fis = strings.Repeat(" ", 9)
	}
	if fail {
		fas = " Yes  "
		cms = err.Error()
	} else {
		fas = strings.Repeat(" ", 6)
		if skip {
			cms = "Skipped"
		} else {
			cms = ""
		}
	}
	fmt.Printf("%s %s %s %s %s %s %s %s\n", key, its, ses, nos, cos, fis, fas, cms)
}
