# Update
This project has been [discontinued](https://medium.com/@hummerb/evo-by-klokare-new-library-same-concept-9eff96126ec0#.rywgvow3a) here and continued under a new [project called EVO](https://github.com/klokare/evo) by klokare. 



RedQ.NEAT
==========
This a Go implementation of NeuralEvolution of Augmenting Topologies (NEAT). From the [NEAT F.A.Q](http://www.cs.ucf.edu/~kstanley/neat.html#FAQ1).

*NEAT stands for NeuroEvolution of Augmenting Topologies. It is a method for evolving artificial neural networks with a genetic algorithm. NEAT implements the idea that it is most effective to start evolution with small, simple networks and allow them to become increasingly complex over generations. That way, just as organisms in nature increased in complexity since the first cell, so do neural networks in NEAT. This process of continual elaboration allows finding highly sophisticated and complex neural networks.*

More information will be provided on the blog [redq.me](http://www.redq.me).

# Installation
```sh
go get github.com/rqme/neat
```

# Usage
The API documentation can be found at [GoDoc](http://godoc.org/github.com/rqme/neat).

The Context and Experiment are the central components of the library. The latter encapsulates everything needed for execution and the former provides access to all the necessary helpers. There are several convenience functions in the starter package.

RedQ.NEAT includes several demonstration experiments, each built at the onset of adding a new feature (like [phased mutation](http://sharpneat.sourceforge.net/phasedsearch.html)) or concept (like [Novelty Search](http://eplex.cs.ucf.edu/noveltysearch/userspage/)). These proof-of-concepts are intended to valid this library with the idea being tested as well as compare different helpers (such as HyperNEAT vs regular NEAT). The experiments are each in their own package in the x/experiments directory.

## Running experiments
Each experiment builds off the trials package which provides a way to compare multiple runs of an experiment against each other. This package provides several command line arguments that are common to all experiments and displays its output in the console window. For example, here is the output of the XOR experiment:

```sh
$ xor --check-stop --trials 40
Run   Iters.   Seconds    Nodes     Conns    Fitness   Fail   Comment 
--- --------- --------- --------- --------- --------- ------ ---------
  0        28     1.339         9        16    16.000        
  1        26     1.192         8        14    15.443        
  2        59     3.384         7        17    16.000        
  3        59     3.609        14        28    16.000        
...
 36        45     2.513        11        20    16.000        
 37        30     1.265         7        12    12.250        
 38        28     1.246         9        17    16.000        
 39        19     0.822         6        12    13.930        

Summary for trials excluding failures (and time for skipped)
      Iters.   Seconds    Nodes     Conns    Fitness
--- --------- --------- --------- --------- ---------
AVG        34     1.769         9        17    14.894        
MED        32     1.625         9        16    15.996        
SDV        13     0.905         2         5     1.637        
MIN         9     0.303         5         7    10.782        
MAX        66     4.503        15        33    16.000  
```

### Common command-line arguments
flag | description | default
-----|-------------|------------
config-name | Common name used as a prefix to archive files | defaults to the name of the executable
config-path | Directory containing the initial configuration file and, if available, state files | Current directory
trials | The number of trial runs to perform | 10
check-stop | Experiments which do not end with an explicit stop are considered to have failed. | false 
show-work | Informs the Evaluator (if it implements Demonstrable) to show its work during evaluation. This is used only for the best genome. | false 
skip-evolve | Skips evolution and only performs summary of archived runs. Best used with --show-work and setting the config-path to the ArchivePath used in the settings file. | false

## Experiments
### XOR
[Exclusive OR](https://en.wikipedia.org/wiki/Exclusive_or), or XOR for short, is the starter experiment to verify the NEAT (called Classic in RedQ.NEAT) functionality. Located in the x/examples/xor directory, the package produces a standalone executable file. A configuration file, xor-config.json, is provided. 

The experiment provides no new command-line arguments but it is recommended to use --check-stop when running to catch trials that do not produce a solution. 

RedQ.NEAT was able to find a solution in 40 out of 40 trials. The median number of nodes and connections were 9 and 16 respectively. The results of this experiment are detailed in the [wiki](https://github.com/rqme/neat/wiki/XOR-experiment-results).

# Background
The core of this library, often called Classic in the code, was written from the ground up using Dr. Kenneth Stanley's [PhD dissertation](http://nn.cs.utexas.edu/keyword?stanley:phd04) as a guide. NEAT has changed a bit since that paper and I have made some adjustments based on the F.A.Q. I have also add some flexibility in the design to allow for growing the library via helpers which will provide for adding HyperNEAT, Novelty Search, etc. to the library without changing the core API.

The library and proof-of-concept experiments utilize SVG to visualize the network of the best genome as well as the experiment's history. This visualization is based on the [NeuroEvolution Visualization Toolkit (NEVT)](http://nevt.sourceforge.net). Each image is output into an .html file for viewing from your desktop or presented through a web server.



