package userinput

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var GET string = "get"
var SET string = "set"
var TESTSET string = "testset"

var legalWord string = "([a-zA-Z0-9_]+)"
var legalGet string = fmt.Sprintf("(%s)\\(%s\\)", GET, legalWord)
var legalSet string = fmt.Sprintf("(%s)\\(%s,%s\\)", SET, legalWord, legalWord)
var legalTestSet string = fmt.Sprintf("(%s)\\(%s,%s,%s\\)", TESTSET, legalWord, legalWord, legalWord)

var legalCommands string = fmt.Sprintf("^(%s|%s|%s)$", legalGet, legalSet, legalTestSet)

type LegalCommand struct {
	Command string
	Args    []string
}

func IsLegalCommand(text string) bool {
	r, _ := regexp.Compile(legalCommands)
	return r.MatchString(text)
}

func ParseCommand(text string) (LegalCommand, error) {
	var parsedCmd LegalCommand
	var err error
	if IsLegalCommand(text) {
		parsedCmd = extractCommand(text)
	} else {
		err = errors.New(fmt.Sprintf("Invalid command: %s", text))
	}
	return parsedCmd, err
}

func extractCommand(text string) LegalCommand {
	splitCmd := strings.Split(text, "(")
	cmdName := splitCmd[0]
	r, _ := regexp.Compile(legalWord)
	cmdArgs := r.FindAllString(splitCmd[1], -1)
	return LegalCommand{cmdName, cmdArgs}
}
