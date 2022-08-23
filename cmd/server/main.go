package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
)

type Config struct {
	Server struct {
		Address string        `default:"0.0.0.0:8080"`
		Timeout time.Duration `default:"5s"`
	}
	Github struct {
		V3APIURL string `default:"https://api.github.com"`
		App      struct {
			IntegrationID int64  `split_words:"true"`
			PrivateKey    string `split_words:"true"`
		}
	}
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger

	_ = godotenv.Load()

	var c Config
	if err := envconfig.Process("app", &c); err != nil {
		logger.Fatal().Err(err).Msg("failed to read envconfig")
	}

	ghConfig := newAppConfig(c)

	cc, err := githubapp.NewDefaultCachingClientCreator(
		ghConfig,
		githubapp.WithClientUserAgent("check-pls/0.0.0"),
		githubapp.WithClientTimeout(c.Server.Timeout),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create client creator")
	}

	ph, crh, prh, csh := initHandlers(cc)
	webhookHandler := githubapp.NewDefaultEventDispatcher(ghConfig, ph, crh, prh, csh)

	http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	logger.Info().Msgf("Starting server on %s...", c.Server.Address)
	err = http.ListenAndServe(c.Server.Address, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("error starting server")
	}
}

func newAppConfig(c Config) githubapp.Config {
	gh := c.Github

	var ghc githubapp.Config
	ghc.V3APIURL = gh.V3APIURL
	ghc.App.IntegrationID = gh.App.IntegrationID
	ghc.App.PrivateKey = gh.App.PrivateKey

	return ghc
}

func initHandlers(cc githubapp.ClientCreator) (*PushHandler, *CheckRunHandler, *PullRequestHandler, *CheckSuiteHandler) {
	return &PushHandler{ClientCreator: cc}, &CheckRunHandler{ClientCreator: cc}, &PullRequestHandler{ClientCreator: cc}, &CheckSuiteHandler{ClientCreator: cc}
}
