package huffman

type (
    hfCode struct {
        Code, Len uint
    }
    hfDump struct {
        Table    []hfCode
        FileSize int64
    }
)
