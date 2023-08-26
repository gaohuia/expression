package parser_test

import (
	"expression/parser"
	"testing"
)

func TestTokenize(t *testing.T) {
	var expression = `"hello" + "world" + "\"aaaa\""`

	tokens := parser.Tokenize(expression)
	t1 := tokens.GetTokenAt(0)

	if tokens.GetTotalTokens() != 5 {
		t.Errorf("Expected 5 tokens but got %d", tokens.GetTotalTokens())
	}

	if string(t1) != "\"hello\"" {
		t.Errorf("Expected 'hello' but got %s", t1)
	}

	t2 := tokens.GetTokenAt(1)
	if string(t2) != "+" {
		t.Errorf("Expected '+' but got %s", t2)
	}

	t3 := tokens.GetTokenAt(2)
	if string(t3) != "\"world\"" {
		t.Errorf("Expected 'world' but got %s", t3)
	}

	t4 := tokens.GetTokenAt(3)
	if string(t4) != "+" {
		t.Errorf("Expected '+' but got %s", t4)
	}

	t5 := tokens.GetTokenAt(4)
	if string(t5) != "\"\\\"aaaa\\\"\"" {
		t.Errorf("Expected '\"aaaa\"' but got %s", t5)
	}
}
