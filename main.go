package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "bytes"

    cmpio "github.com/unkcpz/gocmp/io"
    "github.com/unkcpz/gocmp/crystal"
    "github.com/unkcpz/spglib/go/spglib"
    "gonum.org/v1/gonum/mat"
)
func proc(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    switch r.Method {
    case "GET":
         http.ServeFile(w, r, "form.html")
    case "POST":
        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
        if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        // fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
        fmt.Fprintf(w, "Hello Dr.Yang\n")
        poscar := r.FormValue("poscar")
        symprecStr := r.FormValue("symprec")
        symprec, err := strconv.ParseFloat(symprecStr, 64)
        fmt.Fprintf(w, "Poscar = %s\n", poscar)
        fmt.Fprintf(w, "Symprec = %s\n", symprecStr)

        spacegroup, err := findSpacegroup(poscar, symprec)
        if err != nil {
          fmt.Fprintf(w, "findsym(%s, %v) err: %v", poscar, symprec, err)
        }
        fmt.Fprintf(w, "Spacegroup symbol: %s\n", spacegroup)

        operations, err := findOperations(poscar, symprec)
        if err != nil {
          fmt.Fprintf(w, "findoperations(%s, %v) err: %v", poscar, symprec, err)
        }
        fmt.Fprintf(w, "Operations is:\n %s", operations)

    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }
}

func main() {
    http.HandleFunc("/", proc)

    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":6116", nil); err != nil {
        log.Fatal(err)
    }
}

func findSpacegroup(poscar string, eps float64) (string, error) {
  poscarCell, err := cmpio.ParsePoscar(poscar)
  if err != nil {
    return "", err
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
    return "", err
  }

  ds := spglib.NewDataset(
    MatData(*cell.Lattice),
    MatData(*cell.Positions),
    cell.Types,
    eps,
  )
  return ds.SpaceSymbol, nil
}

func MatData(mat mat.Dense) []float64 {
  blasM := mat.RawMatrix()
  return blasM.Data
}

func findOperations(poscar string, eps float64) (string, error) {
  poscarCell, err := cmpio.ParsePoscar(poscar)
  if err != nil {
    return "", err
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
    return "", err
  }

  ds := spglib.NewDataset(
    MatData(*cell.Lattice),
    MatData(*cell.Positions),
    cell.Types,
    eps,
  )

  rots := ds.Rotations
  trans := ds.Translations

  var buffer bytes.Buffer
  buffer.WriteString(fmt.Sprintf("Number of Operations: %d\n", ds.Nops))
  buffer.WriteString("Rotations:\n")
  for _, rot := range rots {
    for i:=0; i<3; i++ {
      for j:=0; j<3; j++ {
        s := fmt.Sprintf("%d ", int(rot.At(i,j)))
        buffer.WriteString(s)
      }
      buffer.WriteString("\n")
    }
    // buffer.WriteString("\n")
  }

  buffer.WriteString("Translations:\n")
  for _, tran := range trans {
    for i:=0; i<3; i++ {
      s := fmt.Sprintf("%6.4f ", tran.At(i, 0))
      buffer.WriteString(s)
    }
    buffer.WriteString("\n")
  }

  return buffer.String(), nil
}

// func findPermOpMatrix(poscar string, eps float64) (string, error) {
//   poscarCell, err := cmpio.ParsePoscar(poscar)
//   if err != nil {
//     return "", err
//   }
//
//   lattice := poscarCell.Lattice
//   positions := poscarCell.Positions
//   types := make([]int, len(poscarCell.Types), len(poscarCell.Types))
//   for i, _ := range types {
//     types[i] = crystal.SymToNum(poscarCell.Types[i])
//   }
//   var c bool = false
//   if poscarCell.Coordinate == cmpio.Cartesian {
//     c = true
//   }
//   cell, err := crystal.NewCell(lattice, positions, types, c)
//   if err != nil {
//     return "", err
//   }
//
//   ds := spglib.NewDataset(
//     MatData(*cell.Lattice),
//     MatData(*cell.Positions),
//     cell.Types,
//     eps,
//   )
//
//   permTable := getPermOpMatrix(cell, ds)
//
//   var buffer bytes.Buffer
//
//   buffer.WriteString("Permutation Operations are list in each rows following:\n")
//   for i:=0; i<ds.Nops; i++ {
//     for j:=0; j<cell.Natoms; j++ {
//       buffer.WriteString("%d ", permTable[i*cell.Natoms+j])
//     }
//     buffer.WriteString("\n")
//   }
//   return buffer.String(), nil
// }

// func getPermOpMatrix(cell crystal.Cell, ds spglib.Dataset) []int {
//   rots := ds.Rotations
//   trans := ds.Translations
//
//   refinedOriginP := refinePos(cell.Positions)
//   for i, _ := range ds.Nops {
//     rot := rots[i]
//     tran := trans[i]
//     var newP mat.Dense
//     newP.Mul(refinedOriginP, rot)
//     nr, _ := newP.Caps()
//     for i:=0; i<nr; i++ {
//       r := newP.RowView(i)
//       r.Add(r, tran)
//     }
//
//     newP = refinePos(newP)
//
//
//   }
//
// }
