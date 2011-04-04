package common

import "os"


func PanicIf(error os.Error) {
    if error != nil {
        panic(error)
    }
}
