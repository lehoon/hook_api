package service

import (
	"fmt"
	"testing"
)

func TestSequence(t *testing.T) {
	for i := 0; i < 100; i++ {
		sequence_no, err := next_sequece()

		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

		fmt.Printf("%s\n", sequence_no)
	}
}
