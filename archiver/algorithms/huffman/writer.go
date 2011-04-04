package huffman

import (
    "bufio"
    "container/heap"
    "gob"
    "os"
)

import . "common"

const BITS_IN_BYTE byte = 8

type (
    hfNode struct {
        char byte
        freq uint64
        left, right *hfNode
    }
    hfHeap []*hfNode
)

var (
    freq_table map[byte] uint64 = make(map[byte] uint64, 255)
    code_table hfCodeTable = make(hfCodeTable)
    tree hfHeap
)


func countFreq(fobj *os.File) {
    var (
        curr byte = 0
        error os.Error = nil
    )
    reader := bufio.NewReader(fobj)
    for {
        curr, error = reader.ReadByte()
        if error != nil {
            break
        }
        freq_table[curr]++
    }
    _, error = fobj.Seek(0, 0)
    PanicIf(error)
}


func fillCodeTable(node *hfNode, len, code uint) {
    if node.left == nil && node.right == nil {
        code_table[node.char] = hfCode{Len: len, Code: code}
    } else {
        fillCodeTable(node.left,  len + 1, code)
        fillCodeTable(node.right, len + 1, code | 1 << len)
    }
}


func Compress(fin *os.File, fout *os.File) {
    countFreq(fin)

    // create heap
    for ch, freq := range freq_table {
        if freq > 0 {
            tree.Push(&hfNode{char: ch, freq: freq})
        }
    }
    heap.Init(tree)
    for len(tree) > 1 {
        l := heap.Pop(tree).(*hfNode)
        r := heap.Pop(tree).(*hfNode)
        parent := &hfNode{freq: l.freq + r.freq, left: l, right: r}
        heap.Push(tree, parent)
    }

    // save code table
    fillCodeTable(tree[0], 0, 0)
    PanicIf(gob.NewEncoder(fout).Encode(code_table))

    // save source file size
    fi, error := fin.Stat()
    PanicIf(error)
    PanicIf(gob.NewEncoder(fout).Encode(fi.Size))

    // encode
    var (
        outbyte, outlen byte = 0, 0
        i uint = 0
    )
    reader := bufio.NewReader(fin)
    writer := bufio.NewWriter(fout)
    defer writer.Flush()
    for {
        curr, error := reader.ReadByte()
        if error != nil {
            break
        }
        entry := code_table[curr]
        for i < entry.Len {
            for ; i < entry.Len && outlen < BITS_IN_BYTE; i++ {
                if (entry.Code & (1 << i)) != 0 {
                    outbyte |= 1 << outlen
                }
                outlen++
            }
            if outlen == BITS_IN_BYTE {
                writer.WriteByte(outbyte)
                outbyte, outlen = 0, 0
            }
        }
        i = 0
    }
    if outbyte != 0 {
        writer.WriteByte(outbyte)
    }
}
