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
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rqme/neat"
	"github.com/rqme/neat/generator"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/searcher"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

var (
	MazeFile = flag.String("maze", "medium_maze.txt", "Maze file to use in the experiment")
	Steps    = flag.Int("steps", 400, "Number of steps a hero has to solve the maze")
	Novelty  = flag.Bool("novelty", false, "Use novelty search instead of objective fitness")
	RealTime = flag.Bool("realtime", false, "Use the real-time generator and evaluator")
	WorkPath = flag.String("work-path", ".", "Output directory for maze diagrams")
)

type Result struct {
	result.Classic
	behavior []float64
}

func (r *Result) Behavior() []float64 { return r.behavior }

func loadEnv(p string) (e Environment, err error) {
	var f *os.File
	f, err = os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var n, i int
	var parts []string
	for s.Scan() {
		t := s.Text()
		switch i {
		case 0: // Number of lines
			if n, err = strconv.Atoi(t); err != nil {
				return
			}
			e.Lines = make([]Line, 0, n)
		case 1: // Start position
			parts = strings.Split(t, " ")
			if e.Hero.Location.X, err = strconv.ParseFloat(parts[0], 64); err != nil {
				return
			}
			if e.Hero.Location.Y, err = strconv.ParseFloat(parts[1], 64); err != nil {
				return
			}
		case 2: // Initial heading
			if e.Hero.Heading, err = strconv.ParseFloat(t, 64); err != nil {
				return
			}
		case 3: // End position
			parts = strings.Split(t, " ")
			if e.End.X, err = strconv.ParseFloat(parts[0], 64); err != nil {
				return
			}
			if e.End.Y, err = strconv.ParseFloat(parts[1], 64); err != nil {
				return
			}
		default: // Line
			var line Line
			parts := strings.Split(t, " ")
			if line.A.X, err = strconv.ParseFloat(parts[0], 64); err != nil {
				return
			}
			if line.A.Y, err = strconv.ParseFloat(parts[1], 64); err != nil {
				return
			}
			if line.B.X, err = strconv.ParseFloat(parts[2], 64); err != nil {
				return
			}
			if line.B.Y, err = strconv.ParseFloat(parts[3], 64); err != nil {
				return
			}
			e.Lines = append(e.Lines, line)
		}
		i += 1
	}
	return
}

func main() {
	flag.Parse()
	//defer profile.Start(profile.CPUProfile).Stop()
	if *Novelty {
		fmt.Println("Using Novelty search")
	} else {
		fmt.Println("Using Fitness search")
	}
	if *RealTime {
		fmt.Println("Using Real-Time generator")
	} else {
		fmt.Println("Using Classic generator")
	}
	fmt.Println("Using", *Steps, "time steps per evaluation")
	fmt.Println("Loading maze file:", *MazeFile)
	var err error
	orig := Evaluator{}
	if orig.Environment, err = loadEnv(*MazeFile); err != nil {
		log.Fatalf("Could not load maze file %s: %v", *MazeFile, err)
	}

	if err = trials.Run(func(i int) (*neat.Experiment, error) {

		var ctx *starter.Context
		var eval neat.Evaluator
		var gen neat.Generator
		if *RealTime {
			eval = &RTEvaluator{Evaluator: Evaluator{Environment: orig.clone()}}
			gen = &generator.RealTime{RealTimeSettings: ctx}
		} else {
			eval = &Evaluator{Environment: orig.clone()}
			gen = &generator.Classic{ClassicSettings: ctx}
		}
		if *Novelty {
			ctx = starter.NewContext(eval, func(ctx *starter.Context) {
				ctx.SetGenerator(gen)
				ctx.SetSearcher(&searcher.Novelty{NoveltySettings: ctx, Searcher: &searcher.Concurrent{}})
			})
		} else {
			ctx = starter.NewContext(eval, func(ctx *starter.Context) {
				ctx.SetGenerator(gen)
			})
		}

		if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
			return nil, err
		} else {
			//ctx.Settings.ArchivePath = path.Join(ctx.Settings.ArchivePath, "fit")
			//ctx.Settings.ArchiveName = ctx.Settings.ArchiveName + "-fit"
			//ctx.Settings.WebPath = path.Join(ctx.Settings.WebPath, "fit")
			return exp, nil
		}
	}); err != nil {
		log.Fatal("Could not run maze experiment: ", err)
	}
}
