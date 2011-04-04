package main

import (
    "flag"
    "fmt"
    "os"
)

import (
    . "common"
    "algorithms/rle"
    "algorithms/huffman"
)

type mmap map[string] func(*os.File, *os.File)

var (
    fCreate = flag.Bool("c", false, "create archive")
    fExtract = flag.Bool("x", false, "extract files from archive")
    fMethod = flag.String("m", "rle", "compression method")
    fName = flag.String("f", "", "archive file name")

    COMPRESSION_METHODS mmap = make(mmap)
)


func openForRead(name string) *os.File {
    file, error := os.Open(name, os.O_RDONLY, 0444)
    PanicIf(error)
    return file
}


func openForWrite(name string) *os.File {
    file, error := os.Open(name, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
    PanicIf(error)
    return file
}


func main() {
    defer func() {
        if error := recover(); error != nil {
            fmt.Printf("Error: %s", error)
        }
    }()

    COMPRESSION_METHODS["rle.compress"] = rle.Compress
    COMPRESSION_METHODS["rle.decompress"] = rle.Decompress
    COMPRESSION_METHODS["huffman.compress"] = huffman.Compress
    COMPRESSION_METHODS["huffman.decompress"] = huffman.Decompress

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

    compress := COMPRESSION_METHODS[*fMethod + ".compress"]
    decompress := COMPRESSION_METHODS[*fMethod + ".decompress"]
    if compress == nil || decompress == nil {
        panic("Unknown compression method")
    }

    if *fCreate {
        if flag.NArg() == 0 {
            panic("No arguments specified")
        }

        archive := openForWrite(*fName)
        defer archive.Close()
        for _, arg := range flag.Args() {
            fobj := openForRead(arg)
            compress(fobj, archive)
            fobj.Close()
        }
    } else {
        archive := openForRead(*fName)
        result := openForWrite(*fName + ".ex")
        decompress(archive, result)
        archive.Close()
        result.Close()
    }
}
