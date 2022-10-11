package penv

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var Prefix = "ENV_"

type DumpVar struct {
	Name     string
	NameEnv  string
	NameFlag string
	Desc     string
	Type     string
	Value    string
	Secret   bool
}

func generate(name string) string {
	n := ""
	prevBig := true
	for i, char := range name {
		if i > 0 && string(char) == strings.ToUpper(string(char)) && !prevBig {
			n += "_"
		}
		if string(char) == strings.ToUpper(string(char)) {
			prevBig = true
		} else {
			prevBig = false
		}
		n += string(char)
	}
	return n
}

func generateEnvName(name string) string {
	return strings.ToUpper(Prefix + generate(name))
}

func generateFlagName(name string) string {
	return strings.ToLower(generate(name))
}

func stringToInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func stringToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func DumpConfig(config any) map[string]DumpVar {
	res := map[string]DumpVar{}

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		nameEnv := generateEnvName(t.Field(i).Name)
		fieldType := t.Field(i).Type.Kind().String()
		nameFlag := generateFlagName(t.Field(i).Name)
		description := t.Field(i).Tag.Get("description")
		if description == "" {
			description = "No description"
		}
		secret := t.Field(i).Tag.Get("secret")
		if fieldType == "string" {
			res[t.Field(i).Name] = DumpVar{
				Name:     t.Field(i).Name,
				NameEnv:  nameEnv,
				NameFlag: nameFlag,
				Desc:     description,
				Type:     cases.Title(language.English).String(fieldType),
				Value:    *v.Field(i).Addr().Interface().(*string),
				Secret:   secret == "1" || secret == "true",
			}
		} else if fieldType == "int" {
			res[t.Field(i).Name] = DumpVar{
				Name:     t.Field(i).Name,
				NameEnv:  nameEnv,
				NameFlag: nameFlag,
				Desc:     description,
				Type:     cases.Title(language.English).String(fieldType),
				Value:    fmt.Sprintf("%d", *v.Field(i).Addr().Interface().(*int)),
				Secret:   secret == "1" || secret == "true",
			}
		} else if fieldType == "int64" {
			res[t.Field(i).Name] = DumpVar{
				Name:     t.Field(i).Name,
				NameEnv:  nameEnv,
				NameFlag: nameFlag,
				Desc:     description,
				Type:     cases.Title(language.English).String(fieldType),
				Value:    fmt.Sprintf("%d", *v.Field(i).Addr().Interface().(*int64)),
				Secret:   secret == "1" || secret == "true",
			}
		}
	}

	return res
}

func ProcessConfig(config any) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	// Process flags
	for i := 0; i < t.NumField(); i++ {
		nameEnv := generateEnvName(t.Field(i).Name)
		nameFlag := generateFlagName(t.Field(i).Name)
		fieldType := t.Field(i).Type.Kind().String()
		defvalue := t.Field(i).Tag.Get("default")
		description := t.Field(i).Tag.Get("description")
		if description == "" {
			description = "No description"
		}

		if fieldType == "string" {
			value := v.Field(i).Addr().Interface().(*string)
			flag.StringVar(value, nameFlag, defvalue, "Or "+nameEnv+": "+description)
		} else if fieldType == "int" {
			if ndefvalue, err := stringToInt(defvalue); err == nil {
				value := v.Field(i).Addr().Interface().(*int)
				flag.IntVar(value, nameFlag, ndefvalue, "Or "+nameEnv+": "+description)
			} else {
				return err
			}
		} else if fieldType == "int64" {
			if ndefvalue, err := stringToInt64(defvalue); err == nil {
				value := v.Field(i).Addr().Interface().(*int64)
				flag.Int64Var(value, nameFlag, ndefvalue, "Or "+nameEnv+": "+description)
			} else {
				return err
			}
		}
	}
	flag.Parse()

	// Process ENVs
	for i := 0; i < t.NumField(); i++ {
		nameEnv := generateEnvName(t.Field(i).Name)
		fieldType := t.Field(i).Type.Kind().String()

		if os.Getenv(nameEnv) != "" {
			if fieldType == "string" {
				value := v.Field(i).Addr().Interface().(*string)
				*value = os.Getenv(nameEnv)
			} else if fieldType == "int" {
				if nvalue, err := stringToInt(os.Getenv(nameEnv)); err == nil {
					value := v.Field(i).Addr().Interface().(*int)
					*value = nvalue
				} else {
					return err
				}
			} else if fieldType == "int64" {
				if nvalue, err := stringToInt64(os.Getenv(nameEnv)); err == nil {
					value := v.Field(i).Addr().Interface().(*int64)
					*value = nvalue
				} else {
					return err
				}
			}
		}
	}

	return nil
}
