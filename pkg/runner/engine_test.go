package runner

import (
	"testing"
)

func TestAlreadyAdded(t *testing.T) {
	repos := []GHRepo{
		{URL: "https://github.com/zer0yu/ghtopdep"},
		{URL: "https://github.com/projectdiscovery/gologger"},
	}

	if !alreadyAdded("https://github.com/zer0yu/ghtopdep", repos) {
		t.Errorf("alreadyAdded function failed, expected %v, got %v", true, false)
	}

	if alreadyAdded("https://github.com/nonexistent/repo", repos) {
		t.Errorf("alreadyAdded function failed, expected %v, got %v", false, true)
	}
}

func TestSortRepos(t *testing.T) {
	repos := []GHRepo{
		{URL: "https://github.com/zer0yu/ghtopdep", Stars: 10},
		{URL: "https://github.com/projectdiscovery/gologger", Stars: 20},
	}

	sortedRepos := sortRepos(repos, 2)

	if sortedRepos[0].Stars != 20 {
		t.Errorf("sortRepos function failed, expected %v, got %v", 20, sortedRepos[0].Stars)
	}

	if sortedRepos[1].Stars != 10 {
		t.Errorf("sortRepos function failed, expected %v, got %v", 10, sortedRepos[1].Stars)
	}
}

func TestReadableStars(t *testing.T) {
	repos := []GHRepo{
		{Stars: 999},
		{Stars: 1000},
		{Stars: 9999},
		{Stars: 10000},
		{Stars: 999999},
		{Stars: 1000000},
	}

	expected := []string{"999", "1.0K", "10.0K", "10K", "1000K", "1000000"}

	repos4show := readableStars(repos)

	for i, repo := range repos4show {
		if repo.Stars != expected[i] {
			t.Errorf("readableStars function failed, expected %v, got %v", expected[i], repo.Stars)
		}
	}
}
