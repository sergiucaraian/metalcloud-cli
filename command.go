package main

import (
	"flag"
)

//Command defines a command, arguments, description etc
type Command struct {
	Description  string
	Subject      string
	AltSubject   string
	Predicate    string
	AltPredicate string
	FlagSet      *flag.FlagSet
	Arguments    map[string]interface{}
	InitFunc     func(c *Command)
	ExecuteFunc  func(c *Command, client MetalCloudClient) (string, error)
}

func sameCommand(a *Command, b *Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}

const _nilDefaultStr = "__NIL__"
const _nilDefaultInt = -14234
