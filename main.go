package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
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
        fmt.Fprintf(w, "%s", Copyright())
        poscar := r.FormValue("poscar")
        symprecStr := r.FormValue("symprec")
        symprec, err := strconv.ParseFloat(symprecStr, 64)
        fmt.Fprintf(w, "Poscar = \n%s\n", poscar)
        fmt.Fprintf(w, "Symprec = %s\n", symprecStr)

        spacegroup, err := FindSpacegroup(poscar, symprec)
        if err != nil {
          fmt.Fprintf(w, "findsym(%s, %v) err: %v", poscar, symprec, err)
        }
        fmt.Fprintf(w, "Spacegroup: %s\n", spacegroup)

        operations, err := FindSymmetry(poscar, symprec)
        if err != nil {
          fmt.Fprintf(w, "findoperations(%s, %v) err: %v", poscar, symprec, err)
        }
        fmt.Fprintf(w, "Operations is:\n%s", operations)

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
