package services

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
)

//handle sch_gh_tranding

func HandleGHTrending(githubURL string) (*[]string, error) {

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

func (s *ServicesImplementation) DailySchedules(ctx context.Context, interval int) (*[]models.PeriodicResponse, error) {
	dailySchedules, err := s.db.GetOrgSchedulesInterval(ctx, interval)
	if err != nil {
		return nil, err
	}

	//get all the schedule ids first
	scheduleIDs := make([]string, 0)
	for _, schedule := range dailySchedules {
		scheduleIDs = append(scheduleIDs, schedule.ScheduleID)
	}

	//create a map[string][][]string to store the urls against the schedule id
	urlsMap := make(map[string][]string)
	for _, schedule := range dailySchedules {
		urlsMap[schedule.ScheduleID] = make([]string, 0)
	}

	for _, scheduleID := range scheduleIDs {
		urls := make([]string, 0)
		schedule, err := s.db.GetScheduleByID(ctx, scheduleID)
		if err != nil {
			return nil, err
		}

		switch schedule.ScheduleID {
		case "sch_gh_trending_daily", "sch_gh_trending_monthly", "sch_gh_trending_weekly":
			if repos, err := HandleGHTrending(schedule.ScheduleURL); err == nil && repos != nil {
				urls = *repos
			}
		case "sch_ph_leaderboard_daily", "sch_ph_leaderboard_monthly", "sch_ph_leaderboard_weekly":
			phURL := schedule.ScheduleURL + createPhURL(interval)

			if products, err := HandlePHDaily(phURL); err == nil && products != nil {
				urls = *products
			}
		}
		urlsMap[scheduleID] = urls
	}

	//create the response
	periodicData := make(map[string][]string)
	for _, orgSchedule := range dailySchedules {
		periodicData[orgSchedule.OrganizationID] = append(periodicData[orgSchedule.OrganizationID], urlsMap[orgSchedule.ScheduleID]...)
	}

	var urls []models.PeriodicResponse
	for orgID, urlList := range periodicData {
		urls = append(urls, models.PeriodicResponse{
			OrganizationID: orgID,
			URLs:           &urlList,
		})
	}

	return &urls, nil
}
