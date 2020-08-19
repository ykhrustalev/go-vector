package vector

import (
	"errors"
	"sync"
)

var (
	ErrInvalidIndex = errors.New("missing element with the given index")
	ErrSizeDiffers  = errors.New("non equal vectors")
)

const (
	defaultCapacity   = 10
	defaultMultiplier = 2
)

type Vector struct {
	multiplier int
	items      []int

	mx sync.RWMutex
}

func From(items ...int) *Vector {
	v := New()
	v.AppendAll(items...)
	return v
}

func New() *Vector {
	return NewWithCap(defaultCapacity, defaultMultiplier)
}

func NewWithCap(cap int, multiplier int) *Vector {
	return &Vector{
		multiplier: multiplier,
		items:      make([]int, 0, cap),
	}
}

func (v *Vector) withLock(cb func()) {
	v.mx.Lock()
	defer v.mx.Unlock()
	cb()
}

func (v *Vector) withRLock(cb func()) {
	v.mx.RLock()
	defer v.mx.RUnlock()
	cb()
}

func (v *Vector) increaseCapToAtLest(targetCap int) {
	if targetCap <= cap(v.items) {
		return
	}

	// the existing capacity is always > 0
	newCap := cap(v.items) * v.multiplier
	if newCap < targetCap {
		newCap = targetCap
	}

	v.increaseCapTo(newCap)
}

func (v *Vector) increaseCapTo(newCap int) {
	if newCap <= cap(v.items) {
		return
	}
	items := make([]int, 0, newCap)
	v.items = append(items, v.items...)
}

func (v *Vector) checkIndex(index int) error {
	if index < 0 || index > v.len()-1 {
		return ErrInvalidIndex
	}
	return nil
}

func (v *Vector) Add(index int, item int) (err error) {
	v.withLock(func() {
		if index == 0 && v.len() == 0 {
			// append like operation
			v.append(item)
			return
		}

		err = v.checkIndex(index)
		if err != nil {
			return
		}

		v.increaseCapToAtLest(v.len() + 1)

		v.items = append(v.items, 0) // increase
		copy(v.items[index+1:], v.items[index:])
		v.items[index] = item
	})

	return
}

func (v *Vector) Set(index int, item int) (err error) {
	v.withLock(func() {
		err = v.checkIndex(index)
		if err != nil {
			return
		}

		v.items[index] = item
	})

	return
}

func (v *Vector) append(items ...int) {
	v.increaseCapToAtLest(v.len() + len(items))
	v.items = append(v.items, items...)
}

func (v *Vector) Append(item int) {
	v.withLock(func() {
		v.append(item)
	})
}

func (v *Vector) AppendAll(items ...int) {
	v.withLock(func() {
		v.append(items...)
	})
}

func (v *Vector) remove(index int) int {
	r := v.items[index]
	v.items = append(v.items[:index], v.items[index+1:]...)
	return r
}

func (v *Vector) Remove(index int) (r int, err error) {
	v.withLock(func() {
		err = v.checkIndex(index)
		if err != nil {
			return
		}

		r = v.remove(index)
	})
	return
}

func (v *Vector) Peek(index int) (r int, err error) {
	v.withRLock(func() {
		err = v.checkIndex(index)
		if err != nil {
			return
		}

		r = v.items[index]
		return
	})
	return
}

func (v *Vector) IndexOf(item int) (r int) {
	r = -1
	v.withRLock(func() {
		for i, candidate := range v.items {
			if candidate == item {
				r = i
				break
			}
		}
	})
	return
}

func (v *Vector) len() int {
	return len(v.items)
}

func (v *Vector) Len() (r int) {
	v.withRLock(func() {
		r = v.len()
	})
	return
}

func (v *Vector) Cap() (r int) {
	v.withRLock(func() {
		r = cap(v.items)
	})
	return
}

func (v *Vector) each(cb func(index int, item int) bool) {
	for index, item := range v.items {
		if !cb(index, item) {
			return
		}
	}
}
func (v *Vector) Each(cb func(index, item int) bool) {
	v.withRLock(func() {
		v.each(cb)
	})
}

func (v *Vector) Clear() {
	v.withLock(func() {
		v.items = v.items[:0]
	})
}

func (v *Vector) Slice() (r []int) {
	v.withRLock(func() {
		for _, item := range v.items {
			r = append(r, item)
		}
	})
	return
}

func (v *Vector) Clone() (r *Vector) {
	v.withRLock(func() {
		r := NewWithCap(cap(v.items), v.multiplier)
		r.append(v.items...)
	})
	return
}

func (v *Vector) InnerProduct(other *Vector) (r int, err error) {
	v.withRLock(func() {
		if v.len() != other.len() {
			err = ErrSizeDiffers
			return
		}

		for i, val1 := range v.items {
			r += val1 * other.items[i]
		}
	})

	return
}

func (v *Vector) Any(cb func(item int) bool) (r bool) {
	v.withRLock(func() {
		v.each(func(index, item int) bool {
			if cb(item) {
				r = true
				return false
			}
			return true
		})
	})
	return
}

func (v *Vector) All(cb func(item int) bool) (r bool) {
	v.withRLock(func() {
		v.each(func(index, item int) bool {
			if index == 0 {
				r = true // initial
			}
			if !cb(item) {
				r = false
				return false
			}
			return true
		})
	})
	return
}

func (v *Vector) RemoveIf(cb func(item int) bool) {
	v.withLock(func() {
		for i := len(v.items) - 1; i >= 0; i-- {
			item := v.items[i]
			if !cb(item) {
				continue
			}

			v.remove(i)
		}
	})
}

func (v *Vector) Equal(other *Vector) (r bool) {
	v.withRLock(func() {
		if v.len() != other.len() {
			return
		}

		for i, this := range v.items {
			that := other.items[i]
			if that != this {
				return
			}
		}
		r = true
	})
	return
}

func (v *Vector) Accumulate(cb func(a, b int) int) *Vector {
	r := New()

	v.withRLock(func() {
		if v.len() == 0 {
			return
		}

		r.Append(v.items[0])
		for i := 1; i < v.len(); i++ {
			a, _ := r.Peek(r.Len() - 1)
			b := v.items[i]
			r.Append(cb(a, b))
		}
	})
	return r
}

func (v *Vector) Reduce(cb func(a, b int) int) (r int) {
	v.withRLock(func() {
		if v.len() == 0 {
			return
		}

		r = v.items[0]
		for i := 1; i < v.len(); i++ {
			a := v.items[i]
			r = cb(r, a)
		}
	})
	return
}
