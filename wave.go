package main

import (
  "fmt"
  "gear-wheel/ev3Wave/brick"
)

func main() {
  sensors, err := brick.DiscoverSensors()
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(sensors)
  }
}
