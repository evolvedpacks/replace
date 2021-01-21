package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

const version = "1.0.0"

type Args struct {
	Input        string   `arg:"positional" help:"Input data (taken from STDIN when not provided)"`
	Mappings     []string `arg:"-m,--map,separate" help:"Values to be replaced"`
	Replacements []string `arg:"-t,--to,separate" help:"Values to replace with"`
	MapFile      string   `arg:"-f,--mapfile" help:"JSON file to read replacement mappings from"`
}

func (Args) Description() string {
	return "Replace strings in input streams.\n\n" +
		"Either you can provide mappings via command line arguments,\n" +
		"for example like following:\n" +
		"\n" +
		"  cat data.txt | replace --map \"replace this\" --to \"With this\"\n" +
		"\n" +
		"You can chain as many '--map' and '--to' bindings as you want\n" +
		"as long as the same ammount of mappings as of replacements is\n" +
		"provided. The first mappings is replaced with the frist\n" +
		"replacement and so on.\n" +
		"You can also provide a JSON file as mapping which looks like\n" +
		"following, for example:\n" +
		"\n" +
		"  {\n" +
		"    \"Replace this\": \"With this\"\n" +
		"  }\n" +
		"\n"
}

func (Args) Version() string {
	return "replace v" + version
}

type Mapping map[string]string

func (m Mapping) Apply(str string) string {
	for k, v := range m {
		str = strings.ReplaceAll(str, k, v)
	}
	return str
}

func (m Mapping) Merge(m2 Mapping) {
	for k, v := range m2 {
		m[k] = v
	}
}

func main() {
	args := new(Args)
	arg.MustParse(args)

	input := args.Input
	if input == "" {
		input = readInputFromStdin()
	}

	mapping, err := mappingFromArgs(args.Mappings, args.Replacements)
	errIfErrorExit(err)

	if args.MapFile != "" {
		fMapping, err := mappingFromJsonFile(args.MapFile)
		errIfErrorExit(err)
		mapping.Merge(fMapping)
	}

	input = mapping.Apply(input)

	os.Stdout.WriteString(input)
}

func readInputFromStdin() (v string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		v += "\n" + s.Text()
	}
	v = v[1:]
	return
}

func mappingFromArgs(mappings, replacements []string) (Mapping, error) {
	if len(mappings) != len(replacements) {
		return nil, errors.New("Missmatched ammount of arguments for mappings and replacements")
	}

	m := make(Mapping)
	for i, v := range mappings {
		m[v] = replacements[i]
	}

	return m, nil
}

func mappingFromJsonFile(filePath string) (Mapping, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	m := make(Mapping)
	err = json.NewDecoder(f).Decode(&m)

	return m, err
}

func errorExit(msg string) {
	os.Stderr.WriteString("Error: " + msg + "\n")
	os.Exit(1)
}

func fErrorExit(msg string, args ...interface{}) {
	errorExit(fmt.Sprintf(msg, args...))
}

func errIfErrorExit(err error) {
	if err != nil {
		errorExit(err.Error())
	}
}
