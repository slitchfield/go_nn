package main

/*----------------------------------------------------------------------------*/
// Generic Preamble Stuff

import (
  "fmt"
  "github.com/slitchfield/go_nn/readdata"
  "github.com/azer/logger"
  "time"
  "flag"
  "errors"
  "strconv"
  "strings"
  "math"
  "math/rand"
  )

func check(e error) {
  if e != nil {
    panic(e)
  }
}

var log = logger.New("go_nn")
var samp_log = logger.New("samples")

/*----------------------------------------------------------------------------*/
// Declare stuff to do with input parsing

type layer_spec []int

func (l *layer_spec) String() string{
  return fmt.Sprint(*l)
}

func (l *layer_spec) Set(value string) error {
  if len(*l) > 0{
    return errors.New("layer_spec flag already set")
  }
  for _, layer := range strings.Split(value, ",") {
    val, err := strconv.Atoi(layer)
    if err != nil{
      return err
    }
    *l = append(*l, val)
  }
  return nil
}

var dataFile string
var numLayers int
var layerFlag layer_spec = []int{2, 4, 4, 1}

func init() {

  flag.StringVar(&dataFile, "f", "/home/slitchfield3/go/src/github.com/slitchfield/go_nn/training_data/xor.dat", "file containing training samples")
  flag.IntVar(&numLayers, "n", 4, "the number of layers in the neural net")
  flag.Var(&layerFlag, "layers", "comma-separated list of numbers of neurons in each layer")
}

/*----------------------------------------------------------------------------*/
// Neural Net Work



func myActivationFunc(in float32) float32 {
  return float32(math.Tanh(float64(in)))
}

func myGradFunc(in float32) float32 {
  return 1.0 - float32(math.Tanh(float64(in)))
}

type Neuron struct {
  activationFunc func(float32) (float32)
  gradFunc func(float32) (float32)
  inputWeights []float32
  output float32
  myLayer int
}

func (n *Neuron) GetOutput() float32 {
  return n.output
}

func (n *Neuron) Print() {
  fmt.Printf("ActivationFunc: %v\tGradFunc: %v\n", n.activationFunc, n.gradFunc)
  fmt.Printf("InputWeights: %f\n", n.inputWeights)
  fmt.Printf("Output: %f\n", n.output)
  fmt.Printf("My Layer: %d\n", n.myLayer)
}

type Layer []*Neuron

type Net struct {
  shape []int
  layers []Layer
}

func (net *Net) FeedFwd(samp *readdata.Sample) []float32 {

  inputs := samp.GetInputs()

  for i, _ := range(inputs) {
    net.layers[0][i].output = inputs[i]
  }

  // For every layer other than the input...
  for i := 1; i < len(net.layers); i++ {

    // For every neuron in that layer, that isn't the bias neuron
    for j := 0; j < len(net.layers[i]) - 1; j++ {

      // Accumulate!
      accum := float32(0.0)
      for k := 0; k <= net.shape[i-1]; k++ {
        accum += net.layers[i-1][k].GetOutput()*net.layers[i][j].inputWeights[k]
      }
      net.layers[i][j].output = net.layers[i][j].activationFunc(accum)
    }
  }

  // We should have all the outputs set now! Let's collect the final output in a slice
  outputs := make([]float32, len(net.layers[len(net.layers) - 1]))

  for i, _ := range(outputs) {
    outputs[i] = net.layers[len(net.layers) - 1][i].output
  }

  return outputs
}

func main() {

  flag.Parse()
  r := rand.New(rand.NewSource(99))

  log.Info("Starting the NN at %v", time.Now())
  log.Info("Received datafile value of \"%s\"", dataFile)
  log.Info("Received number of layers of \"%d\"", numLayers)
  log.Info("Received layer_spec value of \"%v\"", layerFlag)

  // First, make sure the number of layers accords with the received layerspec

  if numLayers != len(layerFlag) {
    log.Error("Number of layers provided should accord with layer_spec, if provided")
    return
  }

  network := new(Net)
  network.shape = layerFlag // set the shape of the layers
  network.layers = make([]Layer, len(layerFlag))
  for i := 0; i < len(network.shape); i++ {
    network.layers[i] = make([]*Neuron, network.shape[i] + 1)  // Account for bias!
    for j := 0; j <= network.shape[i]; j++ {
      network.layers[i][j] = new(Neuron)
      network.layers[i][j].activationFunc = myActivationFunc
      network.layers[i][j].gradFunc = myGradFunc
      network.layers[i][j].myLayer = i

      // If this is the 0th layer, input weights is a slice of 1
      if i == 0 {
        network.layers[i][j].inputWeights = []float32{1.0}
      } else {
        network.layers[i][j].inputWeights = make([]float32, network.shape[i-1]+1)
        for k, _ := range(network.layers[i][j].inputWeights) {
          network.layers[i][j].inputWeights[k] = r.Float32()
        }
      }

      // If we are the bias neuron, we have no inputs, and our output is one
      if j == network.shape[i] {
        network.layers[i][j].inputWeights = []float32{}
        network.layers[i][j].output = 1.0
      }

    }
    /*
    log.Info("Created layer %d: \"%v\"", i, network.layers[i])
    for j := 0; j < network.shape[i]; j++ {
      network.layers[i][j].Print()
    }
    */
  }

  msgs := make(chan string)
  samps := make(chan *readdata.Sample, 2)

  go readdata.GetTrainingSamps(dataFile, msgs, samps)

  GatherSamps:
    for {
      select {
      case samp := <-samps:
        log.Info("Samp: \"%s\"", samp.Print())
        log.Info("Fed forward! \"%f\"", network.FeedFwd(samp))
        break GatherSamps
      case msg := <-msgs:
        fmt.Println(msg)
        break GatherSamps
      }
    }

}
