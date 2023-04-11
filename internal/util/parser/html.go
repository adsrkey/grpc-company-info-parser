package parser

import (
	"errors"
	"github.com/adsrkey/grpc-company-info-parser/internal/delivery/grpc/dto"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

type Parser struct {
	Info *dto.Info
	Html string
}

func New(info *dto.Info, htmlBody string) *Parser {
	return &Parser{
		Info: info,
		Html: htmlBody,
	}
}

func (p *Parser) Parse() error {
	dom := html.NewTokenizer(strings.NewReader(p.Html))
	domToken := dom.Token()

	hasCompanyName := false
	hasKPP := false

loop:
	for {
		tt := dom.Next()
		switch {
		case tt == html.ErrorToken:
			break loop
		case tt == html.StartTagToken:
			domToken = dom.Token()
			data := domToken.Data
			if data == "h1" {
				continue
			}
			if data == "span" {
				continue
			}
		case tt == html.EndTagToken:
			domToken = dom.Token()
		case tt == html.TextToken:
			data := domToken.Data
			if !hasCompanyName {
				if data == "h1" {
					for _, v := range domToken.Attr {
						if v.Key == "itemprop" {
							if v.Val == "name" {
								hasCompanyName = true

								var re = regexp.MustCompile(`[[:punct:]]`)
								str := re.ReplaceAllString(string(dom.Text()), "")
								p.Info.CompanyName = strings.TrimSpace(str)

								break
							}
						}
					}
				}
			}

			if !hasKPP {
				if data == "span" {
					for _, v := range domToken.Attr {
						if v.Key == "id" {
							if v.Val == "clip_kpp" {
								hasKPP = true
								p.Info.Kpp = string(dom.Text())
								break
							}
						}
					}
				}
			}
			if hasCompanyName && hasKPP {
				break loop
			}
		}
	}

	if !hasCompanyName && !hasKPP {
		return errors.New("no company name and kpp")
	}

	return nil
}
