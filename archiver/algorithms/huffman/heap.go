package huffman


func (heap *hfHeap) Len() int {
    return len(*heap)
}


func (heap *hfHeap) Less(i, j int) bool {
    return (*heap)[i].freq < (*heap)[j].freq
}


func (heap *hfHeap) Swap(i, j int) {
    (*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}


func (heap *hfHeap) Push(x interface{}) {
    (*heap) = append((*heap), x.(*hfNode))
}


func (heap *hfHeap) Pop() interface{} {
    n := heap.Len() - 1
    elm := (*heap)[n]
    (*heap)[n] = nil
    (*heap) = (*heap)[0:n]
    return elm
}
