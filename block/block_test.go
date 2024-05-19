package block

import (
	"testing"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

func TestBlockEncoding(t *testing.T) {
    encoder, err := encoder.NewRSAEncoder(1024)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }
    b := New(encoder, NewHeader())
    b.Data = []byte("Hello world")
    err = b.Encode()
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    if string(b.Data) == "hello world" {
        t.Fatalf("your encoding did not work")
    }
    if b.GetHash() == "" {
        t.Fatalf("hash missing")
    }
}
