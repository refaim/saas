package huffman


func (hfHeap) Len() int {
    return len(tree)
}


func (hfHeap) Less(i, j int) bool {
    return tree[i].freq < tree[j].freq
}


func (hfHeap) Swap(i, j int) {
    tree[i], tree[j] = tree[j], tree[i]
}


func (hfHeap) Push(x interface{}) {
    elm := x.(*hfNode)
    tree = append(tree, elm)
}


func (hfHeap) Pop() interface{} {
    n := len(tree) - 1
    elm := tree[n]
    tree[n] = nil
    tree = tree[0:n]
    return elm
}
