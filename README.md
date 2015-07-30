NEAT for Go
###########

This a Go implementation of NeuralEvolution of Augmenting Topologies (NEAT). From the [NEAT F.A.Q](http://www.cs.ucf.edu/~kstanley/neat.html#FAQ1).

	NEAT stands for NeuroEvolution of Augmenting Topologies. It is a method for evolving artificial neural networks with a genetic algorithm. NEAT implements the idea that it is most effective to start evolution with small, simple networks and allow them to become increasingly complex over generations. That way, just as organisms in nature increased in complexity since the first cell, so do neural networks in NEAT. This process of continual elaboration allows finding highly sophisticated and complex neural networks. 

The core of this library, often called Classic in the code, was written from the ground up using Dr. Kenneth Stanley's [PhD dissertation](http://nn.cs.utexas.edu/keyword?stanley:phd04) as a guide. NEAT has changed a bit since that paper and I have made some adjustments based on the F.A.Q. I have also add some flexibility in the desing to allow for growing the library via helpers which will provide for adding HyperNEAT, Novelty Search, etc. to the library without changing the core API.

The library and proof-of-concept experiments utilizes SVG to visualize the network of the best genome as well as the experiment's history. This visualization is based on the [NeuroEvolution Visualization Toolkit (NEVT)](http://nevt.sourceforge.net). Each image is output into an .html file for viewing from your desktop or presented through a web server.

# How to use

## Installation

```sh
go get github.com/rqme/neat
```

## Run the XOR experiment

# Create a configuration file

```json
{
    "AddConnProbability": 0.025,
    "AddNodeProbability": 0.015,
    "CompatibilityModifier": 0.3,
    "CompatibilityThreshold": 3.0,
    "DisjointCoefficient": 1,
    "EnableProbability": 0.2,
    "ExcessCoefficient": 1,
    "ExperimentName": "xor",
    "FitnessType": 0,
    "HiddenActivation": 2,
    "InterspeciesMatingRate": 0.001,
    "Iterations": 100,
    "MateByAveragingProbability": 0.4,
    "MaxStagnation": 15,
    "MutateOnlyProbability": 0.25,
    "MutateSettingProbability": 0,
    "MutateTraitProbability": 0,
    "MutateWeightProbability": 0.9,
    "NetworkIterations": 1,
    "NumInputs": 2,
    "NumOutputs": 1,
    "OutputActivation": 2,
    "PopulationSize": 150,
    "ReplaceSettingProbability": 0,
    "ReplaceTraitProbability": 0,
    "ReplaceWeightProbability": 0.2,
    "SurvivalThreshold": 0.2,
    "TargetNumberOfSpecies": 15,
    "WeightCoefficient": 0.4,
    "WeightRange": 2.5
}
```

```sh
go build github.com/rqme/neat/x/proof/xor
xor --config-path "." --archive-path "/tmp" --archive-name "xor" --web-path "/tmp"
```

