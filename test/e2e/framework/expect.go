package framework

import (
	"errors"
	"fmt"
	"strings"

	ginkgotypes "github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// FailureError is an error where the error string is meant to be passed to
// ginkgo.Fail directly, i.e. adding some prefix like "unexpected error" is not
// necessary. It is also not necessary to dump the error struct.
type FailureError struct {
	msg            string
	fullStackTrace string
}

func (f FailureError) Error() string {
	return f.msg
}

func (f FailureError) Backtrace() string {
	return f.fullStackTrace
}

func (f FailureError) Is(target error) bool {
	return target == ErrFailure
}

func (f *FailureError) backtrace() {
	f.fullStackTrace = ginkgotypes.NewCodeLocationWithStackTrace(2).FullStackTrace
}

// ErrFailure is an empty error that can be wrapped to indicate that an error
// is a FailureError. It can also be used to test for a FailureError:.
//
//	return fmt.Errorf("some problem%w", ErrFailure)
//	...
//	err := someOperation()
//	if errors.Is(err, ErrFailure) {
//	    ...
//	}
var ErrFailure error = FailureError{}

// ExpectEqual expects the specified two are the same, otherwise an exception raises
func ExpectEqual(actual interface{}, extra interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.Equal(extra), explain...)
}

// ExpectNotEqual expects the specified two are not the same, otherwise an exception raises
func ExpectNotEqual(actual interface{}, extra interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).NotTo(gomega.Equal(extra), explain...)
}

// ExpectError expects an error happens, otherwise an exception raises
func ExpectError(err error, explain ...interface{}) {
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred(), explain...)
}

// ExpectNoError checks if "err" is set, and if so, fails assertion while logging the error.
func ExpectNoError(err error, explain ...interface{}) {
	ExpectNoErrorWithOffset(1, err, explain...)
}

// ExpectNil expects actual is nil
func ExpectNil(actual interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.BeNil(), explain...)
}

// ExpectNotNil expects actual is not nil
func ExpectNotNil(actual interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).NotTo(gomega.BeNil(), explain...)
}

// ExpectTrue expects actual is true
func ExpectTrue(actual interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.BeTrue(), explain...)
}

// ExpectFalse expects actual is false
func ExpectFalse(actual interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).NotTo(gomega.BeTrue(), explain...)
}

// ExpectNoErrorWithOffset checks if "err" is set, and if so, fails assertion while logging the error at "offset" levels above its caller
// (for example, for call chain f -> g -> ExpectNoErrorWithOffset(1, ...) error would be logged for "f").
func ExpectNoErrorWithOffset(offset int, err error, explain ...interface{}) {
	if err == nil {
		return
	}

	// Errors usually contain unexported fields. We have to use
	// a formatter here which can print those.
	prefix := ""
	if len(explain) > 0 {
		if str, ok := explain[0].(string); ok {
			prefix = fmt.Sprintf(str, explain[1:]...) + ": "
		} else {
			prefix = fmt.Sprintf("unexpected explain arguments, need format string: %v", explain)
		}
	}

	// This intentionally doesn't use gomega.Expect. Instead we take
	// full control over what information is presented where:
	// - The complete error object is logged because it may contain
	//   additional information that isn't included in its error
	//   string.
	// - It is not included in the failure message because
	//   it might make the failure message very large and/or
	//   cause error aggregation to work less well: two
	//   failures at the same code line might not be matched in
	//   https://go.k8s.io/triage because the error details are too
	//   different.
	//
	// Some errors include all relevant information in the Error
	// string. For those we can skip the redundant log message.
	// For our own failures we only log the additional stack backtrace
	// because it is not included in the failure message.
	var failure FailureError
	if errors.As(err, &failure) && failure.Backtrace() != "" {
		Logf("Failed inside E2E framework:\n    %s", strings.ReplaceAll(failure.Backtrace(), "\n", "\n    "))
	} else if !errors.Is(err, ErrFailure) {
		Logf("Unexpected error: %s\n%s", prefix, format.Object(err, 1))
	}
	Fail(prefix+err.Error(), 1+offset)
}

// ExpectConsistOf expects actual contains precisely the extra elements.  The ordering of the elements does not matter.
func ExpectConsistOf(actual interface{}, extra interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.ConsistOf(extra), explain...)
}

// ExpectHaveKey expects the actual map has the key in the keyset
func ExpectHaveKey(actual interface{}, key interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.HaveKey(key), explain...)
}

// ExpectEmpty expects actual is empty
func ExpectEmpty(actual interface{}, explain ...interface{}) {
	gomega.ExpectWithOffset(1, actual).To(gomega.BeEmpty(), explain...)
}
