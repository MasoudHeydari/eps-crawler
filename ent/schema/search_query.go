package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// SearchQuery holds the schema definition for the SearchQueries entity.
type SearchQuery struct {
	ent.Schema
}

// Annotations of the SearchQuery.
func (SearchQuery) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "search_query"},
	}
}

// Fields of the SearchQuery.
func (SearchQuery) Fields() []ent.Field {
	return []ent.Field{
		field.String("query").NotEmpty(),
		field.String("location"),
		field.String("language"),
		field.Time("created_at").SchemaType(TimeStampWithTZ).Default(time.Now),
	}
}

// Indexes of the SearchQuery.
func (SearchQuery) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("query", "location", "language").
			Unique(),
	}
}
