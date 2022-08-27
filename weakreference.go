package weakreference

import (
	"sync"
	"time"
)

type WeakReference struct {
	Reference      interface{}
	RecreationData string
	LastAccessed   time.Time
}
type weakReferences struct {
	References             map[string]*WeakReference
	Locker                 *sync.RWMutex
	Recreator              func(string) interface{}
	EvictAfter             time.Duration
	GarbageCollectorActive bool
	ExtraDebug             bool
}

func NewWeakReferences(evictAfter time.Duration, recreator func(string) interface{}) *weakReferences {
	wR := &weakReferences{
		Locker:                 &sync.RWMutex{},
		References:             make(map[string]*WeakReference),
		Recreator:              recreator,
		EvictAfter:             evictAfter,
		GarbageCollectorActive: true,
		ExtraDebug:             false,
	}

	go func() {
		for {
			if wR.GarbageCollectorActive {
				wR.GarbageCollect()
			}
			time.Sleep(time.Second)
		}
	}()

	return wR
}
func (wS *weakReferences) StopGC() {
	wS.GarbageCollectorActive = false
}

func (wS *weakReferences) StartGC() {
	wS.GarbageCollectorActive = true
}

func (wS *weakReferences) AddWeakRef(str string, wR *WeakReference) {
	wS.Locker.Lock()
	defer wS.Locker.Unlock()

	wS.References[str] = wR
}

func (wS *weakReferences) Add(str string, ref interface{}) {
	wS.AddWeakRef(str, &WeakReference{
		Reference:      ref,
		RecreationData: str,
		LastAccessed:   time.Now(),
	})
}

func (wS *weakReferences) Read(str string) interface{} {
	wS.Locker.Lock()
	w, ok := wS.References[str]
	defer wS.Locker.Unlock()

	if !ok {
		w = &WeakReference{}
		w.RecreationData = str
		w.Reference = wS.Recreator(str)
		wS.References[str] = w
	}

	if w.Reference == nil {
		w.Reference = wS.Recreator(str)
	}

	w.LastAccessed = time.Now()

	return w.Reference
}

// PureRead Gets the underlying value, without running Recreator
func (wS *weakReferences) PureRead(str string) interface{} {
	wS.Locker.Lock()
	w, _ := wS.References[str]
	defer wS.Locker.Unlock()

	if w == nil {
		return nil
	}

	w.LastAccessed = time.Now()

	return w.Reference
}

func (wS *weakReferences) InCache(str string) bool {
	return wS.PureRead(str) != nil
}

func (wS *weakReferences) GarbageCollect() {
	wS.Locker.Lock()
	defer wS.Locker.Unlock()

	cTime := time.Now()
	for i, wR := range wS.References {
		if wR.LastAccessed.Before(cTime.Add(-wS.EvictAfter)) && wR.Reference != nil {
			wRCopy := wR
			wRCopy.Reference = nil
			wS.References[i] = wRCopy
		}
	}
}
