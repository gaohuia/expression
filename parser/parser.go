package parser

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

const OpAdd = 1
const OpSub = 2
const OpMul = 3
const OpDiv = 4
const OpVal = 5
const OpFunc = 6

const TokenTypeOperator = 1
const TokenTypeOperand = 2
const TokenTypeQuoteLeft = 3
const TokenTypeQuoteRight = 4
const TokenTypeComma = 5
const TokenTypeIdentifier = 6

type Token string
type Any interface{}
type Number int

var OperatorPriority = map[int]int{
	OpAdd: 1,
	OpSub: 1,
	OpDiv: 2,
	OpMul: 2,
}

func CompareOperatorPriority(op1, op2 int) int {
	return OperatorPriority[op1] - OperatorPriority[op2]
}

var IdentifierPattern = regexp.MustCompile(`[a-zA-Z]\w*`)
var OperandPattern = regexp.MustCompile(`\d+`)

type FunctionCall struct {
	Name string
	Args []*Operator
}

func IsIdentifier(token string) bool {
	return IdentifierPattern.MatchString(token)
}

func (t Token) GetType() int {
	switch t {
	case "+", "-", "*", "/":
		return TokenTypeOperator
	case "(":
		return TokenTypeQuoteLeft
	case ")":
		return TokenTypeQuoteRight
	case ",":
		return TokenTypeComma
	default:

		if IsIdentifier(string(t)) {
			return TokenTypeIdentifier
		} else if OperandPattern.MatchString(string(t)) {
			return TokenTypeOperand
		}

		panic("unknown token type")
	}
}

func (t Token) GetValue() int {
	if opType := t.GetType(); opType == TokenTypeOperator {
		switch t {
		case "+":
			return OpAdd
		case "-":
			return OpSub
		case "*":
			return OpMul
		case "/":
			return OpDiv
		}

		panic("unknown operator")
	} else if opType == TokenTypeOperand {
		if v, err := strconv.Atoi(string(t)); err != nil {
			panic(err)
		} else {
			return v
		}
	} else {
		panic("this token has no value: " + string(t))
	}
}

type Operator struct {
	OpType int
	Left   *Operator
	Right  *Operator
	Value  Any
	Leaf   bool
}

func (op *Operator) IsLeaf() bool {
	// return op.Left == nil && op.Right == nil
	return op.Leaf
}

func (op *Operator) Calculate() int {

	var (
		left  = 0
		right = 0
	)

	if op.Left != nil {
		left = op.Left.Calculate()
	}

	if op.Right != nil {
		right = op.Right.Calculate()
	}

	switch op.OpType {
	case OpAdd:
		return left + right
	case OpSub:
		return left - right
	case OpMul:
		return left * right
	case OpDiv:
		return left / right
	case OpVal:
		if v, ok := op.Value.(int); ok {
			return v
		} else {
			panic("unknown value type")
		}
	case OpFunc:
		if v, ok := op.Value.(*FunctionCall); ok {
			if v.Name == "random" {
				return rand.Int()
			} else if v.Name == "sum" {
				var sum = 0
				for _, arg := range v.Args {
					sum += arg.Calculate()
				}
				return sum
			} else if v.Name == "max" {

				if len(v.Args) == 0 {
					panic("max function requires at least one argument")
				}

				max := v.Args[0].Calculate()

				for _, arg := range v.Args {
					if arg.Calculate() > max {
						max = arg.Calculate()
					}
				}
				return max
			} else if v.Name == "min" {

				if len(v.Args) == 0 {
					panic("min function requires at least one argument")
				}

				min := v.Args[0].Calculate()

				for _, arg := range v.Args {
					if arg.Calculate() < min {
						min = arg.Calculate()
					}
				}
				return min
			} else {
				panic("unknown function name: " + v.Name)
			}
		}
	}

	return 0
}

func NewOperator(opType int, left *Operator, right *Operator) *Operator {
	return &Operator{opType, left, right, 0, false}
}

func NewValue(value int) *Operator {
	return &Operator{OpVal, nil, nil, value, true}
}

var ErrorNoMoreTokens = errors.New("no more token")
var ErrorUnexpectedToken = errors.New("unexpected token")

type Tokens struct {
	tokens []Token
	index  int
	total  int
}

func NewTokens(tokens []Token) *Tokens {
	return &Tokens{tokens: tokens, index: 0, total: len(tokens)}
}

func (t *Tokens) GetIndex() int {
	return t.index
}

func (t *Tokens) SetIndex(index int) {
	t.index = index
}

func (t *Tokens) ReturnToken() {
	t.index -= 1
}

func (t *Tokens) GetToken() (Token, error) {
	if t.index < t.total {
		var token = t.tokens[t.index]
		t.index++
		return token, nil
	}

	return "", ErrorNoMoreTokens
}

func (t *Tokens) GetTokenAt(index int) Token {
	if index < 0 || index >= t.total {
		panic("index out of range")
	}

	return t.tokens[index]
}

func (t *Tokens) GetTotalTokens() int {
	return t.total
}

func (t *Tokens) HasMoreTokens() bool {
	return t.index < t.total
}

func (t *Tokens) GetTokenOfType(tokenType int) (int, error) {
	if token, err := t.GetToken(); err != nil {
		return 0, err
	} else if token.GetType() != tokenType {
		t.ReturnToken()
		return 0, ErrorUnexpectedToken
	} else {
		return token.GetValue(), nil
	}
}

func GetQuotedOperand(tokens *Tokens) (*Operator, error) {
	index := tokens.GetIndex()
	token, err := tokens.GetToken()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tokens.SetIndex(index)
		}
	}()

	if token.GetType() == TokenTypeQuoteLeft {
		operator, err := GetExpression(tokens)

		if err != nil {
			return nil, err
		}

		token, err := tokens.GetToken()

		if err != nil {
			return nil, errors.New("syntax error, expected token ')', but " + err.Error())
		}

		if token.GetType() != TokenTypeQuoteRight {
			return nil, errors.New("unexpected token: " + string(token))
		}

		operator.Leaf = true
		return operator, nil
	} else {
		return nil, errors.New("syntax error, unexpected token: " + string(token))
	}
}

func GetArgumentList(tokens *Tokens) ([]*Operator, error) {
	var args = make([]*Operator, 0)

	for {
		operator, err := GetExpression(tokens)

		if err != nil {
			break
		}

		args = append(args, operator)

		token, err := tokens.GetToken()
		if err != nil {
			return nil, errors.New("syntax error, expected token ',' or ')'")
		}

		if token.GetType() == TokenTypeComma {
			continue
		} else {
			tokens.ReturnToken()
			break
		}
	}

	return args, nil
}

func GetFunctionCallOperand(tokens *Tokens) (op *Operator, err error) {

	index := tokens.GetIndex()

	defer func() {
		if op == nil {
			tokens.SetIndex(index)
		}
	}()

	token, err := tokens.GetToken()
	if err != nil {
		return
	}

	if token.GetType() == TokenTypeIdentifier {
		name := string(token)
		token, err = tokens.GetToken()

		if err != nil {
			return nil, errors.New("syntax error, expected token '('")
		}

		if token.GetType() != TokenTypeQuoteLeft {
			return nil, errors.New("syntax error, expected token '('")
		}

		var args []*Operator
		args, err = GetArgumentList(tokens)
		if err != nil {
			return nil, err
		}

		token, err = tokens.GetToken()
		if err != nil {
			return nil, errors.New("syntax error, expected token ')'")
		}

		if token.GetType() != TokenTypeQuoteRight {
			return nil, errors.New("syntax error, expected token ')'")
		}

		op = &Operator{OpFunc, nil, nil, &FunctionCall{name, args}, true}
		return
	} else {
		return nil, errors.New("syntax error, expected function name")
	}
}

func GetOperand(tokens *Tokens) (*Operator, error) {
	token, err := tokens.GetToken()
	if err != nil {
		return nil, err
	}

	if token.GetType() == TokenTypeQuoteLeft {
		tokens.ReturnToken()
		return GetQuotedOperand(tokens)
	}

	if token.GetType() == TokenTypeIdentifier {
		tokens.ReturnToken()
		return GetFunctionCallOperand(tokens)
	}

	if token.GetType() == TokenTypeOperand {
		return NewValue(token.GetValue()), nil
	}

	tokens.ReturnToken()
	return nil, ErrorUnexpectedToken
}

func GetExpression(tokens *Tokens) (*Operator, error) {
	operand, err := GetOperand(tokens)

	if err != nil {
		return nil, err
	}

	root := operand

	for {
		index := tokens.GetIndex()
		var operator int
		operator, err = tokens.GetTokenOfType(TokenTypeOperator)

		if err != nil {
			break
		}

		// Operand can be either a quoted expression or numeric literal
		operand, err := GetOperand(tokens)

		if err != nil {
			tokens.SetIndex(index)
			break
		}

		if root.IsLeaf() {
			root = NewOperator(operator, root, operand)
		} else {
			// append the new operator to the right most
			var path []*Operator
			rightMost := root

			// Find the right most leaf node
			for {
				if rightMost.IsLeaf() {
					break
				}

				path = append(path, rightMost)
				rightMost = rightMost.Right
			}

			// The operator of the right most node
			lastOperator := path[len(path)-1]

			// Add the operator
			lastOperator.Right = NewOperator(operator, lastOperator.Right, operand)

			// The initial value is the newly added operator node
			current := lastOperator.Right

			for i := len(path) - 1; i >= 0; i-- {
				parent := path[i]

				// Left rotate
				if CompareOperatorPriority(parent.OpType, current.OpType) >= 0 {

					// Top
					isRoot := i == 0

					parent.Right = current.Left
					current.Left = parent

					if isRoot {
						root = current
					} else {
						gParent := path[i-1]
						gParent.Right = current
					}

				} else {
					break
				}
			}
		}
	}

	return root, nil
}

func BuildCalculatorTreeFromTokens(tokens *Tokens) (*Operator, error) {
	operator, err := GetExpression(tokens)
	if err != nil {
		return nil, err
	}

	if tokens.HasMoreTokens() {
		return nil, errors.New("syntax error, unexpected token: " + string(tokens.tokens[tokens.index]))
	}

	return operator, nil
}

func BuildCalculatorTree(expression string) (*Operator, error) {
	var tokens = Tokenize(expression)
	return BuildCalculatorTreeFromTokens(tokens)
}

func Tokenize(expression string) *Tokens {

	var pattern = regexp.MustCompile(`[a-zA-Z]\w*|,|\d+|\+|\-|\*|\/|\(|\)|"`)

	var result = make([]Token, 0)
	var index = 0
	var length = len(expression)

	for index < length {
		if match := pattern.FindStringIndex(expression[index:]); match != nil {
			s := expression[match[0]+index : match[1]+index]

			if s != "\"" {
				result = append(result, Token(s))
				index += match[1]
			} else {
				index = match[0] + index
				quoteStart := index

				// Find the next quote, of which the previous character is not a backslash
				for {
					i := strings.Index(expression[index+1:], "\"")
					if i == -1 {
						panic("unterminated string")
					}

					if expression[index+i] != '\\' {
						result = append(result, Token(expression[quoteStart:index+i+2]))
						index += i + 2
						break
					} else {
						index += i + 1
					}
				}
			}

		} else {
			break
		}
	}

	return NewTokens(result)
}
