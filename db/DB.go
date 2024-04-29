package db

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/ent"
	"github.com/karust/openserp/ent/serp"
	_ "github.com/lib/pq"
)

func NewDB() (*ent.Client, error) {
	client, err := ent.Open(dialect.Postgres, "host=127.0.0.1 port=5432 user=eps_user dbname=eps_db password=eps_password sslmode=disable")
	if err != nil {
		return nil, err
	}
	err = client.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func InsertBulk(ctx context.Context, db *ent.Client, results []core.SearchResult) error {
	b := make([]*ent.SERPCreate, 0, len(results))
	for _, result := range results {
		b = append(b,
			db.SERP.Create().
				SetTitle(result.Title).
				SetDescription(result.Description).
				SetURL(result.URL).
				SetLocation("Germany").
				SetKeyWords([]string{"N/A"}).
				SetContactInfo([]string{"N/A"}),
		)
	}
	_, err := db.SERP.CreateBulk(b...).Save(ctx)
	return err
}

func GetAllResult(ctx context.Context, db *ent.Client, location string) ([]SERP, error) {
	entSERPs, err := db.SERP.Query().
		Where(
			serp.IsRead(false),
			serp.Location(location),
		).All(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]SERP, 0, len(entSERPs))
	for _, entSERP := range entSERPs {
		results = append(results,
			SERP{
				URL:         entSERP.URL,
				Title:       entSERP.Title,
				Description: entSERP.Description,
				Location:    entSERP.Location,
				ContactInfo: entSERP.ContactInfo,
				Keywords:    entSERP.KeyWords,
				IsRead:      entSERP.IsRead,
				CreatedAt:   entSERP.CreatedAt,
			},
		)
	}
	return results, nil
}
