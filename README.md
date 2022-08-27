# Go WeakReference
Loosely based on [C# Style WeakReference](https://docs.microsoft.com/en-us/dotnet/api/system.weakreference?view=net-6.0), an automatic cache eviction and re-creation system using last access times to evict files.

## Usage:

```go
package main

import (
	"weakreference" // or "github.com/Snaddyvitch-Dispenser/go-WeakReference"
	"time"
)

// Recreate the WeakReference if it has been garbage collected.
// This should recreate the object from the key and return it.
func Recreate(key string) interface{} {
	return key
}

func main() {
	// Create a new WeakReference Holder With a Cache Eviction Time of 1 Second (Too low for production) 
	// and a Recreation function that will recreate the object from the key.
	var weakRefs = weakreference.NewWeakReferences(time.Second, Recreate)

	// Stop the garbage collector
	// Can be used to stop the garbage collector, so you can use your own garbage collector.
	weakRefs.StopGC()
	
	// Restart the garbage collector (it runs by default)
	weakRefs.StartGC()

	// Add an item with a value
	weakRefs.Add("test", "123")
	
	// Read an item from the cache / Recreate() if it doesn't exist or has been evicted.
	var _ = weakRefs.Read("test") // value is "123"
	_ = weakRefs.Read("something") // value is "something" as Recreate("something") returns "something"

	// PureRead lets you get the underlying value bypassing Recreate()
	_ = weakRefs.PureRead("else") // value is nil (as doesn't exist)
	_ = weakRefs.PureRead("test") // value is "123" (as it is in the cache still)
	
	time.Sleep(time.Second * 2)
    // As timeout is one second, the cache will be empty now
	
	_ = weakRefs.PureRead("test") // returns nil (as it has been evicted)
	_ = weakRefs.Read("test") // returns Recreate("test") -> "test"
}
```