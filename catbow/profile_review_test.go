package catbow

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateProfileReviewSummary(t *testing.T) {
	now := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	repos := []RepoSummary{
		{
			Name:        "alpha",
			URL:         "https://github.com/example/alpha",
			Description: "cli tool",
			Language:    "Go",
			Stars:       3,
			Forks:       1,
			UpdatedAt:   now.AddDate(0, -1, 0),
		},
		{
			Name:      "beta",
			URL:       "https://github.com/example/beta",
			Language:  "Python",
			Stars:     0,
			Forks:     0,
			UpdatedAt: now.AddDate(0, -8, 0),
		},
	}

	review := generateProfileReview(repos, now)

	assert.Equal(t, 2, review.Signals.TotalRepos)
	assert.Equal(t, 1, review.Signals.RecentRepos)
	assert.Equal(t, 2, review.Signals.ActiveRepos)
	assert.Equal(t, 3, review.Signals.TotalStars)
	assert.Equal(t, 1, review.Signals.TotalForks)
	assert.Equal(t, 1, review.Signals.DescribedRepos)
	assert.Equal(t, []string{"Go", "Python"}, review.Signals.UniqueLanguages)
	assert.NotEmpty(t, review.Summary)
	assert.NotEmpty(t, review.EmployabilitySignal)
}

func TestFormatProfileReviewIncludesRepoDetails(t *testing.T) {
	review := ProfileReview{
		Score:               72,
		Grade:               "B",
		CompetitiveSignal:   "strong",
		EmployabilitySignal: "Based on public repo signals alone, it would be somewhat surprising if you were having difficulty finding a job.",
		Summary:             "Strengths: 1 recently updated repos",
		Signals: ProfileSignals{
			TotalRepos:      1,
			RecentRepos:     1,
			ActiveRepos:     1,
			TotalStars:      2,
			TotalForks:      0,
			DescribedRepos:  1,
			UniqueLanguages: []string{"Go"},
		},
		RepoReviews: []RepoReview{
			{
				Repo: RepoSummary{
					Name:        "catbow",
					URL:         "https://github.com/example/catbow",
					Description: "rainbow text",
					Language:    "Go",
					Stars:       2,
					Forks:       0,
					UpdatedAt:   time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
				},
				Score: 10,
				Notes: []string{"updated recently"},
			},
		},
	}

	output := FormatProfileReview("jeremysball", review)

	assert.True(t, strings.Contains(output, "GitHub public repo review for jeremysball"))
	assert.True(t, strings.Contains(output, "Total public repos analyzed: 1"))
	assert.True(t, strings.Contains(output, "catbow (Go)"))
	assert.True(t, strings.Contains(output, "Notes: updated recently"))
	assert.True(t, strings.Contains(output, "Note: This review only considers public repo metadata"))
}
