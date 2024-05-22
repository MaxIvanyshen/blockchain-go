package chain

import (
	"blockchain/block"
	"testing"

	"github.com/MaxIvanyshen/block-encryption/encoder"
)

func TestAddingBlockToChain(t *testing.T) {
    encoder, err :=  encoder.NewRSAEncoder(2048)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }
    
    chain := New(encoder, 256)
    if chain.Length != 0 {
        t.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
    } 

    b := block.New(encoder, block.NewHeader())
    err = chain.addBlock(b) 
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if chain.Length != 1 {
        t.Fatalf("chain length should be 1 after adding one block but was %d", chain.Length)
    } 
}

func TestSavingDataToChainThenReadingFromItsTail(t *testing.T) {
    chainEncoder, err :=  encoder.NewRSAEncoder(3100)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }
    
    chain := New(chainEncoder, 256)
    if chain.Length != 0 {
        t.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
    } 

    data := make([]byte, 0)
    str := []string{"hello"," world"}

    for i := 0; i < 520; {
        if i % 2 == 0 {
            data = append(data, []byte(str[0])...)
            i += len(str[0])
        } else {
            data = append(data, []byte(str[1])...)
            i += len(str[1]) - 1
        }
    }    

    err = chain.SaveBytes(data)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if chain.Length != 3 {
        t.Fatalf("chain length should be 3 but was %d", chain.Length)
    }
    
    read, err := chain.ReadBytes()
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if string(data) != string(read) {
        t.Fatal("input and output are now equal")
    }
}