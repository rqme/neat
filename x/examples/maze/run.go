package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

var (
	MazeFile = flag.String("maze", "medium_maze.txt", "Maze file to use in the experiment")
	Steps    = flag.Int("steps", 400, "Number of steps a hero has to solve the maze")
	Novelty  = flag.Bool("novelty", false, "Use novelty search instead of objective fitness")
)

type Result struct {
	result.Classic
	behavior []float64
}

func (r *Result) Behavior() []float64 { return r.behavior }

type Evaluator struct {
	Environment
	show     bool
	workPath string

	useTrial bool
	trialNum int
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

func (e Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	// Clone the environemnt
	var env *Environment
	cln := e.Environment.clone()
	env = &cln
	env.init()
	h := &env.Hero

	// Iterate the maze
	var err error
	paths := make([]Line, 0, *Steps)
	stop := false
	for i := 0; i < *Steps; i++ {

		// Note the start of the path
		a := h.Location

		// Update the hero's location
		var outputs []float64
		inputs := generateNeuralInputs(*env)
		if outputs, err = p.Activate(inputs); err != nil {
			break
		}
		interpretOutputs(env, outputs[0], outputs[1])
		update(env)
		b := h.Location
		paths = append(paths, Line{A: a, B: b})

		// Look for solution
		stop = distanceToTarget(env) < 5.0
		if stop {
			break
		}
	}
	f := 300.0 - distanceToTarget(env)         // fitness
	b := []float64{h.Location.X, h.Location.Y} // behavior
	r = &Result{Classic: result.New(p.ID(), f, err, stop), behavior: b}

	// Output the maze
	if e.show {
		// Write the file
		var t string
		if e.useTrial {
			t = fmt.Sprintf("-%d", e.trialNum)
		}
		if err := showMaze(fmt.Sprintf("maze%s-%d.svg", t, p.ID()), *env, paths); err != nil {
			log.Println("Could not output maze run:", err)
		}
	}
	return
}

func showMaze(p string, e Environment, hist []Line) error {

	// Determine image size
	var h, w float64
	for _, line := range e.Lines {
		if line.A.X > w {
			w = line.A.X
		}
		if line.A.Y > h {
			h = line.A.Y
		}
		if line.B.X > w {
			w = line.B.X
		}
		if line.B.Y > h {
			h = line.B.Y
		}
	}

	// Create the image
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	img := svg.New(f)
	img.Start(int(w), int(h))
	defer img.End()

	// Add the maze
	img.Circle(int(hist[0].A.X), int(hist[0].A.Y), 4, `fill="green"`) // start
	img.Circle(int(e.End.X), int(e.End.Y), 4, `fill="red"`)

	for _, line := range e.Lines {
		img.Path(fmt.Sprintf("M %f %f L %f %f", line.A.X, line.A.Y, line.B.X, line.B.Y), `stroke-width="1" stroke="black" fill="none"`)
	}

	for _, line := range hist {
		img.Path(fmt.Sprintf("M %f %f L %f %f", line.A.X, line.A.Y, line.B.X, line.B.Y), `stroke-width="1" stroke="blue" fill="none"`)

	}
	return nil
}

func (e *Evaluator) ShowWork(s bool) {
	e.show = s
}

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
	fmt.Println("Using", *Steps, "time steps per evaluation")
	fmt.Println("Loading maze file:", *MazeFile)
	var err error
	orig := Evaluator{}
	if orig.Environment, err = loadEnv(*MazeFile); err != nil {
		log.Fatalf("Could not load maze file %s: %v", *MazeFile, err)
	}

	if err = trials.Run(func(i int) (*neat.Experiment, error) {

		eval := &Evaluator{Environment: orig.clone()}
		var ctx *starter.Context
		if *Novelty {
			ctx = starter.NewNoveltyContext(eval)
		} else {
			ctx = starter.NewClassicContext(eval)
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
