package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/manager"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/modules/dorkers"
	"github.com/xdumpgo/XDG/utils"
	"github.com/braintree/manners"
	"github.com/foolin/goview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type StatusRefresh struct {
	Status string
	Runtime string
	Threads int
	Proxies int
	Requests int
	Errors int
	RPS int
	Dorks int
	Urls int
	Injectables int
	Tables int
	Columns int
	Rows int
	Progress int
	Index int
	End int
}

const (
	userkey = "user"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.Redirect(302, "/auth/login")
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}


func StartWebUI(address string) *manners.GracefulServer  {

	fmt.Println("Starting up webui")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	server := manners.NewServer()

	gv := goview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		//Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})
	goview.Use(gv)

	r.Static("/assets", "./static/assets")

	r.Use(sessions.Sessions("xdgsess", sessions.NewCookieStore([]byte("secret"))))

	authRoutes := r.Group("/auth/")
	authRoutes.GET("/login", func(c *gin.Context) {
		if err := goview.Render(c.Writer, 200, "login.html", goview.M{}); err != nil {
			utils.LogError(err.Error())
		}
	})
	authRoutes.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		password := c.PostForm("password")

		if len(username) == 0 || len(password) == 0 {
			c.Set("failedLogin", true)
			c.Redirect(200, "/auth/login")
			return
		}

		if resp := auth.Login(username, password); resp.Status == "success" {
			session.Set(userkey, username) // In real world usage you'd set this to the users ID
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			}
			injection.Init()


			c.Redirect(302, "/")
		}
	})

	authorized := r.Group("/", AuthRequired)

	authorized.GET("/", func(c *gin.Context) {
		var runtime string
		if utils.Module == "Idle" {
			runtime = "00:00:00"
		} else {
			runtime = utils.FmtDuration(time.Now().Sub(utils.StartTime))
		}

		if err := goview.Render(c.Writer, 200, "index", goview.M{
			"CurrentTime": time.Now().Format("3:04PM"),
			"Runtime": runtime,
			"Version": auth.Version,
			"User": auth.ClientUsername,
			"Errors": utils.ErrorCounter,
			"Requests": utils.RequestCounter,
			"Module": utils.Module,
			"RPS": utils.RateCounter.Rate(),

			"Urls": len(modules.Scraper.Urls),
			"Injectables": len(modules.Exploiter.Injectables),
			"Rows": modules.Dumper.Rows,
		}); err != nil {
			utils.LogError(err.Error())
		}
	})

	type Inc struct {
		Message string `json:"message"`
	}
	authorized.POST("/chat", func(c *gin.Context) {
		var a *Inc
		c.BindJSON(&a)

		client.XDGAPI.SendMessage(a.Message)
		c.Status(200)
	})

	authorized.GET("/chat", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		ticker := time.NewTicker(pingPeriod)
		go readLoop(conn)
		for {
			select {
			case message := <- client.XDGAPI.Incoming():
				_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				buf, err := json.Marshal(map[string]interface{}{
					"username": message.Name,
					"message": message.Message,
				})

				_, err = w.Write(buf)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	})

	authorized.GET("/ws", func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		ticker := time.NewTicker(pingPeriod)
		go readLoop(conn)
		for {
			select {
			case <- time.After(1000 * time.Millisecond):
				err = conn.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				var runtime string
				if utils.Module == "Idle" {
					runtime = "00:00:00"
				} else {
					runtime = utils.FmtDuration(time.Now().Sub(utils.StartTime))
				}

				var end int
				switch utils.Module {
				case "Idle":
					end = 0
				case "Scraper":
					end = len(modules.Scraper.Dorks)
				case "Exploiter":
					end = len(modules.Scraper.Urls)
				case "Dumper":
					fallthrough
				case "AutoSheller":
					end = len(modules.Exploiter.Injectables)
				}

				var per int
				if !(utils.GlobalIndex == 0 || end == 0) {
					per = (utils.GlobalIndex / end) * 100
				} else {
					per = 0
				}

				buf, err := json.Marshal(&StatusRefresh{
					Status:      utils.Module,
					Runtime:	 runtime,
					Threads:     len(utils.GlobalSem),
					Requests:    utils.RequestCounter,
					Errors:      utils.ErrorCounter,
					RPS:         int(utils.RateCounter.Rate()),
					Urls:        len(modules.Scraper.Urls),
					Injectables: len(modules.Exploiter.Injectables),
					Tables:      modules.Dumper.Tables,
					Columns:     modules.Dumper.Columns,
					Rows:        modules.Dumper.Rows,
					Progress:    per,
					Index:       utils.GlobalIndex,
					End:         end,
				})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				_, err = w.Write(buf)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	})

	authorized.GET("/settings", func(c *gin.Context) {
		raw, _ := ioutil.ReadFile("config.json")
		var i interface{}
		json.Unmarshal(raw, &i)
		fmt.Printf("%s\n\n%#v\n", string(raw), i)
		err := goview.Render(c.Writer, 200, "settings", i)
		if err != nil {
			utils.LogError(err.Error())
		}
	})
	
	authorized.GET("/generator", func(c *gin.Context) {
		var i interface{}
		viper.UnmarshalKey("generator", &i)
		if err := goview.Render(c.Writer, 200, "generator", i); err != nil {
			utils.LogError(err.Error())
		}
	})

	authorized.GET("/stats", func(c *gin.Context) {
		if err := goview.Render(c.Writer, 200, "stats", goview.M{
			"Version": auth.Version,
			"User": auth.ClientUsername,
			"Errors": utils.ErrorCounter,
			"Requests": utils.RequestCounter,
			"Threads": len(utils.GlobalSem),
			"Module": utils.Module,
			"RPS": utils.RateCounter.Rate(),
			"Urls": len(modules.Scraper.Urls),
			"Dorks": len(modules.Scraper.Dorks),
			"Proxies": len(manager.PManager.Proxies),
			"Injectables": len(modules.Exploiter.Injectables),
			"Tables": modules.Dumper.Tables,
			"Columns": modules.Dumper.Columns,
			"Rows": modules.Dumper.Rows,
		}); err != nil {
			utils.LogError(err.Error())
		}
	})

	authorized.GET("/profile", func(c *gin.Context) {
		err := goview.Render(c.Writer, 200, "profile", goview.M{
			"User": auth.ClientUsername,
			"Version": auth.Version,
			"Expiry": auth.Expiry.Format("Mon Jan 2 15:04:05 2006"),
		})
		if err != nil {
			utils.LogError(err.Error())
		}
	})

	authorized.POST("/settings/main", func(c *gin.Context) {
		viper.Set("core.Threads", c.PostForm("threads"))
		viper.Set("core.Timeouts", c.PostForm("timeouts"))
		at := c.PostForm("autoThreads")
		if at == "on" {
			at = "true"
		} else {
			at = "false"
		}
		//viper.Set("", at)
		td := c.PostForm("targetedDump")
		if td == "on" {
			viper.Set("dumper.Targeted", true)
		} else {
			viper.Set("dumper.Targeted", false)
		}

		bm := c.PostForm("batchMode")
		if bm == "on" {
			viper.Set("core.BatchMode", true)
		} else {
			viper.Set("core.BatchMode", false)
		}

		viper.WriteConfig()
		c.Redirect(301, "/settings")
	})

	authorized.POST("/settings/dorkers", func(c *gin.Context) {
		for _, k := range dorkers.Dorkers {
			o := c.PostForm(k.Name)
			if o == "on" {
				k.Enabled = true
			} else {
				k.Enabled = false
			}
			viper.Set(fmt.Sprintf("scraper.Engines.%s", k.Name), k.Enabled)
		}

		viper.WriteConfig()
		c.Redirect(301, "/settings")
	})

	authorized.POST("/settings/tech", func(c *gin.Context) {
		te := c.PostForm("techError")
		if te == "on" {
			viper.Set("exploiter.Techniques.Error", true)
		} else {
			viper.Set("exploiter.Techniques.Error", false)
		}
		tu := c.PostForm("techUnion")
		if tu == "on" {
			viper.Set("exploiter.Techniques.Union", true)
		} else {
			viper.Set("exploiter.Techniques.Union", false)
		}
		tb := c.PostForm("techBlind")
		if tb == "on" {
			viper.Set("exploiter.Techniques.Blind", true)
		} else {
			viper.Set("exploiter.Techniques.Blind", true)
		}

		viper.Set("exploiter.Intensity", c.PostForm("intensityLevel"))

		viper.WriteConfig()
		c.Redirect(301, "/settings")
	})

	lists := authorized.Group("lists")

	lists.GET("/", func(c *gin.Context) {
		err := goview.Render(c.Writer, 200, "lists", goview.M{
			"Version": auth.Version,
			"User": auth.ClientUsername,
			"Dorks": modules.Scraper.Dorks,
			"Proxies": manager.PManager.ProxyStringArray(),
			"Urls": modules.Scraper.Urls,
			"Injectables": modules.Exploiter.InjectableStringArray(),
		})
		if err != nil {
			utils.LogError(err.Error())
		}
	})

	lists.POST("/:name", func(c *gin.Context) {
		defer c.Redirect(301, "/lists")

		switch c.Param("name") {
		case "proxies":
			file, err := c.FormFile("proxiesFilename")
			if err != nil {
				utils.LogDebug(err.Error())
				c.String(500, `{"status":"error", "message":"%s"}`, err.Error())
				return
			}
			pType := c.PostForm("proxyType")

			i, err := strconv.Atoi(pType)
			if err != nil {
				utils.LogDebug(err.Error())
				c.String(500, `{"status":"error", "message":"%s"}`, err.Error())
				return
			}

			manager.PManager.LoadFile(file.Filename, i)
		case "dorks":
			file, err := c.FormFile("dorksFilename")
			if err != nil {
				utils.LogDebug(err.Error())
				return
			}

			f, err := file.Open()
			if err != nil {
				utils.LogDebug(err.Error())
				c.String(500, `{"status":"error", "message":"%s"}`, err.Error())
				return
			}
			defer f.Close()
			modules.Scraper.Dorks = []string{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				modules.Scraper.Dorks = append(modules.Scraper.Dorks, scanner.Text())
			}
		case "urls":
			f, err := c.FormFile("urlsFilename")
			if err != nil {
				return
			}
			file, err := f.Open()
			if err != nil {
				return
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				modules.Scraper.Urls = append(modules.Scraper.Urls, scanner.Text())
			}
			file.Close()
		case "injectables":
			f, err := c.FormFile("injectablesFilename")
			if err != nil {
				return
			}
			file, err := f.Open()
			if err != nil {
				return
			}

			scanner := bufio.NewScanner(file)
			modules.Exploiter.Injectables = make(map[string]*injection.Injection)
			for scanner.Scan() {
				i := injection.ParseUrlString(scanner.Text())
				modules.Exploiter.Injectables[i.Base.Hostname()] = i
			}
			file.Close()
		case "cust":
			fmt.Println(manager.PManager.LoadCustomerProxies(nil))
		case "clean":
			var domains []string
			var output []string
			var err error
			for _, u := range modules.Scraper.Urls {
				var h *url.URL
				if h, err = url.Parse(u); err != nil {
					continue
				}
				if !utils.StrInArr(domains, h.Hostname()) {
					domains = append(domains, h.Hostname())
				} else {
					continue
				}
				/*if blacklist == "on" {
					if utils.HasAny(h.Hostname(), blu) {
						continue
					}
				}*/

				if !(strings.Contains(u, "?") && strings.Contains(u, "=")) {
					continue
				}

				output = append(output, u)
				c.Redirect(301, "/lists")
			}
			modules.Scraper.Urls = output
			ioutil.WriteFile("urls-cleaned.txt", []byte(strings.Join(output, "\r\n")), 0644)
		}
	})

	type UpdateList struct {
		List []string `json:"list"`
		ProxyType int
	}
	
	lists.PUT("/:name", func(c *gin.Context) {
		defer c.Redirect(301, "/lists")

		var resp *UpdateList
		c.BindJSON(&resp)
		switch c.Param("name") {
		case "proxies":
			manager.PManager.Proxies = make([]*manager.Proxy, len(resp.List))
			for i, proxy := range resp.List {
				parts := strings.Split(proxy, ":")
				var p *manager.Proxy
				if len(parts) == 2 {
					p = &manager.Proxy{
						Address:   strings.Join(parts[:2], ":"),
						Type:      resp.ProxyType,
						Alive:     true,
					}
				} else {
					p = &manager.Proxy{
						Address:   strings.Join(parts[:2], ":"),
						Type:      resp.ProxyType,
						Alive:     true,
						Auth:      strings.Join(parts[2:], ":"),
					}
				}
				manager.PManager.Proxies[i] = p
			}
		case "dorks":
			modules.Scraper.Dorks = resp.List
		case "urls":
			modules.Scraper.Urls = resp.List
		case "injectables":
			modules.Exploiter.Injectables = make(map[string]*injection.Injection)
			for _, injstr := range resp.List {
				i := injection.ParseUrlString(injstr)
				modules.Exploiter.Injectables[i.Base.Hostname()] = i
			}
		}
	})

	lists.DELETE("/:name", func(c *gin.Context) {
		//defer c.String(200, `{"status":"success", "message":"Unloaded %s!"}`, c.Param("name"))
		switch c.Param("name") {
		case "proxies":
			manager.PManager.Proxies = []*manager.Proxy{}
		case "dorks":
			modules.Scraper.Dorks = []string{}
		case "urls":
			modules.Scraper.Urls = []string{}
		case "injectables":
			modules.Exploiter.Injectables = make(map[string]*injection.Injection)
		}

		c.JSON(200, map[string]interface{}{
			"status": "success",
			"message": fmt.Sprintf("Unloaded %s!", c.Param("name")),
		})
	})

	authorized.POST("/module/:modulename", func(context *gin.Context) {
		if context.Param("modulename") == "stop" {
			if utils.Module == "Stopping" {
				context.String(200, `{"status":"success", "message":"Killing..."}`)
				manager.PManager.CancelAll()
				if !utils.IsClosed(utils.Kill) {
					close(utils.Kill)
				}
			} else if utils.Module == "Idle" {
				context.String(200, `{"status":"error", "message":"Nothing is running."}`)
			} else {
				context.String(200, `{"status":"success", "message":"Stopping..."}`)
				if !utils.IsClosed(utils.Done) {
					close(utils.Done)
				}
				utils.Module = "Stopping"
			}
		} else if !utils.ArrContains([]string {"Idle", "Stopping"}, utils.Module) {
			context.String(500,  `{"status":"error", "message":, "Module '%s' already running."}"`, utils.Module)
			return
		}

		switch context.Param("modulename") {
		case "scraper":
			if len(modules.Scraper.Dorks) == 0 {
				context.String(200, `{"status":"error", "message":"Please load dorks."}`)
				return
			}
			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				utils.Module = "Scraper"
				modules.Scraper.Start()
			}()
			context.String(200, `{"status":"success", "message":"Started Scraper"}`)
			break
		case "exploiter":
			if len(modules.Scraper.Urls) == 0 {
				context.String(200, `{"status":"error", "message":"Please load or scrape urls first."}`)
				return
			}
			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				utils.Module = "Exploiter"
				modules.Exploiter.Start(modules.Scraper.Urls)
			}()
			context.String(200, `{"status":"success", "message":"Started Exploiter"}`)
			break
		case "dumper":
			if len(modules.Exploiter.Injectables) == 0 {
				context.String(200, `{"status":"error", "message":"Please load or test for injectables first."}`)
				return
			}
			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				utils.Module = "Dumper"
				var arr []*injection.Injection
				for _, inj := range modules.Exploiter.Injectables {
					arr = append(arr, inj)
				}

				modules.Dumper.Start(arr)
			}()
			context.String(200, `{"status":"success", "message":"Started Dumper"}`)
			break
		case "parse":
			break
		case "single":
			_url := context.PostForm("url")
			utils.LogDebug(_url)
			/*go func(u string) {
				inj := injection.FindInjection(u)
				if inj != nil {
					utils.LogDebug(fmt.Sprintf("%#v", inj))
					a, _ := injection.BuildInjectionRaw(inj.Vector, "[t]", inj.Method, inj.Prefix, inj.Suffix, inj.Technique, inj.UCount, inj.UInj)
					CustomUpdate <- &SingleshotRefresh{
						Message: fmt.Sprintf("Exploitable. Url: %s %s", inj.Base, a),
					}
				} else {
					CustomUpdate <- &SingleshotRefresh{
						Message: "Not Exploitable.",
					}
				}
			}(_url)*/
			context.String(200, `{"status":"success","message":"Started single shot test."}`)
		}
	})

	authorized.GET("/exit", func(c *gin.Context) {
		c.String(200, "Thanks for using!")
		server.Close()
	})

	server.Handler = r
	server.Addr = address
	err := server.ListenAndServe()
	if err != nil {
		utils.LogError(err.Error())
	}

	return server
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}