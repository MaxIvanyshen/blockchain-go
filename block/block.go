package block

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

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
