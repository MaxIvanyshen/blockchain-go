package block

import (
	"errors"
	"fmt"
	"os"
	"testing"
    //"bytes"

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
    if b.Hash == "" {
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

func TestWritingBlockToFileAndReadingFromIT(t *testing.T) {
    input := "hello"
    rsa, err := encoder.NewRSAEncoder(2048)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }
    b := New(rsa, NewHeader())
    b.Data = []byte(input)

    parent := New(rsa, NewHeader())
    parent.Data = []byte(input)
    err = parent.Encode()

    b.Parent = parent.Hash
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    err = b.Encode()
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    path := "./"
    blockEncoder, err := encoder.NewRSAEncoder(5000)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    err = SaveToFile(blockEncoder, b, path)
    if err != nil {
        t.Fatalf("encountered an error: %v", err)
    }

    filepath := path + string(b.Hash)

    if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
        t.Fatalf("didn't save it to file :(")
    }

    readedBlock, err := ReadFromFile(filepath, blockEncoder)
    if err != nil {
        t.Fatalf("error occured: %v", err)
    }

    if readedBlock.Hash != b.Hash {
        t.Fatalf("Hash is not equal. want %s got %s", b.Hash, readedBlock.Hash)
    }

    blockData, err := DecodeBlockData(readedBlock, rsa)
    if err != nil {
        t.Fatalf("error occured: %v", err)
    }

    if string(blockData) != input {
        t.Fatalf("block data is not equal. want %s got %s", input, string(blockData))
    }

    if readedBlock.Parent != parent.Hash {
        t.Fatalf("Parent Hash is not equal. want %s got %s", input, string(blockData))
    }

    if err = os.Remove(filepath); err != nil {
        fmt.Println("couldn't delete test file")
    }
}
