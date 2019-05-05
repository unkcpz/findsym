package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"

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
        fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
        poscar := r.FormValue("poscar")
        symprecStr := r.FormValue("symprec")
        symprec, err := strconv.ParseFloat(symprecStr, 64)
        fmt.Fprintf(w, "Poscar = %s\n", poscar)
        fmt.Fprintf(w, "Symprec = %s\n", symprecStr)

        symdata, err := findsym(poscar, symprec)
        if err != nil {
          fmt.Fprintf(w, "findsym(%s, %v) err: %v", poscar, symprec, err)
        }
        fmt.Fprintf(w, "Spacegroup symbol: %s", symdata)
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

func findsym(poscar string, eps float64) (string, error) {
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
