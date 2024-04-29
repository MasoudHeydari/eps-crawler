package cmd

import (
	"context"
	"fmt"
	"github.com/karust/openserp/db"
	"os"
	"strings"
	"time"

	"github.com/karust/openserp/baidu"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/google"
	"github.com/karust/openserp/yandex"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchCMD = &cobra.Command{
	Use:     "search",
	Aliases: []string{"find"},
	Short:   "Search results using chosen web search engine (google, yandex, baidu)",
	Args:    cobra.MatchAll(cobra.OnlyValidArgs, cobra.ExactArgs(2)),
	Run:     search,
}

// ./eps search --offset "SEARCH QUERY"

func search(cmd *cobra.Command, args []string) {
	var err error
	engineType := args[0]
	offset, err := cmd.Flags().GetInt("offset")
	if err != nil {
		logrus.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Connect to DB
	client, err := db.NewDB()
	if err != nil {
		logrus.Errorf("Failed to connect to DB, error: %v", err)
		return
	}
	fmt.Println(client.Schema)

	query := core.Query{
		Text:     args[1],
		LangCode: "en",
		Location: "NL",
		Limit:    10,
		Offset:   offset,
	}
	//var results []core.SearchResult
	engine := buildEngine(engineType)
	if engine == nil {
		logrus.Errorf("Failed to build Engine, No `%s` search engine found", engineType)
		return
	}

	//if config.App.IsRawRequests {
	//	results, err = searchRaw(engineType, query)
	//} else {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Inserted all found results")
			os.Exit(0)
		default:
			results, err := searchBrowser(engine, query)
			if err != nil {
				logrus.Error(err)
				return
			}

			// Save found records into the DB
			err = db.InsertBulk(ctx, client, results)
			if err != nil {
				logrus.Errorf("failed to insert results to DB, error: %v", err)
				return
			}
			query.NextPage()
		}
	}

	//b, err := json.MarshalIndent(results, "", " ")
	//if err != nil {
	//	logrus.Error(err)
	//	return
	//}
	//
	//fmt.Println(string(b))
}

func searchBrowser(engine core.SearchEngine, query core.Query) ([]core.SearchResult, error) {
	return engine.Search(query)
}

//func searchRaw(engineType string, query core.Query) ([]core.SearchResult, error) {
//	logrus.Warn("Browserless results are very inconsistent or may not even work!")
//
//	switch strings.ToLower(engineType) {
//	case "yandex":
//		return yandex.Search(query)
//	case "google":
//		return google.Search(query)
//	case "baidu":
//		return baidu.Search(query)
//	default:
//		logrus.Infof("No `%s` search engine found", engineType)
//	}
//	return nil, nil
//}

func buildEngine(engineType string) core.SearchEngine {
	opts := core.BrowserOpts{
		IsHeadless:    !config.App.IsBrowserHead, // Disable headless if browser head mode is set
		IsLeakless:    config.App.IsLeakless,
		Timeout:       time.Second * time.Duration(config.App.Timeout),
		LeavePageOpen: config.App.IsLeaveHead,
	}

	if config.App.IsDebug {
		opts.IsHeadless = false
	}

	browser, err := core.NewBrowser(opts)
	if err != nil {
		logrus.Error(err)
	}
	var engine core.SearchEngine
	switch strings.ToLower(engineType) {
	case "yandex":
		engine = yandex.New(*browser, config.YandexConfig)
	case "google":
		engine = google.New(*browser, config.GoogleConfig)
	case "baidu":
		engine = baidu.New(*browser, config.BaiduConfig)
	default:
	}
	return engine
}

func init() {
	RootCmd.AddCommand(searchCMD)
	searchCMD.Flags().IntP("offset", "o", 0, "set offset for your search query")
}
