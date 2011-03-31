package rle

import (
    "bufio"
    "os"
)

func Compress(in *bufio.Reader, out *bufio.Writer) {
    var (
        curr, prev, count byte = 0, 0, 0
        found bool = false
        error os.Error = nil
    )
    for {
        curr, error = in.ReadByte()
        if error != nil {
            break
        }
        if found {
            if curr == prev && count < 255 {
                count++
            } else {
                out.WriteByte(count)
                out.WriteByte(curr)
                count = 0
                found = false
            }
        } else {
            out.WriteByte(curr)
            found = curr == prev
        }
        prev = curr
    }
    if count > 0 {
        out.WriteByte(count)
    }
}
