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

ATTRIBUTIONS and other notes:
* This visualizer is an adaptation of the NeuroEvolution Visualization Toolkit which can be
  found at http://sourceforge.net/projects/nevt/ which was released under LGPL v2 license. All
  functionality derived from NEVT retains the original copyright.

* The SVG creation is made easier with SVGo, github.com/ajstarks/svgo, which was released
  under the Creative Commons license.

* The statistics creation takes advantage of the stats library, github.com/montanaflynn/stats
*/

package visualizer

import (
	svg "github.com/ajstarks/svgo"
	"github.com/montanaflynn/stats"
	. "github.com/rqme/errors"
	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/network"

	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// Visualizes the population by creating web pages
type Web struct {

	// Descriptive name of the experiment
	ExperimentName string "neat:config"

	// Path to output images
	WebPath string "neat:config"

	// Decoder to use for creating network from best genome
	neat.Decoder

	// History
	WebFitness    [][3]float64  "neat:state"
	WebComplexity [][4]float64  "neat:state"
	WebSpecies    [][]int       "neat:state"
	WebBest       []neat.Genome "neat:state"
}

// Sets the configuration of the helper
func (d *Web) Configure(cfg string) error {
	return neat.Configure(cfg, d)
}

// Resets the history
func (d *Web) Reset() {
	d.WebFitness = make([][3]float64, 0, 100)
	d.WebComplexity = make([][4]float64, 0, 100)
	d.WebSpecies = make([][]int, 0, 100)
	d.WebBest = make([]neat.Genome, 0, 100)
}

// Creates visuals of the population which can be displayed in a browser
func (v *Web) Visualize(pop neat.Population) error {

	// Ensure the directory
	if _, err := os.Stat(v.WebPath); os.IsNotExist(err) {
		if err = os.Mkdir(v.WebPath, os.ModePerm); err != nil {
			return fmt.Errorf("Could not create web path %s: %v", v.WebPath, err)
		}
	}

	// Add the population to the history
	updateFitness(v, pop)
	updateComplexity(v, pop)
	updateSpecies(v, pop)
	updateBest(v, pop)

	// Create the visual components
	errs := new(Errors)
	if err := visualizeFitness(v); err != nil {
		errs.Add(err)
	}
	if err := visualizeComplexity(v); err != nil {
		errs.Add(err)
	}
	if err := visualizeSpecies(v); err != nil {
		errs.Add(err)
	}
	if err := visualizeBest(v); err != nil {
		errs.Add(err)
	}
	return errs.Err()
}

func updateFitness(v *Web, pop neat.Population) {
	// Build fitness slice
	x := make([]float64, len(pop.Genomes))
	for i, g := range pop.Genomes {
		x[i] = g.Fitness
	}

	// Append the record
	min, _ := stats.Min(x)
	max, _ := stats.Max(x)
	mean, _ := stats.Mean(x)
	v.WebFitness = append(v.WebFitness, [3]float64{
		min,
		mean,
		max,
	})
}

func updateComplexity(v *Web, pop neat.Population) {
	// Build complexity slice
	x := make([]float64, len(pop.Genomes))
	for i, g := range pop.Genomes {
		x[i] = float64(g.Complexity())
	}

	var b neat.Genome
	max := -1.0
	for _, g := range pop.Genomes {
		if g.Fitness > max {
			b = g
			max = g.Fitness
		}
	}

	// Append the record
	min, _ := stats.Min(x)
	max, _ = stats.Max(x)
	mean, _ := stats.Mean(x)

	v.WebComplexity = append(v.WebComplexity, [4]float64{
		min,
		mean,
		max,
		float64(b.Complexity()),
	})
}

func updateSpecies(v *Web, pop neat.Population) {
	cnt := make([]int, len(pop.Species))
	for _, g := range pop.Genomes {
		cnt[g.SpeciesIdx] += 1
	}
	v.WebSpecies = append(v.WebSpecies, cnt)
}

func updateBest(v *Web, pop neat.Population) {
	var best neat.Genome
	for _, g := range pop.Genomes {
		if g.Fitness > best.Fitness {
			best = g
		}
	}
	v.WebBest = append(v.WebBest, best)
}

func visualizeFitness(v *Web) error {
	// Create the file
	f, err := os.Create(path.Join(v.WebPath, "fitness.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Create the image
	img := svg.New(f)
	img.Start(575, 375)

	// Draw and label horizontal axis
	img.Path("M 40 340 L 540 340", `id="generation" stroke-width="1" stroke="black" fill="none"`)
	img.Textpath("Generation", "#generation", `fill="blue" font-size="15" font-family="Verdana" dy="30" startOffset="40%"`)

	img.Path("M 140 345 L 140 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 240 345 L 240 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 340 345 L 340 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 440 345 L 440 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 540 345 L 540 335", `stroke-width="1" stroke="black" fill="none"`)

	generations := len(v.WebFitness)
	img.Text(132, 355, fmt.Sprintf("%d", generations/5*1), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(232, 355, fmt.Sprintf("%d", generations/5*2), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(332, 355, fmt.Sprintf("%d", generations/5*3), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(432, 355, fmt.Sprintf("%d", generations/5*4), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(532, 355, fmt.Sprintf("%d", generations/5*5), `fill="black" font-size="11" font-family="Verdana"`)

	// Draw and label veritical axis
	img.Path("M 40 340 L 40 40", `id="fitness" stroke-width="1" strok="black" fill="none"`)
	img.Textpath("Fitness", "#fitness", `fill="blue" font-size="15" font-family="Verdana" dy="-25" startOffset="40%"`)

	var fitness_max, fitness_min float64
	fitness_min = 9e10
	for _, generation := range v.WebFitness {
		if generation[0] < fitness_min {
			fitness_min = generation[0]
		}
		if generation[2] > fitness_max {
			fitness_max = generation[2]
		}
	}
	fitness_range := fitness_max - fitness_min

	img.Path("M 40 40 L 540 40", `id="fitness1" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%2f", fitness_range+fitness_min), "#fitness1", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="90%"`)

	img.Path("M 40 100 L 540 100", `id="fitness2" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%2f", (fitness_range*0.8+fitness_min)), "#fitness2", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="70%"`)

	img.Path("M 40 160 L 540 160", `id="fitness3" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%2f", (fitness_range*0.6+fitness_min)), "#fitness3", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="50%"`)

	img.Path("M 40 220 L 540 220", `id="fitness4" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%2f", (fitness_range*0.4+fitness_min)), "#fitness4", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="30%"`)

	img.Path("M 40 280 L 540 280", `id="fitness5" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%2f", (fitness_range*0.2+fitness_min)), "#fitness5", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="10%"`)

	img.Textpath(fmt.Sprintf("%2f", (fitness_min)), "#generation", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8" startOffset="0%"`)

	for i, generation := range v.WebFitness {

		fitnessMax := generation[2]
		fitnessMin := generation[0]
		fitnessAvg := generation[1]

		xplot := 500/generations*i + 40
		yplotMax := 340 - 300/fitness_range*(fitnessMax-fitness_min)
		yplotMin := 340 - 300/fitness_range*(fitnessMin-fitness_min)
		yplotAvg := 340 - 300/fitness_range*(fitnessAvg-fitness_min)

		img.Circle(xplot, int(yplotAvg), 1, fmt.Sprintf(`id="%d" fill="black"`, i))
		img.Path(fmt.Sprintf("M %v %v L %v %v", xplot, yplotMin, xplot, yplotMax), `stroke-width="0.5" stroke="blue" fill="none"`)
	}
	return nil
}

func visualizeComplexity(v *Web) error {
	// Create the file
	f, err := os.Create(path.Join(v.WebPath, "complexity.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Create the image
	img := svg.New(f)
	img.Start(575, 375)

	// Draw and label horizontal axis
	img.Path("M 40 340 L 540 340", `id="generation" stroke-width="1" stroke="black" fill="none"`)
	img.Textpath("Generation", "#generation", `fill="blue" font-size="15" font-family="Verdana" dy="30" startOffset="40%"`)

	img.Path("M 140 345 L 140 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 240 345 L 240 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 340 345 L 340 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 440 345 L 440 335", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 540 345 L 540 335", `stroke-width="1" stroke="black" fill="none"`)

	generations := len(v.WebComplexity)
	img.Text(132, 355, fmt.Sprintf("%d", generations/5*1), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(232, 355, fmt.Sprintf("%d", generations/5*2), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(332, 355, fmt.Sprintf("%d", generations/5*3), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(432, 355, fmt.Sprintf("%d", generations/5*4), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(532, 355, fmt.Sprintf("%d", generations/5*5), `fill="black" font-size="11" font-family="Verdana"`)

	// Draw and label veritical axis
	img.Path("M 40 340 L 40 40", `id="complexity" stroke-width="1" strok="black" fill="none"`)
	img.Textpath("Number of Genes", "#complexity", `fill="blue" font-size="15" font-family="Verdana" dy="-25" startOffset="40%"`)

	var complexity_max, complexity_min float64
	complexity_min = 9e10
	for _, generation := range v.WebComplexity {
		if generation[0] < complexity_min {
			complexity_min = generation[0]
		}
		if generation[2] > complexity_max {
			complexity_max = generation[2]
		}
	}
	complexity_range := complexity_max - complexity_min

	img.Path("M 40 40 L 540 40", `id="complexity1" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%v", complexity_range+complexity_min), "#complexity1", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="90%"`)

	img.Path("M 40 100 L 540 100", `id="complexity2" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%v", int(complexity_range*0.8+complexity_min)), "#complexity2", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="70%"`)

	img.Path("M 40 160 L 540 160", `id="complexity3" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%v", int(complexity_range*0.6+complexity_min)), "#complexity3", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="50%"`)

	img.Path("M 40 220 L 540 220", `id="complexity4" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%v", int(complexity_range*0.4+complexity_min)), "#complexity4", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="30%"`)

	img.Path("M 40 280 L 540 280", `id="complexity5" stroke-width="0.5" stroke="green" fill="none"`)
	img.Textpath(fmt.Sprintf("%v", int(complexity_range*0.2+complexity_min)), "#complexity5", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8", startOffset="10%"`)

	img.Textpath(fmt.Sprintf("%v", int(complexity_min)), "#generation", `fill="green" fill-opacity="1.0" font-size="9" font-family="Verdana" dy="8" startOffset="0%"`)

	for i, generation := range v.WebComplexity {
		complexityMax := generation[2]
		complexityMin := generation[0]
		complexityAvg := generation[1]
		complexityChamp := generation[3]

		xplot := 500/generations*i + 40
		yplotMax := 340 - 300/complexity_range*(complexityMax-complexity_min)
		yplotMin := 340 - 300/complexity_range*(complexityMin-complexity_min)
		yplotAvg := 340 - 300/complexity_range*(complexityAvg-complexity_min)
		yplotChamp := 340 - 300/complexity_range*(complexityChamp-complexity_min)

		img.Path(fmt.Sprintf("M %v %v L %v %v L %v %v z", xplot-2, yplotChamp-2, xplot, yplotChamp-6, xplot+2, yplotChamp-2), fmt.Sprintf(`id="%d" fill="red"`, i))
		img.Circle(xplot, int(yplotAvg), 1, fmt.Sprintf(`id="%d" fill="black"`, i))
		img.Path(fmt.Sprintf("M %v %v L %v %v", xplot, yplotMin, xplot, yplotMax), `stroke-width="0.5" stroke="blue" fill="none"`)
	}
	return nil
}

func visualizeSpecies(v *Web) error {

	// Create the file
	f, err := os.Create(path.Join(v.WebPath, "species.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Identify the max population size
	popSize := 0
	for _, h := range v.WebSpecies {
		cnt := 0
		for _, s := range h {
			cnt += s
		}
		if cnt > popSize {
			popSize = cnt
		}
	}
	if popSize == 0 {
		return nil
	}

	// Create the image
	img := svg.New(f)
	img.Start(popSize*2+400, 460)
	//img.Text(10, 10, fmt.Sprintf("ID=%d Time/Date=%v", ?, ?), `style="font-size:12"`)
	img.Text(10, 25, fmt.Sprintf("PopSize=%d NumGenerations=%d", popSize, len(v.WebSpecies)), `style="font-size:10"`)
	img.Path("M 40 340 L 540 340", `id="generation" stroke-width="1" stroke="black" fill="none"`)
	img.Textpath("Generation", "#generation", `fill="blue" font-size="12" font-family="Verdana" dy="30" startOffset="25%"`)
	img.Path("M 140 345 L 140 340", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 240 345 L 240 340", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 340 345 L 340 340", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 440 345 L 440 340", `stroke-width="1" stroke="black" fill="none"`)
	img.Path("M 540 345 L 540 340", `stroke-width="1" stroke="black" fill="none"`)
	img.Text(132, 355, fmt.Sprintf("%d", len(v.WebSpecies)/5*1), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(232, 355, fmt.Sprintf("%d", len(v.WebSpecies)/5*2), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(332, 355, fmt.Sprintf("%d", len(v.WebSpecies)/5*3), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(432, 355, fmt.Sprintf("%d", len(v.WebSpecies)/5*4), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(532, 355, fmt.Sprintf("%d", len(v.WebSpecies)/5*5), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(5, 75, "# of Individuals", `style="writing-mode: tb; glyph-orientation-vertical:0; fill: blue; font-size: 10; font-family: Verdana;"`)
	img.Text(15, 285, fmt.Sprintf("%d", int(float64(popSize)*0.2)), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(15, 225, fmt.Sprintf("%d", int(float64(popSize)*0.4)), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(15, 165, fmt.Sprintf("%d", int(float64(popSize)*0.6)), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(15, 105, fmt.Sprintf("%d", int(float64(popSize)*0.8)), `fill="black" font-size="11" font-family="Verdana"`)
	img.Text(15, 45, fmt.Sprintf("%d", int(float64(popSize)*1.0)), `fill="black" font-size="11" font-family="Verdana"`)

	popIncrement := 300 / popSize
	for i, generation := range v.WebSpecies {
		xplot := 500/len(v.WebSpecies)*i + 42

		for j, speciesCount := range generation {
			var speciesColor string
			switch j % 3 {
			case 0:
				speciesColor = "CornflowerBlue"
			case 1:
				speciesColor = "Yellow"
			case 2:
				speciesColor = "Plum"
			default:
				speciesColor = "Chartreuse"
			}
			yplotFrom := 340.0
			for k := 0; k < j; k++ {
				yplotFrom -= float64(generation[k] * popIncrement)
			}

			yplotTo := yplotFrom - float64(speciesCount*popIncrement) + 0.5
			img.Path(fmt.Sprintf("M %v %v L %v %v", xplot, yplotFrom, xplot, yplotTo), fmt.Sprintf(`stroke-width="3" stroke="%s" fill="none"`, speciesColor))
		}
	}
	return nil
}

func visualizeBest(v *Web) error {
	// Create the file
	f, err := os.Create(path.Join(v.WebPath, "network.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	best := v.WebBest[len(v.WebBest)-1]
	net0, err := v.Decode(best)
	if err != nil {
		return err
	}
	var net *network.Classic
	p, ok := net0.(decoder.Phenome)
	if !ok {
		return errors.New("Web visualizer only knows the decoder package's phenome")
	}
	net, ok = p.Network.(*network.Classic)
	if !ok {
		return errors.New("Web visualizer only knows the Clasic network")
	}

	// Create the image
	img := svg.New(f)
	w, h := 1024.0, 1280.0
	img.Start(int(w)+30, int(h)+30)

	// Write out the title
	img.Text(10, 10, fmt.Sprintf("Best Genome is %d", best.ID), `font-size="12"`)

	// Define connection heads
	img.Def()
	img.Marker("triangle_black", 0, 10, 8, 6, `viewBox="0 0 20 20" markerUnits="strokeWidth" orient="auto"`)
	img.Path("M 0 0 L 20 10 L 0 20 z", `fill="black" fill-opacity="0.8"`)
	img.MarkerEnd()
	img.Marker("triangle_red", 0, 10, 8, 6, `viewBox="0 0 20 20" markerUnits="strokeWidth" orient="auto"`)
	img.Path("M 0 0 L 20 10 L 0 20 z", `fill="red" fill-opacity="0.8"`)
	img.MarkerEnd()
	img.DefEnd()

	// Draw neurons
	for i, neuron := range net.Neurons {
		var node_color, font_color string
		switch neuron.NeuronType {
		case neat.Bias:
			node_color = "black"
			font_color = "white"
		case neat.Input:
			node_color = "paleturquoise"
			font_color = "black"
		case neat.Hidden:
			node_color = "palegreen"
			font_color = "black"
		case neat.Output:
			node_color = "thistle"
			font_color = "black"
		}
		cx := int(neuron.X*w) + 15
		cy := int((1.0-neuron.Y)*h) + 15
		img.Circle(cx, cy, 10, fmt.Sprintf(`fill="%s" stroke="black" stroke-width="1"`, node_color))
		img.Text(cx-3, cy+3, fmt.Sprintf(`%d`, i), fmt.Sprintf(`font-size="5pt" font-color=%s`, font_color))
	}

	// Draw synapses
	for _, synapse := range net.Synapses {
		src := net.Neurons[synapse.Source]
		tgt := net.Neurons[synapse.Target]
		fromX := int(src.X*w) + 15
		fromY := int((1.0-src.Y)*h) + 15
		toX := int(tgt.X*w) + 15
		toY := int((1.0-tgt.Y)*h) + 15

		var line_color, triangle_color string
		if synapse.Weight >= 0 {
			line_color = "black"
			triangle_color = "#triangle_black"
		} else {
			line_color = "red"
			triangle_color = "#triangle_red"
		}

		var opacity, strokewidth string
		switch {
		case synapse.Weight < 1.0 && synapse.Weight >= 0.5, synapse.Weight > -1.0 && synapse.Weight <= -0.5:
			opacity = "0.8"
			strokewidth = "0.3"
		case synapse.Weight >= 1.0, synapse.Weight <= -1.0:
			opacity = "1.0"
			strokewidth = "0.3"
		default:
			opacity = "0.5"
			strokewidth = "1.0"
		}
		img.Path(fmt.Sprintf("M %v %v L %v %v", fromX, fromY, toX, toY), fmt.Sprintf(`fill="none" stroke="%s" stroke-width="%s" stroke-opacity="%s" marker-end="%s"`, line_color, strokewidth, opacity, triangle_color))
	}

	f.WriteString(fmt.Sprintf("<P>%s</P>", strings.Replace(best.String(), "\n", "<BR/>", -1)))
	f.WriteString(fmt.Sprintf("<P>%s</P>", strings.Replace((*net).String(), "\n", "<BR/>", -1)))
	return nil
}
