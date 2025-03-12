package services

import (
	"context"
	"errors"
	_url "net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/go-shiori/go-readability"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"
	"github.com/rajnandan1/smaraka/utils"
)

func (s *ServicesImplementation) GetContentEasy(url string) (*models.URLStore, error) {

	// call url and get html
	html, err := utils.FetchHTML(url)
	if err != nil {
		return nil, err
	}

	bookmark, err := utils.ParseSEOFromHTML(html)
	if err != nil {
		logger.LogError("Error parsing SEO from HTML", err)
		return nil, err
	}

	bookmark.ImageLarge = utils.ProperImageURL(url, bookmark.ImageLarge)
	bookmark.ImageSmall = utils.ProperImageURL(url, bookmark.ImageSmall)

	reader := strings.NewReader(html)
	if parsedURL, err := _url.Parse(url); err == nil {
		article, err := readability.FromReader(reader, parsedURL)
		if err == nil {

			if article.Excerpt != "" {
				bookmark.Excerpt = article.Excerpt
			}

			if article.Favicon != "" {
				bookmark.ImageSmall = article.Favicon
			}

			if article.Title != "" {
				bookmark.Title = article.Title
			}

		}
	}
	//https://github.com/stars/GreatMach/repositories?direction=desc&filter=all&page=2&sort=created
	if bookmark.Title == "" {
		bookmark.Title = url
	}
	if bookmark.AccentColor == "" {
		bookmark.AccentColor = utils.MurmurHashToRange(url)
	}

	bookmark.URL = url
	bookmark.Status = constants.BookmarkStatusPending

	//get domain
	if domain, err := utils.GetDomain(url); err == nil {

		//replace all . with dot
		bookmark.Domain = domain
	}

	return bookmark, nil

}
func (s *ServicesImplementation) DoContentCompleteByID(url string) (*models.URLStore, error) {

	//add a timeout for context cancel 60 secs
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	urlStore, err := s.db.GetURLStoreByURL(ctx, url)
	if err != nil {
		logger.LogError("Error getting URLStore by URL", err)
		return nil, err
	}
	logger.LogInfo("Fetching inner HTML", urlStore.URL)
	htmlText, err := utils.FetchInnerHTMLChrome(urlStore.URL)
	if err != nil {
		logger.LogError("Error fetching inner HTML", err)
		return nil, err
	}

	result := strings.Trim(s.policy.Sanitize(htmlText), " ")
	result = utils.StripHTML(result)

	if newBookmark, newBookmarkErr := utils.ParseSEOFromHTML(htmlText); newBookmarkErr == nil {
		if newBookmark.Title != "" {
			urlStore.Title = newBookmark.Title
		}
		if newBookmark.Excerpt != "" {
			urlStore.Excerpt = newBookmark.Excerpt
		}
		if newBookmark.AccentColor != "" {
			urlStore.AccentColor = newBookmark.AccentColor
		}

		urlStore.ImageLarge = utils.ProperImageURL(urlStore.URL, urlStore.ImageLarge)
		urlStore.ImageSmall = utils.ProperImageURL(urlStore.URL, urlStore.ImageSmall)
	}

	urlStore.FullText = result
	urlStore.Status = constants.BookmarkStatusComplete

	updatedUrlStore, err := s.db.UpdateURLStoreByID(ctx, urlStore.ID, *urlStore)
	if err != nil {
		logger.LogError("Error updating bookmark", err)
		return nil, err
	}

	return updatedUrlStore, nil

}

func (s *ServicesImplementation) BulkLightAndFullJob(validURLs []string, orgId string) error {
	jobStartAt := time.Now()
	logger.LogInfo("BulkLightAndFullJob for count: ", len(validURLs))

	//loop through urls
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(constants.HeadlessUserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Duration(len(validURLs)*15+60)*time.Second)
	defer cancel()

	for _, validURL := range validURLs {
		s.db.InsertJobQueue(ctx, orgId, validURL)
	}

	for _, validURL := range validURLs {

		s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusQueued)
		_, getURLStoreByURLErr := s.db.GetURLStoreByURL(ctx, validURL)
		if getURLStoreByURLErr != nil {
			urlStore, getContentEasyErr := s.GetContentEasy(validURL)
			if getContentEasyErr != nil {
				logger.LogError("Error getting content", getContentEasyErr)
				s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
				continue
			}
			urlStore.ID = s.db.NewID("url")
			if _, insertNewURLStoreErr := s.db.InsertNewURLStore(ctx, *urlStore); insertNewURLStoreErr != nil {
				logger.LogError("Error inserting url store", insertNewURLStoreErr)
				s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
				continue
			}
		}
		urlStore, getURLStoreByURLErr := s.db.GetURLStoreByURL(ctx, validURL)
		if getURLStoreByURLErr != nil {
			logger.LogError("Error getting url store", getURLStoreByURLErr)
			s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
			continue
		}

		urlOrg := models.URLOrganizations{
			ID:             s.db.NewID("url_org"),
			URLID:          urlStore.ID,
			OrganizationID: orgId,
			Status:         constants.URLStatusActive,
		}

		if _, insertNewURLOrganizationErr := s.db.InsertNewURLOrganization(ctx, urlOrg); insertNewURLOrganizationErr != nil {
			// logger.LogError("Error inserting url org", insertNewURLOrganizationErr)
			if strings.Contains(insertNewURLOrganizationErr.Error(), "duplicate key value violates unique constraint") {
				s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusComplete)
			}
			s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
			continue
		}

		if urlStore.Status != constants.BookmarkStatusPending {
			continue
		}

		logger.LogInfo("Fetching inner HTML", urlStore.URL)

		var htmlText string

		chromeErr := chromedp.Run(ctx,
			chromedp.Navigate(validURL),
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
			chromedp.ActionFunc(func(ctx context.Context) error {
				node, chromeErr2 := dom.GetDocument().Do(ctx)
				if chromeErr2 != nil {
					return chromeErr2
				}
				htmlText, chromeErr2 = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				return chromeErr2
			}),
		)
		if chromeErr != nil {
			if errors.Is(chromeErr, context.DeadlineExceeded) {
				logger.LogError("timeout while fetching page:", chromeErr)

			}
			logger.LogError("error fetching page:", chromeErr)
			s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
			continue
		}

		result := strings.Trim(s.policy.Sanitize(htmlText), " ")
		result = utils.StripHTML(result)

		if newBookmark, newBookmarkErr := utils.ParseSEOFromHTML(htmlText); newBookmarkErr == nil {
			urlStore.Title = newBookmark.Title
			urlStore.Excerpt = newBookmark.Excerpt
			if newBookmark.AccentColor != "" {
				urlStore.AccentColor = newBookmark.AccentColor
			}

			urlStore.ImageLarge = utils.ProperImageURL(urlStore.URL, urlStore.ImageLarge)
			urlStore.ImageSmall = utils.ProperImageURL(urlStore.URL, urlStore.ImageSmall)
		}

		urlStore.FullText = result
		urlStore.Status = constants.BookmarkStatusComplete

		_, updateURLStoreByIDErr := s.db.UpdateURLStoreByID(ctx, urlStore.ID, *urlStore)
		if updateURLStoreByIDErr != nil {
			logger.LogError("Error updating bookmark", updateURLStoreByIDErr)
			s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusFailed)
			return updateURLStoreByIDErr
		}

		s.db.UpdateJobQueueStatus(ctx, orgId, validURL, constants.JobQueueStatusComplete)
	}

	jobEndAt := time.Now()

	logger.LogInfo("BulkLightAndFullJob for count: ", len(validURLs), " took: ", jobEndAt.Sub(jobStartAt))

	return nil
}
