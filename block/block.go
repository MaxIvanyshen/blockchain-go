package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

type hash string

type Block struct {
    encoder encoder.Encoder
    Header *Header
    Data []byte
    Hash hash
    ParentHash hash
}

func New(encoder encoder.Encoder, header *Header) *Block {
    return &Block {
        encoder: encoder,
        Header: header,
    }
}

func (b *Block) Encode() error {
    hasher := sha256.New()
    hasher.Write([]byte(string(b.Data) + string(b.ParentHash)))
    blockHash := hash(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))

    encoded, err := b.encoder.Encode(b.Data)
    if err != nil {
        return fmt.Errorf("could not encode block: %v", err)
    }

    b.Data = encoded
    b.Hash= blockHash

    return nil
}

func DecodeBlockData(block *Block, decoder encoder.Encoder) ([]byte, error) {
    decoded, err := decoder.Decode(block.Data)
    if err != nil {
        return make([]byte, 0), fmt.Errorf("encountered an error while decoding block's data: %v", err)
    }
    return decoded, nil
}

var NotAllBytesWritten = errors.New("not all block bytes were written to file")

func SaveToFile(encoder encoder.Encoder, b *Block, dir string) error {
    if b.Hash == "" {
        fmt.Println("Block does not have hash. Encoding it first....")
        err := b.Encode()
        if err != nil {
            return err
        }
    }

    buffer := bytes.Buffer{}
    enc := gob.NewEncoder(&buffer)
    err := enc.Encode(b)
    if err != nil {
        return fmt.Errorf("encountered an error while saving block to file: %v", err)
    }
    blockBytes, err := encoder.Encode(buffer.Bytes())
    if err != nil {
        return fmt.Errorf("encountered an errorwhile saving block to file: %v", err)
    }

    if !strings.HasSuffix(dir, "/") {
        dir += "/"
    }
    filepath := dir + string(b.Hash)
    file, err := os.Create(filepath)  
    n, err := file.Write(blockBytes)     
    if n < len(blockBytes) {
        return fmt.Errorf("encountered an error file while saving block to file: %v", NotAllBytesWritten)
    }

    file.Close()

    return nil
}

var NoBlockFileError = errors.New("block file with this name does not exists")

func ReadFromFile(filepath string, decoder encoder.Encoder) (*Block, error) {
    if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
        return &Block{}, fmt.Errorf("couldn't read block from file: %v", NoBlockFileError)
    }

    encodedBlockBytes, err := os.ReadFile(filepath)
    if err != nil {
        return &Block{}, fmt.Errorf("couldn't read block from file: %v", NoBlockFileError)
    }

    blockBytes, err := decoder.Decode(encodedBlockBytes) 
    if err != nil {
        return &Block{}, fmt.Errorf("couldn't decode encoded block: %v", err)
    }

    buffer := bytes.Buffer{}
    buffer.Write(blockBytes)
    
    var block Block
    err = gob.NewDecoder(&buffer).Decode(&block)
    if err != nil {
        return &Block{}, fmt.Errorf("couldn't convert bytes to block: %v", err)
    }

    return &block, nil
}
