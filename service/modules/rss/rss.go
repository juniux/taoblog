package rss

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/movsb/taoblog/auth"
	"github.com/movsb/taoblog/config"
	"github.com/movsb/taoblog/modules/datetime"
	"github.com/movsb/taoblog/protocols"
	"github.com/movsb/taoblog/service"
	"github.com/movsb/taoblog/themes/modules/handle304"
)

const tmpl = `
<rss version="2.0">
<channel>
	<title>{{ .Name }}</title>
	<link>{{ .Link }}</link>
	<description>{{ .Description }}</description>
	{{- range .Articles -}}
	<item>
		<title>{{ .Title }}</title>
		<link>{{ .Link }}</link>
		<pubDate>{{ .Date }}</pubDate>
		<description>{{ .Content }}</description>
	</item>
	{{ end }}
</channel>
</rss>
`

// Article ...
type Article struct {
	*protocols.Post
	Date    string
	Content template.HTML
	Link    string
}

// RSS ...
type RSS struct {
	Name        string
	Description string
	Link        string
	Articles    []*Article
	Config      *config.Config

	tmpl *template.Template
	svc  *service.Service
	auth *auth.Auth
}

// New ...
func New(cfg *config.Config, svc *service.Service, auth *auth.Auth) *RSS {
	r := &RSS{
		Config:      cfg,
		Name:        svc.Name(),
		Description: svc.Description(),
		Link:        svc.HomeURL(),
		svc:         svc,
		auth:        auth,
	}

	r.tmpl = template.Must(template.New(`rss`).Parse(tmpl))

	return r
}

// ServeHTTP ...
func (r *RSS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handle304.ArticleRequest(w, req, r.svc.LastArticleUpdateTime()) {
		return
	}

	user := r.auth.AuthCookie2(req)

	articles := r.svc.GetLatestPosts(
		user.Context(nil),
		"id,title,date,content",
		int64(r.Config.Site.RSS.ArticleCount),
	)

	rssArticles := make([]*Article, 0, len(articles))
	for _, article := range articles {
		rssArticle := Article{
			Post:    article,
			Date:    datetime.My2Feed(datetime.Proto2My(*article.Date)),
			Content: template.HTML("<![CDATA[" + strings.Replace(string(article.Content), "]]>", "]]]]><!CDATA[>", -1) + "]]>"),
			Link:    fmt.Sprintf("%s/%d/", r.svc.HomeURL(), article.Id),
		}
		rssArticles = append(rssArticles, &rssArticle)
	}

	cr := *r
	cr.Articles = rssArticles

	handle304.ArticleResponse(w, r.svc.LastArticleUpdateTime())

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))

	if err := cr.tmpl.Execute(w, cr); err != nil {
		panic(err)
	}
}
