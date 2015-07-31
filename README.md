NEAT for Go
###########

This a Go implementation of NeuralEvolution of Augmenting Topologies (NEAT). From the [NEAT F.A.Q](http://www.cs.ucf.edu/~kstanley/neat.html#FAQ1).

*NEAT stands for NeuroEvolution of Augmenting Topologies. It is a method for evolving artificial neural networks with a genetic algorithm. NEAT implements the idea that it is most effective to start evolution with small, simple networks and allow them to become increasingly complex over generations. That way, just as organisms in nature increased in complexity since the first cell, so do neural networks in NEAT. This process of continual elaboration allows finding highly sophisticated and complex neural networks.*

The core of this library, often called Classic in the code, was written from the ground up using Dr. Kenneth Stanley's [PhD dissertation](http://nn.cs.utexas.edu/keyword?stanley:phd04) as a guide. NEAT has changed a bit since that paper and I have made some adjustments based on the F.A.Q. I have also add some flexibility in the design to allow for growing the library via helpers which will provide for adding HyperNEAT, Novelty Search, etc. to the library without changing the core API.

The library and proof-of-concept experiments utilize SVG to visualize the network of the best genome as well as the experiment's history. This visualization is based on the [NeuroEvolution Visualization Toolkit (NEVT)](http://nevt.sourceforge.net). Each image is output into an .html file for viewing from your desktop or presented through a web server.

# How to use

## Installation

```sh
go get github.com/rqme/neat
```

## Proof-of-concept experiments

Inside the github.com/rqme/neat/x/proofs direcory are a series of experiments. I have tried to include at least one for each new feature of the library, usually from (or based on) the one the feature's creator used. Each experiment is set up to run as a series of indpendent trials with the results displayed in the console. 

Feature         | Experiment  | Use check-stop flag (see below)
----------------|-------------|--------------------------------
NEAT            | XOR         | yes
NEAT            | Double Pole | yes
Phased Mutation | OCR         | no

### To build

```sh
go build github.com/rqme/neat/x/proof/xor
```

###To run

```sh
xor --config-path "." --archive-path "/tmp" --archive-name "xor" --web-path "/tmp" --check-stop
```
There is a configuration file in each. Place this in the archive-path or, preferrably, config-path directory.

#### Command-line flags

Flag         | Default | Description
-------------|---------|------------------------------------------------------------------------------------------
archive-path | ""      | the directory to which generational settings and state will be written
archive-name | ""      | prefix for the archive files
config-path  | ""      | overrides the archive-path. used to restore settings and state from a different location
web-path     | ""      | the directory where html files from the web visualizer will be written
trials       | 10      | the number of trials to run
check-stop   | false   | consider not meeting the stop condition a failure. 
duration     | 90      | maximum number of minutes to run a trial. used by OCR only
velocity     | false   | include velocity (Markov) in inputs. used by double pole only

Note: as the archiving process writes out all settings, including zero files, it is advisable to use a config-path to store an initial settings file to ensure it is not overwritten. This is especially important if setting traits are used as original settings will be overwritten during evolution.

