package main

import (
	"github.com/rqme/neat"
	"github.com/rqme/neat/decoder"
)

type HyperNEATSettings struct {
	decoder.HyperNEATSettings
	layers []decoder.SubstrateNodes
}

func (h HyperNEATSettings) SubstrateLayers() []decoder.SubstrateNodes { return h.layers }

// The solution substrate is configured as a state-space sandwich that includes two sheets: (1) The
// visual field is a two-dimensional array of sensors that are either on or off (i.e. black or
// white); (2) The target field is an equivalent two-dimensional array of outputs that are
// activated at variable intensity between zero and one. (Stanley, p.15)
func newHyperNEAT(cfg decoder.HyperNEATSettings) (hns HyperNEATSettings) {
	hns.HyperNEATSettings = cfg

	// Create the substrate layers
	r := *Resolution
	var ilayer decoder.SubstrateNodes = make([]decoder.SubstrateNode, 0, r*r)
	var olayer decoder.SubstrateNodes = make([]decoder.SubstrateNode, 0, r*r)
	for x := 0; x < r; x++ {
		for y := 0; y < r; y++ {
			px := float64(x)/float64(r-1)*2.0 - 1.0
			py := float64(y)/float64(r-1)*2.0 - 1.0
			ilayer = append(ilayer, decoder.SubstrateNode{Position: []float64{px, py}, NeuronType: neat.Input})
			olayer = append(olayer, decoder.SubstrateNode{Position: []float64{px, py}, NeuronType: neat.Output})
		}
	}
	hns.layers = []decoder.SubstrateNodes{ilayer, olayer}
	return
}
