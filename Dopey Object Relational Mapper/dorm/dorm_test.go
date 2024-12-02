package dorm

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func connectSQL() *sql.DB {
	conn, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	if err != nil {
		panic(err)
	}
	return conn
}

func createUserTable(conn *sql.DB) {
	_, err := conn.Exec(`CREATE TABLE user (
		first_name TEXT,
		last_name TEXT
	)`)

	if err != nil {
		panic(err)
	}
}

func insertUsers(conn *sql.DB, users []User) {
	for _, uc := range users {
		_, err := conn.Exec(`INSERT INTO user (first_name, last_name)
		VALUES (?, ?)`, uc.FirstName, uc.LastName)

		if err != nil {
			panic(err)
		}
	}
}

type User struct {
	FirstName string
	LastName  string
}

var MockUsers = []User{
	{FirstName: "Test", LastName: "User1"},
	{FirstName: "Nicholas", LastName: "Nguyen"},
	{FirstName: "Rocky", LastName: "Nguyen"},
}

//func TestColumnNames(t *testing.T) {
//	conn := connectSQL()
//	createUserTable(conn)
//	insertUsers(conn, MockUsers)
//
//	db := NewDB(conn)
//	defer db.Close()
//
//	colNames := ColumnNames(&User{})
//	expected := []string{"first_name", "last_name"}
//
//	for i, names := range expected {
//		if colNames[i] != names {
//			t.Errorf("Expected: %s; Output: %s", names, colNames[i])
//		}
//	}
//}

func TestFind(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	results := []User{}
	db.Find(&results)

	for i, user := range results {
		fmt.Print("The results after calling Find():\n")
		fmt.Printf("{FirstName: %s, LastName: %s}\n", user.FirstName, user.LastName)
		fmt.Print("The expected results:\n")
		fmt.Printf("{FirstName: %s, LastName: %s}\n", MockUsers[i].FirstName, MockUsers[i].LastName)
	}
}

func TestFirst(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	result := &User{}

	_ = db.First(result)
	fmt.Printf("{FirstName: %s, LastName: %s}\n", result.FirstName, result.LastName)
	fmt.Printf("{FirstName: %s, LastName: %s}\n", MockUsers[0].FirstName, MockUsers[0].LastName)
}

func TestCamelToUnderscore(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel_case"},
		{"EMail", "e_mail"},
		{"COSFiles", "cos_files"},
		{"camelCase", "camel_case"},
		{"OldCOSFiles", "old_cos_files"},
		{"COSFilesX", "cos_files_x"},
	}

	for _, test := range tests {
		output := CamelToUnderscore(test.input)
		if output != test.expected {
			t.Errorf("input: %s; output: %s; expected: %s", test.input, output, test.expected)
		}
	}
}

func TestUnderscoreToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"camel_case", "CamelCase"},
		{"e_mail", "EMail"},
		{"cos_files", "CosFiles"},
		{"old_cos_files", "OldCosFiles"},
		{"cos_files_x", "CosFilesX"},
	}

	for _, test := range tests {
		output := UnderscoreToCamel(test.input)
		if output != test.expected {
			t.Errorf("input: %s; output: %s; expected: %s", test.input, output, test.expected)
		}
	}
}
