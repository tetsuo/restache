package stache_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/tetsuo/stache"
)

type testCase struct {
	line     int
	data     string
	expected string
}

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("test error")
}

func TestTokenizer(t *testing.T) {
	const file = "testdata/lexer_testcases.txt"
	for _, tc := range buildTestcases(t, file) {
		t.Run(fmt.Sprintf("%s L%d", file, tc.line), func(t *testing.T) {
			r := strings.NewReader(tc.data)
			z := stache.NewTokenizer(r)

			var ops []string
			for {
				tt := z.Next()
				if tt == stache.ErrorToken {
					if err := z.Err(); err != io.EOF {
						t.Errorf("expected io.EOF, got %v", err)
					}
					break
				}
				var op string

				switch tt {
				case stache.StartTagToken, stache.SelfClosingTagToken:
					tagName, hasAttr := z.TagName()
					if tt == stache.StartTagToken {
						op += "open(" + string(tagName)
					} else {
						op += "openclose(" + string(tagName)
					}
					if hasAttr {
						op += ", "
						var (
							attrKey, attrValue []byte
							isExpr, moreAttr   bool
						)
						moreAttr = true
						for moreAttr {
							attrKey, attrValue, isExpr, moreAttr = z.TagAttr()
							op += string(attrKey) + "="
							if isExpr {
								op += "expr("
							} else {
								op += "text("
							}
							op += string(attrValue) + ")"
							if moreAttr {
								op += " "
							}
						}
					}
					op += ")"
				case stache.EndTagToken:
					tagName, _ := z.TagName()
					op += "close(" + string(tagName) + ")"
				case stache.WhenToken:
					controlName := z.ControlName()
					op += "when(" + string(controlName) + ")"
				case stache.UnlessToken:
					controlName := z.ControlName()
					op += "unless(" + string(controlName) + ")"
				case stache.RangeToken:
					controlName := z.ControlName()
					op += "range(" + string(controlName) + ")"
				case stache.EndControlToken:
					controlName := z.ControlName()
					op += "endctl(" + string(controlName) + ")"
				case stache.VariableToken:
					varName := z.Raw()
					op += "expr(" + string(varName) + ")"
				case stache.CommentToken:
					comment := z.Comment()
					op += "comment(" + string(comment) + ")"
				case stache.TextToken:
					text := z.Raw()
					if len(bytes.TrimSpace(text)) == 0 {
						continue
					}
					op += "text(" + string(text) + ")"
				}
				ops = append(ops, op)
			}

			actual := strings.TrimSpace(strings.Join(ops, "\n"))
			if actual != tc.expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v\n", "\t"+strings.Join(strings.Split(tc.expected, "\n"), "\n\t"), "\t"+strings.Join(strings.Split(actual, "\n"), "\n\t"))
			}
		})
	}

	t.Run("call after error", func(t *testing.T) {
		z := stache.NewTokenizer(strings.NewReader(""))
		for range 2 {
			tt := z.Next()
			if tt != stache.ErrorToken {
				t.Errorf("expected ErrorToken, got %v", tt)
			}
			if err := z.Err(); err != io.EOF {
				t.Errorf("expected io.EOF, got %v", err)
			}
		}
	})

	t.Run("faulty reader", func(t *testing.T) {
		z := stache.NewTokenizer(&errorReader{})
		tt := z.Next()
		if tt != stache.ErrorToken {
			t.Errorf("expected ErrorToken, got %v", tt)
		}
		if err := z.Err().Error(); err != "test error" {
			t.Errorf("expected test error, got %v", err)
		}
	})
}

func buildTestcases(t *testing.T, filename string) []testCase {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var testCases []testCase
	var dataBuilder, expectedBuilder strings.Builder
	scanner := bufio.NewScanner(file)

	reading := true

	lineStart := 0
	lineEnd := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineEnd += 1
		if strings.TrimSpace(line) == "%" {
			if !reading {
				testCases = append(testCases, testCase{
					line:     lineStart,
					data:     strings.TrimSpace(dataBuilder.String()),
					expected: strings.TrimSpace(expectedBuilder.String()),
				})
				lineStart = lineEnd
				dataBuilder.Reset()
				expectedBuilder.Reset()
			}
			// Flip: data -> expected, or expected -> next data
			reading = !reading
			continue
		}

		if reading {
			dataBuilder.WriteString(line + "\n")
		} else {
			expectedBuilder.WriteString(line + "\n")
		}
	}

	if dataBuilder.Len() > 0 && expectedBuilder.Len() > 0 {
		testCases = append(testCases, testCase{
			line:     lineStart,
			data:     strings.TrimSpace(dataBuilder.String()),
			expected: strings.TrimSpace(expectedBuilder.String()),
		})
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	return testCases
}
