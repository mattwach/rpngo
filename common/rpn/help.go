package rpn

import (
	"fmt"
	"sort"
)

func (rpn *RPN) initHelp() {
	conceptHelp := map[string]string{
		"base": "Different number bases are supported:\n" +
			"Examples:\n" +
			"  1.1 -1 0 # complex floats\n" +
			"  2+3i # also a complex float\n" +
			"  32d # a decimal integer\n" +
			"  ffx # hexidecimal ff (255 in decimal)\n" +
			"  55o # octal 55 (45 in decimal)\n" +
			"  1001b # binary (9 in decimal)\n" +
			"Numbers can be converted to different formats:\n" +
			"  bin # convert to binary\n" +
			"  float # convert to complex float\n" +
			"  hex # convert to hecidecimal\n" +
			"  int # cnvert to integer\n" +
			"  oct # cnvert to octal\n" +
			"  str # cnvert to string",

		"basics": "- Enter numbers to push them to the stack\n" +
			"- Numbers can be separated by spaces or newlines\n" +
			"- Enter an operator to replace numbers on the stack with a result\n" +
			"- For example: 2 3 +",

		"conditionals": "The operators, >, >=, <, <=, =, != can be used to\n" +
			"compare numbers.  Note that > and friends ignore the complex\n" +
			"part of numbers.\n" +
			"Example: 3 5 > # this will put false on the stack",

		"control": "The if, ifelse, and for operators can be used to provide\n" +
			"support for simple programmng\n" +
			"Example (print 1-50): 1 'println 1 + $0 50 <' for",

		"complex": "Enter a complex value as i, -i, 3+i or 3-i\n" +
			"Do not use spaces.",

		"conversions": rpn.conv.Help(),

		"keymacros": "The variables .f1 to .f12 can be set to a string.\n" +
			"Pressing the corresponding function key will execute the string\n" +
			"as a macro.",

		"macros": "Execute a variable as @name.  Execute a string with just @\n" +
			"Convert any value to a string with str.\n" +
			"Example:\n" +
			"'. 3.14159 * *' cirarea=\n" +
			"5 @cirarea\n" +
			"See Also: keymacros, variables",

		"printing": "There are various printing functions that print values\n" +
			"at the head of the stack. These include:\n" +
			"  - print : print the value at the head of the stack\n" +
			"  - printx : pop and print the value at the head of the stack\n" +
			"  - prints : print with a space\n" +
			"  - printsx : printx with a space\n" +
			"  - println : print with a newline\n" +
			"  - printlnx : printx with a newline",

		"stack": "Operators are provided to manipute the stack to set up calculations\n" +
			"    If things are getting complex, consider using variables.\n" +
			"Examples:\n" +
			"  $x $0 + 1 /  # Uses $0 to copy the stack head and execute a/(a+1)'\n" +
			"  '$0 cos 1> sin' pplot # parametric plot of a circle using sw to swap elements\n" +
			"   0/ # drop the element at the head of the stack\n" +
			"   1/ # drop the second element from the stack\n" +
			"   X # clear the stack\n" +
			"   2> # Moves the third element in the stack to the head\n" +
			"   2< # Moves the head element two backwards\n" +
			"   $2 # Copies the third element in the stack to the head\n" +
			"\n",

		"strings": "Enter a string value as 'example 1' or \"example 2\"",

		"variables": "Set a variable as name=\n" +
			"Use a variable with $name\n" +
			"Example: 5 x= $x $x *\n" +
			"\n" +
			"Clear a variable with a trailing /. e.g. x/\n" +
			"Choose the first value of the stack: $0\n" +
			"\n" +
			"Variables can hold stacks of values:\n" +
			"  x> Pops a value from x and pushes it to the stack. Removes x if it is empty.\n" +
			"  x< Pops a value from the stack and pushes it to x, creating it if needed.\n" +
			"  x>> Pops all x values to the stack and removes x\n" +
			"  x<< Pushes all of the stack to x, creating it if needed\n" +
			"  x== Sets x equal to to all stack values (and clears stack)\n" +
			"  $$x Appends all values of x to the stack without changing x\n" +
			"See Also: macros, stack",
	}
	rpn.help = map[string]map[string]string{CatConcepts: conceptHelp}
}

func (r *RPN) RegisterConceptHelp(helpmap map[string]string) {
	for concept, help := range helpmap {
		r.help[CatConcepts][concept] = help
	}
}

func (r *RPN) printHelp(topic string) error {
	if len(topic) == 0 {
		r.listCommands()
		return nil
	}
	var help string
	for cat := range r.help {
		help = r.help[cat][topic]
		if help != "" {
			break
		}
	}
	if help == "" {
		return fmt.Errorf("use ? to list all: %w", ErrNotFound)
	}
	r.Print("\n")
	r.Println(help)
	return nil
}

func (r *RPN) listCommands() {
	var cats []string
	for cat := range r.help {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		r.dumpMap(cat, r.help[cat])
	}
}

const colWidth = 32

func (r *RPN) dumpMap(title string, m map[string]string) {
	r.Println(title)
	var topics []string
	for k := range m {
		topics = append(topics, k)
	}
	sort.Strings(topics)
	line := []byte("  ")
	nextCol := len(line) + colWidth
	for _, t := range topics {
		line = append(line, []byte(t)...)
		for len(line) < nextCol {
			line = append(line, ' ')
		}
		nextCol += colWidth
		if nextCol > r.TextWidth {
			r.Println(string(line))
			line = line[:2]
			nextCol = len(line) + colWidth
		}
	}
	if len(line) > 0 {
		r.Println(string(line))
	}
}
