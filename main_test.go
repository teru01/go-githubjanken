package main

import (
	"testing"
)

func TestFetchContributions(t *testing.T) {
	url := "https://github.com/users/teru01/contributions"
	contri, err := fetchContributions(url)
	if err != nil {
		t.Error(err)
	}
	if contri != 602 {
		t.Error("expected 602, got", contri)
	}
}

func TestFetchContributionsZero(t *testing.T) {
	url := "https://github.com/users/ksoo/contributions"
	contri, err := fetchContributions(url)
	if err != nil {
		t.Error(err)
	}
	if contri != 0 {
		t.Error("expected 0, got", contri)
	}
}

func TestFetchError(t *testing.T) {
	url := "https://github.com/_"
	_, err := fetchContributions(url)
	if err != nil {
		switch err := err.(type) {
		case *statusError:
			if err.code != 404 {
				t.Error("expected 404, got", err.code)
			}
			return
		default:
			t.Error(err)
		}
	}
	t.Error("err did not returned")
}
