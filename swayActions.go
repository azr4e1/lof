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
	ClassCriteria                     = "class"
	ConIDCriteria                     = "con_id"
	ConMarkCriteria                   = "con_mark"
	FloatingCriteria                  = "floating"
	IdCriteria                        = "id"
	InstanceCriteria                  = "instance"
	PIDCriteria                       = "pid"
	ShellCriteria                     = "shell"
	TilingCriteria                    = "tiling"
	TitleCriteria                     = "title"
	UrgentCriteria                    = "urgent"
	WindowRoleCriteria                = "window_role"
	WindowTypeCriteria                = "window_type"
	WorkspaceCriteria                 = "workspace"
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

func OutputWindows(tree *Node, criteria FilterCriteria) FlattenedNodes {
	return tree.Flatten().Filter(criteria)
}

func Focus(wCriteria WindowCriteria, value string) error {
	swayMsg := exec.Command("swaymsg", fmt.Sprintf("[%s=%s]", string(wCriteria), value), "focus")

	_, err := swayMsg.CombinedOutput()

	return err
}

func Launch(cmd string) error {
	fields := strings.Fields(cmd)
	if len(fields) < 1 {
		return errors.New("You must provide a command")
	}

	path := fields[0]
	c := exec.Command(path, fields[1:]...)

	_, err := c.CombinedOutput()

	return err
}
