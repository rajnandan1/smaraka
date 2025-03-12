package services

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
)

//handle sch_gh_tranding

func HandleGHTrending(githubURL string, interval int) (*[]string, error) {

	if interval == 1 {
		githubURL = githubURL + "/trending?since=daily"
	} else if interval == 7 {
		githubURL = githubURL + "/trending?since=weekly"
	} else if interval == 30 {
		githubURL = githubURL + "/trending?since=monthly"
	}

	html, err := utils.FetchHTML(githubURL)
	if err != nil {
		logger.LogError("Error fetching Github repos", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.LogError("Error parsing Github repos", err)
		return nil, err
	}
	results := make([]string, 0)
	doc.Find("div.Box article.Box-row").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Find("a").Attr("href"); ok {
			// extract the actual path
			if strings.HasPrefix(href, "/login?return_to=") {
				parts := strings.SplitN(href, "return_to=", 2)
				if len(parts) == 2 {
					if unescaped, err := url.QueryUnescape(parts[1]); err == nil {
						href = unescaped
						results = append(results, "https://github.com"+href)
					}
				}

			}

		}

	})

	if len(results) == 0 {

		return nil, nil
	}

	return &results, nil

}

func HandleHackerNews(url string) (*[]string, error) {
	html, err := utils.FetchHTML(url)
	if err != nil {
		logger.LogError("Error fetching HackerNews", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.LogError("Error parsing HackerNews", err)
		return nil, err
	}

	results := make([]string, 0)
	doc.Find(".submission .titleline > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			results = append(results, href)
		}
	})

	if len(results) == 0 {
		logger.LogInfo("No HackerNews results found", nil)
		return nil, nil
	}
	return &results, nil
}

// create phurl given interval
func createPhURL(interval int) string {
	if interval == 1 {
		yesterday := time.Now().AddDate(0, 0, -1)
		return "/" +
			strconv.Itoa(yesterday.Year()) + "/" +
			strconv.Itoa(int(yesterday.Month())) + "/" +
			strconv.Itoa(yesterday.Day())
	} else if interval == 7 {
		lastWeek := time.Now().AddDate(0, 0, -7)
		year, week := lastWeek.ISOWeek()
		return "/" + strconv.Itoa(year) + "/" + strconv.Itoa(week)
	} else if interval == 30 {
		lastMonth := time.Now().AddDate(0, -1, 0)
		return "/" +
			strconv.Itoa(lastMonth.Year()) + "/" +
			strconv.Itoa(int(lastMonth.Month()))
	}
	return ""
}

func HandlePHDaily(phURL string) (*[]string, error) {
	html, err := utils.FetchHTML(phURL)
	if err != nil {
		logger.LogError("Error fetching Product Hunt daily", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.LogError("Error parsing Product Hunt daily", err)
		return nil, err
	}

	results := make([]string, 0)
	doc.Find("section[data-test^='post-item-']").Each(func(i int, s *goquery.Selection) {

		if href, ok := s.Find("a").Attr("href"); ok {
			if strings.HasPrefix(href, "/posts/") {
				productURL := "https://www.producthunt.com" + href
				results = append(results, productURL)
			}
		}

	})

	if len(results) == 0 {
		logger.LogInfo("No Product Hunt results found", nil)
		return nil, nil
	}

	return &results, nil
}

func (s *ServicesImplementation) RunSchedule(ctx context.Context, interval int) (*[]models.PeriodicResponse, error) {
	dailySchedules, err := s.db.GetActiveSchedulesWithInterval(ctx, interval)
	if err != nil {
		return nil, err
	}
	return s.CreateRun(dailySchedules)
}

func (s *ServicesImplementation) PlaySchedule(ctx context.Context, schedule_ids []string, org_id string) (*[]models.PeriodicResponse, error) {
	dailySchedules, err := s.db.GetSchedulesByIDsAndOrgIDs(ctx, schedule_ids, org_id)
	if err != nil {
		return nil, err
	}
	return s.CreateRun(dailySchedules)
}

func (s *ServicesImplementation) CreateRun(dailySchedules *[]models.Schedule) (*[]models.PeriodicResponse, error) {

	if dailySchedules == nil {
		return nil, nil
	}

	orgMap := make(map[string][]string)
	for _, schedule := range *dailySchedules {

		if _, ok := orgMap[schedule.OrganizationID]; !ok {
			orgMap[schedule.OrganizationID] = make([]string, 0)
		}

		switch schedule.ScheduleType {
		case constants.ScheduleTypeGHTrending:
			if repos, err := HandleGHTrending(schedule.ScheduleURL, schedule.IntervalDays); err == nil && repos != nil {
				orgMap[schedule.OrganizationID] = append(orgMap[schedule.OrganizationID], *repos...)
			}
		case constants.ScheduleTypeHNTrending:
			if urls, err := HandleHackerNews(schedule.ScheduleURL); err == nil && urls != nil {
				orgMap[schedule.OrganizationID] = append(orgMap[schedule.OrganizationID], *urls...)
			}
		case constants.ScheduleTypeGHStarredRepo:
			repos, err := s.ImportGithubStars(schedule.ScheduleURL)
			if err == nil && repos != nil {
				for _, repo := range *repos {
					orgMap[schedule.OrganizationID] = append(orgMap[schedule.OrganizationID], repo.URL)
				}
			}

		case constants.ScheduleTypePHLeaderBoard:
			phURL := schedule.ScheduleURL + createPhURL(schedule.IntervalDays)
			if products, err := HandlePHDaily(phURL); err == nil && products != nil {
				orgMap[schedule.OrganizationID] = append(orgMap[schedule.OrganizationID], *products...)
			}
		}
	}

	//create the response by iterating over the orgMap
	periodicResponse := make([]models.PeriodicResponse, 0)
	for orgID, urls := range orgMap {
		periodicResponse = append(periodicResponse, models.PeriodicResponse{
			OrganizationID: orgID,
			URLs:           urls,
		})
	}

	return &periodicResponse, nil
}
