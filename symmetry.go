package main

import (
  "fmt"
  "bytes"

  cmpio "github.com/unkcpz/gocmp/io"
  "github.com/unkcpz/gocmp/crystal"
)

func GetCell(poscar string) (*crystal.Cell, error) {
  poscarCell, err := cmpio.ParsePoscar(poscar)
  if err != nil {
    return nil, err
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
    return nil, err
  }
  return cell, nil
}

func FindSpacegroup(poscar string, eps float64) (string, error) {
  cell, err := GetCell(poscar)
  if err !=  nil {
    return "", err
  }

  return cell.Spacegroup(eps), nil
}

func FindSymmetry(poscar string, eps float64) (string, error) {
  cell, err := GetCell(poscar)
  if err !=  nil {
    return "", err
  }

  _, rots, trans := cell.Symmetry(eps)
  var buffer bytes.Buffer
  for _, r := range rots {
    buffer.WriteString(fmt.Sprintf("%v\n", r))
  }
  for _, t := range trans {
    buffer.WriteString(fmt.Sprintf("%v\n", t))
  }
  return buffer.String(), nil
}
