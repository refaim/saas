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


func deserializeMetaInfo(fin, fout *os.File) int64 {
    var (
        dump hfDump
    )
    PanicIf(gob.NewDecoder(fin).Decode(&dump))
    for i := range dump.Table {
        entry := &dump.Table[i]
        if entry.Len != 0 {
            codeptr := new(hfReverseCode)
            *codeptr = hfReverseCode{char: byte(i), len: entry.Len}
            rct[entry.Code] = codeptr
        }
    }
    return dump.FileSize
}


func Decompress(fin, fout *os.File) {
    var (
        code, code_len uint = 0, 0
        outptr *hfReverseCode = nil
        cursize, filesize int64 = 0, 0
    )
    filesize = deserializeMetaInfo(fin, fout)

    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()

    for cursize <= filesize {
        curr, error := reader.ReadByte()
        if error != nil {
            break
        }
        for i := 0; i < int(BITS_IN_BYTE); i++ {
            if (curr & (1 << uint(i))) != 0 {
                code |= 1 << code_len
            }
            code_len++
            outptr = rct[code]
            if outptr != nil && outptr.len == code_len && cursize < filesize {
                writer.WriteByte(outptr.char)
                cursize++
                code, code_len = 0, 0
            }
        }
    }
}
