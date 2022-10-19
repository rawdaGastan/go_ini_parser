package main

import (
	"fmt"
	"parser/ini"
)

func main() {
	parser := ini.Parser{}

	sampleContent := map[string]string{
		"valid": "[owner]\nname=John\norganization = threefold\n\n[database]\nserver = 192.0.2.62\nport = 143\npassword = 123456\nprotected = true\nversion = 12.6\n\n",
	}

	parser.FromString(sampleContent["valid"])

	parser.SaveToFile("data.txt")
	parser.FromFile("data.txt")

	sections, _ := parser.GetSections()
	options, _ := parser.GetOptions(sections[0])

	fmt.Println("Debug map: ", parser.GetParsedMap())
	fmt.Println("sections: ", sections)
	fmt.Println("options: ", options)
}
