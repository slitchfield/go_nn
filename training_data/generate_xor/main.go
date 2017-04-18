package main


import "fmt"
import "math/rand"

func main() {

  num_samps := 1000
  var_string := "x1 x2: y"

  r := rand.New(rand.NewSource(99))

  fmt.Printf("%d\n", num_samps)
  fmt.Println(var_string)

  for i := 0; i < num_samps; i++ {
    x1 := r.Intn(2)
    x2 := r.Intn(2)
    y := x1^x2

    fmt.Printf("%d %d: %d\n", x1, x2, y)
  }
}
