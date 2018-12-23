package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/movsb/taoblog/admin"
	"github.com/movsb/taoblog/front"
	"github.com/movsb/taoblog/gateway"
	"github.com/movsb/taoblog/modules/datetime"
	"github.com/movsb/taoblog/modules/memory_cache"
	"github.com/movsb/taoblog/protocols"
	"github.com/movsb/taoblog/service"
	"github.com/movsb/taoblog/service/modules/file_managers"
)

type xConfig struct {
	base     string
	listen   string
	username string
	password string
	database string
	key      string
	files    string
	fileHost string
	mail     string
}

var gkey string
var config xConfig
var gdb *sql.DB
var tagmgr *TagManager
var postmgr *PostManager
var optmgr *OptionManager
var auther *GenericAuth
var uploadmgr *FileUpload
var cmtmgr *CommentManager
var postcmtsmgr *PostCommentsManager
var fileredir *FileRedirect
var catmgr *CategoryManager
var memcch *memory_cache.MemoryCache
var theFront *front.Front
var theAdmin *admin.Admin
var implServer protocols.IServer
var cacheServer protocols.IServer
var theGateway *gateway.Gateway

func auth(c *gin.Context, finish bool) bool {
	if auther.AuthHeader(c) || auther.AuthCookie(c) {
		return true
	}
	if finish {
		EndReq(c, false, "auth error")
	}
	return false
}

func main() {
	flag.StringVar(&config.listen, "listen", "127.0.0.1:2564", "the port to which the server listen")
	flag.StringVar(&config.username, "username", "taoblog", "the database username")
	flag.StringVar(&config.password, "password", "taoblog", "the database password")
	flag.StringVar(&config.database, "database", "taoblog", "the database name")
	flag.StringVar(&config.key, "key", "", "api key")
	flag.StringVar(&config.base, "base", ".", "taoblog directory")
	flag.StringVar(&config.files, "files", ".", "the files folder")
	flag.StringVar(&config.fileHost, "file-host", "//localhost", "the backup file host")
	flag.StringVar(&config.mail, "mail", "//", "example.com:465/user@example.com/password")
	flag.Parse()

	if config.key == "" {
		panic("invalid key")
	}

	var err error
	dataSource := fmt.Sprintf("%s:%s@/%s", config.username, config.password, config.database)
	gdb, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	gdb.SetMaxIdleConns(10)

	defer gdb.Close()

	tagmgr = NewTagManager()
	postmgr = NewPostManager()
	optmgr = newOptionsModel()
	auther = &GenericAuth{}
	auther.SetLogin(optmgr.GetDef(gdb, "login", "x"))
	auther.SetKey(config.key)
	uploadmgr = NewFileUpload(file_managers.NewLocalFileManager(config.files))
	cmtmgr = newCommentManager()
	postcmtsmgr = newPostCommentsManager()
	fileredir = NewFileRedirect(config.base, config.files, config.fileHost)
	catmgr = NewCategoryManager()
	memcch = memory_cache.NewMemoryCache(time.Minute * 10)
	defer memcch.Stop()
	implServer = service.NewImplServer(gdb, auther)

	router := gin.Default()

	theAdmin = admin.NewAdmin(implServer, &router.RouterGroup)
	theFront = front.NewFront(implServer, &router.RouterGroup)

	routerV1(router)

	v2 := router.Group("/v2")
	v2.Use(func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					if err == sql.ErrNoRows {
						c.Status(404)
						return
					}
				}
				panic(e)
			}
		}()
		c.Next()
	})
	theGateway = gateway.NewGateway(v2, implServer)

	router.Run(config.listen)
}

func toInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		//panic(fmt.Errorf("expect number: %s", s))
	}
	return n
}

func routerV1(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	posts := v1.Group("/posts")

	posts.POST("", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		var post Post
		if err := c.ShouldBindJSON(&post); err != nil {
			EndReq(c, err, err)
			return
		}
		if err := txCall(gdb, func(tx Querier) error {
			return postmgr.CreatePost(tx, &post)
		}); err != nil {
			EndReq(c, err, err)
			return
		}
		EndReq(c, nil, &post)
	})

	posts.GET("/:parent/files/*name", func(c *gin.Context) {
		referrer := strings.ToLower(c.GetHeader("referer"))
		if strings.Contains(referrer, "://blog.csdn.net") {
			c.Redirect(302, "/1.jpg")
			return
		}
		parent := toInt64(c.Param("parent"))
		name := c.Param("name")
		if strings.Contains(name, "/../") {
			c.String(400, "bad file")
			return
		}
		logged := auth(c, false)
		path := fileredir.Redirect(logged, fmt.Sprintf("%d/%s", parent, name))
		c.Redirect(302, path)
	})

	posts.GET("/:parent/comments:count", func(c *gin.Context) {
		parent := toInt64(c.Param("parent"))
		count := postmgr.GetCommentCount(gdb, parent)
		EndReq(c, true, count)
	})

	posts.GET("/:parent/comments", func(c *gin.Context) {
		var err error

		parent := toInt64(c.Param("parent"))
		offset := toInt64(c.Query("offset"))
		count := toInt64(c.Query("count"))
		order := c.DefaultQuery("order", "asc")

		cmts, err := postcmtsmgr.GetPostComments(gdb, 0, offset, count, parent, order == "asc")

		if err != nil {
			EndReq(c, err, nil)
			return
		}

		var loggedin = auth(c, false)

		for _, c := range cmts {
			c.private = loggedin
		}

		EndReq(c, true, cmts)
	})

	posts.POST("/:parent/comments", func(c *gin.Context) {
		var err error
		var cmt Comment
		var loggedin bool

		loggedin = auth(c, false)

		cmt.PostID = toInt64(c.Param("parent"))
		cmt.Parent = toInt64(c.DefaultPostForm("parent", "0"))
		cmt.Author = c.DefaultPostForm("author", "")
		cmt.Email = c.DefaultPostForm("email", "")
		cmt.URL = c.DefaultPostForm("url", "")
		cmt.IP = c.ClientIP()
		cmt.Date = datetime.MyLocal()
		cmt.Content = c.DefaultPostForm("content", "")

		tx, err := gdb.Begin()
		if err != nil {
			EndReq(c, err, nil)
			return
		}

		if has, err := postmgr.Has(tx, cmt.PostID); err != nil || !has {
			EndReq(c, errors.New("找不到文章"), nil)
			tx.Rollback()
			return
		}

		if !loggedin {
			{
				notAllowedEmails := strings.Split(optmgr.GetDef(tx, "not_allowed_emails", ""), ",")
				if adminEmail := optmgr.GetDef(tx, "email", ""); adminEmail != "" {
					notAllowedEmails = append(notAllowedEmails, adminEmail)
				}

				log.Println(notAllowedEmails)

				// TODO use regexp to detect equality.
				for _, email := range notAllowedEmails {
					if email != "" && cmt.Email != "" && strings.EqualFold(email, cmt.Email) {
						EndReq(c, errors.New("不能使用此邮箱地址"), nil)
						tx.Rollback()
						return
					}
				}
			}
			{
				notAllowedAuthors := strings.Split(optmgr.GetDef(tx, "not_allowed_authors", ""), ",")
				if adminName := optmgr.GetDef(tx, "nickname", ""); adminName != "" {
					notAllowedAuthors = append(notAllowedAuthors, adminName)
				}

				for _, author := range notAllowedAuthors {
					if author != "" && cmt.Author != "" && strings.EqualFold(author, string(cmt.Author)) {
						EndReq(c, errors.New("不能使用此昵称"), nil)
						tx.Rollback()
						return
					}
				}
			}
		}

		if err = cmtmgr.CreateComment(tx, &cmt); err != nil {
			EndReq(c, err, nil)
			tx.Rollback()
			return
		}

		count := cmtmgr.GetAllCount(tx)
		optmgr.Set(tx, "comment_count", count)

		postcmtsmgr.UpdatePostCommentsCount(tx, cmt.PostID)

		retCmt := c.DefaultQuery("return_cmt", "0") == "1"

		if !retCmt {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
				EndReq(c, err, nil)
				return
			}
			EndReq(c, nil, gin.H{
				"id": cmt.ID,
			})
		} else {
			cmts, err := postcmtsmgr.GetPostComments(tx, cmt.ID, 0, 1, cmt.PostID, true)
			if err != nil || len(cmts) == 0 {
				EndReq(c, errors.New("error get comment"), nil)
				tx.Rollback()
				return
			}
			if err = tx.Commit(); err != nil {
				tx.Rollback()
				EndReq(c, err, nil)
				return
			}
			cmts[0].private = loggedin
			EndReq(c, err, cmts[0])
		}

		doNotify(gdb, &cmt) // TODO use cmts[0]
	})

	posts.DELETE("/:parent/comments/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		var err error

		parent := toInt64(c.Param("parent"))
		id := toInt64(c.Param("name"))

		// TODO check referrer
		_ = parent

		tx, err := gdb.Begin()
		if err != nil {
			panic(err)
		}
		err = postcmtsmgr.DeletePostComment(tx, id)
		if err = tx.Commit(); err != nil {
			tx.Rollback()
		}
		EndReq(c, err, nil)
	})

	posts.GET("/:parent/files", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		files, err := uploadmgr.List(c)
		EndReq(c, err, files)
	})

	posts.POST("/:parent/files/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		err := uploadmgr.Upload(c)
		EndReq(c, err, nil)
	})

	posts.DELETE("/:parent/files/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		err := uploadmgr.Delete(c)
		EndReq(c, err, nil)
	})

	posts.POST("/:parent/tags", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		var tags []string
		if err := c.ShouldBindJSON(&tags); err != nil {
			EndReq(c, err, nil)
			return
		}

		pid := toInt64(c.Param("parent"))
		if has, err := postmgr.Has(gdb, pid); true {
			if err != nil {
				EndReq(c, err, nil)
				return
			} else if !has {
				EndReq(c, fmt.Errorf("post not found: %v", pid), nil)
				return
			}
		}

		tx, err := gdb.Begin()
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		tagmgr.UpdateObjectTags(tx, pid, tags)
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			EndReq(c, err, nil)
			return
		}
		EndReq(c, nil, nil)
	})

	v1.GET("/posts!manage", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		posts, err := postmgr.GetPostsForManagement(gdb)
		EndReq(c, err, posts)
	})

	v1.GET("/posts!rss", func(c *gin.Context) {
		if ifModified := c.GetHeader("If-Modified-Since"); ifModified != "" {
			if modified := optmgr.GetDef(gdb, "last_post_time", ""); modified != "" {
				if ifModified == datetime.Local2Gmt(modified) {
					c.Status(http.StatusNotModified)
					return
				}
			}
		}

		var rss string

		if s, ok := memcch.Get("rss"); ok {
			rss = s.(string)
		} else {
			s, err := theFeed(gdb)
			if err != nil {
				EndReq(c, err, nil)
				return
			}
			rss = s
			memcch.Set("rss", s)
		}

		c.Header("Content-Type", "application/xml")
		if modified := optmgr.GetDef(gdb, "last_post_time", ""); modified != "" {
			c.Header("Last-Modified", datetime.Local2Gmt(modified))
		}
		c.String(http.StatusOK, "%s", rss)
	})

	v1.GET("/posts!all", func(c *gin.Context) {
		var posts []*PostForArchiveQuery
		if p, ok := memcch.Get("posts:all"); ok {
			posts = p.([]*PostForArchiveQuery)
		} else {
			p, err := postmgr.ListAllPosts(gdb)
			if err != nil {
				EndReq(c, err, posts)
				return
			}
			memcch.Set("posts:all", p)
			posts = p
		}
		EndReq(c, nil, posts)
	})

	archives := v1.Group("/archives")

	archives.GET("/categories/:name", func(c *gin.Context) {
		id := toInt64(c.Param("name"))
		ps, err := postmgr.GetPostsByCategory(gdb, id)
		EndReq(c, err, ps)
	})

	archives.GET("/tags/:name", func(c *gin.Context) {
		tag := c.Param("name")
		ps, err := postmgr.GetPostsByTags(gdb, tag)
		EndReq(c, err, ps)
	})

	archives.GET("/dates/:year/:month", func(c *gin.Context) {
		year := toInt64(c.Param("year"))
		month := toInt64(c.Param("month"))

		ps, err := postmgr.GetPostsByDate(gdb, year, month)
		EndReq(c, err, ps)
	})

	tools := v1.Group("/tools")

	tools.POST("/aes2htm", func(c *gin.Context) {
		aes2htm(c)
	})

	v1.Group("/sitemap.xml").GET("", func(c *gin.Context) {
		host := "https://" + optmgr.GetDef(gdb, "home", "localhost")
		maps, err := createSitemap(gdb, host)
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		c.Header("Content-Type", "application/xml")
		c.String(200, "%s", maps)
	})

	tagsV1(v1)
}

func tagsV1(routerV1 *gin.RouterGroup) {
	tagsV1 := routerV1.Group("/tags")

	tagsV1.GET("", func(c *gin.Context) {
		tags, err := tagmgr.ListTags(gdb)
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		EndReq(c, nil, tags)
		return
	})

	tagsV1.POST("/:parent", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		tagID := toInt64(c.Param("parent"))

		var tag Tag

		if err := c.ShouldBindJSON(&tag); err != nil {
			EndReq(c, err, nil)
			return
		}

		tag.ID = tagID

		tx, err := gdb.Begin()
		if err != nil {
			EndReq(c, err, nil)
			return
		}

		err = tagmgr.UpdateTag(tx, &tag)
		if err != nil {
			tx.Rollback()
			EndReq(c, err, nil)
			return
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			EndReq(c, err, nil)
			return
		}
	})
}
