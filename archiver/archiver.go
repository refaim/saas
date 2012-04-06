package main

import (
    "encoding/gob"
    "flag"
    "fmt"
    "os"
)

import (
    "algorithms/huffman"
    "algorithms/rle"
    . "common"
)

type (
    cmap  map[string]func(*os.File, *os.File)
    dcmap map[string]func(*os.File, *os.File) int64
)

var (
    fCreate  = flag.Bool("c", false, "create archive")
    fExtract = flag.Bool("x", false, "extract files from archive")
    fMethod  = flag.String("m", "rle", "compression method")
    fName    = flag.String("f", "", "archive file name")

    COMPRESSORS   cmap  = make(cmap)
    DECOMPRESSORS dcmap = make(dcmap)
)

func init() {
    COMPRESSORS["rle.compress"] = rle.Compress
    DECOMPRESSORS["rle.decompress"] = rle.Decompress
    COMPRESSORS["huffman.compress"] = huffman.Compress
    DECOMPRESSORS["huffman.decompress"] = huffman.Decompress
}

func main() {
    defer func() {
        /*if error := recover(); error != nil {
            fmt.Printf("Error: %s", error)
        }*/
    }()

    flag.Parse()
    if !(*fCreate || *fExtract) {
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

        archive, error := os.Create(*fName)
        PanicIf(error)
        defer archive.Close()

        compress := COMPRESSORS[*fMethod+".compress"]
        if compress == nil {
            panic("Unknown compression method")
        }
        PanicIf(gob.NewEncoder(archive).Encode(*fMethod))

        PanicIf(gob.NewEncoder(archive).Encode(flag.Args()))
        for _, arg := range flag.Args() {
            fobj, error := os.Open(arg)
            PanicIf(error)
            compress(fobj, archive)
            fobj.Close()
        }
    } else {
        dir := *fName + ".ex"
        PanicIf(os.MkdirAll(dir, 0666))

        archive, error := os.Open(*fName)
        PanicIf(error)
        defer archive.Close()

        var method string
        PanicIf(gob.NewDecoder(archive).Decode(&method))
        decompress := DECOMPRESSORS[method+".decompress"]
        if decompress == nil {
            panic("Unknown compression method")
        }

        var filenames []string = make([]string, 0)
        PanicIf(gob.NewDecoder(archive).Decode(&filenames))
        gob.NewDecoder(archive).Decode(&filenames)
        for _, arg := range filenames {
            file_begin_pos := GetFilePos(archive)

            result, error := os.Create(fmt.Sprintf("%s/%s", dir, arg))
            PanicIf(error)
            bytes_read := decompress(archive, result)
            result.Close()

            // workaround for buffered reader
            SafeSeek(archive, file_begin_pos+bytes_read, 0)
        }
    }
}
