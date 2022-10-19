package ini

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

// Helpers
func contains(array []string, element string) bool {
	for _, x := range array {
		if x == element {
			return true
		}
	}
	return false
}

type Entry = map[string]string
type INIMap = map[string]Entry

// Parser
type Parser struct {
	parsedMap INIMap
}

/////////////////////
// inner functions //
/////////////////////

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

////////////////////
// user functions //
////////////////////

func (p *Parser) GetParsedMap() INIMap {
	return p.parsedMap
}

func (p *Parser) ToString() string {
	var parsedStr string = ""

	for parent, dict := range p.parsedMap {
		parsedStr += fmt.Sprintf("[%s]\n", parent)
		for key, val := range dict {
			parsedStr += fmt.Sprintf("%s = %s\n", key, val)
		}
		parsedStr += "\n"
	}

	return parsedStr
}

func (p *Parser) FromFile(filename string) error {
	file, err := os.ReadFile(filename)

	if err == nil {
		p.FromString(string(file))
	}

	return err
}

// TODO: SaveToFile, read from file
func (p *Parser) SaveToFile(path string) error {
	// saves the parsed dict as ini file in the specified path after converting to string
	file, err := os.Create(path)
	defer file.Close()

	if err == nil {
		_, err = file.WriteString(p.ToString())
	}

	return err
}

func (p *Parser) FromString(content string) error {
	// parse the ini content as :
	// [section1: map[key1: value1 key2: value2 ...] section2: map[key1: value1 key2: value2 ...] ...]

	var err error

	p.parsedMap = INIMap{}
	var key string
	var val string
	var parent string

	// for parents
	var newParent bool = false

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
					err = InvalidSection
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
				err = InvalidContent
			}

			// parse empty
		} else if s.Trim(line, " ") == "" {
			continue

			// invlid content
		} else {
			err = InvalidContent
		}
	}

	return err
}

func (p *Parser) ToDict() INIMap {
	// Get the parsed dictionary as :
	// [section1: map[key1: value1 key2: value2 ...] section2: map[key1: value1 key2: value2 ...] ...]

	return p.parsedMap
}

func (p *Parser) GetSections() ([]string, error) {
	// returns a list of all sections --> [section1 section2 ...]

	var err error
	var sections []string

	for parent := range p.parsedMap {
		sections = append(sections, parent)
	}

	// if no data in the parsed dict
	if len(sections) == 0 {
		err = NoParsedData
	}

	return sections, err
}

func (p *Parser) GetSection(sectionKey string) (Entry, error) {
	// returns a dictionary for the section given --> map[key1: val1 key2: val2 ....]

	sections, err := p.GetSections()
	var section Entry

	if err == nil && !contains(sections, sectionKey) {
		err = SectionNotExist
	} else {
		section = p.parsedMap[sectionKey]
	}
	return section, err
}

func (p *Parser) GetOptions(sectionKey string) ([]string, error) {
	// returns all options of the given section

	section, err := p.GetSection(sectionKey)
	var options []string

	if err == nil {
		for key := range section {
			options = append(options, key)
		}

		if len(options) == 0 {
			err = OptionsNotExist
		}
	}

	return options, err
}

func (p *Parser) GetOption(sectionKey string, optionKey string) (string, error) {
	// returns the value of the option key which belongs to the section key given

	section, err := p.GetOptions(sectionKey)

	var options []string
	var option string

	if err == nil {
		for _, value := range section {
			options = append(options, value)
		}

		if !contains(section, optionKey) {
			err = OptionNotExist
		} else {
			option = p.parsedMap[sectionKey][optionKey]
		}
	}

	return option, err
}

func (p *Parser) SetOption(sectionKey string, optionKey string, optionValue string) error {
	// update the option in the given section

	section, err := p.GetOptions(sectionKey)

	var options []string

	if err == nil {
		for _, value := range section {
			options = append(options, value)
		}

		// if no data in the parsed dict
		if !contains(section, optionKey) {
			err = OptionNotExist
		} else {
			p.parsedMap[sectionKey][optionKey] = optionValue
		}
	}

	return err
}

////////////////////
// values getters //
////////////////////

func (p *Parser) GetBool(sectionKey string, optionKey string) (bool, error) {
	// returns the bool value of the option key which belongs to the section key given
	// bool could be : true, false, True, False, yes, no, 1, 0

	boolOption := false

	option, err := p.GetOption(sectionKey, optionKey)

	if err == nil {
		if contains([]string{"true", "True", "yes", "1"}, option) {
			boolOption = true

		} else if contains([]string{"false", "False", "no", "0"}, option) {
			boolOption = false

		} else {
			err = InvalidBoolean
		}
	}

	return boolOption, err
}

func (p *Parser) GetInt(sectionKey string, optionKey string) (int64, error) {
	// returns the int value of the option key which belongs to the section key given

	option, err := p.GetOption(sectionKey, optionKey)
	var parsedInt int64

	if err == nil {
		parsedInt, err = strconv.ParseInt(option, 10, 32)
	}

	return parsedInt, err
}

func (p *Parser) GetFloat(sectionKey string, optionKey string) (float64, error) {
	// returns the float value of the option key which belongs to the section key given

	option, err := p.GetOption(sectionKey, optionKey)
	var parsedFloat float64

	if err == nil {
		parsedFloat, err = strconv.ParseFloat(option, 64)
	}

	return parsedFloat, err
}
