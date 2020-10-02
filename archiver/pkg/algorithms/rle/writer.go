package rle

import (
    "bufio"
    "encoding/gob"
    "os"
)

import . "github.com/refaim/saas/archiver/pkg/common"

func Compress(fin, fout *os.File) {
    var (
        curr, prev, count byte  = 0, 0, 0
        found             bool  = false
        error             error = nil
    )

    PanicIf(gob.NewEncoder(fout).Encode(GetFileSize(fin)))

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
