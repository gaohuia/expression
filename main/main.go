package main

import (
	"expression/parser"
	"log"
	"strings"
)

func Calc(expression string) int {
	tree, err := parser.BuildCalculatorTree(expression)
	if err != nil {
		log.Println(err)
		return -1
	}

	return tree.Calculate()
}

func main() {

	var expressions = []string{
		"1+(2+3)=6",
		"1+(2+3)*4=21",
		"1+2*3=7",
		"1*2*3*4=24",
		"1+2+3+4+5=15",
		"2*3+4=10",
		"4+5*2=14",
		"1+2*3+4=11",
		"1+2*3+4+5+6+7+8+9*10=127",
		"1+((2+3)*4)*5=101",
		"1+(2*3)*((4*5+1)+1)=133",
		"(((111)))=111",
		"min(1, max(1,2,3)+1)=1",
		"max(1,2,3,4,5,6,7,8,9,10)=10",
		"min(1,2,3,4,5,6,7,8,9,10)=1",
		"min(1,2,3,4,5,6,7,8,9,10)+max(1,2,3,4,5,6,7,8,9,10)=11",
		"min(1,2,3,4,5,6,7,8,9,10)*max(1,2,3,4,5,6,7,8,9,10)=10",
		"min(1,2,3,4,5,6,7,8,9,10)+max(1,2,3,4,5,6,7,8,9,10)*min(1,2,3,4,5,6,7,8,9,10)=11",
		"min(1,2,3,4,5,6,7,8,9,10)*max(1,2,3,4,5,6,7,8,9,10)+min(1,2,3,4,5,6,7,8,9,10)=11",
		"sum(1,2,3,4,5,6,7,8,9,10)=55",
		"sum(1+2+3, 4+5+6, 7+8+9, min(1,10))*max(1,2)=92",
		"sum(1+2+3, 2*(5+6), 7+8+9, min(1,10)) * max(1,2)=106",
		"sum(1+22+3, 2*(5+6), 7+8+9, min(1,10)) * max(1,2)=146",
	}

	// Print the header of the table
	log.Printf("%80s %10s %s", "Expression", "Result", "Expected")
	for _, expression := range expressions {

		segments := strings.SplitN(expression, "=", 2)

		result := Calc(segments[0])
		log.Printf("%80s %10s %d", segments[0], segments[1], result)
	}

}
