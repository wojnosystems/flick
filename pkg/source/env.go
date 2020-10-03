package source

import (
	"github.com/wojnosystems/flick/pkg/set_value"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type envReader interface {
	// Get the value of a single environment with the name envNamed
	Get(envNamed string) string
	// Keys get a list of keys that begin with the prefix. If "" is passed, matches all and returns all keys
	Keys(prefix string) []string
}

// Implements the default system environment getters
// this allows us to test this
type envSystem struct {
}

func (s *envSystem) Get(envNamed string) string {
	return os.Getenv(envNamed)
}

func (s *envSystem) Keys(prefix string) (out []string) {
	for _, key := range os.Environ() {
		if strings.HasPrefix(key, prefix) {
			out = append(out, key)
		}
	}
	return
}

// Env reads environment variables and Unmarshalls them into the object provided
// Because of the limitations of environment variables from the shell's perspective,
// we're limited to the following character set for valid environment variable names:
// [a-zA-Z_]+[a-zA-Z0-9_]*
// due to shell restrictions. POSIX says it doesn't care what the name is, as long as it doesn't contain an equal sign:
// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html#tag_08 so we're limited by the shell
// restrictions.
// Thus, the naming conventions for environment variables being mapped to members in structures is structured:
//
// Example variable naming scheme:
//
// type myStruct struct {
//   Name      string     "env:name"
//   PetNames  []string   "env:pet_names"
//   Addresses addrStruct // no tag, assumes "Addresses"
// }
//
// type addrStruct struct {
//   Street string     "env:street"
// }
//
// The following environment variables can create the following data structure:
//
// ```bash
// name=Chris pet_names_0_=Fluffy pet_names_4_=Foxy Address_0_street="742 Evergreen Terrace" \
//   Address_1_street="2001 Creaking Oak Drive" ./my-app
// ```
//
// s := myStruct{
//   Name: "chris",
//   PetNames: []string{
//     "Fluffy",
//     "",
//     "", // Foxy was index 4, so we created 3 blank strings in between as this is an array
//     "",
//     "Foxy",
//   },
//   Addresses: []addrStruct{
//     {
//       Street: "742 Evergreen Terrace",
//     },
//     {
//       Street: "2001 Creaking Oak Drive",
//     },
//   }
// }
//
// Because the structure enforces the name and prevents duplicate names from being used, this scheme is guaranteed
// to produce unique names for structure members and child members as long as you don't use tags to create collisions.
//
type Env struct {
	envs          envReader
	parseRegistry *set_value.Registry
}

func (e *Env) Unmarshall(into interface{}) (err error) {
	root := newTypePair(into)
	err = e.envValidateInto(&root)
	if err != nil {
		return
	}
	err = e.envUnmarshall("", typePair{
		v: root.v.Elem(),
		t: root.v.Elem().Type(),
	})
	if err != nil {
		return
	}
	return nil
}

type typePair struct {
	v reflect.Value
	t reflect.Type
}

func newTypePair(from interface{}) typePair {
	return typePair{
		v: reflect.ValueOf(from),
		t: reflect.TypeOf(from),
	}
}

func newTypePairFromValue(from reflect.Value) typePair {
	return typePair{
		v: from,
		t: from.Type(),
	}
}

func (e *Env) envUnmarshall(parentName string, structRef typePair) (err error) {
	for i := 0; i < structRef.v.NumField(); i++ {
		fieldV := structRef.v.Field(i)
		if fieldV.CanSet() {
			fieldT := structRef.t.Field(i)
			fieldName := envFieldNameOrDefault(fieldT)
			fullPath := envNameFromParent(fieldName, parentName)

			if fieldT.Type.Kind() == reflect.Slice {
				err = e.envSliceUnmarshall(fullPath, fieldV)
				if err != nil {
					return
				}
			} else {
				envValue := e.envs.Get(fullPath)
				if "" != envValue {
					var wasCalled bool
					wasCalled, err = e.parseRegistry.SetValue(fieldV.Addr().Interface(), envValue)
					if err != nil {
						return
					}
					if !wasCalled {
						// fall back
						if fieldT.Type.Kind() == reflect.Struct {
							err = e.envUnmarshall(fullPath, typePair{
								v: fieldV,
								t: fieldV.Type(),
							})
							if err != nil {
								return
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func envFieldNameOrDefault(fieldT reflect.StructField) (fieldName string) {
	fieldName = fieldT.Tag.Get("env")
	if "" == fieldName {
		fieldName = fieldT.Name
	}
	return
}

func envNameFromParent(name string, parent string) string {
	if parent != "" {
		return parent + "." + name
	}
	return name
}

func (e *Env) envValidateInto(root *typePair) (err error) {
	if root.v.IsNil() {
		return NewErrProgramming("'into' argument must be not be nil")
	}
	if root.t.Kind() != reflect.Ptr {
		return NewErrProgramming("'into' argument must be a reference")
	}
	if root.v.Elem().Kind() != reflect.Struct {
		return NewErrProgramming("'into' argument must be a struct")
	}
	return nil
}

func (e *Env) envSliceUnmarshall(path string, sliceValue reflect.Value) (err error) {
	var length int
	length, err = envEnvKeyLength(e.envs, path+"_")
	if err != nil {
		return
	}
	if length > 0 {
		newSlice := reflect.MakeSlice(sliceValue.Type(), length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			err = e.envUnmarshall(path+"_"+strconv.FormatInt(int64(i), 10)+"_", typePair{
				v: sliceElement,
				t: sliceElement.Type(),
			})
			if err != nil {
				return
			}
		}
	}
	return
}

func envEnvKeyLength(env envReader, pathPrefix string) (length int, err error) {
	maxIndex := int64(-1)
	for _, key := range env.Keys(pathPrefix) {
		possibleNumber := envIndexRegexp.FindString(key[len(pathPrefix):])
		if "" != possibleNumber {
			var index int64
			index, err = strconv.ParseInt(possibleNumber, 10, 0)
			if err != nil {
				return
			}
			if index > maxIndex {
				maxIndex = index
			}
		}
	}
	length = int(maxIndex + 1)
	return
}

var envIndexRegexp = regexp.MustCompile(`^(\d+)`)
