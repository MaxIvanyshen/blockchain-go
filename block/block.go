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
    hash hash
    parentHash hash
}

func New(encoder encoder.Encoder, header *Header) *Block {
    return &Block {
        encoder: encoder,
        Header: header,
    }
}

func (b *Block) Encode() error {
    hasher := sha256.New()
    hasher.Write([]byte(string(b.Data) + string(b.parentHash)))
    blockHash := hash(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))

    encoded, err := b.encoder.Encode(b.Data)
    if err != nil {
        return fmt.Errorf("could not encode block: %v", err)
    }

    b.Data = encoded
    b.hash = blockHash

    return nil
}

func DecodeBlockData(block *Block, decoder encoder.Encoder) ([]byte, error) {
    decoded, err := decoder.Decode(block.Data)
    if err != nil {
        return make([]byte, 0), fmt.Errorf("encountered an error while decoding block's data: %v", err)
    }
    return decoded, nil
}

func (b *Block) GetHash() hash {
    return b.hash
}

func (b *Block) GetParentHash() hash {
    return b.parentHash
}

var NotAllBytesWritten = errors.New("not all block bytes were written to file")

func SaveToFile(b *Block, dir string) error {
    if b.GetHash() == "" {
        fmt.Println("Block does not have hash. Encoding it first....")
        err := b.Encode()
        if err != nil {
            return err
        }
    }

    blockEncoder, err := encoder.NewRSAEncoder(4096) 
    if err != nil {
        return fmt.Errorf("encountered an error while saving block to file: %v", err)
    }
    buffer := bytes.Buffer{}
    enc := gob.NewEncoder(&buffer)
    err = enc.Encode(b)
    if err != nil {
        return fmt.Errorf("encountered an error while saving block to file: %v", err)
    }
    blockBytes, err := blockEncoder.Encode(buffer.Bytes())
    if err != nil {
        return fmt.Errorf("encountered an errorwhile saving block to file: %v", err)
    }

    if !strings.HasSuffix(dir, "/") {
        dir += "/"
    }
    filepath := dir + string(b.GetHash())
    file, err := os.Create(filepath)  
    n, err := file.Write(blockBytes)     
    if n < len(blockBytes) {
        return fmt.Errorf("encountered an error file while saving block to file: %v", NotAllBytesWritten)
    }

    return nil
}
