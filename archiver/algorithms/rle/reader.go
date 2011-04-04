package rle

import (
    "bufio"
    "os"
)

func Decompress(fin *os.File, fout *os.File) {
    var (
        curr, prev byte = 0, 0
        found, invalid_prev bool = false, false
        error os.Error = nil
    )
    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()
    for {
        curr, error = reader.ReadByte()
        if error != nil {
            if found {
                panic("Archive corrupted")
            }
            break
        }
        if found {
            for ; curr > 0; curr-- {
                writer.WriteByte(prev)
            }
            prev = 0
            found = false
            invalid_prev = true
        } else {
            writer.WriteByte(curr)
            found = curr == prev && !invalid_prev
            prev = curr
            invalid_prev = false
        }
    }
}
