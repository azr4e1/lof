package lof_test

import (
	"testing"

	"github.com/azr4e1/lof"
	_ "github.com/google/go-cmp/cmp"
)

func TestIsContainer_DetectsAContainerCorrectly(t *testing.T) {
	t.Parallel()

	containers := []lof.ContainerType{lof.SimpleContainer, lof.FloatingContainer}
	nonContainers := []lof.ContainerType{lof.OutputContainer, lof.RootContainer, lof.WorkspaceContainer, lof.ScratchpadContainer}

	for _, con := range containers {
		container := new(lof.BaseNode)

		container.Type = con

		want := true
		got := container.IsContainer()

		if got != want {
			t.Errorf("want %v, got %v in testcase %v", want, got, con)
		}
	}

	for _, con := range nonContainers {
		container := new(lof.BaseNode)

		container.Type = con

		want := false
		got := container.IsContainer()

		if got != want {
			t.Errorf("want %v, got %v in testcase %v", want, got, con)
		}
	}

	var container *lof.BaseNode
	want := false
	got := container.IsContainer()

	if got != want {
		t.Errorf("want %v, got %v in testcase %v", want, got, "nil")
	}
}

func TestIsTrueWindow_DetectsAnActualWindowCorrectly(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		Id     int
		Test   *lof.BaseNode
		Answer bool
	}

	testCases := []TestCase{
		{
			Id:     1,
			Test:   &lof.BaseNode{Type: lof.SimpleContainer, Name: "con"},
			Answer: true,
		},
		{
			Id:     2,
			Test:   &lof.BaseNode{Type: lof.ScratchpadContainer, Name: "con"},
			Answer: false,
		},
		{
			Id:     3,
			Test:   &lof.BaseNode{Type: lof.FloatingContainer},
			Answer: false,
		},
		{
			Id:     4,
			Test:   &lof.BaseNode{Type: lof.FloatingContainer, Name: "con"},
			Answer: true,
		},
		{
			Id:     5,
			Test:   &lof.BaseNode{Type: lof.SimpleContainer},
			Answer: false,
		},
		{
			Id:     6,
			Test:   &lof.BaseNode{Type: lof.WorkspaceContainer},
			Answer: false,
		},
		{
			Id:     7,
			Test:   &lof.BaseNode{Type: lof.RootContainer, Name: "con"},
			Answer: false,
		},
		{
			Id:     8,
			Test:   nil,
			Answer: false,
		},
	}

	for _, tc := range testCases {
		want := tc.Answer
		got := tc.Test.IsTrueWindow()
		if got != want {
			t.Errorf("want %v, got %v for test case %d", want, got, tc.Id)
		}
	}
}
