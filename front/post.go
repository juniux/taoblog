package front

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/movsb/taoblog/service"
	"github.com/movsb/taoblog/service/models"
)

type Post struct {
	*models.Post
	Content      template.HTML
	RelatedPosts []*models.PostForRelated
	server       *service.ImplServer
	_Metas       map[string]interface{}
}

func newPost(post *models.Post, server *service.ImplServer) *Post {
	return &Post{
		Post:    post,
		server:  server,
		Content: template.HTML(post.Content),
	}
}

func newPosts(posts []*models.Post, server *service.ImplServer) []*Post {
	ps := []*Post{}
	for _, p := range posts {
		ps = append(ps, newPost(p, server))
	}
	return ps
}

func (p *Post) Link() string {
	return fmt.Sprintf("/%d/", p.ID)
}

func (p *Post) DateString() string {
	d := strings.Split(strings.Split(p.Date, " ")[0], "-")
	return fmt.Sprintf("%v年%v月%v日", d[0], d[1], d[2])
}

func (p *Post) ModifiedString() string {
	d := strings.Split(strings.Split(p.Modified, " ")[0], "-")
	return fmt.Sprintf("%v年%v月%v日", d[0], d[1], d[2])
}

func (p *Post) TagsString() template.HTML {
	var ts []string
	for _, t := range p.Tags {
		et := html.EscapeString(t)
		ts = append(ts, fmt.Sprintf(`<a href="/tags/%[1]s">%[1]s</a>`, et))
	}
	return template.HTML(strings.Join(ts, " · "))
}

func (p *Post) decodeMetas() {
	if p._Metas == nil {
		json.Unmarshal([]byte(p.Metas), &p._Metas)
	}
}

func (p *Post) CustomHeader() (header string) {
	p.decodeMetas()
	if i, ok := p._Metas["header"]; ok {
		if s, ok := i.(string); ok {
			header = s
		}
	}
	return
}

func (p *Post) CustomFooter() (footer string) {
	p.decodeMetas()
	if i, ok := p._Metas["footer"]; ok {
		if s, ok := i.(string); ok {
			footer = s
		}
	}
	return
}
