package rle

import (
    "bufio"
    "gob"
    "os"
)

import . "common"


func Decompress(fin, fout *os.File) int64 {
    var (
        curr, prev byte = 0, 0
        found, invalid_prev bool = false, false
        error os.Error = nil
        real_size, source_size, read_bytes int64 = 0, 0, 0
    )

    PanicIf(gob.NewDecoder(fin).Decode(&source_size))

    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()

    for real_size < source_size {
        curr, error = reader.ReadByte()
        if error != nil {
            if found {
                panic("Archive corrupted")
            }
            panic(error)
        }
        read_bytes++
        if found {
            real_size += int64(curr)
            for ; curr > 0; curr-- {
                writer.WriteByte(prev)
            }
            prev = 0
            found = false
            invalid_prev = true
        } else {
            real_size++
            writer.WriteByte(curr)
            found = curr == prev && !invalid_prev
            prev = curr
            invalid_prev = false
        }
    }
    return read_bytes
}
