package main

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/movsb/taoblog/server/modules/file_managers"
	"github.com/movsb/taoblog/server/modules/memory_cache"
	"github.com/movsb/taoblog/server/modules/utils/datetime"
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
var backupmgr *BlogBackup
var cmtmgr *CommentManager
var postcmtsmgr *PostCommentsManager
var fileredir *FileRedirect
var catmgr *CategoryManager
var memcch *memory_cache.MemoryCache
var blog *Blog
var admin *Admin
var renderer *Renderer

var templates *template.Template

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
	backupmgr = NewBlogBackup()
	cmtmgr = newCommentManager()
	postcmtsmgr = newPostCommentsManager()
	fileredir = NewFileRedirect(config.base, config.files, config.fileHost)
	catmgr = NewCategoryManager()
	memcch = memory_cache.NewMemoryCache(time.Minute * 10)
	blog = NewBlog()
	admin = NewAdmin()
	renderer = NewRenderer()
	defer memcch.Stop()

	loadTemplates()

	router := gin.Default()

	postapi := router.Group("/posts")

	postapi.POST("/update", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		var err error
		pidstr, has := c.GetPostForm("pid")
		pid, err := strconv.ParseInt(pidstr, 10, 64)
		if !has || err != nil || pid < 0 {
			c.String(400, "expect: pid")
			return
		}
		typ, has := c.GetPostForm("source_type")
		if !has {
			c.String(400, "expect: source_type")
			return
		}
		source, has := c.GetPostForm("source")
		if !has {
			c.String(400, "expect: source")
			return
		}

		tx, err := gdb.Begin()
		if err == nil {
			if err = postmgr.update(tx, pid, typ, source); err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			}
		}

		EndReq(c, err, nil)
	})

	routerV1(router)
	routerInternalV1(router)
	routerBlog(router)
	routerAdmin(router)

	router.Run(config.listen)
}

func toInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func routerInternalV1(router *gin.Engine) {
	v1 := router.Group("/.v1")

	optapi := v1.Group("/options")

	optapi.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		varlue, err := optmgr.Get(gdb, name)
		EndReq(c, err, varlue)
	})
}

func routerV1(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1.GET("/avatar", func(c *gin.Context) {
		GetAvatar(c)
	})

	posts := v1.Group("/posts")

	posts.GET("", func(c *gin.Context) {
		tax, hasTax := c.GetQuery("tax")
		parents, hasParents := c.GetQuery("parents")
		slug, hasSlug := c.GetQuery("slug")
		modified := c.Query("modified")
		if (hasTax || hasParents) && hasSlug {
			if hasParents {
				tax = parents
			}
			post, err := postmgr.GetPostBySlug(gdb, tax, slug, modified, hasParents)
			posts := make([]*Post, 0)
			if err == nil {
				posts = append(posts, post)
			} else if err == sql.ErrNoRows {
				err = nil
			}
			// TODO don't return array.
			EndReq(c, err, posts)
			return
		}
		c.Status(400)
	})

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

	posts.GET("/:parent", func(c *gin.Context) {
		pid := toInt64(c.Param("parent"))
		modified := c.Query("modified")
		post, err := postmgr.GetPostByID(gdb, pid, modified)
		posts := make([]*Post, 0)
		if err == nil {
			posts = append(posts, post)
		} else if err == sql.ErrNoRows {
			err = nil
		}
		// TODO don't return array.
		EndReq(c, err, posts)
	})

	posts.POST("/:parent", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		var post Post
		if err := c.ShouldBindJSON(&post); err != nil {
			EndReq(c, err, err)
			return
		}
		if err := txCall(gdb, func(tx Querier) error {
			return postmgr.UpdatePost(tx, &post)
		}); err != nil {
			EndReq(c, err, err)
			return
		}
		EndReq(c, nil, post.ID)
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
		cmt.EMail = c.DefaultPostForm("email", "")
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
					if email != "" && cmt.EMail != "" && strings.EqualFold(email, cmt.EMail) {
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

	posts.POST("/:parent/files", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		err := uploadmgr.Upload(c)
		EndReq(c, err, nil)
	})

	posts.DELETE("/:parent/files/*name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		err := uploadmgr.Delete(c)
		EndReq(c, err, nil)
	})

	posts.GET("/:parent/tags", func(c *gin.Context) {
		pid := toInt64(c.Param("parent"))
		tags, err := tagmgr.GetObjectTagNames(gdb, pid)
		EndReq(c, err, tags)
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

	posts.POST("/:parent/view", func(c *gin.Context) {
		pid := toInt64(c.Param("parent"))
		postmgr.IncrementPageView(gdb, pid)
	})

	posts.GET("/:parent/relates", func(c *gin.Context) {
		pid := toInt64(c.Param("parent"))
		key := fmt.Sprintf("relates:%d", pid)
		if relates, ok := memcch.Get(key); ok {
			EndReq(c, nil, relates)
			return
		}
		relates, err := postmgr.GetRelatedPosts(gdb, pid)
		memcch.SetIf(err == nil, key, relates)
		EndReq(c, err, relates)
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

		rss, err := theFeed(gdb)
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		c.Header("Content-Type", "application/xml")
		if modified := optmgr.GetDef(gdb, "last_post_time", ""); modified != "" {
			c.Header("Last-Modified", datetime.Local2Gmt(modified))
		}
		c.String(http.StatusOK, "%s", rss)
	})

	v1.GET("/posts!all", func(c *gin.Context) {
		posts, err := postmgr.ListAllPosts(gdb)
		EndReq(c, err, posts)
	})

	v1.GET("/posts!latest", func(c *gin.Context) {
		limit := toInt64(c.Query("limit"))
		posts, err := postmgr.GetLatest(gdb, limit)
		EndReq(c, err, posts)
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

	optionsV1(v1)
	tagsV1(v1)
	backupsV1(v1)
	categoryV1(v1)
}

func optionsV1(routerV1 *gin.RouterGroup) {
	optapi := routerV1.Group("/options")

	optapi.GET("", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		items, err := optmgr.List(gdb)
		EndReq(c, err, items)
	})

	optapi.GET("/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		name := c.Param("name")
		varlue, err := optmgr.Get(gdb, name)
		EndReq(c, err, varlue)
	})

	optapi.POST("/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		name := c.Param("name")
		value, _ := ioutil.ReadAll(c.Request.Body) // WARN: Body is consumed

		tx, err := gdb.Begin()
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		err = optmgr.Set(tx, name, string(value))
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
		EndReq(c, err, nil)
	})

	optapi.DELETE("/:name", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}
		name := c.Param("name")

		tx, err := gdb.Begin()
		if err != nil {
			EndReq(c, err, nil)
			return
		}
		err = optmgr.Del(tx, name)
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
		EndReq(c, err, nil)
	})
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

func backupsV1(routerV1 *gin.RouterGroup) {
	backups := routerV1.Group("/backups")

	backups.GET("", func(c *gin.Context) {
		if !auth(c, true) {
			return
		}

		var sb bytes.Buffer
		err := backupmgr.Backup(&sb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "%s", sb.String())
	})
}

func categoryV1(router *gin.RouterGroup) {
	cats := router.Group("/categories")

	cats.GET("", func(c *gin.Context) {
		cats, err := catmgr.ListCategories(gdb)
		EndReq(c, err, cats)
	})

	router.GET("/categories!tree", func(c *gin.Context) {
		if cats, ok := memcch.Get("/categories!tree"); ok {
			EndReq(c, nil, cats)
			return
		}
		cats, err := catmgr.GetTree(gdb)
		memcch.SetIf(err == nil, "/categories!tree", cats)
		EndReq(c, err, cats)
	})
	router.GET("/categories!parse", func(c *gin.Context) {
		tree := c.Query("tree")
		id, err := catmgr.ParseTree(gdb, tree)
		EndReq(c, err, id)
	})
}

func routerBlog(router *gin.Engine) {
	b := router.Group("/blog")
	b.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		blog.Query(c, path)
	})
}

func routerAdmin(router *gin.Engine) {
	a := router.Group("/admin")
	a.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		switch path {
		case "", "/":
			c.Redirect(302, "/admin/login")
			return
		}
		admin.Query(c, path)
	})
	a.POST("/*path", func(c *gin.Context) {
		path := c.Param("path")
		admin.Post(c, path)
	})
}

func loadTemplates() {
	funcs := template.FuncMap{
		"get_config": func(name string) string {
			return optmgr.GetDef(gdb, name, "")
		},
	}
	templates = template.New("taoblog").Funcs(funcs)
	var err error
	if templates, err = templates.ParseGlob("../theme/*.html"); err != nil {
		panic(err)
	}
	if templates, err = templates.ParseGlob("../admin/*.html"); err != nil {
		panic(err)
	}
}
