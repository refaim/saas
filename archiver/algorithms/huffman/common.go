package huffman

type (
    hfCode struct {
        Code, Len uint
    }
    hfCodeTable map[byte] hfCode
)
