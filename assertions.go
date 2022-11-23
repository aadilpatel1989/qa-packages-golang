package asserts

import (
	"github.com/aadilpatel1989/qa-packages-golang"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Equal(t *testing.T, actual, expected interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.Equal(t, expected, actual, msgAndArgs...) {
		return
	}
	output := fmt.Sprintf("Not equal: \n"+
		"expected: %#v\n"+
		"actual  : %#v", expected, actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func Contains(t *testing.T, array, arrayObject interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.Contains(t, array, arrayObject, msgAndArgs...) {
		return
	}
	output := fmt.Sprintf("Not contains: \n"+
		"arrayObject: %#v\n"+
		"array  : %#v", arrayObject, array)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func NotEqual(t *testing.T, actual, expected interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.NotEqual(t, expected, actual, msgAndArgs...) {
		return
	}
	output := fmt.Sprintf("Equal: \n"+
		"expected: %#v\n"+
		"actual  : %#v", expected, actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func NotNil(t *testing.T, actual interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.NotNil(t, actual, msgAndArgs...) {
		return
	}
	output := fmt.Sprintf("Not nil: \n"+
		"actual  : %#v", actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func Nil(t *testing.T, actual interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.Nil(t, actual, msgAndArgs...) {
		return
	}
	output := fmt.Sprintf("Not nil: \n"+
		"actual  : %#v", actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func NotEmpty(t *testing.T, actual interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.NotEmpty(t, actual, msgAndArgs...) {
		return
	}

	output := fmt.Sprintf("Not empty: \n"+
		"actual  : %#v", actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func Empty(t *testing.T, actual interface{}, input string, err error, msgAndArgs ...interface{}) {
	if assert.Empty(t, actual, msgAndArgs...) {
		return
	}

	output := fmt.Sprintf("Empty: \n"+
		"actual  : %#v", actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}

func Error(t *testing.T, actual interface{}, input string, err error, msgAndArgs ...interface{}) {
	if !assert.Error(t, err) {
		return
	}

	output := fmt.Sprintf("Empty: \n"+
		"actual  : %#v", actual)
	reports.MarkTestStatus("Fail", t.Name(), input, output, err)
	t.FailNow()
}
