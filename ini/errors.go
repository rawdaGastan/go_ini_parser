package ini

import "errors"

var NoParsedData = errors.New("There is no parsed data")
var InvalidSection = errors.New("Invalid section! Please make sure that you have one '[' and one ']'")
var InvalidContent = errors.New("Not a valid ini content")
var SectionNotExist = errors.New("This section does not exist")
var OptionsNotExist = errors.New("This section does not have any options")
var OptionNotExist = errors.New("This option does not exist in the given section")

var InvalidBoolean = errors.New("Not a valid boolean")
