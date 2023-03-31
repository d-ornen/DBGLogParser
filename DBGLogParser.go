package DBGLogParser

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

// Execute ONLY if there is no possibility to recover from error.
// If there is, prefer to return error structure instead
func check(e error) {
  if e != nil {
    panic(e)
  }
}

type stepSnapshot struct {
  instructionMnemonic []string
  memoryOps           []string
  registerOps         []string
}

// data  stored in JSON header of tracefile
type TraceJson struct{
  Ver int
  Arch string
  HashAlgorithm string
  Hash string
  Compression string
  Path string
}

type traceCtx struct {
  registerNames  []string
  registerValues []int
  parsedSteps    []stepSnapshot
  jsonData       *TraceJson
}

// wtf is bblock??
type bBlock struct{
  blockType uint8
  registerChanges uint8
  memoryAccesses uint8
  blockFlagsAndOpcodeSize uint8
  threadID uint32
  opcode []uint8
  registerChangeNewData []uint32
  memoryAccessFlags []uint8
  memoryAccessAddress []uint32
  memoryAccessOldData []uint32
  memoryAcessNewData []uint32
}

func checkMagic(data []byte) (err error) {
  // Probably in future we will expand list of trace files, so in future we will have to move magic variable to separate list with "magic constants".
  magic := []byte{'T', 'R', 'A', 'C'}
  if !bytes.Equal(data[0:4], magic) {
    err = &ERR_NOT_A_TRACE{}
  }
  return err
}

func checkJSONLength(data []byte) (length uint32, err error) {
  length = binary.LittleEndian.Uint32(data[4:8])
  if length == 0 {
    return 0, &ERR_ZERO_TRACE_LENGTH{}
  }
  // check that length of magick + dword(jsonLength) + json is >= length
  if uint32(len(data)+8) < length {
    return length, &ERR_JSON_FORMAT_ERROR{}
  }
  // If length is not zero and not less that length of data, then assume there is no errors in json
  return length, nil
}

func Parse(filename string) (ctx traceCtx, err error) {
  data, err := os.ReadFile(filename)
  check(err)
  err = checkMagic(data)
  check(err) //TOFIX: recover from error
  length, err := checkJSONLength(data)
  check(err) //TOFIX: recover from error
  jsonData := extractJSON(data, length, 8)
  ctx.jsonData = readJSON(jsonData)

  return ctx, nil
}

//func parseBinaryTraceBlock(data []byte, blocksParsed int64) (binaryBlock *bBlock) {
//  //TODO: parse binary block
//}

func extractJSON(raw []byte, length uint32, offset uint8) (json []byte) {
  // does not contain \0 at end
  json = raw[offset:length+uint32(offset)]
  return json
}

func readJSON(jsonData []byte) (*TraceJson) {
  var trace TraceJson
  err := json.Unmarshal(jsonData, &trace)

  check(err) //TOFIX: recover from error
  return &trace
}
