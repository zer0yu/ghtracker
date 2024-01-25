package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v58/github"
	"github.com/olekukonko/tablewriter"
	"github.com/projectdiscovery/gologger"
	"github.com/zer0yu/ghtracker/pkg/config"
	"github.com/zer0yu/ghtracker/pkg/utils/textwrap"
	"io"
	"math"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (r *GHRunner) GHCrawlEngine(writers []io.Writer) error {
	destination := "repository"
	destinations := "repositories"
	repositories := r.options.Repositories // Assuming repositories is a list of strings

	if !repositories {
		destination = "package"
		destinations = "packages"
	}
	//
	var repos []GHRepo
	moreThanZeroCount := 0
	totalReposCount := 0

	pageURL, _ := r.getPageUrl(destination)
	//fmt.Println(pageURL)
	for {
		resp, err := r.reqClient.R().Get(pageURL)
		if err != nil {
			gologger.Error().Msgf("Req pageURL Error!")
		}
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		dependents := doc.Find(config.ItemSelector)
		totalReposCount += dependents.Length()

		dependents.Each(func(i int, s *goquery.Selection) {
			repoStarsList := s.Find(config.StarsSelector)
			if repoStarsList.Length() > 0 {
				repoStars := strings.Replace(repoStarsList.First().Text(), ",", "", -1)
				repoStarsNum, _ := strconv.Atoi(strings.TrimSpace(repoStars))

				if repoStarsNum != 0 {
					moreThanZeroCount++
				}

				if repoStarsNum >= r.options.MinStar {
					relativeRepoURL, _ := s.Find(config.RepoSelector).First().Attr("href")
					repoURL := fmt.Sprintf("%s%s", config.GithubURL, relativeRepoURL)

					if !alreadyAdded(repoURL, repos) && repoURL != pageURL {
						if r.options.Description {
							repoDescription, _ := r.fetchDescription(r.ghClient.ghctx, relativeRepoURL)
							repos = append(repos, GHRepo{URL: repoURL, Stars: repoStarsNum,
								Description: repoDescription})
						} else {
							repos = append(repos, GHRepo{URL: repoURL, Stars: repoStarsNum})
						}
					}
				}
			}
		})

		node := doc.Find(config.NextButtonSelector)
		if node.Length() == 2 {
			pageURL, _ = node.Eq(1).Attr("href")
		} else if node.Length() == 0 || node.First().Text() == "Previous" {
			break
		} else if node.First().Text() == "Next" {
			pageURL, _ = node.First().Attr("href")
		}
	}

	sortedRepos := sortRepos(repos, totalReposCount)

	if r.options.Search != "" {
		for _, repo := range repos {
			u, _ := url.Parse(repo.URL)
			repoPath := strings.TrimPrefix(u.Path, "/")
			query := fmt.Sprintf("%s repo:%s", r.options.Search, repoPath)
			results, _, _ := r.ghClient.gitClient.Search.Code(r.ghClient.ghctx, query, nil)

			for _, s := range results.CodeResults {
				gologger.Info().Msgf("%s with %d stars\n", *s.HTMLURL, repo.Stars)
			}

			outputWriter := NewOutputWriter(true)
			for _, writer := range writers {
				err := outputWriter.writeJSONResults(results.CodeResults, writer)
				if err != nil {
					gologger.Error().Msgf("Could not write results for %s: %s\n", query, err)
				}
			}

		}
	} else {
		r.showResult(sortedRepos, totalReposCount, moreThanZeroCount, destinations, writers)
	}

	return nil
}

type GHPackage struct {
	Count     int
	PackageID string
}

type GHRepo struct {
	URL         string
	Stars       int
	Description string
}

type GHRepo4Show struct {
	URL         string
	Stars       string
	Description string
}

func (r *GHRunner) getPageUrl(destination string) (string, error) {
	pageURL := fmt.Sprintf("%s/network/dependents?dependent_type=%s", r.options.URL, strings.ToUpper(destination))
	resp, err := r.reqClient.R().Get(pageURL)
	if err != nil {
		gologger.Error().Msgf("Req pageURL Error!")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil { // Append raw dump content to error message if goquery parse failed to help troubleshoot.
		gologger.Error().Msgf("failed to parse html: %s, raw content:\n%s", err.Error(), resp.Dump())
	}
	link := doc.Find(".select-menu-item")
	if link.Length() > 0 {
		packages := make([]GHPackage, 0)
		link.Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			repoURL := fmt.Sprintf("https://github.com/%s", href)
			resp, err := r.reqClient.R().Get(repoURL)
			if err != nil {
				// 处理错误
				gologger.Error().Msgf("Req repoURL Error!")
			}
			parsedItem, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				gologger.Error().Msgf("failed to parse html: %s, raw content:\n%s", err.Error(), resp.Dump())
			}
			hrefURL, _ := url.Parse(repoURL)
			packageID := strings.Split(hrefURL.RawQuery, "=")[1]

			countStr := parsedItem.Find(".table-list-filters a:first-child").First().Text()
			countStr = strings.Split(countStr, " ")[0]
			countStr = strings.ReplaceAll(countStr, ",", "")
			count, _ := strconv.Atoi(countStr)

			packages = append(packages, GHPackage{
				Count:     count,
				PackageID: packageID,
			})
		})
		sort.Slice(packages, func(i, j int) bool {
			return packages[i].Count > packages[j].Count
		})

		mostPopularPackageID := packages[0].PackageID
		pageURL = fmt.Sprintf("%s/network/dependents?dependent_type=%s&package_id=%s", r.options.URL,
			strings.ToUpper(destination), mostPopularPackageID)
	}
	return pageURL, nil
}

// fetchDescription fetches the description of a repository
func (r *GHRunner) fetchDescription(ctx context.Context, relativeURL string) (string, error) {
	parts := strings.Split(relativeURL, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid relative URL")
	}
	owner, repo := parts[0], parts[1]
	repoinfo, _, err := r.ghClient.gitClient.Repositories.Get(ctx, owner, repo)
	if _, ok := err.(*github.RateLimitError); ok {
		gologger.Error().Msgf("hit rate limit")
	}
	//fmt.Println(repoinfo.GetDescription())
	if repoinfo.GetDescription() != "" {
		return textwrap.Shorten(repoinfo.GetDescription(), 60), nil
	}

	return " ", nil
}

func alreadyAdded(repoURL string, repos []GHRepo) bool {
	for _, repo := range repos {
		if repo.URL == repoURL {
			return true
		}
	}
	return false
}

func sortRepos(repos []GHRepo, rows int) []GHRepo {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})

	if rows > len(repos) {
		rows = len(repos)
	}
	return repos[:rows]
}

// showResult shows the result of the search or fmt the result in table format
func (r *GHRunner) showResult(repos []GHRepo, totalReposCount int, moreThanZeroCount int, destinations string, writers []io.Writer) {
	if r.options.Table {
		if len(repos) > 0 {
			repos4show := readableStars(repos)
			tw := tablewriter.NewWriter(os.Stdout)
			tw.SetHeader([]string{"URL", "Stars", "Description"})
			for _, repo := range repos4show {
				tw.Append([]string{repo.URL, repo.Stars, repo.Description})
			}
			tw.Render()
			fmt.Printf("found %d %s others %s are private\n", totalReposCount, destinations, destinations)
			fmt.Printf("found %d %s with more than zero star\n", moreThanZeroCount, destinations)
		} else {
			fmt.Printf("Doesn't find any %s that match search request\n", destinations)
		}
	} else {
		reposJSON, _ := json.Marshal(readableStars(repos))
		fmt.Println(string(reposJSON))
	}

	outputWriter := NewOutputWriter(true)
	for _, writer := range writers {
		err := outputWriter.writeJSONResults(readableStars(repos), writer)
		if err != nil {
			gologger.Error().Msgf("Could not write results for %s: %s\n", r.options.OutputFile, err)
		}
	}

}

func humanize(num int) string {
	if num < 1000 {
		return fmt.Sprintf("%d", num)
	} else if num < 10000 {
		return fmt.Sprintf("%.1fK", math.Round(float64(num)/100.0)/10.0)
	} else if num < 1000000 {
		return fmt.Sprintf("%.0fK", math.Round(float64(num)/1000.0))
	} else {
		return fmt.Sprintf("%d", num)
	}
}

func readableStars(repos []GHRepo) []GHRepo4Show {
	repos4show := make([]GHRepo4Show, 0)
	for i := range repos {
		repos4show = append(repos4show, GHRepo4Show{URL: repos[i].URL, Stars: humanize(repos[i].Stars), Description: repos[i].Description})
	}
	return repos4show
}
