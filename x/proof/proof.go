package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/x/starter"
)

const trials = 10

func main2() {
	// create the genome
	g := neat.Genome{ID: 725, Nodes: make(map[int]neat.Node), Conns: make(map[int]neat.Connection)}
	g.Nodes[1] = neat.Node{Innovation: 1, X: 0, Y: 0, NeuronType: neat.Bias, ActivationType: neat.Direct}
	g.Nodes[2] = neat.Node{Innovation: 2, X: 0.5, Y: 0, NeuronType: neat.Input, ActivationType: neat.Direct}
	g.Nodes[3] = neat.Node{Innovation: 3, X: 1.0, Y: 0, NeuronType: neat.Input, ActivationType: neat.Direct}
	g.Nodes[4] = neat.Node{Innovation: 4, X: 0.5, Y: 1.0, NeuronType: neat.Output, ActivationType: neat.SteependSigmoid}
	g.Nodes[64] = neat.Node{Innovation: 64, X: 0.5, Y: 0.5, NeuronType: neat.Hidden, ActivationType: neat.SteependSigmoid}
	g.Nodes[234] = neat.Node{Innovation: 234, X: 0.5, Y: 0.25, NeuronType: neat.Hidden, ActivationType: neat.SteependSigmoid}
	g.Conns[5] = neat.Connection{Innovation: 5, Source: 1, Target: 4, Enabled: true, Weight: 2.4734613949471784}
	g.Conns[6] = neat.Connection{Innovation: 6, Source: 2, Target: 4, Enabled: false, Weight: -5.64995113868551}
	g.Conns[7] = neat.Connection{Innovation: 7, Source: 3, Target: 4, Enabled: true, Weight: 0.8428774810124069}
	g.Conns[65] = neat.Connection{Innovation: 65, Source: 2, Target: 64, Enabled: false, Weight: 1.6538892148212083}
	g.Conns[66] = neat.Connection{Innovation: 66, Source: 64, Target: 4, Enabled: true, Weight: -2.9999856850356714}
	g.Conns[84] = neat.Connection{Innovation: 84, Source: 3, Target: 64, Enabled: true, Weight: 3.9055228465697596}
	g.Conns[159] = neat.Connection{Innovation: 159, Source: 1, Target: 64, Enabled: true, Weight: 0.350296980355459}
	g.Conns[235] = neat.Connection{Innovation: 235, Source: 2, Target: 234, Enabled: true, Weight: 6.111340339072617}
	g.Conns[236] = neat.Connection{Innovation: 236, Source: 234, Target: 64, Enabled: true, Weight: -2.1660946640041074}
	g.Conns[335] = neat.Connection{Innovation: 335, Source: 234, Target: 4, Enabled: true, Weight: -0.338797031578537}
	g.Conns[463] = neat.Connection{Innovation: 463, Source: 1, Target: 234, Enabled: true, Weight: -3.4407146147284307}

	// decode the genome
	d := &decoder.Classic{}
	p, _ := d.Decode(g)

	// Evaluate
	e := &NEATEval{}
	r := e.Evaluate(p)
	fmt.Println(r)
}
func main() {
	log.Println("Running proofs. Success rates under 100% are OK.")
	run("neat", neatContext)
	run("phased", phasedContext)
	run("hyperneat", hyperneatContext)
}

func run(name string, f func() *starter.Context) {
	wg := new(sync.WaitGroup)
	ch := make(chan float64, trials)
	for i := 0; i < trials; i++ {
		wg.Add(1)
		go func() {
			ctx := f()
			exp := &neat.Experiment{ExperimentSettings: ctx}
			exp.SetContext(ctx)
			if err := neat.Run(exp); err != nil {
				log.Fatalf("Fatal error in %s: %v\n", name, err)
			}
			if exp.Stopped() {
				ch <- 0.0
			} else {
				ch <- 1.0
			}
			wg.Done()
		}()
	}
	wg.Wait()

	sum := 0.0
	for i := 0; i < trials; i++ {
		sum += <-ch
	}
	log.Println(name, "success rate:", sum/10.0)
}
