package config

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Init the Env Setup
func initEnv(envName ...string) {
	var env string
	if len(envName) > 0 {
		env = envName[0]
	} else {
		env = ".env"
	}

	cwd, err := os.Getwd()

	if err != nil {
		log.Println("Error getting current working directory:", err)
	}

	LoadEnv(cwd + "/" + env)
	log.Println("EEEEEE", os.Getenv("APP_ENV"))
}

// LoadEnv loads environment variables from a .env file.
func LoadEnv(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        log.Println("Could not load env file: ", err)
		return err
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        // Ignore comments and empty lines
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }
        // Split the line into key and value
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        // Set the environment variable
        err := os.Setenv(key, value)
        if err != nil {
            return err
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

// GetEnv retrieves the value of an environment variable.
// If the variable is not set, it returns the provided default value (if any).
func GetEnv(key string, defaultValue ...interface{}) interface{} {
    value := os.Getenv(key)
    if value == "" {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return "" // Return empty string if no default value is provided
    }

    // If a default value is provided, determine its type and convert the env var to that type
    if len(defaultValue) > 0 {
        defaultValueType := reflect.TypeOf(defaultValue[0])
        switch defaultValueType.Kind() {
        case reflect.String:
            return value
        case reflect.Int:
            intValue, err := strconv.Atoi(value)
            if err != nil {
                log.Println("Error converting value to integer: ", err)
                return defaultValue[0]
            }
            return intValue
        case reflect.Bool:
            boolValue, err := strconv.ParseBool(value)
            if err != nil {
                log.Println("Error converting value to boolean: ", err)
                return defaultValue[0]
            }
            return boolValue
        default:
            return value // If default value is not a string or int, return the raw string value
        }
    }

    return value // If no default value is provided, return the raw string value
}
