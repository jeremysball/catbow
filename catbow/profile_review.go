package catbow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const GitHubAPIBaseURL = "https://api.github.com"

const (
	recentMonths = 6
	activeMonths = 12
)

type RepoSummary struct {
	Name        string
	URL         string
	Description string
	Language    string
	Stars       int
	Forks       int
	UpdatedAt   time.Time
	Archived    bool
	Fork        bool
	Topics      []string
}

type RepoReview struct {
	Repo  RepoSummary
	Score int
	Notes []string
}

type ProfileSignals struct {
	TotalRepos      int
	RecentRepos     int
	ActiveRepos     int
	TotalStars      int
	TotalForks      int
	DescribedRepos  int
	UniqueLanguages []string
}

type ProfileReview struct {
	Score               int
	Grade               string
	CompetitiveSignal   string
	EmployabilitySignal string
	Summary             string
	Signals             ProfileSignals
	RepoReviews         []RepoReview
}

type githubRepo struct {
	Name            string    `json:"name"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	Stars           int       `json:"stargazers_count"`
	Forks           int       `json:"forks_count"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Archived        bool      `json:"archived"`
	Fork            bool      `json:"fork"`
	Topics          []string  `json:"topics"`
	DefaultBranch   string    `json:"default_branch"`
	Visibility      string    `json:"visibility"`
	OpenIssuesCount int       `json:"open_issues_count"`
}

func NewGitHubHTTPClient() *http.Client {
	return &http.Client{Timeout: 12 * time.Second}
}

func FetchPublicRepos(ctx context.Context, client *http.Client, baseURL, user string) ([]RepoSummary, error) {
	if strings.TrimSpace(user) == "" {
		return nil, errors.New("github username is required")
	}
	if client == nil {
		return nil, errors.New("http client is required")
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = GitHubAPIBaseURL
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	page := 1
	var repos []RepoSummary
	for {
		reqURL := fmt.Sprintf("%s/users/%s/repos?per_page=100&type=public&sort=updated&page=%d", baseURL, url.PathEscape(user), page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build github request: %w", err)
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "catbow-profile-review")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("github request failed: %w", err)
		}
		body := resp.Body
		if body == nil {
			return nil, errors.New("github response body was empty")
		}
		if resp.StatusCode != http.StatusOK {
			payload, _ := io.ReadAll(io.LimitReader(body, 1024))
			_ = body.Close()
			return nil, fmt.Errorf("github api error: %s", strings.TrimSpace(string(payload)))
		}

		var pageRepos []githubRepo
		if err := json.NewDecoder(body).Decode(&pageRepos); err != nil {
			_ = body.Close()
			return nil, fmt.Errorf("decode github response: %w", err)
		}
		if err := body.Close(); err != nil {
			return nil, fmt.Errorf("close github response: %w", err)
		}
		for _, repo := range pageRepos {
			updatedAt := repo.PushedAt
			if updatedAt.IsZero() {
				updatedAt = repo.UpdatedAt
			}
			repos = append(repos, RepoSummary{
				Name:        repo.Name,
				URL:         repo.HTMLURL,
				Description: strings.TrimSpace(repo.Description),
				Language:    strings.TrimSpace(repo.Language),
				Stars:       repo.Stars,
				Forks:       repo.Forks,
				UpdatedAt:   updatedAt,
				Archived:    repo.Archived,
				Fork:        repo.Fork,
				Topics:      repo.Topics,
			})
		}

		if !hasNextLink(resp.Header.Get("Link")) {
			break
		}
		page++
	}
	return repos, nil
}

func GenerateProfileReview(repos []RepoSummary) ProfileReview {
	return generateProfileReview(repos, time.Now().UTC())
}

func FormatProfileReview(user string, review ProfileReview) string {
	var builder strings.Builder
	username := strings.TrimSpace(user)
	if username == "" {
		username = "(unknown)"
	}
	fmt.Fprintf(&builder, "GitHub public repo review for %s\n", username)
	fmt.Fprintf(&builder, "Total public repos analyzed: %d\n", review.Signals.TotalRepos)
	fmt.Fprintf(&builder, "Overall score: %d/100 (%s)\n", review.Score, review.Grade)
	fmt.Fprintf(&builder, "Competitive signal: %s\n", review.CompetitiveSignal)
	fmt.Fprintf(&builder, "Employability note: %s\n", review.EmployabilitySignal)
	if review.Summary != "" {
		fmt.Fprintf(&builder, "Summary: %s\n", review.Summary)
	}

	fmt.Fprintln(&builder, "\nSignals:")
	fmt.Fprintf(&builder, "- Recent activity (last %d months): %d repos\n", recentMonths, review.Signals.RecentRepos)
	fmt.Fprintf(&builder, "- Active in last %d months: %d repos\n", activeMonths, review.Signals.ActiveRepos)
	fmt.Fprintf(&builder, "- Total stars: %d\n", review.Signals.TotalStars)
	fmt.Fprintf(&builder, "- Total forks: %d\n", review.Signals.TotalForks)
	fmt.Fprintf(&builder, "- Descriptions present: %d/%d repos\n", review.Signals.DescribedRepos, review.Signals.TotalRepos)
	if len(review.Signals.UniqueLanguages) > 0 {
		fmt.Fprintf(&builder, "- Languages: %s\n", strings.Join(review.Signals.UniqueLanguages, ", "))
	} else {
		fmt.Fprintln(&builder, "- Languages: none listed")
	}

	fmt.Fprintln(&builder, "\nRepo-by-repo review:")
	if len(review.RepoReviews) == 0 {
		fmt.Fprintln(&builder, "(no public repos found)")
		return builder.String()
	}

	for i, repoReview := range review.RepoReviews {
		repo := repoReview.Repo
		fmt.Fprintf(&builder, "%d. %s", i+1, repo.Name)
		if repo.Language != "" {
			fmt.Fprintf(&builder, " (%s)", repo.Language)
		}
		fmt.Fprintln(&builder)
		if repo.URL != "" {
			fmt.Fprintf(&builder, "   URL: %s\n", repo.URL)
		}
		if repo.Description != "" {
			fmt.Fprintf(&builder, "   Description: %s\n", repo.Description)
		}
		updated := "unknown"
		if !repo.UpdatedAt.IsZero() {
			updated = repo.UpdatedAt.Format("2006-01-02")
		}
		fmt.Fprintf(&builder, "   Stars: %d, Forks: %d, Updated: %s\n", repo.Stars, repo.Forks, updated)
		if len(repo.Topics) > 0 {
			fmt.Fprintf(&builder, "   Topics: %s\n", strings.Join(repo.Topics, ", "))
		}
		if len(repoReview.Notes) > 0 {
			fmt.Fprintf(&builder, "   Notes: %s\n", strings.Join(repoReview.Notes, "; "))
		}
	}
	fmt.Fprintln(&builder, "\nNote: This review only considers public repo metadata and is not a hiring decision.")
	return builder.String()
}

func generateProfileReview(repos []RepoSummary, now time.Time) ProfileReview {
	signals := ProfileSignals{TotalRepos: len(repos)}
	languageSet := make(map[string]struct{})

	reviews := make([]RepoReview, 0, len(repos))
	for _, repo := range repos {
		repoScore, notes, recent, active := scoreRepo(repo, now)
		reviews = append(reviews, RepoReview{Repo: repo, Score: repoScore, Notes: notes})

		signals.TotalStars += repo.Stars
		signals.TotalForks += repo.Forks
		if repo.Description != "" {
			signals.DescribedRepos++
		}
		if recent {
			signals.RecentRepos++
		}
		if active {
			signals.ActiveRepos++
		}
		if repo.Language != "" {
			languageSet[repo.Language] = struct{}{}
		}
	}

	signals.UniqueLanguages = sortedKeys(languageSet)

	score := overallScore(signals)
	grade := gradeFromScore(score)
	competitive := competitiveSignal(score)
	employability := employabilitySignal(score)
	summary := summaryFromSignals(signals)

	return ProfileReview{
		Score:               score,
		Grade:               grade,
		CompetitiveSignal:   competitive,
		EmployabilitySignal: employability,
		Summary:             summary,
		Signals:             signals,
		RepoReviews:         reviews,
	}
}

func scoreRepo(repo RepoSummary, now time.Time) (int, []string, bool, bool) {
	var notes []string
	isRecent := false
	isActive := false
	repoScore := 0

	if repo.Description != "" {
		repoScore += 2
		notes = append(notes, "has description")
	} else {
		notes = append(notes, "missing description")
	}

	if repo.Language != "" {
		repoScore++
	}

	if repo.Stars > 0 {
		repoScore += minInt(5, repo.Stars)
		notes = append(notes, fmt.Sprintf("%d stars", repo.Stars))
	}

	if repo.Forks > 0 {
		repoScore += minInt(3, repo.Forks)
		notes = append(notes, fmt.Sprintf("%d forks", repo.Forks))
	}

	if !repo.UpdatedAt.IsZero() {
		recentCutoff := now.AddDate(0, -recentMonths, 0)
		activeCutoff := now.AddDate(0, -activeMonths, 0)
		if repo.UpdatedAt.After(recentCutoff) {
			repoScore += 5
			isRecent = true
			isActive = true
			notes = append(notes, "updated recently")
		} else if repo.UpdatedAt.After(activeCutoff) {
			repoScore += 2
			isActive = true
			notes = append(notes, "updated within a year")
		}
		if !isRecent && !isActive {
			notes = append(notes, "stale update")
		}
	}

	if repo.Archived {
		repoScore -= 2
		notes = append(notes, "archived")
	}
	if repo.Fork {
		repoScore--
		notes = append(notes, "forked repo")
	}

	return repoScore, notes, isRecent, isActive
}

func overallScore(signals ProfileSignals) int {
	repoCountScore := minInt(30, signals.TotalRepos*3)
	activityScore := 0
	if signals.RecentRepos > 0 {
		activityScore = 20
	} else if signals.ActiveRepos > 0 {
		activityScore = 10
	}

	popularity := signals.TotalStars*2 + signals.TotalForks
	popularityScore := int(math.Round(math.Min(30, math.Log1p(float64(popularity))*8)))

	docScore := 0
	if signals.TotalRepos > 0 {
		docScore = int(math.Round(float64(signals.DescribedRepos) / float64(signals.TotalRepos) * 10))
	}

	languageScore := minInt(10, len(signals.UniqueLanguages)*2)

	score := repoCountScore + activityScore + popularityScore + docScore + languageScore
	if score > 100 {
		return 100
	}
	return score
}

func gradeFromScore(score int) string {
	switch {
	case score >= 85:
		return "A"
	case score >= 70:
		return "B"
	case score >= 55:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}

func competitiveSignal(score int) string {
	switch {
	case score >= 70:
		return "strong"
	case score >= 55:
		return "moderate"
	case score >= 40:
		return "emerging"
	default:
		return "limited"
	}
}

func employabilitySignal(score int) string {
	switch {
	case score >= 70:
		return "Based on public repo signals alone, it would be somewhat surprising if you were having difficulty finding a job."
	case score >= 55:
		return "Based on public repo signals alone, it would not be surprising either way; the portfolio shows promise but still has gaps."
	case score >= 40:
		return "Based on public repo signals alone, it would not be surprising if finding a job is difficult; strengthening activity and documentation could help."
	default:
		return "Based on public repo signals alone, it is not surprising if finding a job is difficult; consider adding clearer documentation and recent work."
	}
}

func summaryFromSignals(signals ProfileSignals) string {
	if signals.TotalRepos == 0 {
		return "No public repositories were found to evaluate."
	}

	var strengths []string
	var gaps []string

	if signals.RecentRepos > 0 {
		strengths = append(strengths, fmt.Sprintf("%d recently updated repos", signals.RecentRepos))
	} else if signals.ActiveRepos > 0 {
		strengths = append(strengths, fmt.Sprintf("%d repos updated within a year", signals.ActiveRepos))
	} else {
		gaps = append(gaps, "no recent updates in the last year")
	}

	if signals.TotalStars > 0 {
		strengths = append(strengths, fmt.Sprintf("%d total stars", signals.TotalStars))
	} else {
		gaps = append(gaps, "no starred repos yet")
	}

	if signals.DescribedRepos == signals.TotalRepos {
		strengths = append(strengths, "every repo has a description")
	} else if signals.DescribedRepos > 0 {
		gaps = append(gaps, fmt.Sprintf("%d repos missing descriptions", signals.TotalRepos-signals.DescribedRepos))
	} else {
		gaps = append(gaps, "descriptions are missing")
	}

	if len(signals.UniqueLanguages) > 1 {
		strengths = append(strengths, fmt.Sprintf("languages used: %s", strings.Join(signals.UniqueLanguages, ", ")))
	} else if len(signals.UniqueLanguages) == 1 {
		strengths = append(strengths, fmt.Sprintf("primary language: %s", signals.UniqueLanguages[0]))
	} else {
		gaps = append(gaps, "languages not detected")
	}

	var parts []string
	if len(strengths) > 0 {
		parts = append(parts, "Strengths: "+strings.Join(strengths, ", "))
	}
	if len(gaps) > 0 {
		parts = append(parts, "Gaps: "+strings.Join(gaps, ", "))
	}

	return strings.Join(parts, " | ")
}

func hasNextLink(linkHeader string) bool {
	if linkHeader == "" {
		return false
	}
	for _, part := range strings.Split(linkHeader, ",") {
		if strings.Contains(part, "rel=\"next\"") {
			return true
		}
	}
	return false
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
