package rle

import (
    "bufio"
    "os"
)

func Compress(fin *os.File, fout *os.File) {
    var (
        curr, prev, count byte = 0, 0, 0
        found bool = false
        error os.Error = nil
    )
    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()
    for {
        if curr, error = reader.ReadByte(); error != nil {
            break
        }
        if found {
            if curr == prev && count < 255 {
                count++
            } else {
                writer.WriteByte(count)
                writer.WriteByte(curr)
                count = 0
                found = false
            }
        } else {
            writer.WriteByte(curr)
            found = curr == prev
        }
        prev = curr
    }
    if count > 0 {
        writer.WriteByte(count)
    }
}
