package main


import (
  "fmt"
  "github.com/slitchfield/go_nn/readdata"
  "github.com/azer/logger"
  "time"
  "flag"
  "errors"
  "strconv"
  "strings"
  )

func check(e error) {
  if e != nil {
    panic(e)
  }
}

var log = logger.New("go_nn")

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
var layerFlag layer_spec = []int{2, 4, 1}

func init() {

  flag.StringVar(&dataFile, "f", "/home/slitchfield3/go/src/github.com/slitchfield/go_nn/training_data/xor.dat", "file containing training samples")
  flag.IntVar(&numLayers, "n", 3, "the number of layers in the neural net")
  flag.Var(&layerFlag, "layers", "comma-separated list of numbers of neurons in each layer")
}

/*----------------------------------------------------------------------------*/

func main() {

  //init()

  flag.Parse()

  log.Info("Starting the NN at %v", time.Now())
  log.Info("Received datafile value of \"%s\"", dataFile)
  log.Info("Received number of layers of \"%d\"", numLayers)
  log.Info("Received layer_spec value of \"%v\"", layerFlag)

  // First, make sure the number of layers accords with the received layerspec

  if numLayers != len(layerFlag) {
    log.Error("Number of layers provided should accord with layer_spec, if provided")
    return
  }

  msgs := make(chan string)
  samps := make(chan *readdata.Sample, 2)

  go readdata.GetTrainingSamps(dataFile, msgs, samps)

  GatherSamps:
    for {
      select {
      case samp := <-samps:
        log.Info("Samp: \"%s\"", samp.Print())
      case msg := <-msgs:
        fmt.Println(msg)
        break GatherSamps
      }
    }

}
