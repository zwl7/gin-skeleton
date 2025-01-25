package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// ConfigMap ...
type ConfigMap map[string]map[string]string

// DefaultSection ...
var DefaultSection = "default"

// RootPath 获取根目录的绝对路径
func RootPath() string {
	tempSlice := []string{}
	if runtime.GOOS == "windows" {
		//windows系统
		root, _ := exec.LookPath(os.Args[0])
		rootPath, _ := filepath.Abs(root)
		tempSlice = strings.Split(rootPath, `\`)
	} else {
		//其他系统
		root, _ := exec.LookPath(os.Args[0])
		rootPath, _ := filepath.Abs(root)
		tempSlice = strings.Split(rootPath, "/")
	}
	return strings.Join(tempSlice[0:len(tempSlice)-1], "/")
}

// ParseFile ...
func ParseFile(fileName string) (cfg ConfigMap, err error) {
	var file *os.File
	cfg = make(ConfigMap, 0)
	file, err = os.OpenFile(fileName, os.O_RDONLY, 0755)
	if err != nil {
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	//
	var (
		configSection    = regexp.MustCompile("^\\s*\\[\\s*(\\w+)\\s*\\]\\s*$")
		quotedConfigLine = regexp.MustCompile("^\\s*(\\w+)\\s*=\\s*[\"'](.*)[\"']\\s*$")
		configLine       = regexp.MustCompile("^\\s*(\\w+)\\s*=\\s*(.*)\\s*$")
		commentLine      = regexp.MustCompile("^#.*$")
		commentLine2     = regexp.MustCompile("^;.*$")
		blankLine        = regexp.MustCompile("^\\s*$")
	)
	//
	var (
		line           string
		longLine       bool
		currentSection string
		lineBytes      []byte
		isPrefix       bool
	)

	for {
		err = nil
		lineBytes, isPrefix, err = buf.ReadLine()
		if io.EOF == err {
			err = nil
			break
		} else if err != nil {
			break
		} else if isPrefix {
			line += string(lineBytes)

			longLine = true
			continue
		} else if longLine {
			line += string(lineBytes)
			longLine = false
		} else {
			line = string(lineBytes)
		}
		line = strings.TrimPrefix(line, string([]byte{239, 187, 191}))
		//fmt.Println(line)
		if commentLine.MatchString(line) {
			continue
		} else if commentLine2.MatchString(line) {
			continue
		} else if blankLine.MatchString(line) {
			continue
		} else if configSection.MatchString(line) {
			section := configSection.ReplaceAllString(line,
				"$1")
			if section == "" {
				err = fmt.Errorf("invalid structure in file")
				break
			} else if !cfg.SectionInConfig(section) {
				cfg[section] = make(map[string]string, 0)
			}
			currentSection = section
		} else if configLine.MatchString(line) {
			regex := configLine
			if quotedConfigLine.MatchString(line) {
				regex = quotedConfigLine
			}
			if currentSection == "" {
				currentSection = DefaultSection
				if !cfg.SectionInConfig(currentSection) {
					cfg[currentSection] = make(map[string]string, 0)
				}
			}
			key := regex.ReplaceAllString(line, "$1")
			val := regex.ReplaceAllString(line, "$2")
			if key == "" {
				continue
			}
			cfg[currentSection][key] = val
		} else {
			//fmt.Println("!!!!!",line,"!!!!!")
			err = fmt.Errorf("invalid config file")
			break
		}
	}
	return
}

// SectionInConfig ...
func (that *ConfigMap) SectionInConfig(section string) bool {
	for s, _ := range *that {
		if section == s {
			return true
		}
	}
	return false
}

// ListSections ...
func (that *ConfigMap) ListSections() (sections []string) {
	for section, _ := range *that {
		sections = append(sections, section)
	}
	return
}

// WriteFile ...
func (that *ConfigMap) WriteFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, section := range that.ListSections() {
		sName := fmt.Sprintf("[ %s ]\n", section)
		_, err = file.Write([]byte(sName))
		if err != nil {
			return
		}

		for k, v := range (*that)[section] {
			line := fmt.Sprintf("%s = %s\n", k, v)
			_, err = file.Write([]byte(line))
			if err != nil {
				return
			}
		}
		_, err = file.Write([]byte{0x0a})
		if err != nil {
			return
		}
	}
	return
}

// AddSection ...
func (that *ConfigMap) AddSection(section string) {
	if nil != (*that)[section] {
		(*that)[section] = make(map[string]string, 0)
	}
}

// AddKeyVal ...
func (that *ConfigMap) AddKeyVal(section, key, val string) {
	if "" == section {
		section = DefaultSection
	}
	if nil == (*that)[section] {
		that.AddSection(section)
	}
	(*that)[section][key] = val
}

// GetValue ...
func (that *ConfigMap) GetValue(section, key string) (val string, present bool) {
	if that == nil {
		return
	}
	if section == "" {
		section = DefaultSection
	}
	cm := *that
	_, ok := cm[section]
	if !ok {
		return
	}
	val, present = cm[section][key]
	return
}

// Has ...
func (that *ConfigMap) Has(sectionAndKey string) bool {
	section := DefaultSection
	key := sectionAndKey
	_slice := strings.Split(sectionAndKey, ".")
	if len(_slice) > 1 {
		section = _slice[0]
		key = strings.Join(_slice[1:], ".")
	}
	cm := *that
	if _, _ok := cm[section]; !_ok {
		return false
	}
	if _, _ok := cm[section][key]; !_ok {
		return false
	}
	return true
}

// StringDefault ...
func (that *ConfigMap) StringDefault(sectionAndKey, defaultValue string) string {
	section := DefaultSection
	key := sectionAndKey
	_slice := strings.Split(sectionAndKey, ".")
	if len(_slice) > 1 {
		section = _slice[0]
		key = strings.Join(_slice[1:], ".")
	}
	cm := *that
	_, ok := cm[section]
	if !ok {
		return defaultValue
	}
	_v, _ok := cm[section][key]
	if !_ok {
		return defaultValue
	}
	return _v
}

// IntDefault ...
func (that *ConfigMap) IntDefault(sectionAndKey string, defaultValue int) int {
	section := DefaultSection
	key := sectionAndKey
	_slice := strings.Split(sectionAndKey, ".")
	if len(_slice) > 1 {
		section = _slice[0]
		key = strings.Join(_slice[1:], ".")
	}
	cm := *that
	_, ok := cm[section]
	if !ok {
		return defaultValue
	}
	_v, _ok := cm[section][key]
	if !_ok {
		return defaultValue
	}
	return ToInt(_v)
}

// Int64Default ...
func (that *ConfigMap) Int64Default(sectionAndKey string, defaultValue int64) int64 {
	section := DefaultSection
	key := sectionAndKey
	_slice := strings.Split(sectionAndKey, ".")
	if len(_slice) > 1 {
		section = _slice[0]
		key = strings.Join(_slice[1:], ".")
	}
	cm := *that
	_, ok := cm[section]
	if !ok {
		return defaultValue
	}
	_v, _ok := cm[section][key]
	if !_ok {
		return defaultValue
	}
	return ToInt64(_v)
}

// BoolDefault ...
func (that *ConfigMap) BoolDefault(sectionAndKey string, defaultValue bool) bool {
	section := DefaultSection
	key := sectionAndKey
	_slice := strings.Split(sectionAndKey, ".")
	if len(_slice) > 1 {
		section = _slice[0]
		key = strings.Join(_slice[1:], ".")
	}
	cm := *that
	_, ok := cm[section]
	if !ok {
		return defaultValue
	}
	_v, _ok := cm[section][key]
	if !_ok {
		return defaultValue
	}
	switch strings.TrimSpace(_v) {
	case "0":
		return false
	case "false":
		return false
	case "FALSE":
		return false
	default:
		return true
	}
}
