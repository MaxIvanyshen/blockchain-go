package block

import (
    "fmt"
    "errors"
)

type Header struct {
    store map[string][]byte
}

func NewHeader() *Header {
    return &Header{store: make(map[string][]byte)}
}

func (h *Header) Add(key string, val []byte) {
    h.store[key] = val
}

var NoHeaderParamError = errors.New("No header param")

func (h *Header) Get(key string) ([]byte, error) {
    val := h.store[key]
    if len(val) == 0 {
        return val, fmt.Errorf(
            "%v: header param with name '%s' does not exist",
            NoHeaderParamError,
            key,
        )
    }
    return val, nil
}
