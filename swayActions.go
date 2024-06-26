package lof

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type WindowCriteria string

const (
	AppIDCriteria      WindowCriteria = "app_id"
	ClassCriteria      WindowCriteria = "class"
	ConIDCriteria      WindowCriteria = "con_id"
	ConMarkCriteria    WindowCriteria = "con_mark"
	FloatingCriteria   WindowCriteria = "floating"
	X11IDCriteria      WindowCriteria = "id"
	InstanceCriteria   WindowCriteria = "instance"
	PIDCriteria        WindowCriteria = "pid"
	ShellCriteria      WindowCriteria = "shell"
	TilingCriteria     WindowCriteria = "tiling"
	TitleCriteria      WindowCriteria = "title"
	UrgentCriteria     WindowCriteria = "urgent"
	WindowRoleCriteria WindowCriteria = "window_role"
	WindowTypeCriteria WindowCriteria = "window_type"
	WorkspaceCriteria  WindowCriteria = "workspace"
)

// TODO: to test
func GetTree() (*Node, error) {
	swayMsg := exec.Command("swaymsg", "--raw", "-t", "get_tree")
	data, err := swayMsg.CombinedOutput()
	if err != nil {
		return nil, err
	}

	node := new(Node)
	err = json.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// TODO: to test
func GetWindows(tree *Node, criteria FilterCriteria) FlattenedNodes {
	return tree.Flatten().Filter(criteria)
}

// TODO: to test
func Focus(wCriteria WindowCriteria, value string) error {
	swayMsg := exec.Command("swaymsg", fmt.Sprintf("[%s=%s] focus", string(wCriteria), value))

	err := swayMsg.Run()

	return err
}

// TODO: to test
func Launch(cmd string) error {
	fields := strings.Fields(cmd)
	if len(fields) < 1 {
		return errors.New("You must provide a command")
	}

	path := fields[0]
	c := exec.Command(path, fields[1:]...)

	err := c.Run()

	return err
}

// TODO: to test
func Close(wCriteria WindowCriteria, value string) error {
	c := exec.Command("swaymsg", fmt.Sprintf("[%s=%s] kill", string(wCriteria), value))

	err := c.Run()

	return err
}

// TODO: to test
func ToggleMark(mark string, wCriteria WindowCriteria, value string) error {
	swayMsg := exec.Command("swaymsg", fmt.Sprintf("[%s=%s] mark --toggle %s", string(wCriteria), value, mark))
	err := swayMsg.Run()

	return err
}

// TODO: to test
func AddMark(mark string, wCriteria WindowCriteria, value string) error {
	swayMsg := exec.Command("swaymsg", fmt.Sprintf("[%s=%s] mark --add %s", string(wCriteria), value, mark))
	err := swayMsg.Run()

	return err
}

func RemoveMark(mark string, wCriteria WindowCriteria, value string) error {
	swayMsg := exec.Command("swaymsg", fmt.Sprintf("[%s=%s] unmark %s", string(wCriteria), value, mark))
	err := swayMsg.Run()

	return err
}

// TODO: to test
func GetIdFromMark(fn FlattenedNodes, mark string) (int, error) {
	window := new(BaseNode)
outer:
	for _, bn := range fn {
		if bn == nil {
			continue
		}
		for _, curMark := range bn.Marks {
			if curMark == mark {
				window = bn
				break outer
			}
		}
	}

	if window == nil {
		return 0, fmt.Errorf("No window with mark %s", mark)
	}

	return window.Id, nil
}

// TODO: to test
func GetIdFromFocused(fn FlattenedNodes) (int, error) {
	window := new(BaseNode)
	for _, bn := range fn {
		if bn == nil {
			continue
		}
		if bn.Focused {
			window = bn
			break
		}
	}

	if window == nil {
		return 0, fmt.Errorf("No window in focus")
	}

	return window.Id, nil
}
