package services

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

// ParseUploadFileFirefox parses the file content based on the file extension
func parseUploadFileFirefox(html string) ([]models.FileUploadResponse, error) {
	// Parse the file content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	response := make([]models.FileUploadResponse, 0)

	//get all the links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			data := models.FileUploadResponse{
				Name: s.Text(),
				URL:  href,
			}
			if addDate, ok := s.Attr("add_date"); ok {
				data.AddedOn = addDate
			}
			if iconURI, ok := s.Attr("icon_uri"); ok {
				data.Icon = iconURI
			}
			if data.Icon == "" {
				if icon, ok := s.Attr("icon"); ok {
					data.Icon = icon
				}
			}

			response = append(response, data)
		}
	})

	return response, nil
}

func (s *ServicesImplementation) ParseUploadFile(fileObj models.FileUpload) ([]models.FileUploadResponse, error) {
	// Parse the file content
	switch fileObj.ImportType {
	case constants.Browser:
		return parseUploadFileFirefox(fileObj.FileContent)
	default:
		return nil, nil
	}
}
