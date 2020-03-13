package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivym/zhihu-spider/internal/db"
	"github.com/vivym/zhihu-spider/internal/nlp"
	"github.com/vivym/zhihu-spider/internal/spiders"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	v, p := viper.New(), pflag.NewFlagSet(friendlyAppName, pflag.ExitOnError)

	configure(v, p)

	p.String("config", "", "Configuration file")
	p.Bool("version", false, "Show version information")

	_ = p.Parse(os.Args[1:])

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", friendlyAppName, version, commitHash, buildDate)

		os.Exit(0)
	}

	if c, _ := p.GetString("config"); c != "" {
		v.SetConfigFile(c)
	}

	err := v.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if !configFileNotFound {
		log.Panic("failed to read configuration", err)
	}

	var config configuration
	err = v.Unmarshal(&config)
	if err != nil {
		log.Panic("failed to unmarshal configuration", err)
	}

	if configFileNotFound {
		log.Println("configuration file not found")
	}

	fmt.Printf("%+v\n", config)

	if err := db.SetupDB(config.DB.URI, config.DB.DBName); err != nil {
		log.Fatalf("mongodb error: %v\n", err)
	}

	var nlpToolkit *nlp.NLPToolkit
	nlpToolkit, err = nlp.New(config.NLP)
	if err != nil {
		log.Fatalf("nlp error: %v\n", err)
	}

	spider := spiders.New(config.Spider, nlpToolkit)
	if err := spider.Go(); err != nil {
		log.Fatalf("spider error: %v\n", err)
	}
	log.Println("done.")
}
