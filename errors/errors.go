package errors

import (
	"fmt"
	"strings"
)

type TestResults struct {
	msgs []string
}

func (e *TestResults) AppendFailure(msg string) {
	if msg != "" {
		e.msgs = append(e.msgs, msg)
	}
}

func (e *TestResults) Err() error {
	if len(e.msgs) == 0 {
		return nil
	}
	var sb strings.Builder
	_, _ = sb.WriteString("failed the following test cases:\n")
	for _, msg := range e.msgs {
		_, _ = sb.WriteString("\t")
		_, _ = sb.WriteString(msg)
		_, _ = sb.WriteString("\n")
	}
	return fmt.Errorf(sb.String())
}
