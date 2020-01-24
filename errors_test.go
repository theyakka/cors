package cors_test

import (
	"fmt"
	"github.com/theyakka/cors"
	"testing"
)

func TestError(t *testing.T) {
	code := 55
	message := "this is a validation error"
	err := cors.ValidationError{
		Code:          code,
		Message:       message,
		OriginalError: nil,
	}
	errMessage := err.Error()
	computedMessage := fmt.Sprintf("%s [%d]", message, code)
	if errMessage != computedMessage {
		t.Error("the expected error message was not outputted")
	}
}
