package config

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig(flags *flag.FlagSet) {
	initConfig(flags, *GetConfig(), "")
}

func initConfig(flags *flag.FlagSet, config interface{}, prefix string) {
	vp := reflect.ValueOf(config)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(config)) {
		envKey := field.Tag.Get("env")
		usage := field.Tag.Get("usage")
		if envKey == "-" {
			continue // skip
		}
		if len(envKey) == 0 {
			envKey = field.Name
		}
		kebabKey := strcase.ToKebab(field.Name)

		// apply prefix
		if len(prefix) > 0 {
			envKey = fmt.Sprintf("%v_%v", prefix, envKey)
			kebabKey = fmt.Sprintf("%v-%v", strings.ToLower(prefix), kebabKey)
		}

		switch field.Type.Kind() {
		case reflect.String:
			value := vp.FieldByName(field.Name)
			valueIf := value.Interface().(string)
			flags.StringVar(&valueIf, kebabKey, value.String(), usage)
			viper.RegisterAlias(field.Name, kebabKey)
			_ = viper.BindEnv(kebabKey, envKey)
		case reflect.Int:
			value := vp.FieldByName(field.Name)
			valueIf := value.Interface().(int)
			flags.IntVar(&valueIf, kebabKey, int(value.Int()), usage)
			viper.RegisterAlias(field.Name, kebabKey)
			_ = viper.BindEnv(kebabKey, envKey)
		case reflect.Int64:
			value := vp.FieldByName(field.Name)
			valueIf := value.Interface().(int64)
			flags.Int64Var(&valueIf, kebabKey, value.Int(), usage)
			viper.RegisterAlias(field.Name, kebabKey)
			_ = viper.BindEnv(kebabKey, envKey)
		case reflect.Bool:
			value := vp.FieldByName(field.Name)
			valueIf := value.Interface().(bool)
			flags.BoolVar(&valueIf, kebabKey, value.Bool(), usage)
			viper.RegisterAlias(field.Name, kebabKey)
			_ = viper.BindEnv(kebabKey, envKey)
		case reflect.Slice:
			value := vp.FieldByName(field.Name)
			valueIf := value.Interface().([]string)
			flags.StringSliceVar(&valueIf, kebabKey, value.Interface().([]string), usage)
			viper.RegisterAlias(field.Name, kebabKey)
			_ = viper.BindEnv(kebabKey, envKey)
		case reflect.Struct:
			value := vp.FieldByName(field.Name)
			initConfig(flags, value.Interface(), envKey)
		case reflect.Pointer:
			panic("pointer not allowed in config field")
		}
	}
}

func UpdateConfig() {
	log.Infof("update config")
	vp := reflect.ValueOf(GetConfig()).Elem()
	updateConfig(vp, reflect.TypeOf(*GetConfig()), "")
	updateCommand()
}

var privateFieldRegex = regexp.MustCompile("^[a-z].+$")

func updateConfig(vp reflect.Value, structType reflect.Type, prefix string) {
	for _, field := range reflect.VisibleFields(structType) {
		fieldVp := vp.FieldByName(field.Name)
		if field.Type.Kind() == reflect.Struct {
			envKey := field.Tag.Get("env")
			updateConfig(fieldVp, field.Type, envKey)
		} else {
			if privateFieldRegex.Match([]byte(field.Name)) {
				continue
			}
			kebabKey := strcase.ToKebab(field.Name)
			// apply prefix
			if len(prefix) > 0 {
				kebabKey = fmt.Sprintf("%v-%v", strings.ToLower(prefix), kebabKey)
			}

			if viper.IsSet(kebabKey) {
				value := viper.Get(kebabKey)
				switch field.Type.Kind() {
				case reflect.Int:
					value = viper.GetInt(kebabKey)
				case reflect.Int64:
					value = viper.GetInt64(kebabKey)
				case reflect.Bool:
					value = viper.GetBool(kebabKey)
				case reflect.Slice:
					if field.Type.String() == "[]string" {
						value = TrimAndCreateSlices(StringToStringSlice(kebabKey))
					} else {
						value = viper.GetStringSlice(kebabKey)
					}
				}
				fieldVp.Set(reflect.ValueOf(value))
			}
		}
	}
}

func updateCommand() {
	// init command(find in 'PATH')
	goarch := runtime.GOARCH // os.Getenv("GOARCH")
	goos := runtime.GOOS     // os.Getenv("GOOS")
	goroot, _ := os.Getwd()
	log.Debugf("GOARCH - %v, GOOS - %v, current dir - %v", goarch, goos, goroot)
}

func whichCommand(command string) string {
	cmd := exec.Command("which", command)

	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("dcmtk command(%s) is not in path - %v", command, err))
	}
	return strings.TrimSpace(string(out))
}

func StringToStringSlice(key string) []string {
	in := viper.Get(key)
	if in == nil {
		return nil
	}
	switch reflect.TypeOf(in).String() {
	case "string":
		val, _ := ReadAsCSV(in.(string))
		return val
	case "[]string":
		return in.([]string)
	}
	return in.([]string)
}

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func StringToRuneSlice(s string) []rune {
	var r []rune
	for _, runeValue := range s {
		r = append(r, runeValue)
	}
	return r
}

func initLookupTable(row, column int) [][]bool {
	lookup := make([][]bool, row)
	for i := range lookup {
		lookup[i] = make([]bool, column)
	}
	return lookup
}

// Function that matches input str with given wildcard pattern
func WildcardPatternMatch(str, pattern string) bool {
	s := StringToRuneSlice(str)
	p := StringToRuneSlice(pattern)

	// empty pattern can only match with empty string
	if len(p) == 0 {
		return len(s) == 0
	}

	// lookup table for storing results of subproblems
	// zero value of lookup is false
	lookup := initLookupTable(len(s)+1, len(p)+1)

	// empty pattern can match with empty string
	lookup[0][0] = true

	// Only '*' can match with empty string
	for j := 1; j < len(p)+1; j++ {
		if p[j-1] == '*' {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	// fill the table in bottom-up fashion
	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(p)+1; j++ {
			if p[j-1] == '*' {
				// Two cases if we see a '*'
				// a) We ignore ‘*’ character and move
				//    to next  character in the pattern,
				//     i.e., ‘*’ indicates an empty sequence.
				// b) '*' character matches with ith
				//     character in input
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]

			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				// Current characters are considered as
				// matching in two cases
				// (a) current character of pattern is '?'
				// (b) characters actually match
				lookup[i][j] = lookup[i-1][j-1]

			} else {
				// If characters don't match
				lookup[i][j] = false
			}
		}
	}

	return lookup[len(s)][len(p)]
}

func TrimAndCreateSet(vs []string) map[string]bool {
	vsf := map[string]bool{}
	for _, v := range vs {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			vsf[v] = true
		}
	}
	return vsf
}

func TrimAndCreateSlices(vs []string) []string {
	var vsf []string
	for _, v := range vs {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
