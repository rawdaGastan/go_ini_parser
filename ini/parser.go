package ini

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

type Entry = map[string]string
type INIMap = map[string]Entry

// Parser
type Parser struct {
	parsedMap INIMap
}

func (p *Parser) addParent(parent string) {
	// check if the parent does not exist in the parsed dict
	if p.parsedMap[parent] == nil {
		p.parsedMap[parent] = Entry{}
	}
}

func (p *Parser) add(parent string, key string, value string) {
	// add the parent first
	p.addParent(parent)

	p.parsedMap[parent][key] = value
}

// GetParsedMap return the parsed map of the struct
func (p *Parser) GetParsedMap() INIMap {
	return p.parsedMap
}

// String converts the parsed map to ini string and returns it
func (p *Parser) String() string {
	formattedStr := ""

	for parent, dict := range p.parsedMap {
		formattedStr += fmt.Sprintf("[%s]\n", parent)
		for key, val := range dict {
			formattedStr += fmt.Sprintf("%s = %s\n", key, val)
		}
		formattedStr += "\n"
	}

	return formattedStr
}

// FromFile read an ini file and converts it to a parsed map
func (p *Parser) FromFile(filename string) error {
	file, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	p.FromString(string(file))
	return nil
}

// SaveToFile saves the parsed map as ini file in the specified path after converting to string
func (p *Parser) SaveToFile(path string) error {

	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.WriteString(p.String())

	return err
}

// FromString parses the ini content as :
// [section1: map[key1: value1 key2: value2 ...] section2: map[key1: value1 key2: value2 ...] ...]
func (p *Parser) FromString(content string) error {

	p.parsedMap = INIMap{}
	key := ""
	val := ""
	parent := ""

	// for parents
	newParent := false

	// read content lines
	scanner := bufio.NewScanner(s.NewReader(content))
	for scanner.Scan() {
		s.ReplaceAll(scanner.Text(), "\n", "")
		s.ReplaceAll(scanner.Text(), "\r", "")

		line := scanner.Text()

		if len(line) > 0 {
			// parse sections
			if string(line[0]) == "[" && string(line[len(line)-1]) == "]" {
				// check number of opened and closed sections []
				if s.Count(line, "[") == 1 && s.Count(line, "]") == 1 {
					parent = line[1 : len(line)-1]
					p.addParent(parent)
					newParent = true
				} else {
					return fmt.Errorf("invalid section %v! please make sure that you have one '[' and one ']'", parent)
				}

				// parse sections values
			} else if newParent && s.Count(line, "=") == 1 && (!contains([]string{"", "=", " "}, string(line[0])) && !contains([]string{"", "=", " "}, string(line[len(line)-1]))) {
				if s.Contains(line, " = ") {
					splitted := s.Split(line, " = ")
					key = splitted[0]
					val = splitted[1]
				} else if s.Contains(line, "=") {
					splitted := s.Split(line, "=")
					key = splitted[0]
					val = splitted[1]
				}

				p.add(parent, key, val)

				// parse comment lines
			} else if string(line[0]) == ";" {
				continue

				// invlid content
			} else {
				return fmt.Errorf("not a valid ini content")
			}

			// parse empty
		} else if s.Trim(line, " ") == "" {
			continue

			// invlid content
		} else {
			return fmt.Errorf("not a valid ini content")
		}
	}

	return nil
}

// GetSections returns a list of all sections --> [section1 section2 ...]
func (p *Parser) GetSections() []string {

	var sections []string

	for parent := range p.parsedMap {
		sections = append(sections, parent)
	}

	return sections
}

// GetSection returns a dictionary for the section given --> map[key1: val1 key2: val2 ....]
func (p *Parser) GetSection(sectionKey string) (Entry, error) {

	sections := p.GetSections()
	var section Entry

	if !contains(sections, sectionKey) {
		return section, fmt.Errorf("section %v does not exist", sectionKey)
	} else {
		section = p.parsedMap[sectionKey]
	}
	return section, nil
}

// GetOptions returns all options of the given section
func (p *Parser) GetOptions(sectionKey string) []string {

	section, _ := p.GetSection(sectionKey)
	var options []string

	for key := range section {
		options = append(options, key)
	}

	return options
}

// GetOption returns the value of the option key which belongs to the section key given
func (p *Parser) GetOption(sectionKey string, optionKey string) (string, error) {

	option := ""

	if option, exists := p.parsedMap[sectionKey][optionKey]; exists {
		return option, nil
	}

	return option, fmt.Errorf("option %v does not exist in the given section %v", optionKey, sectionKey)
}

// SetOption updates the option value in the given section
// If the section key does not exist it inserts the section key in the map with its option and value
func (p *Parser) SetOption(sectionKey string, optionKey string, optionValue string) error {

	section := p.GetOptions(sectionKey)

	var options []string

	for _, value := range section {
		options = append(options, value)
	}

	p.parsedMap[sectionKey][optionKey] = optionValue

	return nil
}

// GetBool returns the bool value of the option key which belongs to the section key given
// Bool could be : true, false, True, False, yes, no, 1, 0
func (p *Parser) GetBool(sectionKey string, optionKey string) (bool, error) {

	boolOption := false

	option, err := p.GetOption(sectionKey, optionKey)

	if err != nil {
		return boolOption, err
	}

	if contains([]string{"true", "True", "yes", "1"}, option) {
		boolOption = true

	} else if contains([]string{"false", "False", "no", "0"}, option) {
		boolOption = false

	} else {
		return boolOption, fmt.Errorf("%v is not a valid boolean", option)
	}

	return boolOption, nil
}

// GetInt returns the int value of the option key which belongs to the section key given
func (p *Parser) GetInt(sectionKey string, optionKey string) (int64, error) {

	option, err := p.GetOption(sectionKey, optionKey)
	var parsedInt int64

	if err != nil {
		return parsedInt, err
	}

	parsedInt, err = strconv.ParseInt(option, 10, 32)

	return parsedInt, err
}

// GetFloat returns the float value of the option key which belongs to the section key given
func (p *Parser) GetFloat(sectionKey string, optionKey string) (float64, error) {

	option, err := p.GetOption(sectionKey, optionKey)
	var parsedFloat float64

	if err != nil {
		return parsedFloat, err
	}

	parsedFloat, err = strconv.ParseFloat(option, 64)

	return parsedFloat, err
}
