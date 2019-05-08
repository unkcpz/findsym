package main

import (
  "fmt"
)

func Copyright() string {
  cr :=
`----------------------------------------------------------------------
This package using BSD-3-Clause

Copyright (c) 2019, Jason Eu
All rights reserved.
----------------------------------------------------------------------`
  return fmt.Sprintf("%s\n", cr)
}
