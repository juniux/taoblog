package data

import (
	"html/template"
	"io"

	"github.com/movsb/taoblog/auth"
	"github.com/movsb/taoblog/config"
	"github.com/movsb/taoblog/protocols"
	"github.com/movsb/taoblog/service"
)

// Data holds all data for rendering the site.
type Data struct {
	// all configuration.
	Config *config.Config

	// current login user, non-nil.
	User *auth.User

	// The response writer.
	Writer io.Writer

	// The template
	Template *template.Template

	// Metadata
	Meta *MetaData

	// If it is home page.
	Home *HomeData

	// If it is a post.
	Post *PostData

	// If it is the RSS.
	Rss *RssData

	// If it is the Sitemap.
	Sitemap *SitemapData

	// If it is the Search.
	Search *SearchData

	// If it is the Posts.
	Posts *PostsData

	// If it is the Tags.
	Tags *TagsData
}

// MetaData ...
type MetaData struct {
	Title string
}

// HomeData ...
type HomeData struct {
	Posts    []*Post
	Comments []*Comment

	PostCount    int64
	PageCount    int64
	CommentCount int64
}

// PostData ...
type PostData struct {
	Post *Post
}

// NewDataForPost ...
func NewDataForPost(cfg *config.Config, user *auth.User, service *service.Service, post *protocols.Post) *Data {
	d := &Data{
		Config: cfg,
		User:   user,
		Meta:   &MetaData{},
	}
	p := &PostData{
		Post: newPost(post),
	}
	d.Post = p
	if cfg.Site.ShowRelatedPosts {
		p.Post.Related = service.GetRelatedPosts(post.Id)
	}
	p.Post.Tags = service.GetPostTags(p.Post.Id)
	return d
}