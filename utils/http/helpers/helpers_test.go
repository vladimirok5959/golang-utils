package helpers_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var _ = Describe("helpers", func() {
	Context("CurlGetStatusError", func() {
		It("recognize error", func() {
			err := error(&helpers.CurlGetStatusError{Expected: http.StatusOK, Received: http.StatusServiceUnavailable})
			Expect(errors.Is(err, helpers.ErrCurlGetStatus)).To(BeTrue())

			err = error(&helpers.CurlGetStatusError{Expected: http.StatusOK, Received: http.StatusBadGateway})
			Expect(errors.Is(err, helpers.ErrCurlGetStatus)).To(BeTrue())

			err = fmt.Errorf("Some error")
			Expect(errors.Is(err, helpers.ErrCurlGetStatus)).To(BeFalse())

			Expect(errors.Is(fs.ErrNotExist, helpers.ErrCurlGetStatus)).To(BeFalse())
		})

		It("generate error message", func() {
			err := error(&helpers.CurlGetStatusError{Expected: http.StatusOK, Received: http.StatusBadGateway})
			Expect(err.Error()).To(Equal("CurlGet: expected 200, received 502"))
		})
	})

	Context("ClientIP", func() {
		It("return client IP", func() {
			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "127.0.0.1",
			})).To(Equal("127.0.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.1,127.0.0.1",
			})).To(Equal("192.168.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.1, 127.0.0.1",
			})).To(Equal("192.168.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.50,192.168.0.1,127.0.0.1",
			})).To(Equal("192.168.0.50"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.50, 192.168.0.1, 127.0.0.1",
			})).To(Equal("192.168.0.50"))
		})
	})

	Context("ClientIPs", func() {
		It("return array of client IPs", func() {
			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "127.0.0.1",
			})).To(ConsistOf(
				"127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.1,127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.1, 127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.50,192.168.0.1,127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.50", "192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.50, 192.168.0.1, 127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.50", "192.168.0.1", "127.0.0.1",
			))
		})

		It("return array of client IPs without port", func() {
			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "[::1]:544722",
			})).To(ConsistOf(
				"[::1]",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "127.0.0.1:8080",
			})).To(ConsistOf(
				"127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.1:80,127.0.0.1:443",
			})).To(ConsistOf(
				"192.168.0.1", "127.0.0.1",
			))
		})
	})

	Context("Handles", func() {
		var srv *httptest.Server
		var client *http.Client
		var resp *http.Response
		var err error

		Context("HandleAppStatus", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleAppStatus())
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle app status", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(MatchRegexp(`{"memory":{"alloc":[0-9]+,"num_gc":[0-9]+,"sys":[0-9]+,"total_alloc":[0-9]+},"routines":[0-9]+}`))
			})
		})

		Context("HandleFile", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleFile("MyContent", "my/type"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle file", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("my/type"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleImageGif", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleImageGif("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle image gif", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("image/gif"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleImageJpeg", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleImageJpeg("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle image jpeg", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("image/jpeg"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleImagePng", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleImagePng("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle image png", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("image/png"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleTextCss", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleTextCss("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle text css", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("text/css"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleTextJavaScript", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleTextJavaScript("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle text javascript", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("text/javascript"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("HandleTextPlain", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleTextPlain("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle text plain", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("text/plain"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("InArrayInt", func() {
			Expect(helpers.InArrayInt(3, []int{1, 2, 3, 4, 5})).To(BeTrue())
			Expect(helpers.InArrayInt(9, []int{1, 2, 3, 4, 5})).To(BeFalse())
		})

		Context("InArrayInt64", func() {
			Expect(helpers.InArrayInt64(3, []int64{1, 2, 3, 4, 5})).To(BeTrue())
			Expect(helpers.InArrayInt64(9, []int64{1, 2, 3, 4, 5})).To(BeFalse())
		})

		Context("InArrayStr", func() {
			Expect(helpers.InArrayStr("3", []string{"1", "2", "3"})).To(BeTrue())
			Expect(helpers.InArrayStr("9", []string{"1", "2", "3"})).To(BeFalse())
		})

		Context("HandleTextXml", func() {
			BeforeEach(func() {
				srv = httptest.NewServer(helpers.HandleTextXml("MyContent"))
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle text xml", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("Content-Type")).To(Equal("text/xml"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(Equal("MyContent"))
			})
		})

		Context("IntToStr", func() {
			It("convert int to string", func() {
				Expect(helpers.IntToStr(1)).To(Equal("1"))
				Expect(helpers.IntToStr(5)).To(Equal("5"))
				Expect(helpers.IntToStr(50)).To(Equal("50"))
				Expect(helpers.IntToStr(100)).To(Equal("100"))
				Expect(helpers.IntToStr(1000)).To(Equal("1000"))
			})
		})

		Context("IntToStr64", func() {
			It("convert int to string", func() {
				Expect(helpers.IntToStr64(1)).To(Equal("1"))
				Expect(helpers.IntToStr64(5)).To(Equal("5"))
				Expect(helpers.IntToStr64(50)).To(Equal("50"))
				Expect(helpers.IntToStr64(100)).To(Equal("100"))
				Expect(helpers.IntToStr64(1000)).To(Equal("1000"))
			})
		})

		Context("RespondAsBadRequest", func() {
			BeforeEach(func() {
				var handler = func() http.HandlerFunc {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						helpers.RespondAsBadRequest(w, r, fmt.Errorf("MyError"))
					})
				}

				srv = httptest.NewServer(handler())
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle bad request", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(string(body)).To(MatchRegexp(`{"error":"MyError"}`))
			})
		})

		Context("RespondAsInternalServerError", func() {
			BeforeEach(func() {
				var handler = func() http.HandlerFunc {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						helpers.RespondAsInternalServerError(w, r)
					})
				}

				srv = httptest.NewServer(handler())
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle method not allowed", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("RespondAsMethodNotAllowed", func() {
			BeforeEach(func() {
				var handler = func() http.HandlerFunc {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						helpers.RespondAsMethodNotAllowed(w, r)
					})
				}

				srv = httptest.NewServer(handler())
				client = srv.Client()
				resp, err = client.Get(srv.URL + "/")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				Expect(resp.Body.Close()).To(Succeed())
				srv.Close()
			})

			It("handle method not allowed", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

	Context("MinifyHtmlCode", func() {
		It("minify Html code", func() {
			Expect(helpers.MinifyHtmlCode(`
				<!doctype html>
				<html lang="uk">
					<head>
						<meta charset="utf-8" />
						<script>
							// Comment
							var LangTexts = {
								"empty": "Empty",
								"full": "Full"
							};

							/* Comment */

							/*
								Comment line 1
								Comment line 2
							*/
							function Func(value) {
								console.log(value);
							};
						</script>
					</head>
					<body>
						Index
					</body>
				</html>
			`)).To(Equal(`<!doctype html><html lang="uk"><head><meta charset="utf-8" /><script>var LangTexts={"empty":"Empty","full":"Full"};function Func(value){console.log(value);};</script></head><body>Index</body></html>`))

			Expect(helpers.MinifyHtmlCode(`
				<div>
					<a href="#">Link 1</a>, <a href="#">Link 2</a>
				</div>
			`)).To(Equal(`<div><a href="#">Link 1</a>, <a href="#">Link 2</a></div>`))

			Expect(helpers.MinifyHtmlCode(`
				<div>
					<b>Contacts:</b> <a href="#">Link 1</a>, <a href="#">Link 2</a>
				</div>
			`)).To(Equal(`<div><b>Contacts:</b> <a href="#">Link 1</a>, <a href="#">Link 2</a></div>`))
		})
	})

	Context("StrToInt", func() {
		It("convert int to string", func() {
			Expect(helpers.StrToInt("")).To(Equal(0))
			Expect(helpers.StrToInt("1")).To(Equal(1))
			Expect(helpers.StrToInt("5")).To(Equal(5))
			Expect(helpers.StrToInt("50")).To(Equal(50))
			Expect(helpers.StrToInt("100")).To(Equal(100))
			Expect(helpers.StrToInt("1000")).To(Equal(1000))
		})
	})

	Context("StrToInt64", func() {
		It("convert int to string", func() {
			Expect(helpers.StrToInt64("")).To(Equal(int64(0)))
			Expect(helpers.StrToInt64("1")).To(Equal(int64(1)))
			Expect(helpers.StrToInt64("5")).To(Equal(int64(5)))
			Expect(helpers.StrToInt64("50")).To(Equal(int64(50)))
			Expect(helpers.StrToInt64("100")).To(Equal(int64(100)))
			Expect(helpers.StrToInt64("1000")).To(Equal(int64(1000)))
		})
	})

	Context("FakeResponseWriter", func() {
		It("write data to fake response writer", func() {
			var someHandleFunc = func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("body"))
			}

			writer := helpers.NewFakeResponseWriter()
			someHandleFunc(writer)

			Expect(writer.Body).To(Equal([]byte("body")))
			Expect(writer.Headers).To(Equal(http.Header{
				"Content-Type": []string{"application/json"},
			}))
			Expect(writer.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "helpers")
}
