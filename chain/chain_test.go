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

func TestSavingDataToChain(t *testing.T) {
    chainEncoder, err :=  encoder.NewRSAEncoder(3076)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }
    
    chain := New(chainEncoder, 256)
    if chain.Length != 0 {
        t.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
    } 

    data := make([]byte, 520)
    for i := 0; i < 520; i++ {
        data[i] = byte(i * 3)
    }    

    blocks, err := chain.SaveToBytes(data)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if chain.Length != 3 {
        t.Fatalf("chain length should be 3 but was %d", chain.Length)
    }

    if string(blocks[2]) != string(chain.Tail.Hash) {
        t.Fatalf("wrong tail hash")
    }
}
