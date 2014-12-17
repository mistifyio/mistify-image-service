package imageservice

import (
	"net/http"
	"github.com/bakins/net-http-recover"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
    "github.com/mistifyio/mistify-agent/log"
	"os"
    "strings"
    "runtime"
    "fmt"
    "encoding/json"
)

type (
    
    HttpRequest struct {
		ResponseWriter http.ResponseWriter
		Request        *http.Request
		Context        *Context
		vars           map[string]string
	}

	HttpErrorMessage struct {
		Message string   `json:"message"`
		Code    int      `json:"code"`
		Stack   []string `json:"stack"`
	}

    Chain struct {
        alice.Chain
        ctx *Context
    }
)

func Run(ctx *Context, address string) error {
	router := mux.NewRouter()
	router.StrictSlash(true)

	chain := Chain {
		ctx: ctx,
		Chain: alice.New(
			func(h http.Handler) http.Handler {
				return handlers.CombinedLoggingHandler(os.Stdout, h)
			},
			handlers.CompressHandler,
			func(h http.Handler) http.Handler {
				return recovery.Handler(os.Stderr, h, true)
			}),
	}

	router.HandleFunc("/images", chain.RequestWrapper(listImages)).Methods("GET")
	router.HandleFunc("/images", chain.RequestWrapper(putImage)).Methods("POST")
	router.HandleFunc("/images/{id}", chain.RequestWrapper(getImage)).Methods("GET")
	router.HandleFunc("/images/{id}", chain.RequestWrapper(deleteImage)).Methods("DELETE")

	server := &http.Server{
		Addr: address,
		Handler: router,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}

func (c *Chain) RequestWrapper(fn func(*HttpRequest) *HttpErrorMessage) http.HandlerFunc {
	return c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := HttpRequest{
			Context:        c.ctx,
			ResponseWriter: w,
			Request:        r,
		}
		if err := fn(&req); err != nil {
			log.Error("%s\n\t%s\n", err.Message, strings.Join(err.Stack, "\t\n\t"))
			req.JSON(err.Code, err)
		}
	})).ServeHTTP
}

func (r *HttpRequest) SetHeader(key, val string) {
	r.ResponseWriter.Header().Set(key, val)
}

func (r *HttpRequest) NewError(err error, code int) *HttpErrorMessage {
	if code <= 0 {
		code = 500
	}
	msg := HttpErrorMessage{
		Message: err.Error(),
		Code:    code,
		Stack:   make([]string, 0, 4),
	}
	for i := 1; ; i++ { //
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		msg.Stack = append(msg.Stack, fmt.Sprintf("%s:%d (0x%x)", file, line, pc))
	}
	return &msg
}

func (r *HttpRequest) JSON(code int, obj interface{}) *HttpErrorMessage {
	r.SetHeader("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(code)
	encoder := json.NewEncoder(r.ResponseWriter)
	if err := encoder.Encode(obj); err != nil {
		return r.NewError(err, 500)
	}
	return nil
}

func listImages(req *HttpRequest) *HttpErrorMessage {
    return nil
}

func putImage(req *HttpRequest) *HttpErrorMessage {
    return nil
}

func getImage(req *HttpRequest) *HttpErrorMessage {
    return nil
}

func deleteImage(req *HttpRequest) *HttpErrorMessage {
    return nil
}

