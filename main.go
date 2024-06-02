package main

import (
	"blockchain/block"
    "fmt"
    "encoding/gob"
    "bytes"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

func main() {
    rsa, err := encoder.NewRSAEncoder(5000)
    if err != nil {
        panic(err)
    }

    header := block.NewHeader()
    header.Store["dir"] = []byte("huy")
    b := block.New(rsa, block.NewHeader())

    b.Data = []byte("hello world")

    err = b.Encode()
    if err != nil {
        panic(err)
    }

    buffer := bytes.Buffer{}
    enc := gob.NewEncoder(&buffer)
    err = enc.Encode(b)
    if err != nil {
        panic(err)
    }

    blockBytes := buffer.Bytes()
    half := len(blockBytes) / 2
    half1 := blockBytes[:half]
    half2 := blockBytes[half:]

    encodedHalf1, err := rsa.Encode(half1)
    if err != nil {
        panic(err)
    }
    encodedHalf2, err := rsa.Encode(half2)
    if err != nil {
        panic(err)
    }

    decodedHalf1, err := rsa.Decode(encodedHalf1)
    if err != nil {
        panic(err)
    }
    decodedHalf2, err := rsa.Decode(encodedHalf2)
    if err != nil {
        panic(err)
    }
    
    var decodedBlock block.Block
    decodedBlockBytes := bytes.NewBuffer(append(decodedHalf1, decodedHalf2...))
    err = gob.NewDecoder(decodedBlockBytes).Decode(&decodedBlock)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(decodedBlock.Data) == string(b.Data))
    fmt.Println(string(decodedBlock.Header.Store["dir"]) == string(b.Header.Store["dir"]))
}
