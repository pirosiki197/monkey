package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/pirosiki197/monkey/evaluator"
	"github.com/pirosiki197/monkey/lexer"
	"github.com/pirosiki197/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	evaluator := evaluator.New()
	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(os.Stdout, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			out.Write([]byte{'\n'})
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
