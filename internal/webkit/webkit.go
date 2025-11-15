//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package webkit

import (
	"net/url"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/playwright-community/playwright-go"
)

type WebKit struct {
	service *playwright.Playwright
	browser playwright.Browser
	html2md *converter.Converter
	ua      string
}

type Config struct {
	AutoConfig bool
	DriverDir  string
	UserAgent  string
}

const (
	UserAgentCurl   = "curl/8.7.1"
	UserAgentSafari = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.6 Safari/605.1.15"
)

func New(cfg Config) (web *WebKit, err error) {
	if cfg.AutoConfig {
		if err = autoconfig(cfg); err != nil {
			return nil, err
		}
	}

	web = &WebKit{}

	web.ua = cfg.UserAgent
	if web.ua == "" {
		web.ua = UserAgentSafari
	}

	web.html2md = converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	web.html2md.Register.TagType("img", converter.TagTypeRemove, converter.PriorityStandard)

	web.service, err = playwright.Run(
		&playwright.RunOptions{
			DriverDirectory: cfg.DriverDir,
		},
	)
	if err != nil {
		return nil, err
	}

	web.browser, err = web.service.WebKit.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
		},
	)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (api *WebKit) Close() error {
	if err := api.browser.Close(); err != nil {
		return err
	}
	if err := api.service.Stop(); err != nil {
		return err
	}
	return nil
}

func (api *WebKit) Extract(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	host := url.Scheme + "://" + url.Host

	bcxt, err := api.browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(api.ua),
	})
	if err != nil {
		return "", err
	}

	page, err := bcxt.NewPage()
	if err != nil {
		return "", err
	}

	if _, err = page.Goto(uri); err != nil {
		return "", err
	}

	html, err := page.Content()
	if err != nil {
		return "", err
	}

	md, err := api.html2md.ConvertString(html,
		converter.WithDomain(host),
	)
	if err != nil {
		return "", err
	}

	return md, nil
}

func autoconfig(conf Config) error {
	return playwright.Install(
		&playwright.RunOptions{
			Browsers:        []string{"webkit"},
			DriverDirectory: conf.DriverDir,
		},
	)
}
