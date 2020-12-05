package machine

type Heap map[int]int

func (h Heap) Store(addr, value int) {
	h[addr] = value
}
func (h Heap) Load(addr int) (value int) {
	value = h[addr]
	return
}
