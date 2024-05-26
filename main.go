package main

import (
    "testing"
    "fmt"
    "github.com/MaxIvanyshen/block-encryption/encoder"
    "blockchain/chain"
    "flag"
)

func main() {
    ReadWriteBenchmark()
}

func ReadWriteBenchmark() {
    benchWhole := func(b *testing.B) {
        chainEncoder, err :=  encoder.NewRSAEncoder(3100)
        if err != nil {
            b.Fatalf("an error occured: %v", err)
        }
        
        chain := chain.New(chainEncoder, 256)
        if chain.Length != 0 {
            b.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
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
            b.Fatalf("an error occured: %v", err)
        }

        for i := 0; i < b.N; i++ {
            _, err := chain.ReadBytesInChunks()
            if err != nil {
                b.Fatalf("an error occured: %v", err)
            }
        }
    }

    benchChunks := func(b *testing.B) {
        chainEncoder, err :=  encoder.NewRSAEncoder(3100)
        if err != nil {
            b.Fatalf("an error occured: %v", err)
        }
        
        chain := chain.New(chainEncoder, 256)
        if chain.Length != 0 {
            b.Fatalf("chain length should be 0 after initialization but was %d", chain.Length)
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
            b.Fatalf("an error occured: %v", err)
        }

        for i := 0; i < b.N; i++ {
            _, err := chain.ReadBytesInChunks()
            if err != nil {
                b.Fatalf("an error occured: %v", err)
            }
        }
    }

    testing.Init()
    flag.Parse()
    br := testing.Benchmark(benchWhole)
    fmt.Println("decoding whole byte array at once: " + br.String() + br.MemString())
    br = testing.Benchmark(benchChunks)
    fmt.Println("decoding byte array in chunks:  " + br.String() + br.MemString())
}
