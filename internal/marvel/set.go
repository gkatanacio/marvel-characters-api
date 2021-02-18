package marvel

type void struct{}

var member void

type IntSet struct {
	elems map[int]void
}

func NewIntSet() *IntSet {
	elems := make(map[int]void)
	return &IntSet{elems}
}

func (is *IntSet) Add(i int) {
	is.elems[i] = member
}

func (is *IntSet) Remove(i int) {
	delete(is.elems, i)
}

func (is *IntSet) Len() int {
	return len(is.elems)
}

func (is *IntSet) Contains(i int) bool {
	_, exists := is.elems[i]
	return exists
}

func (is *IntSet) ToSlice() []int {
	list := []int{}
	for i := range is.elems {
		list = append(list, i)
	}
	return list
}
