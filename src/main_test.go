package main

import (
	"log/slog"
	"testing"
)

func TestCheckMimetype(t *testing.T) {

	tests := map[string]string{
		"application/json; charset=utf-8": "json",
		"":                                "text",
		"application/pdf":                 "text",
		"application/json":                "json",
	}

	s := StoryConfig{}

	for contentType, expectedResult := range tests {
		res := s.getOutputFormat(contentType)
		if expectedResult != res {
			t.Logf("Expected: %v, got: %v", expectedResult, res)
			t.Fail()
		}
	}

}
