package apiserv

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/vladimirok5959/golang-utils/utils/http/logger"
)

var mParam = regexp.MustCompile(`\{([^/]*)}`)
var TestingMockParams func() []Param = nil

var mParams = &Params{
	list: map[*http.Request][]Param{},
}

type ServeMux struct {
	handlers *Handlers
}

func NewServeMux() *ServeMux {
	s := ServeMux{
		handlers: &Handlers{
			list: map[*regexp.Regexp]Handler{},
		},
	}
	return &s
}

func GetParams(r *http.Request) []Param {
	if TestingMockParams != nil {
		return TestingMockParams()
	}
	mParams.Lock()
	defer mParams.Unlock()
	if v, ok := mParams.list[r]; ok {
		return v
	}
	return []Param{}
}

func (s ServeMux) regexp(pattern string) string {
	pattern = strings.ReplaceAll(pattern, "-", "\\-")
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "/", "\\/")
	pattern = mParam.ReplaceAllStringFunc(pattern, func(str string) string {
		if str == "{i}" {
			return "([\\d]+)"
		} else if str == "{s}" {
			return "([^\\/]+)"
		}
		return str
	})
	return "(?i)^" + pattern + "$"
}

func (s ServeMux) Delete(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodDelete}, handler)
}

func (s ServeMux) Get(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodGet}, handler)
}

func (s ServeMux) Handle(pattern string, methods []string, handler http.Handler) {
	s.handlers.Lock()
	defer s.handlers.Unlock()
	if pattern != "" {
		re := regexp.MustCompile(s.regexp(pattern))
		s.handlers.list[re] = Handler{
			handler: handler,
			methods: methods,
		}
	} else {
		s.handlers.list[nil] = Handler{
			handler: handler,
			methods: methods,
		}
	}
}

func (s ServeMux) NotFound(handler http.Handler) {
	s.Handle("", []string{http.MethodGet}, handler)
}

func (s ServeMux) Options(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodOptions}, handler)
}

func (s ServeMux) Patch(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodPatch}, handler)
}

func (s ServeMux) Post(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodPost}, handler)
}

func (s ServeMux) Put(pattern string, handler http.Handler) {
	s.Handle(pattern, []string{http.MethodPut}, handler)
}

func (s ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handlers.Lock()
	defer s.handlers.Unlock()
	for re, h := range s.handlers.list {
		if re != nil && re.Match([]byte(r.URL.Path)) {
			if h.IsMethod(r.Method) {
				if rs := re.FindAllStringSubmatch(r.URL.Path, 1); len(rs) >= 1 {
					defer func() {
						mParams.Lock()
						delete(mParams.list, r)
						defer mParams.Unlock()
					}()
					mParams.Lock()
					for _, p := range rs[0] {
						if _, ok := mParams.list[r]; !ok {
							mParams.list[r] = []Param{{value: p}}
						} else {
							mParams.list[r] = append(mParams.list[r], Param{value: p})
						}
						select {
						case <-r.Context().Done():
							mParams.Unlock()
							return
						default:
						}
					}
					mParams.Unlock()
				}
				select {
				case <-r.Context().Done():
					return
				default:
				}
				logger.LogRequests(h.handler).ServeHTTP(w, r)
				return
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		}
		select {
		case <-r.Context().Done():
			return
		default:
		}
	}

	// Error 404
	if h, ok := s.handlers.list[nil]; ok {
		logger.LogRequests(h.handler).ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
