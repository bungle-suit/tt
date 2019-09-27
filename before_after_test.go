package tt_test

import (
	"io"
	"testing"

	"github.com/bungle-suit/tt"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	log string
}

func (t *testLogger) Log() string {
	return t.log
}

func (t *testLogger) Append(s string) {
	t.log += s + "\n"
}

func (t *testLogger) Action(msg string) func() {
	return func() {
		t.Append(msg)
	}
}

func (tl *testLogger) Func(msg string) tt.TestFunction {
	return func(t *testing.T) {
		tl.Append(msg)
	}
}

func (tl *testLogger) Panic(msg string) tt.TestFunction {
	return func(t *testing.T) {
		tl.Append(msg)
		panic("foo")
	}
}

type CloseFunc func() error

func (f CloseFunc) Close() error {
	return f()
}

func (tl *testLogger) Closer(before, onAfter string) func() io.Closer {
	return func() io.Closer {
		tl.Append(before)
		return CloseFunc(func() error {
			tl.Append(onAfter)
			return nil
		})
	}
}

func hideFooPanic(f tt.TestFunction) tt.TestFunction {
	return func(t *testing.T) {
		defer func() {
			e := recover()
			if s, ok := e.(string); ok && s == "foo" {
				return
			}

			panic(e)
		}()
		f(t)
	}
}

func TestBefore(t *testing.T) {
	logger := testLogger{}
	t.Run("before", tt.Before(logger.Action("1"), logger.Func("act")))
	assert.Equal(t, "1\nact\n", logger.Log())
}

func TestBeforeFP(t *testing.T) {
	logger := testLogger{}
	fp := tt.BeforeFP(logger.Action("1"))
	t.Run("before", fp(logger.Func("act")))
	assert.Equal(t, "1\nact\n", logger.Log())
}

func TestAfter(t *testing.T) {
	logger := testLogger{}
	t.Run("normal", tt.After(logger.Action("1"), logger.Func("act")))
	assert.Equal(t, "act\n1\n", logger.Log())

	logger = testLogger{}
	t.Run("panic", hideFooPanic(tt.After(logger.Action("1"), logger.Panic("act"))))
	assert.Equal(t, "act\n1\n", logger.Log())
}

func TestAfterFP(t *testing.T) {
	logger := testLogger{}
	fp := tt.AfterFP(logger.Action("1"))
	t.Run("after", fp(logger.Func("act")))
	assert.Equal(t, "act\n1\n", logger.Log())
}

func TestBeforeAfter(t *testing.T) {
	logger := testLogger{}
	t.Run("normal", tt.BeforeAfter(logger.Action("1"), logger.Action("2"), logger.Func("act")))
	assert.Equal(t, "1\nact\n2\n", logger.Log())

	logger = testLogger{}
	t.Run("panic", hideFooPanic(tt.BeforeAfter(logger.Action("1"), logger.Action("2"), logger.Panic("act"))))
	assert.Equal(t, "1\nact\n2\n", logger.Log())
}

func TestBeforeAfterFP(t *testing.T) {
	logger := testLogger{}
	fp := tt.BeforeAfterFP(logger.Action("1"), logger.Action("2"))
	t.Run("FP", fp(logger.Func("act")))
	assert.Equal(t, "1\nact\n2\n", logger.Log())
}

func TestCloser(t *testing.T) {
	logger := testLogger{}
	t.Run("normal", tt.Closer(logger.Closer("1", "2"), logger.Func("act")))
	assert.Equal(t, "1\nact\n2\n", logger.Log())

	logger = testLogger{}
	t.Run("panic", hideFooPanic(tt.Closer(logger.Closer("1", "2"), logger.Panic("act"))))
	assert.Equal(t, "1\nact\n2\n", logger.Log())
}

func TestCloserFP(t *testing.T) {
	logger := testLogger{}
	fp := tt.CloserFP(logger.Closer("1", "2"))
	t.Run("normal", fp(logger.Func("act")))
	assert.Equal(t, "1\nact\n2\n", logger.Log())
}
