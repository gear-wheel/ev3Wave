package brick

import (
  "io"
  "os"
  "errors"
  "fmt"
  "strings"
  "strconv"
)

// should be more than enough
const BufferSize = 128

type SysFSReader struct {
  file *os.File
  path string
}

type SysFSWriter struct {
  file *os.File
  path string
}

func OpenSysFSReader(path string) (SysFSReader, error) {
  if file, err := os.Open(path); err == nil {
    return SysFSReader{file, path}, nil
  } else {
    return SysFSReader{}, err
  }
}

func OpenSysFSWriter(path string) (SysFSWriter, error) {
  if file, err := os.OpenFile(path, os.O_WRONLY, 0); err == nil {
    return SysFSWriter{file, path}, nil
  } else {
    return SysFSWriter{}, err
  }
}

func (sysFSReader * SysFSReader) Read(p []byte) (int, error) {
  sysFSReader.file.Seek(0, io.SeekStart)
  return sysFSReader.file.Read(p)
}

func (sysFSWriter * SysFSWriter) Write(p []byte) (int, error) {
  sysFSWriter.file.Seek(0, io.SeekStart)
  return sysFSWriter.file.Write(p)
}

func ReadOnce(path string) (string, error) {
  if file, err := os.Open(path); err == nil {
    buffer := make([] byte, BufferSize)

    if n, err := file.Read(buffer); err == nil {
      var data string

      if buffer[n-1] == '\n' {
        data = string(buffer[:(n - 1)])
      } else {
        data = string(buffer[:n])
      }

      return data, file.Close()
    } else {
      file.Close()
      return "", errors.New(fmt.Sprintf("[reading %s]: %v", path, err.Error()))
    }
  } else {
    return "", errors.New(fmt.Sprintf("[reading %s]: %v", path, err.Error()))
  }
}

func ReadIntOnce(path string) (int, error) {
  value, err := ReadOnce(path)
  if err == nil {
    return strconv.Atoi(value), nil
  } else {
    return 0, err
  }
}

func WriteOnce(path string, data string) error {
  if file, err := os.OpenFile(path, os.O_WRONLY, 0); err == nil {
    if data[len(data) - 1] != '\n' {
      data = data + "\n"
    }
    _, err := file.WriteString(data + "\n")
    if err != nil {
      file.Close()
      return err
    } else {
      return file.Close()
    }
  } else {
    return errors.New(fmt.Sprintf("[writing %s]: %v", path, err.Error()))
  }
}

func WriteIntOnce(path string, value int) error {
  return WriteOnce(path, fmt.Sprintf("%v\n", value))
}

