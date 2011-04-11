package huffman

import (
    "bufio"
    "gob"
    "os"
)

import . "common"

type (
    hfReverseCode struct {
        char byte
        len uint
    }
    hfReverseCodeTable map[uint] *hfReverseCode
)

var rct hfReverseCodeTable = make(hfReverseCodeTable)


func deserializeMetaInfo(fin *os.File) int64 {
    var (
        dump hfDump
    )
    PanicIf(gob.NewDecoder(fin).Decode(&dump))
    for i, record := range dump.Table {
        if record.Len != 0 {
            rct[record.Code] = &hfReverseCode{char: byte(i), len: record.Len}
        }
    }
    return dump.FileSize
}


func Decompress(fin, fout *os.File) int64 {
    var (
        code, code_len uint = 0, 0
        outptr *hfReverseCode = nil
        real_size, source_size, read_bytes int64 = 0, 0, 0
        i byte = 0
    )
    source_size = deserializeMetaInfo(fin)

    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()

    for real_size < source_size {
        curr, error := reader.ReadByte()
        if error != nil {
            break
        }
        read_bytes++
        i = 0
        for ; i < BITS_IN_BYTE; i++ {
            if (curr & (1 << i)) != 0 {
                code |= 1 << code_len
            }
            code_len++
            outptr = rct[code]
            if outptr != nil && outptr.len == code_len && real_size < source_size {
                writer.WriteByte(outptr.char)
                real_size++
                code, code_len = 0, 0
            }
        }
    }
    return read_bytes
}
