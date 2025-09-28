package util

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// binds a JSON configuration file to the provided destination struct.
func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Unmarshal the config into the provided struct
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal config: %v", err)
		return err
	}
	return nil
}

// sets environment variables from a Consul KV store using the provided Viper instance.
func SetEnvFromConsulKV(v *viper.Viper) error {
	val := make(map[string]any)

	err := v.Unmarshal(&val)
	if err != nil {
		logrus.Errorf("failed to unmarshal config from consul: %v", err)
		return err
	}

	// iterate over the map to convert all values to string and set them as environment variables
	for key, value := range val {
		var (
			valOf = reflect.ValueOf(value)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32, reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			panic("unsupported type")
		}

		// set the environment variable
		err = os.Setenv(key, val)
		if err != nil {
			logrus.Errorf("failed to set env from consul kv: %v", err)
			return err
		}
	}

	return nil
}

// binds configuration from a Consul KV store to the provided destination struct.
func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote config: %v", err)
		return err
	}
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal config: %v", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to set env from consul kv: %v", err)
		return err
	}

	return nil
}
