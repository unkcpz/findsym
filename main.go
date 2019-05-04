package main

import (
  "io/ioutil"
  "fmt"

  cmpio "github.com/unkcpz/gocmp/io"
  "github.com/unkcpz/gocmp/crystal"
  "github.com/unkcpz/spglib/go/spglib"
  "gonum.org/v1/gonum/mat"
)

func main() {
  data, err := ioutil.ReadFile("./POSCAR")
  if err != nil {
    panic(err)
  }

  s := string(data)
  poscarCell, err := cmpio.ParsePoscar(s)
  if err != nil {
    panic(err)
  }

  lattice := poscarCell.Lattice
  positions := poscarCell.Positions
  types := make([]int, len(poscarCell.Types), len(poscarCell.Types))
  for i, _ := range types {
    types[i] = crystal.SymToNum(poscarCell.Types[i])
  }
  var c bool = false
  if poscarCell.Coordinate == cmpio.Cartesian {
    c = true
  }
  cell, err := crystal.NewCell(lattice, positions, types, c)
  if err != nil {
    panic(err)
  }

  ds := spglib.NewDataset(
    MatData(*cell.Lattice),
    MatData(*cell.Positions),
    cell.Types,
    1e-5,
  )
  fmt.Printf("SpaceSymbol: %s\n", ds.SpaceSymbol)
}

func MatData(mat mat.Dense) []float64 {
  blasM := mat.RawMatrix()
  return blasM.Data
}
