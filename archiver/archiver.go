package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
)

var (
    fCreate = flag.Bool("c", false, "create archive")
    fExtract = flag.Bool("x", false, "extract files from archive")
    fMethod = flag.String("m", "rle", "compression method")
    fName = flag.String("f", "", "archive file name")
)


func RLECompress(in *bufio.Reader, out *bufio.Writer) {
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


func RLEUncompress(in *bufio.Reader, out *bufio.Writer) {
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


func openForRead(name string) *os.File {
    file, error := os.Open(name, os.O_RDONLY, 0777)
    if error != nil {
        panic(error)
    }
    return file
}


func openForWrite(name string) *os.File {
    file, error := os.Open(name, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
    if error != nil {
        panic(error)
    }
    return file
}


func main() {
    defer func() {
        if error := recover(); error != nil {
            fmt.Printf("Error: %s", error)
        }
    }()

    flag.Parse()
    if !*fCreate && !*fExtract {
        panic("Missing operation")
    }
    if *fCreate && *fExtract {
        panic("Only one operation must be specified")
    }
    if *fName == "" {
        panic("Missing archive file name")
    }

    if *fCreate {
        if flag.NArg() == 0 {
            panic("No arguments specified")
        }

        archive := openForWrite(*fName)
        fout := bufio.NewWriter(archive)
        defer func() {
            fout.Flush()
            archive.Close()
        }()

        for _, arg := range flag.Args() {
            fobj := openForRead(arg)
            defer fobj.Close()
            fin := bufio.NewReader(fobj)
            RLECompress(fin, fout)
        }
    } else {
        archive := openForRead(*fName)
        fin := bufio.NewReader(archive)
        result := openForWrite(*fName + ".ex")
        fout := bufio.NewWriter(result)
        defer func() {
            fout.Flush()
            archive.Close()
            result.Close()
        }()
        RLEUncompress(fin, fout)
    }
}
