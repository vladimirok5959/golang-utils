# golang-utils

Different kind of functions for second reusage

## utils/http/apiserv

Web server, which can use regular expressions in path. Designed to be simple and for usage for project with API and etc

```go
package server

import (
    "context"
    "fmt"
    "net/http"
    ...
)

// External package
type base.Handler struct {
    Ctx      context.Context
    DB       *database.DataBase
    Shutdown context.CancelFunc
}

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
