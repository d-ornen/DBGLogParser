package DBGLogParser

import (
  "os"
  "testing"
)

func TestWrongMagick(t* testing.T){
  wrongMagic := []byte{
    // providing elf file magick instead of TRAC
    0x7F, 0x45, 0x4C, 0x46,
    0x02, 0x01, 0x01, 0x00,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
  }
  if (checkMagic(wrongMagic) != &ERR_NOT_A_TRACE{}) {
    t.Fatalf("Check wrong magic fail.")
  }

}
func TestWrongLength(t *testing.T)  {
  wrongLength:= []byte{
    // this time correct magic file, but json length and delcared length are different.
    // As json length I consider len(file) - len(TRAC)+len(declaredJsonlength)
    'T', 'R', 'A', 'C',
    0x12, 0x1d, 0x1c, 0x22,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
  }
  _, err := checkJSONLength(wrongLength)
  if ( err != &ERR_JSON_FORMAT_ERROR{}) {
    t.Fatalf("Check wrong length fail.")
  }
}

func TestZeroLength(t *testing.T)  {
  zeroLength:= []byte {
    // this time correct magic file, but json length and delcared length are different.
    // As json length I consider len(file) - len(TRAC)+len(declaredJsonlength)
    'T', 'R', 'A', 'C',
    0x0, 0x0, 0x0, 0x0,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
  }
  _, err := checkJSONLength(zeroLength)
  if (err != &ERR_ZERO_TRACE_LENGTH{}) {
    t.Fatalf("Check zero length fail.")
  }

}

func TestFormatCheck(t *testing.T)  {
  correctFile := []byte{
    // This file contains only headers, i.e. magic, dword representing length and json with compression and other info like path to executable.
    0x54, 0x52, 0x41, 0x43, 0x8F, 0x00, 0x00, 0x00, 0x7B, 0x22, 0x76, 0x65, 0x72, 0x22, 0x3A, 0x31, 
    0x2C, 0x22, 0x61, 0x72, 0x63, 0x68, 0x22, 0x3A, 0x22, 0x78, 0x36, 0x34, 0x22, 0x2C, 0x22, 0x68, 
    0x61, 0x73, 0x68, 0x41, 0x6C, 0x67, 0x6F, 0x72, 0x69, 0x74, 0x68, 0x6D, 0x22, 0x3A, 0x22, 0x6D, 
    0x75, 0x72, 0x6D, 0x75, 0x72, 0x68, 0x61, 0x73, 0x68, 0x22, 0x2C, 0x22, 0x68, 0x61, 0x73, 0x68, 
    0x22, 0x3A, 0x22, 0x30, 0x78, 0x39, 0x30, 0x35, 0x38, 0x37, 0x31, 0x36, 0x32, 0x43, 0x43, 0x35, 
    0x39, 0x39, 0x33, 0x34, 0x39, 0x22, 0x2C, 0x22, 0x63, 0x6F, 0x6D, 0x70, 0x72, 0x65, 0x73, 0x73, 
    0x69, 0x6F, 0x6E, 0x22, 0x3A, 0x22, 0x22, 0x2C, 0x22, 0x70, 0x61, 0x74, 0x68, 0x22, 0x3A, 0x22, 
    0x5A, 0x3A, 0x5C, 0x5C, 0x68, 0x6F, 0x6D, 0x65, 0x5C, 0x5C, 0x6A, 0x75, 0x73, 0x5C, 0x5C, 0x44, 
    0x6F, 0x63, 0x75, 0x6D, 0x65, 0x6E, 0x74, 0x73, 0x5C, 0x5C, 0x43, 0x72, 0x61, 0x63, 0x6B, 0x4D, 
    0x65, 0x2E, 0x65, 0x78, 0x65, 0x22, 0x7D, 0x00, 0xAC, 0x01, 0x83, 0x64, 0x01, 0x00, 0x00, 0x48, 
    0x89, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  }

  length, err := checkJSONLength(correctFile)
  if ( err != nil || length != 0x8f) {
    t.Fatalf("Correct file parse error.")
  }

}
func TestExtractJson(t *testing.T) {
  testFile, err := os.ReadFile("test/test1.trace64")
  check(err)
  length, err := checkJSONLength(testFile)
  check(err)
  jsonBytes := extractJSON(testFile, length, 8)
  if jsonBytes[0] != '{' && jsonBytes[len(jsonBytes)-1] != '}'{
    t.Fatalf("json format error.")
  }
}

func TestReadJson(t *testing.T) {
  testFile, err := os.ReadFile("test/test1.trace64")
  check(err)
  length, err := checkJSONLength(testFile)
  check(err)
  jsonBytes := extractJSON(testFile, length, 8)
  ctx := readJSON(jsonBytes)
  if ctx.Ver != 1.{
    t.Fatalf("Version parsing error. Expected 1, got %d", ctx.Ver)
  }
  if ctx.Arch != "x64"{
    t.Fatalf("Architecture parsing error. Expected 'x64', got %s", ctx.Arch)
  }
  if ctx.Compression != ""{
    t.Fatalf("Compression parsing error.")
  }
  if ctx.Hash != "0x3A6C6FA546C8A44E"{
    t.Fatalf("Hash value parsing error.")
  }
  if ctx.HashAlgorithm != "murmurhash"{
    t.Fatalf("Hash algorithm name parsing error.")
  }
  if ctx.Path != "F:\\Factorio\\bin\\x64\\factorio.exe"{
    t.Fatalf("Path parsing error.")
  }
}

func TestParse(t *testing.T)  {
  trace, err := Parse("test/test1.trace64")
  check(err)
  if trace.jsonData == nil || trace.jsonData.Hash != "0x3A6C6FA546C8A44E"{
    t.Fatalf("Failed to parse tracefile.")
  }
  if trace.jsonData.Ver != 1 {
    t.Fatalf("Failed to parse verison.")
  }
  if trace.jsonData.Arch != "x64"{
    t.Fatalf("Failed to parse architecture.")
  }
  if trace.jsonData.HashAlgorithm != "murmurhash" {
    t.Fatalf("Failed to parse hash algorithm.")
  }
  if trace.jsonData.Compression != "" {
    t.Fatalf("Failed to parse compression.")
  }
  if trace.jsonData.Path != "F:\\Factorio\\bin\\x64\\factorio.exe" {
    t.Fatalf("Failed to parse path.")
  }
}
