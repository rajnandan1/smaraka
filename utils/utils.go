package utils

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/logger"
	"github.com/rajnandan1/smaraka/models"

	"github.com/spaolacci/murmur3"
)

func GetDomain(inputURL string) (string, error) {
	// Parse the URL
	parsedURL, err := _url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	// Get the host (domain) part
	domain := parsedURL.Host

	// Check if the domain contains a port, remove it if present
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}

	return domain, nil
}

// create md5 hash of a string
func CreateMD5Hash(input string) string {
	// Create a new MD5 hash instance
	hash := md5.New()

	// Write the string to the hash
	hash.Write([]byte(input))

	// Get the hash result and return it
	return hex.EncodeToString(hash.Sum(nil))
}

func FetchHTML(url string) (string, error) {
	// Create a new context

	// Perform the HTTP GET request with a custom User-Agent header
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("Error creating request", err)
		return "", err
	}
	req.Header.Set("User-Agent", constants.HeadlessUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		logger.LogError("Error fetching URL", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError("Error reading response body", err)
		return "", err
	}

	// Print the HTML content
	return string(body), nil
}

func ProperImageURL(url string, imageURL string) string {

	if imageURL == "" {
		return ""
	}

	if strings.HasPrefix(imageURL, "http") {
		return imageURL
	}

	if strings.HasPrefix(imageURL, "//") {
		return "https:" + imageURL
	}

	if strings.HasPrefix(imageURL, "/") {
		parsedURL, _ := _url.Parse(url)
		return parsedURL.Scheme + "://" + parsedURL.Host + imageURL
	}

	return url + imageURL
}

// write a func to encode string to base64
func EncodeStringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func MurmurHashToRange(s string) string {
	// Create a new Murmur3 hash instance (32-bit)
	hash := murmur3.New32()

	// Write the string to the hash
	hash.Write([]byte(s))

	// Get the hash result and apply modulus to fit within 24-bit range
	hashInt := hash.Sum32() % 16777216
	hexStr := fmt.Sprintf("%06X", hashInt)

	return "#" + hexStr
}

func ParseFetch(htmlText string) (string, error) {

	return htmltomarkdown.ConvertString(htmlText)

}

func ParseSEOFromHTML(html string) (*models.URLStore, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var title string = ""
	var excerpt string = ""
	var accentColor string = ""
	var imageSmall string = ""
	var imageLarge string = ""

	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})

	title = doc.Find("title").First().Text()

	//check og:title set to title if present
	doc.Find("meta[property='og:title']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			title = content
		}
	})

	//check site_name if present prepend to title and not already present in title
	doc.Find("meta[property='og:site_name']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			if !strings.Contains(title, content) {
				title = content + " - " + title
			}
		}
	})

	//find <meta name="description">
	doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			excerpt = content
		}
	})

	//check og:description set to excerpt if present
	doc.Find("meta[property='og:description']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			excerpt = content
		}
	})

	//check twitter:description set to excerpt if present
	doc.Find("meta[name='twitter:description']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			excerpt = content
		}
	})

	// check if name="theme-color" set to accentColor if present
	doc.Find("meta[name='theme-color']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			accentColor = content
		}
	})

	//check rel="shortcut icon" set to imageSmall if present
	doc.Find("link[rel='shortcut icon']").Each(func(i int, s *goquery.Selection) {
		if href, _ := s.Attr("href"); href != "" {
			imageSmall = href
		}
	})

	//check rel="icon" set to imageSmall if present
	doc.Find("link[rel='icon']").Each(func(i int, s *goquery.Selection) {
		if href, _ := s.Attr("href"); href != "" {
			imageSmall = href
		}
	})

	//check if og:image set to imageLarge if present
	doc.Find("meta[property='og:image']").Each(func(i int, s *goquery.Selection) {
		if content, _ := s.Attr("content"); content != "" {
			imageLarge = content
		}
	})

	//if imageLarge is not set check for first <img> tag
	if imageLarge == "" {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			if src, _ := s.Attr("src"); src != "" {
				imageLarge = src
				return
			}
		})
	}

	bookmark := &models.URLStore{
		Title:       title,
		Excerpt:     excerpt,
		AccentColor: accentColor,
		ImageSmall:  imageSmall,
		ImageLarge:  imageLarge,
	}
	return bookmark, nil
}

func RemoveNonNumeric(s string) string {
	// Compile a regular expression that matches non-numeric characters
	re := regexp.MustCompile("[^0-9]")
	// Replace all non-numeric characters with an empty string
	return re.ReplaceAllString(s, "")
}
func RemoveAllNewlines(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	s1 := re.ReplaceAllString(s, " ")
	//replace all double space with single
	re = regexp.MustCompile(`\s+`)
	return strings.Trim(re.ReplaceAllString(s1, " "), " ")
}

func FetchInnerHTMLChrome(url string) (string, error) {
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var htmlText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			htmlText, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("timeout while fetching page: %w", err)
		}
		return "", fmt.Errorf("error fetching page: %w", err)
	}

	return htmlText, nil
}

// given array of strings, break the array n subArrays
func ChunkStringArray(input []string, numOfSubArrays int) [][]string {
	var output [][]string
	chunkSize := (len(input) + numOfSubArrays - 1) / numOfSubArrays
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input)
		}
		output = append(output, input[i:end])
	}
	return output
}

func GetStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}
