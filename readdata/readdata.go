package readdata

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  "strings"
  )

type Sample struct {
  input_len int
  output_len int
  inputs []float32
  outputs []float32
}

func (s *Sample) GetInputs() []float32 {
  return s.inputs
}

func (s *Sample) Print () string {
  return fmt.Sprintf("%d Inputs: %f|%d Outputs: %f", s.input_len, s.inputs, s.output_len, s.outputs)
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func Line2Samp(line string, input_len int, output_len int) *Sample {

  // Declare an initialize the sample
  cur_samp := new(Sample)
  cur_samp.input_len = input_len
  cur_samp.output_len = output_len
  cur_samp.inputs = make([]float32, cur_samp.input_len)
  cur_samp.outputs = make([]float32, cur_samp.output_len)

  // Parse the line according to the specified number of inputs/outputs
  // For now, anything in a sample is guaranteed to be float32, or interpretable that way

  // First split the line on ':'
  splits := strings.Split(line, ":")

  // Now split what we know as the inputs.
  inputs := strings.Split(strings.Trim(splits[0], " "), " ")
  for i := 0; i < cur_samp.input_len; i++ {
    input, err := strconv.ParseFloat(inputs[i], 32)
    check(err)
    cur_samp.inputs[i] = float32(input)
  }

  // Split what we know as the outputs, and stick them in the struct
  outputs := strings.Split(strings.Trim(splits[1], " "), " ")
  for i := 0; i < cur_samp.output_len; i++ {
    output, err := strconv.ParseFloat(outputs[i], 32)
    check(err)
    cur_samp.outputs[i] = float32(output)
  }

  return cur_samp

}

func GetTrainingSamps(filename string, msgs chan<- string, samps chan<- *Sample) {
  // Open the data file, and defer the close
  file, err := os.Open(filename)
  check(err)
  defer file.Close()

  // Declare our scanner, and read the number of samples
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  num_samps, err := strconv.Atoi(scanner.Text())
  check(err)


  // Read the input/output definition line, and parse it
  scanner.Scan()
  splits := strings.Split(scanner.Text(), ":")
  inputs := strings.Split(splits[0], " ")
  input_len := len(inputs)
  outputs := strings.Split(strings.Trim(splits[1], " "), " ")
  output_len := len(outputs)

  // Start reading samples!
  for i := 0; i < num_samps; i++ {
    scanner.Scan()
    line := scanner.Text()
    sample := Line2Samp(line, input_len, output_len)
    samps <- sample
  }

  msgs <- "Goroutine is done!"
}
