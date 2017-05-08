package brick

import (
  "io"
  "os"
)

// should be more than enough
const bufferSize = 128

type SysFSReader struct {
  io.ReadSeeker
  path string
}

type SysFSWriter struct {
  io.WriteSeeker
  path string
}

type SysFSReaderWriter struct {
  io.ReadWriteSeeker
  path string
}

func OpenSysFSReader(path string) (SysFSReader, error) {
  if file, err := os.Open(path); err == nil {
    return SysFSReader{file, path : path}, nil
  } else {
    return nil, err
  }
}

func OpenSysFSWriter(path string) (SysFSWriter, error) {
  if file, err := os.OpenFile(path, os.O_WRONLY, 0); err == nil {
    return SysFSWriter{file, path : path}, nil
  } else {
    return nil, err
  }
}

func (sysFSWriter * SysFSWriter) Write(p []byte) (int, error) {
  sysFSWriter.Seek(0, io.SeekStart)
  return io.Writer(sysFSWriter).Write(p)
}

func OpenSysFSReaderWriter(path string) (SysFSReaderWriter, error) {
  if file, err := os.OpenFile(path, os.O_RDWR, 0); err == nil {
    return SysFSReaderWriter{file, path : path}, nil
  } else {
    return nil, err
  }
}

func (sysFSReaderWriter * SysFSReaderWriter) Read(p []byte) (int, error) {
  return SysFSReader(sysFSReaderWriter).Read(p)
}

func (sysFSReaderWriter * SysFSReaderWriter) Write(p []byte) (int, error) {
  return SysFSWriter{}(sysFSReaderWriter).Write(p)
}

func readOnce(path string) (string, error) {
  if file, err := os.Open(path); err == nil {
    buffer := make([] byte, bufferSize)

    if n, err := file.Read(buffer); err == nil {
      return string(buffer[:n]), nil
    } else {
      return nil, err
    }
  } else {
    return nil, err
  }
}


