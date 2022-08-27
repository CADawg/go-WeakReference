package weakreference

import (
	"fmt"
	"testing"
	"time"
)

func BasicRecreator(str string) interface{} {
	return str
}
func TestAddItemPlain(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	weakRefs.Add("test", "works")

	if weakRefs.Read("test") == "works" {
		t.Logf("Success!")
	} else {
		t.Errorf("Adding item by value failed")
	}
}
func TestRecreator(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	if weakRefs.Read("test") != "test" {
		t.Errorf("Recreator failed")
	} else {
		t.Logf("Success!")
	}
}
func TestCacheEviction(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	weakRefs.Add("test", "works")

	time.Sleep(time.Second * 2)

	if weakRefs.PureRead("test") == nil {
		t.Logf("Success!")
	} else {
		t.Errorf("Cache was not evicted automatically.")
	}
}
func TestPureReadNil(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	if weakRefs.PureRead("test") == nil {
		t.Logf("Success!")
	} else {
		t.Errorf("PureRead Nil Test Failed")
	}
}
func TestInCacheEmpty(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	if !weakRefs.InCache("test") {
		t.Logf("Success!")
	} else {
		t.Errorf("Cache Empty Test Failed")
	}
}
func TestInCacheExists(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	// also tests that Read sets stuff correctly
	weakRefs.Read("test")

	if weakRefs.InCache("test") {
		t.Logf("Success!")
	} else {
		t.Errorf("Cache Empty Test Failed")
	}
}
func TestCacheWontClearIfAccessed(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second*7, BasicRecreator)

	// also tests that Read sets stuff correctly
	weakRefs.Add("test", 123)

	time.Sleep(4)

	// test read refresh cache time
	if weakRefs.Read("test") != 123 {
		t.Errorf("Wasn't available 4 seconds later (cache life: 7s)")
	}

	time.Sleep(4)

	// test pureread refresh cache time
	weakRefs.PureRead("test")

	time.Sleep(4)

	if weakRefs.PureRead("test") == 123 {
		t.Logf("Success!")
	} else {
		t.Errorf("Cache Was accessed recently but didn't keep data")
	}
}
func TestEnableDisableGC(t *testing.T) {
	var weakRefs = NewWeakReferences(time.Second, BasicRecreator)

	weakRefs.StopGC()

	// also tests that Read sets stuff correctly
	weakRefs.Add("test", "123")

	time.Sleep(time.Second * 2)

	// check gc was disabled
	if weakRefs.Read("test") != "123" {
		t.Errorf("Wasn't available 4 seconds later (cache life: 7s)")
	}

	weakRefs.StartGC()

	time.Sleep(time.Second * 2)

	fmt.Println(weakRefs.PureRead("test"))

	if weakRefs.PureRead("test") == nil {
		t.Logf("Success!")
	} else {
		t.Errorf("GC Re-enabled but didn't clear un-accessed file")
	}
}
