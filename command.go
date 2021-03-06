package main

import (
	"flag"
	"fmt"
	"strconv"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

//CommandExecuteFunc a function type a command can take for executing the content
type CommandExecuteFunc = func(c *Command, client interfaces.MetalCloudClient) (string, error)

//CommandInitFunc a function type a command can take for initializing the command
type CommandInitFunc = func(c *Command)

//Command defines a command, arguments, description etc
type Command struct {
	Description  string
	Subject      string
	AltSubject   string
	Predicate    string
	AltPredicate string
	FlagSet      *flag.FlagSet
	Arguments    map[string]interface{}
	InitFunc     CommandInitFunc
	ExecuteFunc  CommandExecuteFunc
	Endpoint     string
}

func sameCommand(a *Command, b *Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}

const _nilDefaultStr = "__NIL__"
const _nilDefaultInt = -14234

//confirms command
func confirmCommand(c *Command, f func() string) (bool, error) {

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		return true, nil
	}

	return requestConfirmation(f())
}

//getPtrValueIfExistsOk returns a string or an int from a map of pointers if the key exists
func getPtrValueIfExistsOk(m map[string]interface{}, key string) (interface{}, bool) {

	if v := m[key]; v != nil {
		switch v.(type) {
		case *int:
			if *v.(*int) != _nilDefaultInt {
				return *v.(*int), true
			}
		case *string:
			if *v.(*string) != _nilDefaultStr {
				return *v.(*string), true
			}
		}
	}
	return nil, false
}

//getIDFromStringOk returns the id and true if valid number
func getIDFromStringOk(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	return i, err == nil
}

//verifyParam returns error if param is not present
func getParam(c *Command, label string, name string) (interface{}, error) {
	v := c.Arguments[label]
	if v == nil {
		return nil, fmt.Errorf("-%s cannot be nil", name)
	}
	switch v.(type) {
	case *int:
		if *v.(*int) <= 0 {
			return nil, fmt.Errorf("-%s cannot be <=0", name)
		}
		if *v.(*int) == _nilDefaultInt {
			return nil, fmt.Errorf("-%s is required", name)
		}
	case *string:
		if *v.(*string) == "" {
			return nil, fmt.Errorf("-%s cannot be empty", name)
		}
		if *v.(*string) == _nilDefaultStr {
			return nil, fmt.Errorf("-%s is required", name)
		}
	}
	return v, nil
}

func idOrLabelString(v string) (int, string, bool) {
	if i, ok := getIDFromStringOk(v); ok {
		return i, "", true
	}
	return 0, v, false
}

//idOrLabel returns an int or a string contained in the interface. The last param is true if int is returned.
func idOrLabel(v interface{}) (int, string, bool) {
	switch v.(type) {
	case *int:
		return *v.(*int), "", true
	case *string:
		if i, ok := getIDFromStringOk(*v.(*string)); ok {
			return i, "", true
		}
		return 0, *v.(*string), false
	}
	return -1, "", false
}

type getIDOrDoFunc func(i string) (int, error)

func getIDOrDo(idOrLabel string, f getIDOrDoFunc) (int, error) {
	id, label, isID := idOrLabelString(idOrLabel)
	if !isID {
		return f(label)
	}
	return id, nil
}

func getIntParam(v interface{}) int {
	if v != nil && *v.(*int) != _nilDefaultInt {
		return *v.(*int)
	}
	return 0
}

func getStringParam(v interface{}) string {
	if v != nil && *v.(*string) != _nilDefaultStr {
		return *v.(*string)
	}
	return ""
}

func getBoolParam(v interface{}) bool {
	return v != nil && *v.(*bool)
}

func getStringParamOk(v interface{}) (string, bool) {
	if v != nil && *v.(*string) != _nilDefaultStr {
		return *v.(*string), true
	}
	return "", false
}

func getIntParamOk(v interface{}) (int, bool) {
	if v != nil && *v.(*int) != _nilDefaultInt {
		return *v.(*int), true
	}
	return 0, false
}

func getBoolParamOk(v interface{}) (bool, bool) {
	if v == nil {
		return false, false
	}
	return v != nil && *v.(*bool), true
}

func updateIfIntParamSet(v interface{}, p *int) {
	if v, ok := getIntParamOk(v); ok {
		*p = v
	}
}

func updateIfStringParamSet(v interface{}, p *string) {
	if v, ok := getStringParamOk(v); ok {
		*p = v
	}
}

func updateIfBoolParamSet(v interface{}, p *bool) {
	if v, ok := getBoolParamOk(v); ok {
		*p = v
	}
}
