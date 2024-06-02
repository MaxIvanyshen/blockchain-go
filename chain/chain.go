package chain

import (
	"blockchain/block"
	"crypto/sha256"
	"encoding/base64"
    /*
	"bytes"
	"encoding/gob"
	"os"
	"strings"
    */
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

type hash string

type Chain struct {
    encoder encoder.Encoder
    Length uint
    blockSize uint
    Tail string
    Hash string
    Timestamp int64
    Blocks []string
}

func New(encoder encoder.Encoder, blockSize uint) *Chain {
    timestamp := time.Now().Unix()
    return &Chain{
        encoder: encoder,
        blockSize: blockSize,
        Timestamp: timestamp,
        Blocks: make([]string, 0),
    }
}

var UnableToAddBlockToChainError = errors.New("unable to add block to chain")

func (c *Chain) addBlock(b *block.Block) error {
    if c.Length != 0 {
        b.Parent = c.Tail
    }
    c.checkHash()
    path := "./" + c.Hash //TODO: make path get variable from env
    err := block.SaveToFile(c.encoder, b, path) 
    if err != nil {
        return fmt.Errorf("%v: %v", UnableToAddBlockToChainError, err)
    }

    c.Tail = b.Hash
    c.Blocks = append(c.Blocks, b.Hash)
    c.Length += 1

    return nil
}

func (c *Chain) checkHash() {
    if c.Hash == "" {
        hasher := sha256.New()
        hasher.Write([]byte(strconv.Itoa(int(c.blockSize)) + strconv.Itoa(int(c.Timestamp))))
        c.Hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    }
}

var ChainSavingError = errors.New("an error occured while saving chain")

func (c *Chain) WriteBytes(data []byte) error {
    done := make(chan interface{})
    defer close(done)
    chunkStream := writeChunkStream(done, data, int(c.blockSize))

    c.checkHash()

    for chunk := range chunkStream {
        header := block.NewHeader()
        header.Store["chain"] = []byte(c.Hash)
        b := block.New(c.encoder, header) 
        b.Data = chunk
        err := c.addBlock(b)
        if err != nil {
            return fmt.Errorf("%v: could not add block to chain: %v", ChainSavingError, err)
        }
    }

    return nil
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

var ZeroLengthError = errors.New("the chain's length is 0")


//this function reads data in chain and decodes its chunks in one thread
func (c *Chain) ReadBytes() ([]byte, error) {
    if c.Length == 0 {
        return make([]byte, 0), fmt.Errorf("could now decode chain data: %v", ZeroLengthError)
    }
    out := make([]byte, 0)

    chunks := make([][]byte, 0)
    current := getBlock(
    for current != nil {
        chunk, err := block.DecodeBlockData(current, c.encoder)
        if err != nil {
            return out, fmt.Errorf("could not decode chain data because of error while reading block: %v", err)
        }
        chunks = append(chunks, chunk)
        current = c.Blocks[i]
    }

    for i := len(chunks) - 1; i >= 0; i-- {
        out = append(out, chunks[i]...)
    }

    return out, nil
}

/*
    This structure is used to decode chain data by chunks,
    specifically to concatenate chunks in the right order
    and not think of passing around data byte and chunks index
*/
type chunk struct {
    bytes []byte
    idx int
}

/*
    This functions does the same as the previous one, but it decodes chunks 
    concurrently, which at times appears to be faster. This also makes us
    later use O(n^2) loop to concatenate all the decoded chunks
*/
func (c *Chain) ReadBytesInChunks() ([]byte, error) {
    if c.Length == 0 {
        return make([]byte, 0), fmt.Errorf("could now decode chain data: %v", ZeroLengthError)
    }
    out := make([]byte, 0)

    chunks := make([]chunk, 0)
    current := c.Tail    

    var wg sync.WaitGroup

    for i := 0; current != nil; i++ {

        wg.Add(1)
        go func(current *block.Block) {
                defer wg.Done()
                decoded, err := block.DecodeBlockData(current, c.encoder)
                if err != nil {
                    panic(err)
                }

                chunks = append(chunks, chunk{bytes: decoded, idx: i})
        }(current)

        current = current.Parent
    }

    wg.Wait()

    for i := len(chunks) - 1; i >= 0; i-- {
        out = append(out, findChunk(chunks, i)...)
    }

    return out, nil
}

func findChunk(chunks []chunk, neededIdx int) []byte {
    for i := 0; i < len(chunks); i++ {
        if chunks[i].idx == neededIdx {
            return chunks[i].bytes
        }
    }
    return make([]byte, 0)
}

var SavingChainError = errors.New("could not save chain to files")
var NotAllBytesWritten = errors.New("not all chain bytes were written to file")

/*
func (c *Chain) SaveToFiles(path string) ([]hash, error) {
    dto := NewChainDTO(c)

    blocks := make([]hash,  c.Length, c.Length)

    current := c.Tail
    dto.Tail = hash(c.Tail.Hash)
    
    if c.Hash == "" {
        hasher := sha256.New()
        hasher.Write([]byte(strconv.Itoa(int(c.blockSize)) + strconv.Itoa(int(c.Timestamp))))
        c.Hash = hash(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
    }

    for i := c.Length - 1; current != nil; i-- {
        err := block.SaveToFile(c.encoder, current, path)
        if err != nil {
            return make([]hash, 0), fmt.Errorf(
                "%v: error while saving block with hash '%s': %v",
                SavingChainError,
                current.Hash,
                err,
            )
        }
        blocks[i] = hash(current.Hash)
        current = current.Parent
    }

    buffer := bytes.Buffer{}
    enc := gob.NewEncoder(&buffer)
    err := enc.Encode(c))
    if err != nil {
        return make([]hash, 0), fmt.Errorf("encountered an error while saving block to file: %v", err)
    }
    
    chainBytes, err := c.encoder.Encode(buffer.Bytes())
    if err != nil {
        return make([]hash, 0), fmt.Errorf("encountered an errorwhile saving block to file: %v", err)
    }

    if !strings.HasSuffix(path, "/") {
        path += "/"
    }
    filepath := path + string(c.Hash) + ".chain"
    file, err := os.Create(filepath)  
    n, err := file.Write(chainBytes)     
    if n < len(chainBytes) {
        return make([]hash, 0), fmt.Errorf("encountered an error file while saving block to file: %v", block.NotAllBytesWritten)
    }

    file.Close()

    return blocks, nil
}
*/

func (c *Chain) SaveToFiles(path string) error {

    return nil
}
