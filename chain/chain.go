package chain

import (
	"blockchain/block"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

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
        b.ParentHash = c.Tail.Hash
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

func (c *Chain) SaveBytes(data []byte) ([]hash, error) {
    blocksCount := int(math.Floor(float64(len(data) / int(c.BlockSize))))
    if len(data) != int(c.BlockSize) {
        blocksCount += 1
    }

    done := make(chan interface{})
    defer close(done)
    chunkStream := chunkStream(done, data, int(c.BlockSize))

    blocks := make([]hash, 0)

    for chunk := range chunkStream {
        header := block.NewHeader()
        header.Store["chain"] = []byte(c.Hash)
        b := block.New(c.encoder, header) 
        b.Data = chunk
        err := c.addBlock(b)
        if err != nil {
            return make([]hash, 0), fmt.Errorf("%v: could not add block to chain: %v", ChainSavingError, err)
        }
        blocks = append(blocks, hash(b.Hash))
    }

    if c.Hash == "" {
        hasher := sha256.New()
        hasher.Write([]byte(string(c.Tail.Data) + strconv.Itoa(int(c.Timestamp))))
        c.Hash = hash(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
    }

    return blocks, nil
}

func chunkStream(done <-chan interface{}, data []byte, chunkSize int) <-chan []byte {
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
