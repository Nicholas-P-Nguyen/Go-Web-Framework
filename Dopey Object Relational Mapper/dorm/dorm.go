package dorm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"
)

// DB handle connection pool to SQL database
type DB struct {
	inner *sql.DB
}

// NewDB returns a new DB using the provided `conn`,
// an sql database connection.
// This function is provided for you. You DO NOT need to modify it.
func NewDB(conn *sql.DB) DB {

	return DB{inner: conn}
}

// Close closes db's database connection.
// This function is provided for you. You DO NOT need to modify it.
func (db *DB) Close() error {
	return db.inner.Close()
}

func CamelToUnderscore(field string) string {
	var components []string

	for i := 0; i < len(field); {
		currChar := rune(field[i])
		var j int
		if unicode.IsUpper(currChar) {
			j = i + 1
			// upper to upper: COSF iles
			for j < len(field) && unicode.IsUpper(rune(field[j])) {
				if j+1 < len(field) && unicode.IsLower(rune(field[j+1])) {
					break
				}
				j++
			}
			// upper to lower: Camel Case
			for j < len(field) && unicode.IsLower(rune(field[j])) {
				j++
			}
			components = append(components, field[i:j])
			i = j
		} else {
			// If not uppercase, move to the next character
			j = i + 1
			for j < len(field) && unicode.IsLower(rune(field[j])) {
				j++
			}
			components = append(components, field[i:j])
			i = j
		}
	}

	var builder strings.Builder
	for i, component := range components {
		if i > 0 {
			builder.WriteRune('_')
		}
		builder.WriteString(strings.ToLower(component))
	}
	return builder.String()
}

// UnderscoreToCamel   camel_case -> CamelCase
func UnderscoreToCamel(field string) string {
	var builder strings.Builder
	capitalizeNext := true

	for _, char := range field {
		if char == '_' {
			capitalizeNext = true
			continue
		}

		if capitalizeNext {
			builder.WriteRune(unicode.ToUpper(char))
			capitalizeNext = false
		} else {
			builder.WriteRune(unicode.ToLower(char))
		}
	}

	return builder.String()
}

// ColumnNames analyzes a struct, v, and returns a list of strings,
// *[]User -> []User -> User
// *User{} -> User{}
func ColumnNames(v interface{}) []string {
	structType := reflect.TypeOf(v).Elem()
	if structType.Kind() == reflect.Slice {
		structType = structType.Elem()
	}

	var columnNames []string
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath == "" {
			convertedField := CamelToUnderscore(field.Name)
			columnNames = append(columnNames, convertedField)
		}
	}
	return columnNames
}

// TableName analyzes a struct, v, and returns a single string, equal
// to the name of that struct's type, converted to underscore_case.
func TableName(result interface{}) string {
	structType := reflect.TypeOf(result).Elem()
	if structType.Kind() == reflect.Slice {
		structType = structType.Elem()
	}
	return CamelToUnderscore(structType.Name())
}

// Find queries a database for all rows in a given table,
// and stores all matching rows in the slice provided as an argument.
// The argument result will be a pointer to an empty slice of models.
// Find should panic if the table doesn't exist.
func (db *DB) Find(result interface{}) {
	resultSlice := reflect.ValueOf(result)
	tableName := TableName(result)
	queryStatement := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.inner.Query(queryStatement)
	if err != nil {
		log.Printf("Table names: %s", tableName)
		log.Panicf("An error has occurred: %s", err)
	}
	defer rows.Close()

	resultValue := reflect.ValueOf(result).Elem().Type().Elem()

	for rows.Next() {
		tempStruct := reflect.New(resultValue).Elem()
		fieldPointers := make([]interface{}, resultValue.NumField())

		for i := 0; i < resultValue.NumField(); i++ {
			fieldPointers[i] = tempStruct.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fieldPointers...); err != nil {
			log.Panic(err)
		}
		resultSlice.Elem().Set(reflect.Append(resultSlice.Elem(), tempStruct))
	}
	if err := rows.Err(); err != nil {
		log.Panic(err)
	}
}

// First queries a database for the first row in a table,
// and stores the matching row in the struct provided as an argument.
// If no such entry exists, First returns false; else it returns true.
// First should panic if the table doesn't exist.
func (db *DB) First(result interface{}) bool {
	tableName := TableName(result)
	queryStatement := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	row := db.inner.QueryRow(queryStatement)

	resultValue := reflect.ValueOf(result).Elem()
	resultType := reflect.TypeOf(result).Elem()
	fieldPointers := make([]interface{}, resultType.NumField())
	for i := 0; i < resultType.NumField(); i++ {
		fieldPointers[i] = resultValue.Field(i).Addr().Interface()
	}

	err := row.Scan(fieldPointers...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		log.Panic(err)
	}
	return true
}

// Create adds the specified model to the appropriate database table.
// The table for the model *must* already exist, and Create() should
// panic if it does not.
func (db *DB) Create(model interface{}) {
	tableName := TableName(model)
	// *User{} -> User{}
	modelValue := reflect.ValueOf(model).Elem()
	modelType := reflect.TypeOf(model).Elem()

	var colNames []string
	var placeHolder []string
	var values []interface{}
	var primaryField reflect.Value
	var hasPrimaryKey bool

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldVal := modelValue.Field(i)

		if tag := field.Tag.Get("dorm"); tag == "primary_key" {
			primaryField = fieldVal
			hasPrimaryKey = true
			continue
		}
		colNames = append(colNames, UnderscoreToCamel(field.Name))
		values = append(values, fieldVal.Interface())
		placeHolder = append(placeHolder, "?")
	}
	queryStatement := fmt.Sprintf(
		"INSERT INTRO %s (%s) VALUES (%s",
		tableName,
		strings.Join(colNames, ", "),
		strings.Join(placeHolder, ", "),
	)

	result, err := db.inner.Exec(queryStatement, values...)
	if err != nil {
		log.Panic(err)
	}

	if hasPrimaryKey {
		id, err := result.LastInsertId()
		if err != nil {
			log.Panic(err)
		}
		primaryField.SetInt(id)
	}
}
