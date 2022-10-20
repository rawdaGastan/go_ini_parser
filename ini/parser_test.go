package ini

import (
	"reflect"
	"testing"
)

var sampleContent = map[string]string{
	//valid options
	"valid": "[owner]\nname=John\norganization = threefold\n\n[database]\nserver = 192.0.2.62\nport = 143\npassword = 123456\nprotected = true\nversion = 12.6",

	"valid_comment": ";comment",
	"valid_empty":   "",

	//invalid options
	"invalid":                     "[owner]\n--",
	"invalid_section":             "[[owner]",
	"invalid_unclosed_section":    "[owner\nkey=val\n",
	"invalid_unopened_section":    "owner]\nkey=val\n",
	"invalid_no_equal":            "[owner]\nkeyval",
	"invalid_no_value":            "[owner]\nkeyval=",
	"invalid_no_key":              "[owner]\n=keyval",
	"invalid_more_than_one_equal": "[owner]\nkey==val",

	"invalid_no_sections": "",
	"invalid_no_options":  "[owner]",
}

func TestValidParser(t *testing.T) {

	t.Run("testValid", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("Content is no valid")
		}
	})

	t.Run("test_valid_comment", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid_comment"])

		if err != nil {
			t.Errorf("Content is no valid")
		}
	})

	t.Run("test_valid_empty", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid_empty"])

		if err != nil {
			t.Errorf("Content is no valid")
		}
	})

	t.Run("test_value", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("function shouldn't throw errors")
		}

		if got, _ := parser.GetOption("owner", "name"); got != "John" {
			t.Errorf("Got %v, want John", got)
		}

		if got, _ := parser.GetOption("owner", "organization"); got != "threefold" {
			t.Errorf("Got %v, want threefold", got)
		}

		if got, _ := parser.GetOption("database", "server"); got != "192.0.2.62" {
			t.Errorf("Got %v, want 192.0.2.62", got)
		}

		if got, _ := parser.GetOption("database", "port"); got != "143" {
			t.Errorf("Got %v, want 143", got)
		}

		if got, _ := parser.GetOption("database", "password"); got != "123456" {
			t.Errorf("Got %v, want 123456", got)
		}

		if got, _ := parser.GetOption("database", "protected"); got != "true" {
			t.Errorf("Got %v, want true", got)
		}

		if got, _ := parser.GetOption("database", "version"); got != "12.6" {
			t.Errorf("Got %v, want 12.6", got)
		}

		if got, _ := parser.GetBool("database", "protected"); got != true {
			t.Errorf("Got %v, want true", got)
		}

		if got, _ := parser.GetInt("database", "port"); got != 143 {
			t.Errorf("Got %v, want 143", got)
		}

		if got, _ := parser.GetFloat("database", "port"); got != 143 {
			t.Errorf("Got %v, want 143", got)
		}

		if got, _ := parser.GetInt("database", "password"); got != 123456 {
			t.Errorf("Got %v, want 123456", got)
		}

		if got, _ := parser.GetFloat("database", "version"); got != 12.6 {
			t.Errorf("Got %v, want 12.6", got)
		}
	})

	t.Run("test_parsed_sections", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := []string{"owner", "database"}
		got := parser.GetSections()

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})

	t.Run("test_parsed_section", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := map[string]string{"name": "John", "organization": "threefold"}
		got, _ := parser.GetSection("owner")

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})

	t.Run("test_parsed_options", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := []string{"name", "organization"}
		got := parser.GetOptions("owner")

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})

	t.Run("test_parsed_option", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got, _ := parser.GetOption("owner", "name"); got != "John" {
			t.Errorf("Got %v, want John", got)
		}
	})

	t.Run("test_set_option", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parser.SetOption("owner", "name", "Ali")

		if got, _ := parser.GetOption("owner", "name"); got != "Ali" {
			t.Errorf("Got %v, want Ali", got)
		}
	})

	t.Run("test_parsed_functions", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parsed := parser.GetParsedMap()

		// parsed str
		testParsedStr := parser.String()

		// parsed map

		err = parser.FromString(testParsedStr)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		testParsedDict := parser.GetParsedMap()

		if !reflect.DeepEqual(testParsedDict, parsed) {
			t.Errorf("Got %v, want %v", testParsedDict, parsed)
		}
	})

	t.Run("test_no_sections", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_no_sections"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		sections := parser.GetSections()

		if len(sections) > 0 {
			t.Errorf("Content have sections")
		}
	})

	t.Run("test_no_options", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_no_options"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		options := parser.GetOptions("owner")

		if len(options) > 0 {
			t.Errorf("owner section has options")
		}
	})

	t.Run("test_set_option_old_option", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parser.SetOption("owner", "name", "Ali")

		if want, _ := parser.GetOption("owner", "name"); "John" == want {
			t.Errorf("Got John, want Ali")
		}
	})

	t.Run("test_set_option_no_option", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parser.SetOption("owner", "age", "30")

		if want, _ := parser.GetOption("owner", "age"); "30" != want {
			t.Errorf("Got error, want 30")
		}
	})

}

func TestInValidParser(t *testing.T) {

	t.Run("test_invalid", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid"])

		if err == nil {
			t.Errorf("Content is valid")
		}
	})

	t.Run("test_invalid_section", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_section"])

		if err == nil {
			t.Errorf("Section is valid")
		}
	})

	t.Run("test_unclosed_section", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_unclosed_section"])

		if err == nil {
			t.Errorf("Section is closed")
		}
	})

	t.Run("test_unopened_section", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_unopened_section"])

		if err == nil {
			t.Errorf("Section is opened")
		}
	})

	t.Run("test_no_equal", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_no_equal"])

		if err == nil {
			t.Errorf("Section has equal")
		}
	})

	t.Run("test_no_value", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_no_value"])

		if err == nil {
			t.Errorf("Section has value")
		}
	})

	t.Run("test_no_key", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_no_key"])

		if err == nil {
			t.Errorf("Section has key")
		}
	})

	t.Run("test_more_than_one_equal", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["invalid_more_than_one_equal"])

		if err == nil {
			t.Errorf("Section has one equal")
		}
	})
}

func TestWrongValues(t *testing.T) {

	t.Run("test_wrong_section", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		_, err = parser.GetSection("ownerr")

		if err == nil {
			t.Errorf("ownerr section exists")
		}
	})

	t.Run("test_wrong_value", func(t *testing.T) {
		parser := NewParser()
		err := parser.FromString(sampleContent["valid"])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if option, _ := parser.GetOption("owner", "server"); option == "John" {
			t.Errorf("Wrong option value John")
		}
	})

	t.Run("test_wrong_bool", func(t *testing.T) {
		parser := NewParser()
		validErr := parser.FromString(sampleContent["valid"])

		if validErr != nil {
			t.Errorf("unexpected error: %v", validErr)
		}

		_, err := parser.GetBool("database", "server")

		if err == nil {
			t.Errorf("Valid boolean")
		}
	})

	t.Run("test_wrong_int", func(t *testing.T) {
		parser := NewParser()
		validErr := parser.FromString(sampleContent["valid"])

		if validErr != nil {
			t.Errorf("unexpected error: %v", validErr)
		}

		_, err := parser.GetInt("database", "protected")

		if err == nil {
			t.Errorf("Valid ineteger")
		}
	})

	t.Run("test_wrong_float", func(t *testing.T) {
		parser := NewParser()
		validErr := parser.FromString(sampleContent["valid"])

		if validErr != nil {
			t.Errorf("unexpected error: %v", validErr)
		}

		_, err := parser.GetFloat("database", "protected")

		if err == nil {
			t.Errorf("Valid float")
		}
	})
}
