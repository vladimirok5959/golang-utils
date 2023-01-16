# golang-utils

Different kind of functions for second reusage

## utils/http/apiserv

Web server, which can use regular expressions in path. Designed to be simple for projects with API

```go
package base

import (
    "net/http"
    ...
)

type Handler struct {
    Ctx      context.Context
    DB       *database.DataBase
    Shutdown context.CancelFunc
}

type ServerData struct {
    ...
}

func (h Handler) FuncMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
    return template.FuncMap{}
}
```

```go
package page_index

import (
    "net/http"
    ...
)

type Handler struct {
    base.Handler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    sess := h.SessionStart(w, r)
    if sess == nil {
        return
    }
    defer sess.Close()

    // /api/v1/aliases/{i}
    // apiserv.GetParams(r)[1].Integer()

    // /api/v1/aliases/{s}
    // apiserv.GetParams(r)[1].String()

    ...

    data := &base.ServerData{}

    if !render.HTML(w, r, h.FuncMap(w, r), data, web.IndexHtml, http.StatusOK) {
        return
    }
}
```

```go
package server

import (
    "net/http"
    ...
)

func NewMux(ctx context.Context, shutdown context.CancelFunc, db *database.DataBase) *apiserv.ServeMux {
    mux := apiserv.NewServeMux()

    handler := base.Handler{
        Ctx:      ctx,
        DB:       db,
        Shutdown: shutdown,
    }

    // Pages
    mux.Get("/", page_index.Handler{Handler: handler})
    mux.Get("/about", page_about.Handler{Handler: handler})

    // API
    mux.Get("/api/v1/app/health", v1_app_health.Handler{Handler: handler})
    mux.Handle("/api/v1/aliases", []string{http.MethodGet, http.MethodPost}, v1_aliases.Handler{Handler: handler})  
    mux.Handle("/api/v1/aliases/{i}", []string{http.MethodGet, http.MethodPut, http.MethodDelete}, v1_aliases.Handler{Handler: handler})

    // Assets
    mux.Get("/favicon.png", helpers.HandleImagePng(web.FaviconPng))
    mux.Get("/robots.txt", helpers.HandleTextPlain(web.RobotsTxt))
    mux.Get("/sitemap.xml", helpers.HandleTextXml(web.SitemapXml))

    // 404
    mux.NotFound(page_404.Handler{Handler: handler})

    return mux
}

func New(ctx context.Context, shutdown context.CancelFunc, db *database.DataBase) (*http.Server, error) {
    mux := NewMux(ctx, shutdown, db)
    srv := &http.Server{
    Addr:   consts.Config.Host + ":" + consts.Config.Port,
            Handler: mux,
    }
    go func() {
        fmt.Printf("Web server: http://%s:%s/\n", consts.Config.Host, consts.Config.Port)
        if err := srv.ListenAndServe(); err != nil {
            if err != http.ErrServerClosed {
                fmt.Printf("Web server startup error: %s\n", err.Error())
                shutdown()
                return
            }
        }
    }()
    return srv, nil
}
```

## utils/http/helpers

```go
type CurlGetOpts struct {
    ExpectStatusCode int
    Headers          map[string][]string
    Timeout          time.Duration
}

func CurlDownload(ctx context.Context, url string, opts *CurlGetOpts, fileName string, filePath ...string) error
func CurlGet(ctx context.Context, url string, opts *CurlGetOpts) ([]byte, error)
```

```go
func ClientIP(r *http.Request) string
func ClientIPs(r *http.Request) []string
func HandleAppStatus() http.Handler
func HandleFile(data, contentType string) http.Handler
func HandleImageJpeg(data string) http.Handler
func HandleImagePng(data string) http.Handler
func HandleTextCss(data string) http.Handler
func HandleTextJavaScript(data string) http.Handler
func HandleTextPlain(data string) http.Handler
func HandleTextXml(data string) http.Handler
func IntToStr(value int64) string
func Md5Hash(str []byte) string
func MinifyHtmlCode(str string) string
func MinifyHtmlJsCode(str string) string
func RespondAsBadRequest(w http.ResponseWriter, r *http.Request, err error)
func RespondAsInternalServerError(w http.ResponseWriter, r *http.Request)
func RespondAsMethodNotAllowed(w http.ResponseWriter, r *http.Request)
func SessionStart(w http.ResponseWriter, r *http.Request) (*session.Session, error)
func SetLanguageCookie(w http.ResponseWriter, r *http.Request) error
```

## utils/http/logger

Used for write logs, `LogRequests` can be used in handlers chain

```go
func LogError(format string, a ...any)
func LogInfo(format string, a ...any)
func LogInternalError(err error)
func LogRequests(handler http.Handler) http.Handler

// Write logs to file:
func main() {
    logger.AccessLogFile = consts.Config.AccessLogFile
    logger.ErrorLogFile = consts.Config.ErrorLogFile
    ...
}
```

RollBar supports from box:

```go
func main() {
    logger.RollBarEnabled = (consts.Config.Deployment == "production")

    logger.RollBarSkipErrors = []string{
        "context canceled",
    }

    rollbar.SetToken("token")
    rollbar.SetEnvironment(consts.Config.Deployment)
    rollbar.SetCodeVersion(consts.AppVersion)

    ...

    rollbar.Wait()
}
```

## utils/http/redirect

HTTP request redirect helper:

```go
func redirect.Handler(url string) http.Handler
```

```go
// Example:
mux.Get("/old", redirect.Handler("/new"))
```

## utils/http/render

Used for rendering HTML pages and for JSON API responses

```go
func HTML(w http.ResponseWriter, r *http.Request, f template.FuncMap, d interface{}, s string, httpStatusCode int) bool
func JSON(w http.ResponseWriter, r *http.Request, o interface{}) bool
```

## utils/http/servauth

Used for adding basic authentication for any HTTP handler

```go
func BasicAuth(handler http.Handler, username, password, realm string) http.Handler
```

## utils/http/servlimit

Used for limiting count of requests per second for client IP

```go
func ReqPerSecond(handler http.Handler, requests int) http.Handler
```

## utils/http/servtools

HTTP handler, universal template and renderer for config

```go
func ConfigVars(config any) http.Handler
```

```go
package consts

var Config struct {
    Host string `default:"127.0.0.1" description:"Web server IP"`
    Port string `default:"8080" description:"Web server port"`
}

// Example:
mux.Get("/config", servtools.ConfigVars(&consts.Config))
```

## utils/pagination

Pagination helper for build SQL request and HTML template

```go
func New(currentPage, resultsCount, resultsPerPage int64) *Data
```

```go
func (d *Data) CurrentPage() int64
func (d *Data) MaxPages() int64
func (d *Data) ResultsCount() int64
func (d *Data) ResultsPerPage() int64
func (d *Data) Limit() (int64, int64)
```

## utils/penv

Auto-config for parameters and ENV variables

```go
var Config struct {
    Host string `default:"127.0.0.1" description:"Web server IP"`
    Port string `default:"8080" description:"Web server port"`
}

// Example:
func init() {
    if err := penv.ProcessConfig(&Config); err != nil {
        panic(err)
    }
}

func main() {
    fmt.Printf("%s\n", Config.Host)
    fmt.Printf("%s\n", Config.Port)
}

// Auto generated params and ENVs:
// Params will be readed: host, port
// ENV will be readed too: ENV_HOST, ENV_PORT
```

## utils

```go
func RandomString(length int) string
```
