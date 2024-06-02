package chain

import (
	"blockchain/block"
	"testing"
    //"os"
    //"errors"

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

func TestSavingDataToChainThenReadingFromIt(t *testing.T) {
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

    err = chain.WriteBytes(data)
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
        t.Fatal("input and output are not equal")
    }
}

func TestSavingDataToChainThenReadingFromItInChunks(t *testing.T) {
    chainEncoder, err :=  encoder.NewRSAEncoder(3100)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }
    
    chain := New(chainEncoder, 256)
    if chain.Length != 0 {
        t.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
    } 

    data := make([]byte, 0)
    str := []string{"world"," hello"}

    for i := 0; i < 520; {
        if i % 2 == 0 {
            data = append(data, []byte(str[0])...)
            i += len(str[0])
        } else {
            data = append(data, []byte(str[1])...)
            i += len(str[1]) - 1
        }
    }    

    err = chain.WriteBytes(data)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if chain.Length != 3 {
        t.Fatalf("chain length should be 3 but was %d", chain.Length)
    }
    
    read, err := chain.ReadBytesInChunks()
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if string(data) != string(read) {
        t.Fatal("input and output are not equal")
    }
}

/*
func TestWritingChainToFiles(t *testing.T) {

    chainEncoder, err :=  encoder.NewRSAEncoder(5000)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }
    
    chain := New(chainEncoder, 32)
    if chain.Length != 0 {
        t.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
    } 

    data := make([]byte, 0)
    str := []string{"world"," hello"}

    for i := 0; i < 520; {
        if i % 2 == 0 {
            data = append(data, []byte(str[0])...)
            i += len(str[0])
        } else {
            data = append(data, []byte(str[1])...)
            i += len(str[1]) - 1
        }
    }    

    err = chain.WriteBytes(data)
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    if chain.Length < 9 {
        t.Fatalf("chain length should be more than 9 but was %d", chain.Length)
    }

    _, err = chain.SaveToFiles("./chainFiles")
    if err != nil {
        t.Fatalf("an error occured: %v", err)
    }

    files, _ := os.ReadDir("./chainFiles")
    if len(files) < 9 {
        t.Fatalf("wrong number of files. want %d got %d", 9, len(files))
    }

    filepath := "./" + string(chain.Hash) + ".chain"

    if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
        t.Fatalf("didn't save it to file :(")
    }
}
*/
