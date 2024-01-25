package runner

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v58/github"
	"github.com/imroc/req/v3"
	"github.com/projectdiscovery/gologger"
	"github.com/zer0yu/ghtracker/pkg/options"
	"io"
	"os"
	"time"
)

type GHRunner struct {
	options   *options.GHTopDepOptions
	ghClient  *GHClient
	reqClient *req.Client
}

type GHClient struct {
	gitClient *github.Client
	ghctx     context.Context
}

func NewRunner(options *options.GHTopDepOptions) (*GHRunner, error) {
	ghRunner := &GHRunner{options: options}
	if options.Version {
		showBanner()
		os.Exit(0)
	}

	if options.URL == "" {
		gologger.Fatal().Msgf("URL is empty!")
	}

	if (options.Description || options.Search != "") && options.Token != "" {
		rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
		if err != nil {
			gologger.Error().Msgf("Init rateLimiter Error: %v\n", err)
		}

		// Set GitHub authentication information and configure the rate limit
		ghRunner.ghClient = &GHClient{
			gitClient: nil,
			ghctx:     context.Background(),
		}

		ghRunner.ghClient.gitClient = github.NewClient(rateLimiter).WithAuthToken(options.Token)

		// verify token
		_, resp, err := ghRunner.ghClient.gitClient.Users.Get(ghRunner.ghClient.ghctx, "")
		if err != nil {
			//fmt.Printf("\nerror: %v\n", err)
			gologger.Error().Msgf("Token is invalid!")
		}

		// If a Token Expiration has been set, it will be displayed.
		if !resp.TokenExpiration.IsZero() {
			gologger.Error().Msgf("Token Expiration: %v\n", resp.TokenExpiration)
		}
	} else if (options.Description || options.Search != "") && options.Token == "" {
		gologger.Error().Msgf("Please provide token!")
	}

	// set req client
	ghRunner.reqClient = req.C().
		EnableDumpEachRequest().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // Ignore when there is an underlying error, e.g. network error.
				return nil
			}
			// Treat non-successful responses as errors, record raw dump content in error message.
			if !resp.IsSuccessState() { // Status code is not between 200 and 299.
				resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
			}
			return nil
		})
	ghRunner.reqClient.ImpersonateChrome()
	ghRunner.reqClient.R().
		SetRetryCount(15).
		SetRetryBackoffInterval(10*time.Second, 20*time.Second).
		SetRetryFixedInterval(2 * time.Second)

	return ghRunner, nil
}

func (r *GHRunner) RunGHCrawler() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = "\n"
	s.FinalMSG = "Complete!\n"
	s.Start()
	outputs := []io.Writer{}
	if r.options.OutputFile != "" {
		outputWriter := NewOutputWriter(true)
		file, err := outputWriter.createFile(r.options.OutputFile, true)
		if err != nil {
			gologger.Error().Msgf("Could not create file %s: %s\n", r.options.OutputFile, err)
			return err
		}
		err = r.GHCrawlEngine(append(outputs, file))
		if err != nil {
			gologger.Error().Msgf("Run GHCrawlEngine Error!")
			s.Stop()
			return err
		}
		file.Close()
	} else {
		err := r.GHCrawlEngine(outputs)
		if err != nil {
			gologger.Error().Msgf("Run GHCrawlEngine Error!")
			s.Stop()
			return err
		}
	}

	s.Stop()
	return nil
}
