package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
)

import (
    "algorithms/rle"
)

var (
    fCreate = flag.Bool("c", false, "create archive")
    fExtract = flag.Bool("x", false, "extract files from archive")
    fMethod = flag.String("m", "rle", "compression method")
    fName = flag.String("f", "", "archive file name")
)


func openForRead(name string) *os.File {
    file, error := os.Open(name, os.O_RDONLY, 0444)
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
            rle.Compress(fin, fout)
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
        rle.Decompress(fin, fout)
    }
}
