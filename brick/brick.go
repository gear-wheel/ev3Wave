package brick

import (
  "fmt"
  "io/ioutil"
  "strings"
  "strconv"
  "path"
  "time"
)

const SensorRoot = "/sys/class/lego-sensor/"
const MaxSensors = 4

type Sensor struct {
  Path      string
  Name      string
  Port      string
}

func (x * Sensor) String() string {
  return fmt.Sprintf(
    "Sensor %v [%v:%v]",
    x.Name, x.Path, x.Port)
}

func (x * Sensor) OpenValueReader(n int) (SysFSReader, error) {
  p := path.Join(x.Path, strconv.Itoa(n))
  return OpenSysFSReader(p)
}

func (x * Sensor) SetMode(mode string) error {
  return WriteOnce(path.Join(x.Path, "mode"), mode)
}

func (x * Sensor) Mode() (string, error) {
  return ReadOnce(path.Join(x.Path, "mode"))
}

func (x * Sensor) SetPollMs(pollMs int) error {
  return WriteIntOnce(path.Join(x.Path, "poll_ms"), pollMs)
}

func (x * Sensor) PollMs() (int, error) {
  return ReadIntOnce(path.Join(x.Path, "poll_ms"))
}

func (x * Sensor) NumOfValues() (int, error) {
  return ReadIntOnce(path.Join(x.Path, "num_values"))
}

func getSensor(path string) (Sensor, error) {
  port, err := ReadOnce(path + "/address")
  if err != nil { return Sensor{}, err }

  name, err := ReadOnce(path + "/driver_name")
  if err != nil { return Sensor{}, err }

  mode, err := ReadOnce(path + "/mode")
  if err != nil { return Sensor{}, err }

  numValuesString, err := ReadOnce(path + "/num_values")
  if err != nil { return Sensor{}, err }
  numValues, err := strconv.Atoi(numValuesString)
  if err != nil { return Sensor{}, err }

  pollMsString, err := ReadOnce(path + "/pool_ms")
  if err != nil { return Sensor{}, err }
  pollMs, err := strconv.Atoi(pollMsString)
  if err != nil { return Sensor{}, err }

  return Sensor{path, name, port, mode, numValues, pollMs}, nil
}

func DiscoverSensorsInPath(root string) ([]Sensor, error) {
  if files, err := ioutil.ReadDir(root); err == nil {
    fmt.Printf("Sensor discovery in %v;\n", root)
    fmt.Println(files)

    sensors := make([]Sensor, MaxSensors)
    i := 0

    for _, file := range files {
      fmt.Printf("Checking %v;\n", file)

      if strings.HasPrefix(file.Name(), "sensor") {
        sensor_path := path.Join(root, file.Name())

        if sensor, err := getSensor(sensor_path); err == nil {
          sensors[i] = sensor
          i = i + 1
          fmt.Printf("%v recognised as %v:%v", sensor_path, sensor.Name, sensor.Port)
        } else {
          fmt.Printf("Error in sensor detection: %v", err)
        }
      }
    }

    return sensors[:i], nil
  } else {
    return nil, err
  }
}

func DiscoverSensors() ([]Sensor, error) {
  return DiscoverSensorsInPath(SensorRoot)
}

func CaptureSensor(sensor * Sensor,
  mode string, pollMs int,
  data chan<- [] int, errors chan<- error) {

  check_err := func (err error) bool {
    if err != nil {
      errors <- err
      close(data)
      return true
    } else {}
    return false
  }

  buffer := make([] byte, BufferSize)

  var err error

  if check_err(sensor.SetMode(mode)) { return }
  if check_err(sensor.SetPollMs(pollMs)) { return }


  numValues, err := sensor.NumOfValues()
  if check_err(err) { return }

  buffer := make([]byte, BufferSize)
  readings := make([] int, numValues)
  readers := make([] * SysFSReader, numValues)

  for i := 0; i < numValues; i++ {
    readers[i], err = sensor.OpenValueReader(i)

    if err == nil {
      errors <- err
      close(data)
      return
    }
  }

  for err == nil {
    var n int
    for i := 0; i < numValues; i++ {
      n, err = readers[i].Read(buffer)
      if err != nil {
        errors <- err
        close(data)
        return
      }
      readings[i] = strconv.Itoa(string(buffer[n - 1]))
    }
  }
}



