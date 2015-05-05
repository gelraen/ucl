package ucl

import (
	"testing"
)

func TestParseJSON(t *testing.T) {
	good := []string{
		`{}`,
		`[]`,
		`[true, false, null]`,
		`[123]`,
		`{"1":"2"}`,
		`{"a":123}`,
		`{"a":1,"b":"321"}`,
		`{"a":123.45}`,
		`{"a":123.45e67}`,
		`{"a":-123}`,
		`{"a":-123.45}`,
		`{"a":-123.45e-67}`,
		`{"a":["b"]}`,
		`{"a":{"b":{"c":123}}}`,
		`[[[123,"1 ", "1", 123 , {}]]]`,
	}
	for _, s := range good {
		if v, err := parse([]rune(s)); err != nil {
			t.Errorf("Failed to parse '%s': %s", s, err)
		} else {
			t.Logf("'%s' -> %s", s, v)
		}
	}
	bad := []string{
		`""`,
		`["""]`,
		`{}{}`,
		`{{}}`,
		`{[]:{}}`,
		`{{{{{`,
		`}{`,
		`][`,
		`"a"a"`,
	}
	for _, s := range bad {
		_, err := parse([]rune(s))
		if err == nil {
			t.Errorf("Parse succeeded on invalid JSON document '%s'", s)
		} else {
			t.Logf("%s: %s", s, err)
		}
	}
}

func TestUnescape(t *testing.T) {
	cases := map[string]string{
		`a`:         `a`,
		`\n`:        "\n",
		`aa\na`:     "aa\na",
		`aa\\n`:     `aa\n`,
		`aa\\n\\a`:  `aa\n\a`,
		`aa\\u1234`: `aa\u1234`,
		`aa\u1234`:  "aa\u1234",
	}
	for in, expected := range cases {
		if got, err := json_unescape(in); err != nil || got != expected {
			t.Errorf("json_unescape(%q) = (%q, %s), want (%q, nil)", in, got, err, expected)
		}
	}
	bad := []string{
		`\`,
		`\u123G`,
		`\o`,
		`\u12`,
	}
	for _, s := range bad {
		_, err := json_unescape(s)
		if err == nil {
			t.Errorf("Parse succeeded on invalid JSON document '%s'", s)
		} else {
			t.Logf("%s: %s", s, err)
		}
	}
}
