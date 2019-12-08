package set

type IntSet struct {
	values map[int]struct{}
}

func NewIntSet() IntSet {
	return IntSet{values: map[int]struct{}{}}
}

func (set *IntSet) Add(value int) {
	set.values[value] = struct{}{}
}

func (set IntSet) Contains(value int) bool {
	_, ok := set.values[value]
	return ok
}

func (set *IntSet) Remove(value int) {
	delete(set.values, value)
}

func (set IntSet) Values() []int {
	ret := make([]int, 0)
	for v := range set.values {
		ret = append(ret, v)
	}
	return ret
}
