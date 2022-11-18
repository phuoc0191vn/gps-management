package conversion

import (
	"fmt"
	"testing"
)

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
)

func TestSequenceToString(t *testing.T) {
	testCases := []struct {
		input  int64
		expect string
	}{
		{int64(1), "\x00\x00\x00\x00\x00\x00\x00\x01"},
		{int64(127), "\x00\x00\x00\x00\x00\x00\x00\x7f"},
		{int64(255), "\x00\x00\x00\x00\x00\x00\x00\xff"},
		{int64(1023), "\x00\x00\x00\x00\x00\x00\x03\xff"},
	}
	t.Log("Test convert sequence to string")
	for index, testCase := range testCases {
		seqString, _ := SequenceToString(testCase.input)
		if seqString == testCase.expect {
			t.Logf("Test case %d passed %s", index, checkMark)
		} else {
			t.Logf("Test case %d failed: Given %s (expected %s) %s", index, seqString, testCase.expect, ballotX)
		}
	}
}

func TestStringToSequence(t *testing.T) {
	testCases := []struct {
		input  string
		expect int64
	}{
		{"\x00\x00\x00\x00\x00\x00\x00\x01", int64(1)},
		{"\x00\x00\x00\x00\x00\x00\x00\x7f", int64(127)},
		{"\x00\x00\x00\x00\x00\x00\x00\xff", int64(255)},
		{"\x00\x00\x00\x00\x00\x00\x03\xff", int64(1023)},
	}
	t.Log("Test convert sequence to string")

	for index, testCase := range testCases {
		fmt.Printf("Test case %d: %x %U\n", index, testCase.input, testCase.input)
		seqString, _ := StringToSequence(testCase.input)
		if seqString == testCase.expect {
			t.Logf("Test case %d passed %s", index, checkMark)
		} else {
			t.Logf("Test case %d failed: Given %s (expected %s) %s", index, seqString, testCase.expect, ballotX)
		}
	}
}
