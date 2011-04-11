package common

import "os"


func PanicIf(error os.Error) {
    if error != nil {
        panic(error)
    }
}


func GetFileSize(fobj *os.File) int64 {
    fileinfo, error := fobj.Stat()
    PanicIf(error)
    return fileinfo.Size
}


func GetFilePos(fobj *os.File) int64 {
    pos, error := fobj.Seek(0, 1)
    PanicIf(error)
    return pos
}
