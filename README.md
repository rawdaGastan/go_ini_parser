# ini parser

## How to use

- Create a new *.ini file and insert your ini content, for example:
  
```ini
[owner]
name = John
organization = threefold

[database]
version = 12.6
server = 192.0.2.62
port = 143
password = 123456
protected = true
```

- Create a new parser struct

```go
import "github.com/rawdaGastan/go_ini_parser/ini"

parser := NewParser()
```

- You can parse a file

```go
parser.FromFile("INI_FILE_PATH")
```

- You can parse a string

```go
iniContent := "[owner]\nname=John\norganization = threefold"
parser.FromString(iniContent)
```

## Functions

- `parser.GetParsedMap()` &rarr; to get your parsed map
- `parser.String()` &rarr; to convert your parsed map to ini string
- `parser.FromFile( filename )` &rarr; to convert your ini file to a parsed map
- `parser.SaveToFile( path )` &rarr; to save your ini string converted parsed map
- `parser.FromString( content )` &rarr; to convert your ini string to a parsed map
- `parser.GetSections()` &rarr; to get your sections' names
- `parser.GetSection( sectionKey )` &rarr; to get the content of the specified section key
- `parser.GetOptions( sectionKey )` &rarr; to get the options of the specified section key
- `parser.GetOption( sectionKey, optionKey )` &rarr; to get the string value of an option key inside a section
- `parser.SetOption( sectionKey, optionKey )` &rarr; to set the string value of an option inside a section
- `parser.GetBool( sectionKey, optionKey )` &rarr; to set the bool value of an option inside a section
- `parser.GetInt( sectionKey, optionKey )` &rarr; to set the integer value of an option inside a section
- `parser.GetFloat( sectionKey, optionKey )` &rarr; to set the float value of an option inside a section

## Testing

Use this command to run the tests

```bash
go test -v ./...
```
