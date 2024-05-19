package block

type Header struct {
    Store map[string][]byte
}

func NewHeader() *Header {
    return &Header{Store: make(map[string][]byte)}
}
