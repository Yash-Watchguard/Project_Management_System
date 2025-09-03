package config

import (
	"testing"
)

func TestInsertQuery(t *testing.T) {
	table := "users"
	columns := []string{"id", "name", "email"}

	expected := "INSERT INTO users (id, name, email) VALUES (?, ?, ?)"
	result := InsertQuery(table, columns)

	if result != expected {
		t.Errorf("InsertQuery failed, expected %q, got %q", expected, result)
	}
}

func TestSelectQuery(t *testing.T) {
	table := "users"
	cols := []string{"id", "name", "email"}

	// Without conditions
	expected := "SELECT id, name, email FROM users"
	result := SelectQuery(table, cols)
	if result != expected {
		t.Errorf("SelectQuery without conditions failed, expected %q, got %q", expected, result)
	}

	// With one condition
	expected = "SELECT id, name, email FROM users WHERE id = ?"
	result = SelectQuery(table, cols, "id")
	if result != expected {
		t.Errorf("SelectQuery with one condition failed, expected %q, got %q", expected, result)
	}

	// With multiple conditions
	expected = "SELECT id, name, email FROM users WHERE id = ? AND name = ?"
	result = SelectQuery(table, cols, "id", "name")
	if result != expected {
		t.Errorf("SelectQuery with multiple conditions failed, expected %q, got %q", expected, result)
	}
}

func TestDeleteQuery(t *testing.T) {
	table := "users"

	// Without conditions
	expected := "DELETE FROM users"
	result := DeleteQuery(table, []string{})
	if result != expected {
		t.Errorf("DeleteQuery without conditions failed, expected %q, got %q", expected, result)
	}

	// With one condition
	expected = "DELETE FROM users WHERE id = ?"
	result = DeleteQuery(table, []string{"id"})
	if result != expected {
		t.Errorf("DeleteQuery with one condition failed, expected %q, got %q", expected, result)
	}

	// With multiple conditions
	expected = "DELETE FROM users WHERE id = ? AND name = ?"
	result = DeleteQuery(table, []string{"id", "name"})
	if result != expected {
		t.Errorf("DeleteQuery with multiple conditions failed, expected %q, got %q", expected, result)
	}
}

func TestUpdateQuery(t *testing.T) {
	table := "users"
	columns := []string{"name", "email"}

	// With one condition
	expected := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	result := UpdateQuery(table, "id", "", columns)
	if result != expected {
		t.Errorf("UpdateQuery with one condition failed, expected %q, got %q", expected, result)
	}

	// With two conditions
	expected = "UPDATE users SET name = ?, email = ? WHERE id = ? AND status = ?"
	result = UpdateQuery(table, "id", "status", columns)
	if result != expected {
		t.Errorf("UpdateQuery with two conditions failed, expected %q, got %q", expected, result)
	}
}
