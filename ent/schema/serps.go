package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// SERP holds the schema definition for the SERPs entity.
type SERP struct {
	ent.Schema
}

// Annotations of the User.
func (SERP) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "serps"},
	}
}

// Fields of the SERP.
func (SERP) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("description"),
		field.String("location"),
		field.JSON("contact_info", []string{}),
		field.JSON("key_words", []string{}),
		field.Bool("is_read").Default(false),
		field.Time("created_at").SchemaType(map[string]string{dialect.Postgres: "TIMESTAMP(0) WITH TIME ZONE"}).Default(time.Now),
	}
}

// Indexes of the SERP.
func (SERP) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("url", "contact_info").
			Unique(),
	}
}
