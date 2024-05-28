package db

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	"errors"
	"fmt"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/ent"
	"github.com/karust/openserp/ent/searchquery"
	"github.com/karust/openserp/ent/serp"
	_ "github.com/lib/pq"
	"os"
	"time"
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

func InsertBulk(ctx context.Context, db *ent.Client, results []core.SearchResult, loc, lang, searchQ string) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting a transaction: %w", err)
	}
	var sqEnt *ent.SearchQuery
	sqCreate := db.SearchQuery.Create().
		SetLocation(loc).
		SetLanguage(lang).
		SetQuery(searchQ)
	sqCreate.OnConflict().
		DoNothing()
	sqEnt, err = sqCreate.Save(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// already exists
			sqEnt, err = tx.SearchQuery.Query().
				Where(
					searchquery.Query(searchQ),
					searchquery.Language(lang),
					searchquery.Location(loc),
				).
				First(ctx)
			if err != nil {
				return rollback(tx, err)
			}
		default:
			return rollback(tx, err)
		}
	}

	b := make([]*ent.SERPCreate, 0, len(results))
	for _, result := range results {
		b = append(b,
			tx.SERP.Create().
				SetTitle(result.Title).
				SetDescription(result.Description).
				SetURL(result.URL).
				SetKeyWords(result.KeyWords).
				SetNillableContactInfo(result.ContactInfo).
				SetSqID(sqEnt.ID),
		)
	}
	_, err = db.SERP.CreateBulk(b...).Save(ctx)
	if err != nil {
		return rollback(tx, err)
	}
	return tx.Commit()
}

func GetSQID(ctx context.Context, db *ent.Client, loc, lang, searchQ string) (int, error) {
	sq, err := db.SearchQuery.Query().
		Where(
			searchquery.Query(searchQ),
			searchquery.Query(loc),
			searchquery.Query(lang),
		).
		First(ctx)
	if err != nil {
		return -1, err
	}
	return sq.ID, nil
}
func GetAllResult(ctx context.Context, db *ent.Client, sqID int) ([]SERP, error) {
	entSERPs, err := db.SERP.Query().
		Where(
			serp.SqID(sqID),
			serp.IsRead(false),
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
				ContactInfo: entSERP.ContactInfo,
				Keywords:    entSERP.KeyWords,
				IsRead:      entSERP.IsRead,
				CreatedAt:   entSERP.CreatedAt,
			},
		)
	}
	return results, nil
}

func ExportCSV(ctx context.Context, db *ent.Client, sqID int) (csvAbsFilePath string, err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("starting a transaction: %w", err)
	}
	entSq, err := tx.SearchQuery.Query().Where(searchquery.ID(sqID)).First(ctx)
	if err != nil {
		return "", rollback(tx, err)
	}
	basePath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("ExportCSV.Getwd: %w", err)
	}
	csvFileName := fmt.Sprintf("%s-%s.csv", entSq.Query, time.Now().Format("2006-01-02_15:04:05"))
	csvAbsFilePath = fmt.Sprintf("%s/EPS/db/csv/%s", basePath, csvFileName)
	csvAbsFilePathInPsqlContainer := fmt.Sprintf("/opt/eps/%s", csvFileName)
	csvFile, err := os.OpenFile(csvAbsFilePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return "", fmt.Errorf("ExportCSV: %w", err)
	}
	defer func() {
		csvFile.Close()
		if err != nil {
			_ = os.Remove(csvFile.Name())
		}
	}()
	q := fmt.Sprintf(`COPY (
		SELECT
			serps.url,
			serps.title,
			serps.description,
			serps.contact_info,
			serps.key_words,
			search_queries.query,
			search_queries.location,
			search_queries.language,
			serps.created_at
		FROM serps
		JOIN search_queries
		ON serps.sq_id=search_queries.id)
	TO '%s' WITH (FORMAT CSV, HEADER);`, csvAbsFilePathInPsqlContainer)

	_, err = tx.ExecContext(ctx, q) //, csvAbsFilePath)
	if err != nil {
		return "", rollback(tx, err)
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return csvAbsFilePath, nil
}

func GetAllSearchQueries(ctx context.Context, db *ent.Client) ([]SearchQuery, error) {
	entSearchQueries, err := db.SearchQuery.Query().
		Where(searchquery.IsCanceled(false)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	searchQueries := make([]SearchQuery, 0, len(entSearchQueries))
	for _, entSearchQuery := range entSearchQueries {
		searchQueries = append(searchQueries, SearchQuery{
			Id:         entSearchQuery.ID,
			Query:      entSearchQuery.Query,
			Language:   entSearchQuery.Language,
			Location:   entSearchQuery.Location,
			IsCanceled: entSearchQuery.IsCanceled,
			CreatedAt:  entSearchQuery.CreatedAt,
		})
	}
	return searchQueries, nil
}

func CancelSQ(ctx context.Context, db *ent.Client, sqID int) error {
	_, err := db.SearchQuery.UpdateOneID(sqID).
		SetIsCanceled(true).
		Save(ctx)
	return err
}

func rollback(tx *ent.Tx, err error) error {
	if rErr := tx.Rollback(); rErr != nil {
		err = fmt.Errorf("%w: %v", err, rErr)
	}
	return err
}
