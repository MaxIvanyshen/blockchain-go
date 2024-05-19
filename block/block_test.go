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

    if string(b.Data) == "Hello world" {
        t.Fatalf("your encoding did not work")
    }
    if b.GetHash() == "" {
        t.Fatalf("hash missing")
    }
}

func TestBlockEncodingAndDecodingWithSameEncoder_DataShouldBeEqual(t *testing.T) {
    input := "Hello world"
    encoder, err := encoder.NewRSAEncoder(1024)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }
    b := New(encoder, NewHeader())
    b.Data = []byte(input)
    err = b.Encode()
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    decoded, err := DecodeBlockData(b, encoder)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    if string(decoded) != input {
        t.Fatalf("original and decrypted data are not equal. want '%s' got '%s'", input, string(decoded))
    }
}

func TestBlockEncodingAndDecodingWithDifferentEncoder_DataShouldNotBeEqual(t *testing.T) {
    input := "Hello world"
    rsa, err := encoder.NewRSAEncoder(1024)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }
    b := New(rsa, NewHeader())
    b.Data = []byte(input)
    err = b.Encode()
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    decoder, err := encoder.NewRSAEncoder(1024)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }
    decoded, err := DecodeBlockData(b, decoder)
    if err == nil {
        t.Fatal("data decoded successfully with different decoder")
    }

    if string(decoded) == input {
        t.Fatal("original and decrypted data are equal")
    }
}


