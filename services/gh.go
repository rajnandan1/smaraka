package services

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
)

func doGithub(username string, pageNo int) (*[]models.GithubRepo, error) {
	githubURL := "https://github.com/stars/" + username + "/repositories?direction=desc&filter=all&page=" + strconv.Itoa(pageNo) + "&sort=created"
	logger.LogInfo("Fetching Github repos for user", githubURL)
	html, err := utils.FetchHTML(githubURL)
	if err != nil {
		logger.LogError("Error fetching Github repos", err)
		return nil, err
	}
	results := make([]models.GithubRepo, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.LogError("Error parsing Github repos", err)
		return nil, err
	}
	doc.Find("ul.repo-list li").Each(func(i int, s *goquery.Selection) {
		title := utils.RemoveAllNewlines(s.Find("h3").Text())
		//remove all spaces from title
		title = strings.ReplaceAll(title, " ", "")
		url, _ := s.Find("a").Attr("href")
		url = "https://github.com" + strings.ReplaceAll(url, " ", "")

		stars, _ := strconv.Atoi(utils.RemoveNonNumeric(s.Find(`a[href$="stargazers"]`).Text()))
		excerpt := utils.RemoveAllNewlines(s.Find("p").Text())
		repo := models.GithubRepo{
			URL:     url,
			Stars:   stars,
			Excerpt: excerpt,
			Title:   title,
		}
		results = append(results, repo)
	})

	if len(results) == 0 {

		return nil, nil
	}

	return &results, nil
}

func (s *ServicesImplementation) ImportGithubStars(username string) (*[]models.GithubRepo, error) {

	results := make([]models.GithubRepo, 0)
	for i := 1; i < 101; i++ {
		if repos, err := doGithub(username, i); err == nil && repos != nil {
			results = append(results, *repos...)
		} else {
			break
		}
	}
	return &results, nil
}
