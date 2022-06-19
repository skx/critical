//
// This is a trivial "TCL" which allows variable expansion,
// single-character variables only, and inline execution.
//
// We're missing a decent tokenizer, but we do have simple support
// for blocks - which is used for `if` and `while` commands.
//
// If we had a real parser we'd return things like this:
//
//   puts    -> TOKEN
//   "Hello" -> STRING
//   {}      -> BLOCK
//   []      -> EXPAND
//
// Bugs?  Many, but I think the biggest issue is the support
// for blocks, and nested-evaluation.
//
// For example this works:
//
//     set a [ expr 2 + 2 ]
//     // a -> 4
//
// But this doesn't work:
//
//    set a [ expr 2 + [ expr 1 + 1 ] ]
//    // a -> 2, rather than 4 as it should.
//
// We're not handling the nested addition so we get:
//    set a [ 2 + 0/error ]
//    // a -> 2
//

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// vars contains our map of global variables
var vars map[string]string

// set a variable
func set(name string, value string) string {
	if value != "" {
		vars[name] = value
		return value
	}
	return vars[name]

}

// incr increments the contents of the named variable by one.
//
// If the variable doesn't exist it is assumed to have a zero value,
// such that "incr unknown" returns 1.
func incr(name string, value string) string {

	offset, _ := strconv.Atoi(value)
	val, ok := vars[name]
	if ok {
		cur, _ := strconv.Atoi(val)
		cur += offset
		vars[name] = fmt.Sprintf("%d", cur)
	} else {
		vars[name] = "1"
	}
	return vars[name]
}

// decr decrements the contents of the named variable.
//
// If the variable does not exist it is assumed to contain a
// 0-value, so "decr unknown" results in -1.
func decr(name string, value string) string {
	offset, _ := strconv.Atoi(value)

	val, ok := vars[name]
	if ok {
		cur, _ := strconv.Atoi(val)
		cur -= offset
		vars[name] = fmt.Sprintf("%d", cur)
		return vars[name]
	} else {
		vars[name] = "-1"
		return "-1"
	}
}

// puts outputs a string, returning the output too.
func puts(arg string) string {
	out := fmt.Sprintf("%s", arg)
	fmt.Printf("%s\n", out)
	return out
}

// expr handles math and comparison operations.
//
// TODO - Currently we only support the use of integers, we should use
// floating-point.
func expr(a string, op string, b string) string {
	aV, _ := strconv.Atoi(a)
	bV, _ := strconv.Atoi(b)

	switch op {
	case "+":
		return (fmt.Sprintf("%f", float64(aV+bV)))
	case "-":
		return (fmt.Sprintf("%f", float64(aV-bV)))
	case "*":
		return (fmt.Sprintf("%f", float64(aV*bV)))
	case "/":
		return (fmt.Sprintf("%f", float64(aV/bV)))
	case "<":
		if aV < bV {
			return "1"
		} else {
			return "0"
		}
	case "<=":
		if aV <= bV {
			return "1"
		} else {
			return "0"
		}
	case ">":
		if aV > bV {
			return "1"
		} else {
			return "0"
		}
	case ">=":
		if aV >= bV {
			return "1"
		} else {
			return "0"
		}
	}
	return fmt.Sprintf("unknown operation %s %s %s", a, op, b)
}

// Expand variables inside strings, so "puts $a" becomes "puts XXX"
func expand(str string) string {

	ret := ""

	i := 0

	for i < len(str) {

		if str[i] == '$' {

			// Skip past the dollar
			i++

			// We build up the name of the variable to
			// expand
			variable := ""

			// While we've not walked off the end of our
			// string, and we've got a "letter" then we
			// can update our variable name.
			for i < len(str) && isLetter(str[i]) {
				variable += string(str[i])
				i++
			}

			// OK append the variable value to the string
			ret += vars[variable]
		} else {

			// Just append the string
			ret += str[i : i+1]
			i++
		}
	}

	return ret
}

// expand sub-commands, enclosed in square brackets.
//
// For example given:
//    puts [expr 3 + 4]
//
// replace the "[" "]" part with the result of running "3 + 4", so
// the line becomes:
//
//   puts 7
//
func expandReplace(str string) string {
	ret := ""
	i := 0

	for i < len(str) {
		c := str[i]

		// OK we have a nested thing.
		if c == '[' {
			tmp := ""
			closed := false
			i++
			for i < len(str) && !closed {
				c = str[i]
				i++
				if c == ']' {
					out, _ := Eval([]string{tmp})
					ret += out
					closed = true

				} else {
					tmp += string(c)
				}
			}

		} else {
			ret += string(c)
		}
		i++
	}

	return ret
}

// This is a hacky solution which allows us to handle "IF", and "WHILE".
//
// Given a line such as "if {aa} {xxxxx} else {yyyyy}" return
//
//   if
//   {aa}
//   {xxxx}
//   else
//   {yyyy}
//
// This handles whitespace inside the blocks, etc.
func lineToBlocks(input string) []string {
	ret := []string{}

	tmp := ""

	for _, x := range input {
		if x == '{' {
			ret = append(ret, tmp)
			tmp = "{"
		} else if x == '}' {
			tmp += "}"
			ret = append(ret, tmp)
			tmp = ""
		} else {
			tmp += string(x)
		}
	}

	// Ensure we don't have whitespace blcoks
	r := []string{}
	for _, ent := range ret {
		trimmed := strings.TrimSpace(ent)
		if trimmed != "" {
			r = append(r, trimmed)
		}

	}
	return r
}

// IsNumber returns true if the given string is wholly numeric.
//
// We try to parse as an integer, and return true if that succeeds, otherwise
// return false.
func IsNumber(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

// IsString returns true if the input is a quote-wrapped string.
//
// We allow both "string" and 'string'.
func IsString(v string) bool {
	if len(v) > 2 {
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			return true
		}
		if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			return true
		}

	}
	return false

}

// Eval interprets a set of code lines
func Eval(input []string) (string, error) {

	// Every evaluation returns a value.
	out := ""

	// For each line
	for _, line := range input {

		// Strip leading/trailing whitespace
		line = strings.TrimSpace(line)

		// Save a copy of the input line.
		//
		// This is a horrid hack for the `while` command,
		// as noted below.
		original := line

		// Expand variables in the line
		line = expand(line)

		// Now expand inline-commands
		line = expandReplace(line)

		// get the first token, which is a command.
		i := 0
		for i < len(line) {
			if isWhitespace(line[i]) {
				break
			}
			i++
		}

		// So we have the first part, and the fields
		// of arguments.
		first := line[0:i]
		fields := strings.Split(line, " ")

		switch first {
		case "set":
			if len(fields) == 3 {
				out = set(fields[1], fields[2])
			} else {
				out = set(fields[1], "")
			}
		case "incr":
			if len(fields) > 2 {
				out = incr(fields[1], fields[2])
			} else {
				out = incr(fields[1], "1")
			}

		case "decr":
			if len(fields) > 2 {
				out = decr(fields[1], fields[2])
			} else {
				out = decr(fields[1], "1")
			}

		case "puts":
			// A bit horrid
			tmp := line[i+1:]
			tmp = strings.TrimSpace(tmp)
			tmp = strings.TrimPrefix(tmp, "\"")
			tmp = strings.TrimSuffix(tmp, "\"")

			// output the string
			out = puts(tmp)

		case "expr":
			out = expr(fields[1], fields[2], fields[3])

		case "while":

			// This is all horrid as we reparse the original
			// line.
			//
			// We need to do this because:
			//
			//   while { expr $i < 10 } { puts ... }
			//
			// Will become:
			//
			//   while { expr 1 < 10 } { ...		}
			//
			// As the value of "$i" will be whatever it was set to
			// on the first time this line is encountered.
			//
			blocks := lineToBlocks(original)
			if len(blocks) != 3 {
				return "", fmt.Errorf("wrong number of blocks for 'while'")
			}

			test := blocks[1]
			body := blocks[2]

			test = strings.Trim(test, "{}")
			test = strings.TrimSpace(test)
			body = strings.Trim(body, "{}")
			body = strings.TrimSpace(body)

			var err error

			// Run the test the first time
			tmp, err := Eval([]string{test})
			if err != nil {
				return "", err
			}

			// While the result is not "false"
			for tmp != "" && tmp != "0" {

				// Reparse again.  FFS
				blocks = lineToBlocks(original)
				test = blocks[1]
				test = strings.Trim(test, "{}")
				test = strings.TrimSpace(test)

				body = blocks[2]
				body = strings.Trim(body, "{}")
				body = strings.TrimSpace(body)

				// Run the body
				b := strings.Split(body, "\n")
				out, err = Eval(b)
				if err != nil {
					return "", err
				}

				// Now repeat the test
				tmp, err = Eval([]string{test})
				if err != nil {
					return "", err
				}
			}

		case "if":
			// Parse into a set of blocks, horrid
			blocks := lineToBlocks(line)

			// blocks[0] == "if"
			expr := blocks[1]
			truthy := blocks[2]
			falsey := ""

			if len(blocks) >= 4 && blocks[3] == "else" {
				falsey = blocks[4]
			}

			// remove "{" "}" wrapper around the string
			expr = strings.Trim(expr, "{}")
			expr = strings.TrimSpace(expr)

			// remove "{" "}" wrapper around the string
			truthy = strings.Trim(truthy, "{}")
			truthy = strings.TrimSpace(truthy)

			// remove "{" "}" wrapper around the string
			falsey = strings.Trim(falsey, "{}")
			falsey = strings.TrimSpace(falsey)

			var err error

			out, err = Eval([]string{expr})
			if err != nil {
				fmt.Printf("error evaluating %s:%s", expr, err)
				return "", err
			}

			// true result is non-empty
			if out != "" && out != "0" {

				// Run the true-condition
				b := strings.Split(truthy, "\n")
				out, err = Eval(b)
				if err != nil {
					return "", err
				}
			} else {

				if falsey != "" {
					// Run the false-condition
					b := strings.Split(falsey, "\n")
					out, err = Eval(b)
					if err != nil {
						return "", err
					}
				}
			}

		case "":
			// Empty input - often the case for
			//   if { $UNSET } ...
			//
			// The line will expand to:
			//
			//   if { } ...
			//
			// Which will result in a call to evaluate the function ""
			//
			out = ""

		case "//":
			// Comment
			out = ""
		default:
			if IsNumber(fields[0]) {
				return fields[0], nil
			}
			if IsString(fields[0]) {
				return fields[0], nil
			}
			return "", fmt.Errorf("unknown function %s for line %s", fields[0], line)
		}
	}

	return out, nil
}

// readLines reads the contents of the given file into an array of lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read from the file.
	scanner := bufio.NewScanner(file)

	// Return value
	var lines []string

	//
	// We read into this string.
	//
	line := ""

	//
	// Loop
	//
	for scanner.Scan() {

		//
		// Get the line, and strip leading/trailing space.
		//
		tmp := scanner.Text()
		tmp = strings.TrimSpace(tmp)

		//
		// Append to our existing line.
		//
		line += tmp

		//
		// If the line ends with "\" then we remove
		// that character, and repeat with a newline.
		//
		if strings.HasSuffix(line, "\\") {
			line = strings.TrimSuffix(line, "\\")

			line += "\n"
			continue
		}

		lines = append(lines, line)
		line = ""
	}

	//
	// Was there an error with the scanner?  If so catch it
	// here.  To be honest I'm not sure if anything needs to
	// happen here
	//
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// No error
	return lines, nil
}

func main() {

	// Setup our arguments
	vars = make(map[string]string)

	if len(os.Args) > 1 {
		text, err := readLines(os.Args[1])
		if err != nil {
			fmt.Printf("Error reading %s: %s\n", os.Args[1], err)
			return
		}

		_, err = Eval(text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		return

	}
	a := []string{
		// variable
		"set a [ expr 2 + 2 ]",
		"puts \"A is set to: $a\"",

		// variable expansion comes first.
		"set a pu",
		"set b ts",
		"$a$b \"Hello World\"",

		// expansion
		"puts [set a 4]",
		"puts [set a]",

		// variable
		"set name Steve",
		"puts \"Hello World my name is $name\"",

		// maths
		"puts [expr 3 * 4]   ok world",
		"puts I'm dividing [expr 3 / 4]",

		// conditional
		"if { 1 } { puts \"OK: 1 was ok\" }",
		"if { 0 } { puts \"FAILURE: 0 was regarded as true\" }",
		"if { \"steve\" } { puts \"OK: steve was ok\" } else { puts \"steve was not ok\" }",
		"puts [expr 3 * 4]   ok world",
		"puts \"Still alive\"",
		"if { $a } { puts \"A is set\" } else { puts \"A is NOT set\" } ",
		"if { $x } { puts \"X is set\" } else { puts \"X is NOT set\" } ",
		//
		// Horrid
		//
		"set i 1",
		"set max 10",
		"while { expr $i <= $max } { puts \"Loop $i\"\n incr i }",
	}

	_, err := Eval(a)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

}

// isWhitespace returns true if the given character is whitespace.
func isWhitespace(ch byte) bool {
	return ch == byte(' ') || ch == byte('\t') || ch == byte('\n') || ch == byte('\r')
}

// isLetter returns true if the given character is suitable for use as a variable-name.
func isLetter(ch byte) bool {
	return ((ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z'))
}
