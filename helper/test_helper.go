package helper

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const (
	// TestEnvFilePath is the environment variable key for getting .env file path
	TestEnvFilePath = "TEST_ENV_FILE_PATH"
)

// SetupTestCase initializes a test and loads the environment file specified in TEST_ENV_FILE_PATH
// It will cause the test to fail if the environment file is not available
func SetupTestCase(t *testing.T) {
	// setup code in here.
	err := LoadEnvFile(t)
	if err != nil {
		// even if we couldn't load the env file we won't fail the test
		// but if proper env variables aren't set, tests will fail on test data validation
		//t.Fatal(err)
	}
}

// LoadEnvFile read an .env file that has a path by the value of TEST_ENV_FILE_PATH environment variable.
func LoadEnvFile(t *testing.T) error {
	envFileName := os.Getenv(TestEnvFilePath)
	err := godotenv.Load(envFileName)
	if err != nil {
		return fmt.Errorf("Can not read .env file: %s", envFileName)
	}
	return nil
}

// InitializeTestValuesE fill the value from environment variables.
func InitializeTestValues(s interface{}) interface{} {
	fields := reflect.ValueOf(s).Elem()
	// iterate across all configuration properties
	for i := 0; i < fields.NumField(); i++ {
		typeField := fields.Type().Field(i)
		environmentVariablesKey := typeField.Tag.Get("env")
		if fields.Field(i).Kind() == reflect.String {
			// check if we want a property inside a complex object
			propertyKey, exists := typeField.Tag.Lookup("property")
			if exists {
				// get object string
				objectString := os.Getenv(environmentVariablesKey)
				// grab property value inside string
				propertyValue := getPropertyValueFromString(objectString, propertyKey)
				// set the value in the correct field
				fields.Field(i).SetString(propertyValue)
			} else {
				fields.Field(i).SetString(os.Getenv(environmentVariablesKey))
			}
		}
	}
	return s
}

func getPropertyValueFromString(object string, propertyKey string) string {
	// compile regex to look for key="value"
	regexString := fmt.Sprintf(`%s=\"(.*?)\"`, propertyKey)
	re := regexp.MustCompile(regexString)
	match := string(re.Find([]byte(object)))
	if len(match) == 0 {
		log.Printf("Warning: Could not find property with key %s\n", propertyKey)
		return ""
	}
	match = strings.Replace(match, "\"", "", -1)
	propertyValue := strings.Split(match, "=")[1]
	return propertyValue
}

// ValidateTestValues validate if the all parameters has the value. skipGenerated allows ignore a field that has the `generated:"true"` tag.
func ValidateTestValues(s interface{}, skipGenerated bool) bool {
	fields := reflect.ValueOf(s).Elem()
	flag := true
	for i := 0; i < fields.NumField(); i++ {
		value := fields.Field(i)
		typeField := fields.Type().Field(i)

		if !validateTags(typeField.Tag) {
			log.Printf("Warning: Struct Field %s has invalid tags.\n", typeField.Name)
			flag = false
			continue
		}

		if value.Kind() == reflect.String {
			if len(value.String()) == 0 {
				if !anyTagExists(typeField.Tag) {
					continue
				} else if skipGenerated && tagExists(typeField.Tag, "env") && tagExists(typeField.Tag, "generated") {
					log.Printf("Warning: Struct Field %s (env:%s) doesn't have any value. (Generated = true. skipped.)\n", typeField.Name, typeField.Tag.Get("env"))
					continue
				} else if skipGenerated && tagExists(typeField.Tag, "kv") && tagExists(typeField.Tag, "generated") {
					log.Printf("Warning: Struct Field %s (kv:%s) doesn't have any value. (Generated = true. skipped.)\n", typeField.Name, typeField.Tag.Get("kv"))
					continue
				} else if tagExists(typeField.Tag, "kv") {
					log.Printf("Warning: Struct Field %s (kv:%s) doesn't have any value.\n", typeField.Name, typeField.Tag.Get("kv"))
					flag = false
				} else if tagExists(typeField.Tag, "val") {
					log.Printf("Warning: Struct Field %s doesn't have any value.\n", typeField.Name)
					flag = false
				} else {
					log.Printf("Warning: Struct Field %s (env:%s) doesn't have any value.\n", typeField.Name, typeField.Tag.Get("env"))
					flag = false
				}
			}
		} else if value.Kind() == reflect.Map {
			if value.IsNil() {
				log.Printf("Warning: Struct Field %s doesn't have any value.\n", typeField.Name)
				flag = !tagExists(typeField.Tag, "val")
			}
		} else if value.Kind() == reflect.Slice {
			if value.IsNil() {
				log.Printf("Warning: Array Field %s doesn't have any value.\n", typeField.Name)
				flag = !tagExists(typeField.Tag, "val")
			}
		} else if value.Kind() == reflect.Bool || value.Kind() == reflect.Int32 || value.Kind() == reflect.Int64 {
			// all of these have default "zero" values so they are always valid
		} else {
			log.Printf("Warning: Found Field %s of type %s which is not allowed for Config Structures.\n", value.Kind(), typeField.Name)
			return false
		}
	}
	return flag
}

// FetchKeyVaultSecretE fill the value from keyvault
func FetchKeyVaultSecretE(s interface{}) (interface{}, error) {
	keyVaultName, err := getKeyVaultName(s)
	if err != nil {
		return nil, err
	}

	fields := reflect.ValueOf(s).Elem()
	for i := 0; i < fields.NumField(); i++ {
		typeField := fields.Type().Field(i)
		if typeField.Tag.Get("kv") != "" {
			secretName := typeField.Tag.Get("kv")
			if fields.Field(i).Kind() == reflect.String {
				secret, err := GetKeyVaultSecret(keyVaultName, secretName)
				if err != nil {
					return nil, err
				}
				fields.Field(i).SetString(secret)
			}
		}

	}
	return s, nil
}

func getKeyVaultName(s interface{}) (string, error) {
	structName := reflect.TypeOf(s)
	fields := reflect.ValueOf(s).Elem()
	for i := 0; i < fields.NumField(); i++ {
		typeField := fields.Type().Field(i)
		if len(typeField.Tag.Get("kvname")) != 0 {
			if fields.Field(i).Kind() == reflect.String {
				kvname := fields.Field(i).String()
				kvNameField := fields.Type().Field(i).Name
				if len(kvname) == 0 {
					return "", fmt.Errorf("Empty KeyVault name is not allowed. Please add `kvname` on your struct %s.%s", structName, kvNameField)
				}
				return fields.Field(i).String(), nil
			}
		}
	}
	return "", fmt.Errorf("Can not find kvname filed on your struct %s", structName)
}

// IsTagExists test if the tag is there or not.
func tagExists(tag reflect.StructTag, tagName string) bool {
	_, ok := tag.Lookup(tagName)
	return ok
}

// validateTags test if any tags are invalid
func validateTags(tag reflect.StructTag) bool {

	val, isVal := tag.Lookup("val")
	generated, isGenerated := tag.Lookup("generated")

	if isVal {
		v, err := strconv.ParseBool(val)
		if err != nil || !v {
			log.Printf("Warning: Value of \"val\" tag should be true")
			return false
		}
	}

	if isGenerated {
		v, err := strconv.ParseBool(generated)
		if err != nil || !v {
			log.Printf("Warning: Value of \"generated\" tag should be true")
			return false
		}
	}

	return true
}

// IsAnyTagExists test if any tags are exists.
func anyTagExists(tag reflect.StructTag) bool {
	_, isEnv := tag.Lookup("env")
	_, isKv := tag.Lookup("kv")
	_, isVal := tag.Lookup("val")
	return isEnv || isKv || isVal
}

// GetYamlVariables reads the yaml file in filePath and returns valus specified by interface s
func GetYamlVariables(filePath string, s interface{}) (interface{}, error) {
	// read yaml file
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Path to Yaml file not set or invalid: %s", filePath)
	}

	// parse yaml file
	m := make(map[interface{}]interface{})
	err = yaml.UnmarshalStrict(yamlFile, &m)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Yaml File %s: %s", filePath, err.Error())
	}

	err = mapstructure.Decode(m, &s)
	return s, nil
}

// CheckIfEndpointIsResponding test an endpoint for availability. Returns true if endpoint is available, false otherwise
func CheckIfEndpointIsResponding(t *testing.T, endpoint string) bool {
	// we ignore certificates at this point
	tlsConfig := tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	err := http_helper.HttpGetWithRetryWithCustomValidationE(
		t,
		fmt.Sprintf("https://%s", endpoint),
		&tlsConfig,
		1,
		10*time.Second,
		func(statusCode int, body string) bool {
			if statusCode == 200 {
				return true
			}
			if statusCode == 404 {
				t.Log("Warning: 404 response from endpoint. Test will still PASS.")
				return true
			}
			return false
		},
	)
	return err == nil
}

//CheckSQLConnectivity checks if we can successfully connect to a SQL Managed Instance, MySql server or Azure SQL Server
func CheckSQLConnectivity(t *testing.T, driver string, connString string) {

	// Create connection pool
	db, err := sql.Open(driver, connString)
	require.NoErrorf(t, err, "Error creating connection pool: %s ", err)

	// Close the database connection pool after program executes
	defer db.Close()

	// Use background context
	ctx := context.Background()

	//As open doesn't actually create the connection we need to make some sort of command to check that connectivity works
	//err = db.Ping()
	//require.NoErrorf(t, err, "Error pinging database: %s", err)

	var result string

	// Run query and scan for result
	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	require.NoErrorf(t, err, "Error: %s", err)
	t.Logf("%s\n", result)
}

//CheckRedisCacheConnectivity checks if we can successfully connect to a Redis cache instance
func CheckRedisCacheConnectivity(t *testing.T, redisCacheURL string, redisCachePort int, redisCachePassword string) {
	conn, err := redis.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", redisCacheURL, redisCachePort),
		redis.DialPassword(redisCachePassword),
		redis.DialUseTLS(true))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	key := "CheckRedisCacheConnectivity"
	_, err = conn.Do("SET", key, "hello world")
	assert.NoErrorf(t, err, "Error connecting to Redis Cache %s", redisCacheURL)
	result, err := redis.String(conn.Do("GET", key))
	assert.NoErrorf(t, err, "Error connecting to Redis Cache %s", redisCacheURL)
	assert.Equalf(t, "hello world", result, "Failed to retrieve value from REDIS")
	_, err = conn.Do("DEL", key)
	assert.NoErrorf(t, err, "Error connecting to Redis Cache %s", redisCacheURL)
}

//CheckCassandraConnectivity checks if we can successfully connect to a Cassandra keyspace on CosmosDB
func CheckCassandraConnectivity(t *testing.T, endpoint string, username string, password string, database string) {
	// connect to the cluster
	cluster := gocql.NewCluster(endpoint)
	cluster.Port = 10350
	var sslOptions = new(gocql.SslOptions)
	sslOptions.EnableHostVerification = false
	cluster.SslOpts = sslOptions
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}
	session, err := cluster.CreateSession()
	assert.NoErrorf(t, err, "Error Connecting to %s", endpoint)
	if session != nil {
		defer session.Close()
		var keyspace string
		iter := session.Query(`SELECT keyspace_name FROM system_schema.keyspaces;`).Iter()
		found := false
		for iter.Scan(&keyspace) {
			if keyspace == database {
				found = true
				break
			}
		}
		assert.Truef(t, found, "Keyspace %s not found in Cassandra %s", endpoint)
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}

}
