package rle

import (
    "bufio"
    "os"
)

func Decompress(in *bufio.Reader, out *bufio.Writer) {
    var (
        curr, prev byte = 0, 0
        found, invalid_prev bool = false, false
        error os.Error = nil
    )
    for error != os.EOF {
        curr, error = in.ReadByte()
        if error != nil {
            if found {
                panic("Archive corrupted")
            }
            break
        }
        if found {
            for ; curr > 0; curr-- {
                out.WriteByte(prev)
            }
            prev = 0
            found = false
            invalid_prev = true
        } else {
            out.WriteByte(curr)
            found = curr == prev && !invalid_prev
            prev = curr
            invalid_prev = false
        }
    }
}
