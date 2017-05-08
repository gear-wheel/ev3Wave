package brick

import (
  "fmt"
  "io/ioutil"
  "strings"
  "strconv"
)

const SensorRoot = "/sys/class/lego-sensor/"
const MaxSensors = 4

type Sensor struct {
  path string
  name string
  port string
  mode string
  numValues int
  pollMs int
}

func (x * Sensor) String() string {
  return fmt.Sprintf(
    "Sensor %v [%v:%v], mode %v [%v values each %v ms]",
    x.name, x.path, x.port, x.mode, x.numValues, x.pollMs)
}

func getSensor(path string) (Sensor, error) {
  port, err := readOnce(path + "/address")
  if err != nil { return nil, err }

  name, err := readOnce(path + "/driver_name")
  if err != nil { return nil, err }

  mode, err := readOnce(path + "/mode")
  if err != nil { return nil, err }

  numValues, err := readOnce(path + "/num_values")
  if err != nil { return nil, err }

  pollMsString, err := readOnce(path + "/pool_ms")
  if err != nil { return nil, err }
  pollMs, err := strconv.Atoi(pollMsString)
  if err != nil { return nil, err }

  return Sensor{path, name, port, mode, numValues, pollMs}
}

func DiscoverSensorsInPath(root string) ([]Sensor, error) {
  if files, err := ioutil.ReadDir(root); err == nil {
    sensors := make([]Sensor, MaxSensors)
    i := 0

    for _, file := range files {
      if file.IsDir() && strings.HasPrefix(file.Name(), "sensor") {
        if sensor, err := getSensor(file.Name()); err == nil {
          sensors[i] = sensor
          i = i + 1
        }
      }
    }

    return sensors[:i]
  } else {
    return nil, err
  }
}

func DiscoverSensors() ([]Sensor, error) {
  return DiscoverSensorsInPath(SensorRoot)
}

func Senses(modes map[string] string, ) {

}



