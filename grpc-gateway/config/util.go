package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
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
