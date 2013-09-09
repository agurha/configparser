package configparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config interface {
	AddConfigSection(section string) bool
	RemoveConfigSection(section string) bool
	AddSectionLeaf(section string, option string, value string) bool
	RemoveSectionLeaf(section string, option string) bool
}

type ConfigParser struct {
	config map[string]map[string]string
}

var DefaultSection = "default"

func (c *ConfigParser) AddConfigSection(section string) bool {

	if _, ok := c.config[section]; ok {

		return false
	}

	c.config[section] = make(map[string]string)

	return true
}

func (c *ConfigParser) RemoveConfigSection(section string) bool {

	switch _, ok := c.config[section]; {
	case !ok:
		return false

	case section == DefaultSection:
		return false

	default:
		for o, _ := range c.config[section] {
			delete(c.config[section], o)
		}

		delete(c.config, section)

	}
	return true
}

func (c *ConfigParser) AddSectionLeaf(section string, option string, value string) bool {

	// First we should make sure the section to which we would like to add the option exists
	c.AddConfigSection(section)

	if _, ok := c.config[section][option]; ok {

		return false
	}

	c.config[section][option] = value

	return true
}

func (c *ConfigParser) RemoveSectionLeaf(section string, option string) bool {

	if _, ok := c.config[section][option]; !ok {

		return false
	}

	delete(c.config[section], option)

	return true

}

func NewConfigFile() *ConfigParser {

	c := new(ConfigParser)

	c.config = make(map[string]map[string]string)

	c.AddConfigSection(DefaultSection)

	return c
}

func (c *ConfigParser) read(buf *bufio.Reader) error {

	var section, option string

	// Get the buffer reader and return the ConfigParser Pointer with added section to it from the buffer reader

	for {

		l, err := buf.ReadString('\n')

		if err == io.EOF {

			if len(l) == 0 {
				break
			}
		} else if err != nil {
			return err

		}

		l = strings.TrimSpace(l)

		switch {

		case len(l) == 0:
			continue

		case l[0] == '#': // assume comment
			continue

		case l[0] == '[' && l[len(l)-1] == ']': // this is a new section
			option = "" // reset multi-line option value
			section = strings.TrimSpace(l[1 : len(l)-1])
			c.AddConfigSection(section)

		default:
			i := strings.IndexRune(l, '=')
			switch {

			case i > 0:
				i := strings.IndexRune(l, '=')
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(l)
				c.AddSectionLeaf(section, option, value)

			default:
				return errors.New(fmt.Sprintf("Could not parse line %s", l))
			}

		}

	}

	return nil
}

func ReturnString(somestring string) (retstring string) {

	retstring = somestring + "Hello"
	return
}

func ReadConfigFile(filename string) (*ConfigParser, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	c := NewConfigFile()

	if err := c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ConfigParser) GetConfigSections() (sections []string) {

	sections = make([]string, len(c.config))

	i := 0
	for s, _ := range c.config {
		sections[i] = s
		i++
	}

	return sections

}

func (c *ConfigParser) GetSectionLeaves(section string) ([]string, error) {

	if _, ok := c.config[section]; !ok {

		return nil, errors.New(fmt.Sprintf("Section not found %s", section))
	}

	options := make([]string, len(c.config[DefaultSection])+len(c.config[section]))

	i := 0
	for s, _ := range c.config[DefaultSection] {

		options[i] = s
		i++
	}

	for s, _ := range c.config[section] {

		options[i] = s
		i++
	}

	return options, nil
}
