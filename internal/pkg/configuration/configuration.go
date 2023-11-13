// TODO
package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Configuration struct {
	Mode              string `json:"mode" env:"REMEMBRANCE_MODE" cmd:"mode"`
	KeysDirectory     string `json:"keysDirectory" env:"REMEMBRANCE_KEYS_DIRECTORY" cmd:"keysDirectory"`
	RunDirectory      string `json:"runDirectory" env:"REMEMBRANCE_RUN_DIRECTORY" cmd:"runDirectory"`
	NodeIdUuidVersion int    `json:"nodeIdUuidVersion" env:"REMEMBRANCE_NODE_ID_UUID_VERSION" cmd:"nodeIdUuidVersion"`
	NodeIdUuidString  string `json:"nodeIdUuidString" env:"REMEMBRANCE_NODE_ID_UUID_STRING" cmd:"nodeIdUuidString"`
}

const defaultConfigurationFilePath = "./config.json"

var (
	currentConfiguration *Configuration
	loaded               bool = false
	guard                sync.RWMutex
	flagsMap             map[string]*string
)

func GetConfiguration(forceReload bool, configurationFilePath string) (Configuration, error) {
	guard.RLock()

	if !loaded || forceReload {
		guard.RUnlock()
		guard.Lock()
		defer guard.Unlock()

		if !loaded || forceReload {
			temp := Configuration{}
			err := loadConfiguration(&temp, configurationFilePath)
			if err != nil {
				return *currentConfiguration, err
			}

			currentConfiguration = &temp
			loaded = true
		}
	} else {
		defer guard.RUnlock()
	}

	return *currentConfiguration, nil
}

func applyEnvironmentVariableValue(field *reflect.StructField, value *reflect.Value) error {
	variableName, present := field.Tag.Lookup("env")
	if !present {
		variableName = "REMEMBRANCE_" + strings.ToUpper(field.Name)
	}

	variableValue, present := os.LookupEnv(variableName)
	if present {
		fieldValue := value.Elem().FieldByName(field.Name)

		actualValue, err := getValue(field, variableValue)
		if err != nil {
			return err
		}

		fieldValue.Set(*actualValue)
	}

	return nil
}

func applyConfigurationFileValue(field *reflect.StructField, value *reflect.Value, jsonConfiguration map[string]interface{}) error {

	attributeName, present := field.Tag.Lookup("json")
	if !present {
		attributeName = strings.ToLower(field.Name)
	}

	fieldValue := value.Elem().FieldByName(field.Name)

	attributeValue, present := jsonConfiguration[attributeName]
	if present && attributeValue != nil {
		switch attributeType := attributeValue.(type) {
		case string:
			if attributeType != "" {
				fieldValue.Set(reflect.ValueOf(attributeType))
			}

		default:
			s := fmt.Sprint(attributeValue)
			actualValue, err := getValue(field, s)
			if err != nil {
				return err
			}

			fieldValue.Set(*actualValue)
		}
	}

	return nil
}

func prepareFlags(pt *reflect.Type) {
	const blankValue = "blank"

	flagsMap = make(map[string]*string)

	t := *pt
	nFields := t.NumField()
	for i := 0; i < nFields; i++ {
		field := t.Field(i)

		argumentName, present := field.Tag.Lookup("cmd")
		if !present {
			argumentName = strings.ToLower(field.Name)
		}

		usage := fmt.Sprintf("Usage: %s", argumentName)
		argumentValue := flag.String(argumentName, blankValue, usage)
		flagsMap[argumentName] = argumentValue
	}

	if !flag.Parsed() {
		flag.Parse()
	}
}

func applyCommandLineValue(field *reflect.StructField, value *reflect.Value) error {
	const blankValue = "blank"

	argumentName, present := field.Tag.Lookup("cmd")
	if !present {
		argumentName = strings.ToLower(field.Name)
	}

	argumentValue := flagsMap[argumentName]

	if (*argumentValue) != blankValue {
		actualValue, err := getValue(field, *argumentValue)
		if err != nil {
			return err
		}

		fieldValue := value.Elem().FieldByName(field.Name)
		fieldValue.Set(*actualValue)
	}

	return nil
}

func applyDefaultValues(configuration *Configuration) {
	configuration.Mode = ""
	configuration.KeysDirectory = ""
	configuration.RunDirectory = "./run"
	configuration.NodeIdUuidVersion = 4
}

func loadConfiguration(configuration *Configuration, configurationFilePath string) error {
	applyDefaultValues(configuration)

	var actualConfigurationFilePath string
	if configurationFilePath == "" {
		actualConfigurationFilePath = defaultConfigurationFilePath

	} else {
		actualConfigurationFilePath = configurationFilePath
	}

	file, err := os.ReadFile(actualConfigurationFilePath)
	if err != nil {
		return err
	}

	var rawConfiguration interface{}
	marshallError := json.Unmarshal(file, &rawConfiguration)
	if marshallError != nil {
		return marshallError
	}

	jsonConfiguration := rawConfiguration.(map[string]interface{})

	t := reflect.TypeOf(*configuration)
	v := reflect.ValueOf(configuration)

	prepareFlags(&t)

	nFields := t.NumField()
	for i := 0; i < nFields; i++ {
		field := t.Field(i)

		if err := applyConfigurationFileValue(&field, &v, jsonConfiguration); err != nil {
			return fmt.Errorf("error while applying value to field '%s': %s", field.Name, err.Error())
		}

		if err := applyEnvironmentVariableValue(&field, &v); err != nil {
			return fmt.Errorf("error while applying value to field '%s': %s", field.Name, err.Error())
		}

		if err := applyCommandLineValue(&field, &v); err != nil {
			return fmt.Errorf("error while applying value to field '%s': %s", field.Name, err.Error())
		}
	}

	return nil
}

func getValue(field *reflect.StructField, value string) (*reflect.Value, error) {
	switch field.Type.Kind() {
	case reflect.Float32:
		actual, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(float32(actual))
		return &result, nil

	case reflect.Float64:
		actual, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(actual)
		return &result, nil

	case reflect.Int:
		actual, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(int(actual))
		return &result, nil

	case reflect.Int8:
		actual, err := strconv.ParseInt(value, 0, 8)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(int8(actual))
		return &result, nil

	case reflect.Int16:
		actual, err := strconv.ParseInt(value, 0, 16)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(int16(actual))
		return &result, nil

	case reflect.Int32:
		actual, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(int32(actual))
		return &result, nil

	case reflect.Int64:
		actual, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(int64(actual))
		return &result, nil

	case reflect.Uint:
		actual, err := strconv.ParseUint(value, 0, 0)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(uint(actual))
		return &result, nil

	case reflect.Uint8:
		actual, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(uint8(actual))
		return &result, nil

	case reflect.Uint16:
		actual, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(uint16(actual))
		return &result, nil

	case reflect.Uint32:
		actual, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(uint32(actual))
		return &result, nil

	case reflect.Uint64:
		actual, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(uint64(actual))
		return &result, nil

	case reflect.Bool:
		actual, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}

		result := reflect.ValueOf(actual)
		return &result, nil

	case reflect.String:
		result := reflect.ValueOf(value)
		return &result, nil
	}

	return nil, fmt.Errorf("unexpected field type: %s", field.Type.String())
}
