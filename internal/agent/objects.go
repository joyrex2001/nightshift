package agent

import (
	"container/heap"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

// objectspq is the priority queue that contains objects found by the scanners.
// Each scanner has a priority and Objects found with these scanners take the
// same priority. The highest priority takes precedence over earlier scanned
// objects with a lower priority (hence this implementation uses a priority
// queue). If an Object is added with the same priority it will be replaced.
type objectspq []*scanner.Object

// Len returns the length of the priority queue, as required by the heap
// interface.
func (pq objectspq) Len() int {
	return len(pq)
}

// Less compares the scanner.Objects, and determines the order of the priority
// queue, as required by the heap interface.
func (pq objectspq) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

// Swap will swap two scanner.Objects on the priority queue, as required by
// the heap interface.
func (pq objectspq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Pop will remove the scanner.Object with the lowest Priority from the
// priority queue.
func (pq *objectspq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Push will add an Object to the priority queue.
func (pq *objectspq) Push(x interface{}) {
	item := x.(*scanner.Object)
	*pq = append(*pq, item)
}

// Index will return the raw array index of the given scanner.Object.
func (pq objectspq) Index(obj *scanner.Object) int {
	// this could be optimized with a hashmap indexing the actual positions,
	// however, since the queues are expected to have just a few entries, and
	// this code will only be called when an update is received, the
	// implementation is left as-is (choosing readiblity over performance).
	for i, o := range pq {
		if o.Priority == obj.Priority {
			return i
		}
	}
	return -1
}

// InitObjects will initialize the objects. If objects were stored, this method
// will remove these and the objects are re-initialized.
func (a *worker) InitObjects() {
	a.m.Lock()
	defer a.m.Unlock()
	a.objects = map[string]*objectspq{}
}

// GetObjects will go through all object priority queues, and for each object
// found, it will append the result with the highest priority to the Objects
// result map.
func (a *worker) GetObjects() map[string]*scanner.Object {
	a.m.Lock()
	defer a.m.Unlock()
	objs := map[string]*scanner.Object{}
	for _, opq := range a.objects {
		if len(*opq) > 0 {
			obj := (*opq)[0].Copy()
			objs[obj.UID] = obj
		}
	}
	return objs
}

// addObject will add (or replace!) an object to the collection of objects.
// Each object is stored in its own priority queue, and if an object with the
// same priority is to be added, it will replace the object instead.
func (a *worker) addObject(obj *scanner.Object) {
	a.m.Lock()
	defer a.m.Unlock()
	opq, ok := a.objects[obj.UID]
	if !ok {
		// no entries yet, init the heap!
		opq := &objectspq{obj}
		a.objects[obj.UID] = opq
		heap.Init(opq)
		return
	}
	if idx := opq.Index(obj); idx >= 0 {
		// existing entry found for this priority; replace the item in the
		// priority queue
		(*opq)[idx] = obj
		return
	}
	// add the object to the priority queue
	heap.Push(opq, obj)
}

// removeObject will remove an Object from the priority queue.
func (a *worker) removeObject(obj *scanner.Object) {
	a.m.Lock()
	defer a.m.Unlock()
	opq, ok := a.objects[obj.UID]
	if !ok {
		return
	}
	if idx := opq.Index(obj); idx >= 0 {
		heap.Remove(opq, idx)
	}
}
