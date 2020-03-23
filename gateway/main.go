package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/movsb/taoblog/auth"
	"github.com/movsb/taoblog/service"
)

type Gateway struct {
	router  *gin.RouterGroup
	service *service.Service
	auther  *auth.Auth
}

func NewGateway(router *gin.RouterGroup, service *service.Service, auther *auth.Auth) *Gateway {
	g := &Gateway{
		router:  router,
		service: service,
		auther:  auther,
	}

	router.POST("/comments!setPostID", g.auth, g.SetCommentPostID)

	g.routePosts()
	g.routeOthers()
	g.routeOptions()

	return g
}

func (g *Gateway) routeOthers() {
	g.router.GET("/avatar", g.GetAvatar)
	g.router.GET("/backup", g.auth, g.GetBackup)
}

func (g *Gateway) routePosts() {
	c := g.router.Group("/posts")

	// posts
	c.POST("", g.auth, g.CreatePost)
	c.POST("/:name", g.auth, g.UpdatePost)
	c.GET("/:name", g.auth, g.GetPost)

	// comments
	c.GET("/:name/comments!count", g.GetPostCommentCount)
	c.DELETE("/:name/comments/:comment_name", g.auth, g.DeleteComment)

	// files
	c.GET("/:name/files/*file", g.GetFile)
	c.GET("/:name/files", g.auth, g.ListFiles)
	c.POST("/:name/files/*file", g.auth, g.CreateFile)
	c.DELETE("/:name/files/*file", g.auth, g.DeleteFile)

	c.POST("/:name/status", g.auth, g.SetPostStatus)

	// for mirror host
	files := g.router.Group("/files")
	files.GET("/:name/*file", g.GetFile)
}

func (g *Gateway) routeOptions() {
	c := g.router.Group("/options")
	c.GET("/:name", g.auth, g.GetOption)
	c.POST(":name", g.auth, g.SetOption)
}
