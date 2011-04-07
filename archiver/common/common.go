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
