// Code generated by ent, DO NOT EDIT.

package searchquery

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the searchquery type in the database.
	Label = "search_query"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldQuery holds the string denoting the query field in the database.
	FieldQuery = "query"
	// FieldLocation holds the string denoting the location field in the database.
	FieldLocation = "location"
	// FieldLanguage holds the string denoting the language field in the database.
	FieldLanguage = "language"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// Table holds the table name of the searchquery in the database.
	Table = "search_query"
)

// Columns holds all SQL columns for searchquery fields.
var Columns = []string{
	FieldID,
	FieldQuery,
	FieldLocation,
	FieldLanguage,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// QueryValidator is a validator for the "query" field. It is called by the builders before save.
	QueryValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)

// OrderOption defines the ordering options for the SearchQuery queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByQuery orders the results by the query field.
func ByQuery(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldQuery, opts...).ToFunc()
}

// ByLocation orders the results by the location field.
func ByLocation(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocation, opts...).ToFunc()
}

// ByLanguage orders the results by the language field.
func ByLanguage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLanguage, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}
