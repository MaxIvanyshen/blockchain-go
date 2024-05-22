package chain

import (
	"blockchain/block"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"
	"math"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

type hash string

type Chain struct {
    encoder encoder.Encoder
    Length uint
    BlockSize uint
    Tail *block.Block
    Hash hash
    Timestamp int64
}

func New(encoder encoder.Encoder, blockSize uint) *Chain {
    timestamp := time.Now().Unix()
    return &Chain{encoder: encoder, BlockSize: blockSize, Timestamp: timestamp}
}

var UnableToAddBlockToChainError = errors.New("unable to add block to chain")

func (c *Chain) addBlock(b *block.Block) error {
    if c.Length != 0 {
        b.Parent = c.Tail
    }
    err := b.Encode()
    if err != nil {
        return fmt.Errorf("%v: %v", UnableToAddBlockToChainError, err)
    }

    c.Tail = b
    c.Length += 1

    return nil
}

var ChainSavingError = errors.New("an error occured while saving chain")

func (c *Chain) SaveBytes(data []byte) error {
    blocksCount := int(math.Floor(float64(len(data) / int(c.BlockSize))))
    if len(data) != int(c.BlockSize) {
        blocksCount += 1
    }

    done := make(chan interface{})
    defer close(done)
    chunkStream := writeChunkStream(done, data, int(c.BlockSize))

    blocks := make([]hash, 0)

    for chunk := range chunkStream {
        header := block.NewHeader()
        header.Store["chain"] = []byte(c.Hash)
        b := block.New(c.encoder, header) 
        b.Data = chunk
        err := c.addBlock(b)
        if err != nil {
            return fmt.Errorf("%v: could not add block to chain: %v", ChainSavingError, err)
        }
        blocks = append(blocks, hash(b.Hash))
    }

    if c.Hash == "" {
        hasher := sha256.New()
        hasher.Write([]byte(string(c.Tail.Data) + strconv.Itoa(int(c.Timestamp))))
        c.Hash = hash(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
    }

    return nil
}

var ZeroLengthError = errors.New("the chain's length is 0")

func (c *Chain) ReadBytes() ([]byte, error) {
    if c.Length == 0 {
        return make([]byte, 0), fmt.Errorf("could now decode chain data: %v", ZeroLengthError)
    }
    out := make([]byte, 0)

    chunks := make([][]byte, 0)
    current := c.Tail    
    for i := 0; current != nil; i++ {
        chunk, err := block.DecodeBlockData(current, c.encoder)
        if err != nil {
            return out, fmt.Errorf("could not decode chain data because of error while reading block: %v", err)
        }
        chunks = append(chunks, chunk)
        current = current.Parent
    }

    for i := len(chunks) - 1; i >= 0; i-- {
        out = append(out, chunks[i]...)
    }

    return out, nil
}

func (c *Chain) ReadTailBytes() ([]byte, error) {
    if c.Length == 0 {
        return make([]byte, 0), fmt.Errorf("could now decode chain data: %v", ZeroLengthError)
    }
    out, err := block.DecodeBlockData(c.Tail, c.encoder)
    if err != nil {
        return make([]byte, 0), fmt.Errorf("could now read tail block data: %v", err)
    }
    return out, nil
}

func writeChunkStream(done <-chan interface{}, data []byte, chunkSize int) <-chan []byte {
    stream := make(chan []byte)

    go func() {
        defer close(stream)
        for i := 0; i < len(data); i += chunkSize {
            var chunk []byte

            if i + chunkSize >= len(data) {
                chunkSize = len(data) - i
            }

            chunk = data[i:i+chunkSize]

            select {
            case <-done:
                return
            case stream<-chunk:
            }
        }
    }()

    return stream
}
