package cmdline

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/azr4e1/lof"
)

type ErrorCode int

const (
	MissingArgumentError ErrorCode = iota + 1
	ExecutableError
	CriteriaError
	SwaymsgError
	EncodingError
	MarkError
)

const PrevMark = "__lof_prev"

func usage() {
}

type actionValue struct {
	validActions []string
	action       *string
}

// TODO: to test
func (av actionValue) String() string {
	if av.action == nil {
		return ""
	}
	return *av.action
}

// TODO: to test
func (av actionValue) Set(s string) error {
	if av.validActions == nil {
		return errors.New("Must specify valid actions")
	}

	if !slices.Contains(av.validActions, s) {
		return fmt.Errorf("Actions %s is not valid. Valid actions are: %s", s, strings.Join(av.validActions, ", "))
	}

	*(av.action) = s

	return nil
}

var validActions = []string{"launch", "focus", "launch_focus", "switch_prev", "get_windows"}

// TODO: to test
func markCurrentAsPrev() error {
	tree, err := lof.GetTree()
	if err != nil {
		return err
	}
	flatTree := lof.GetWindows(tree, func(bn *lof.BaseNode) bool { return bn.IsTrueWindow() })
	id, err := lof.GetIdFromFocused(flatTree)
	if err != nil {
		return err
	}
	lof.RemoveMark(PrevMark, lof.ConMarkCriteria, PrevMark)
	err = lof.AddMark(PrevMark, lof.ConIDCriteria, strconv.Itoa(id))
	return err
}

func getPrevId() (int, error) {
	tree, err := lof.GetTree()
	if err != nil {
		return 0, err
	}
	flatTree := lof.GetWindows(tree, func(bn *lof.BaseNode) bool { return bn.IsTrueWindow() })
	id, err := lof.GetIdFromMark(flatTree, PrevMark)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Main() ErrorCode {
	actionFlag := actionValue{validActions: validActions, action: new(string)}
	flag.Var(actionFlag, "action", "Action to perform")
	criteria := flag.String("criteria", string(lof.ConIDCriteria), "Criteria to filter windows.")
	identifier := flag.String("identifier", "", "Value of criteria")
	cmd := flag.String("cmd", "", "Command of window to launch")

	flag.Parse()

	switch *actionFlag.action {
	case "launch":
		if cmd == nil || *cmd == "" {
			fmt.Fprintln(os.Stderr, "You need to provide a command")
			return MissingArgumentError
		}

		markCurrentAsPrev()
		err := lof.Launch(*cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExecutableError
		}
	case "focus":
		if identifier == nil || *identifier == "" {
			fmt.Fprintln(os.Stderr, "You need to provide an identifier")
			return MissingArgumentError
		}
		markCurrentAsPrev()
		err := lof.Focus(lof.WindowCriteria(*criteria), *identifier)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return CriteriaError
		}
	case "launch_focus":
		if cmd == nil || *cmd == "" {
			fmt.Fprintln(os.Stderr, "You need to provide a command")
			return MissingArgumentError
		}
		if identifier == nil || *identifier == "" {
			fmt.Fprintln(os.Stderr, "You need to provide an identifier")
			return MissingArgumentError
		}
		markCurrentAsPrev()
		errFocus := lof.Focus(lof.WindowCriteria(*criteria), *identifier)
		if errFocus != nil {
			errLaunch := lof.Launch(*cmd)
			if errLaunch != nil {
				fmt.Fprintln(os.Stderr, errLaunch)
				return ExecutableError
			}
		}
	case "switch_prev":
		prevId, err := getPrevId()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return SwaymsgError
		}
		markCurrentAsPrev()
		err = lof.Focus(lof.ConIDCriteria, strconv.Itoa(prevId))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return CriteriaError
		}
	case "get_windows":
		tree, err := lof.GetTree()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return SwaymsgError
		}

		windows, err := lof.GetWindows(tree, func(bn *lof.BaseNode) bool { return bn.IsTrueWindow() }).ToJSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return EncodingError
		}

		fmt.Fprintln(os.Stdout, string(windows))
	default:

	}

	return 0
}
