package main

import (
    "flag"
    "fmt"
    "gob"
    "os"
)

import (
    . "common"
    "algorithms/rle"
    "algorithms/huffman"
)

type (
    cmap map[string] func(*os.File, *os.File)
    dcmap map[string] func(*os.File, *os.File) int64
)

var (
    fCreate = flag.Bool("c", false, "create archive")
    fExtract = flag.Bool("x", false, "extract files from archive")
    fMethod = flag.String("m", "rle", "compression method")
    fName = flag.String("f", "", "archive file name")

    COMPRESSORS cmap = make(cmap)
    DECOMPRESSORS dcmap = make(dcmap)
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


func init() {
      COMPRESSORS["rle.compress"] = rle.Compress
    DECOMPRESSORS["rle.decompress"] = rle.Decompress
      COMPRESSORS["huffman.compress"] = huffman.Compress
    DECOMPRESSORS["huffman.decompress"] = huffman.Decompress
}

func main() {
    /*defer func() {
        if error := recover(); error != nil {
            fmt.Printf("Error: %s", error)
        }
    }()*/

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

    compress := COMPRESSORS[*fMethod + ".compress"]
    decompress := DECOMPRESSORS[*fMethod + ".decompress"]
    if compress == nil || decompress == nil {
        panic("Unknown compression method")
    }

    if *fCreate {
        if flag.NArg() == 0 {
            panic("No arguments specified")
        }

        archive := openForWrite(*fName)
        defer archive.Close()

        PanicIf(gob.NewEncoder(archive).Encode(flag.Args()))
        for _, arg := range flag.Args() {
            fobj := openForRead(arg)
            compress(fobj, archive)
            fobj.Close()
        }
    } else {
        dir := *fName + ".ex"
        PanicIf(os.MkdirAll(dir, 0666))

        archive := openForRead(*fName)
        defer archive.Close()

        var filenames []string = make([]string, 0)
        PanicIf(gob.NewDecoder(archive).Decode(&filenames))
        for _, arg := range filenames {
            file_begin_pos, error := archive.Seek(0, 1)
            PanicIf(error)

            result := openForWrite(fmt.Sprintf("%s/%s", dir, arg))
            bytes_read := decompress(archive, result)
            result.Close()

            // workaround for buffered reader
            _, error = archive.Seek(file_begin_pos + bytes_read, 0)
            PanicIf(error)
        }
    }
}
