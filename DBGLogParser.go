package DBGLogParser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)
type ERR_NOT_A_TRACE struct{}

func (m *ERR_NOT_A_TRACE) Error() string {
  return "Please, provide a trace file"
}
type ERR_ZERO_TRACE_LENGTH struct{}

func (m *ERR_ZERO_TRACE_LENGTH) Error() string {
  return "Trace file seems to be empty"
}

type ERR_JSON_FORMAT_ERROR struct{}

func (m *ERR_JSON_FORMAT_ERROR) Error() string {
  return "json seems to be corrupted"
}

// Execute ONLY if there is no possibility to recover from error. If there is, prefer to return error structure instead
func check(e error)  {
  if e != nil{
    panic(e)
  }
}

type stepSnapshot struct {
  instructionMnemonic []string
  memoryOps []string
  registerOps []string
}

type traceCtx struct{
  registerNames []string
  registerValues []int
  parsedSteps []stepSnapshot

}

func checkMagic(data []byte) (err error) {
  // Probably in future we will expand list of trace files, so maybe in future we will have to move magic variable to separate list with "magic constants".
  magic := []byte{'T', 'R', 'A', 'C'}
  if !bytes.Equal(data[0:3], magic) {
    err = &ERR_NOT_A_TRACE{}
  }
  return
}

func checkJSONLength(data []byte) (length uint32, err error) {
  length = binary.LittleEndian.Uint32(data[4:7])
  if length==0 {
    return 0, &ERR_ZERO_TRACE_LENGTH{}
  } 
  if (data[length+8]!=0) | (uint32(len(data)) < length+8) {
    // I expect \0 after the end of enclosing bracket
    return length, &ERR_JSON_FORMAT_ERROR{}
  }
  // If length is not zero and \0 after bracket is present, then assume there is no errors in json
  return length, nil
}

func Parse(filename string) (ctx traceCtx, err error){
  data, err := os.ReadFile(filename) 
  check(err)
  if err = checkMagic(data); err != nil {
    return ctx, err
  }

  fmt.Println(jsonLength)
  return ctx, nil
}
