package tap

import "io"
import "os"
import "fmt"
import "testing/quick"

type T struct {
	nextTestNumber int
	outputFile     *os.File
}

// New creates a new Tap value
func New(fileName string) *T {
	output, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return &T{
		nextTestNumber: 1,
		outputFile:     output,
	}
}

// Header displays a TAP header including version number and expected
// number of tests to run.
func (t *T) Header(testCount int) {
	str := fmt.Sprintf("TAP version 13\n")
	io.WriteString(t.outputFile, str)
	str = fmt.Sprintf("1..%d\n", testCount)
	io.WriteString(t.outputFile, str)
}

// Ok generates TAP output indicating whether a test passed or failed.
func (t *T) Ok(test bool, description string) {
	// did the test pass or not?
	ok := "ok"
	if !test {
		ok = "not ok"
	}

	str := fmt.Sprintf("%s %d - %s\n", ok, t.nextTestNumber, description)
	io.WriteString(t.outputFile, str)
	t.nextTestNumber++
}

// Check runs randomized tests against a function just as "testing/quick.Check"
// does.  Success or failure generate appropriate TAP output.
func (t *T) Check(function interface{}, description string) {
	err := quick.Check(function, nil)
	if err == nil {
		t.Ok(true, description)
		return
	}

	fmt.Printf("# %s\n", err)
	t.Ok(false, description)
}

// return number of completed tests
func (t *T) Count() int {
	return t.nextTestNumber - 1
}

// generates plan string based on number of tests ran
func (t *T) AutoPlan() {
	str := fmt.Sprintf("1..%d\n", t.nextTestNumber-1)
	io.WriteString(t.outputFile, str)
}
