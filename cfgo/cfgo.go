package cfgo

import (
	"syscall/js"
	"encoding/json"
	"path"
	"os"
	"github.com/fikisipi/cloudflare-workers-go/cfgo/structs"
)

type Request struct {
	Body string
	Headers map[string]string
	QueryParams map[string]string
	URL string
	Hostname string
	Pathname string
	Method string
}

type Response = structs.Response

// Used to chain .SetStatus(), .AddHeader(), etc.
// together when creating a Response. Example:
//   return BuildResponse().SetStatus(200).SetBody("Hello").Build()
// The final .Build() is not mandatory, it's just for
// reducing ambiguity about the return type.
func BuildResponse() Response {
	reply := new(structs.RawResponse)

	reply.StatusCode = 200
	reply.Headers = make(map[string]string)
	return reply
}

type FetchBody interface {
	Get() js.Value
}

func BodyString(body string) FetchBody {
	return &structs.StringBody{body}
}

func BodyForm(body map[string]string) FetchBody {
	return &structs.FormBody{body}
}

// Fetches any URL using the fetch() Web API. Unlike browsers,
// CloudFlare workers miss some features like credentials.
//
// If you don't need headers or a request body, set them to nil:
//    Fetch(myUrl, "GET", nil, nil)
//
// To create a POST/PUT body, use cfgo.BodyForm() or cfgo.BodyString():
//    Fetch(myURL, "PUT", nil, cfgo.BodyForm(myDict))
func Fetch(url string, method string, headers map[string]string, requestBody FetchBody) string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headersJs := structs.CreateJsMap(headers)

	bodyJs := js.Null()
	if requestBody != nil {
		bodyJs = requestBody.Get()
	}

	out := make(chan string)
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		out <- args[0].String()
		cb.Release()
		return 1
	})
	js.Global().Call("_cfFetch", url, method, headersJs, bodyJs, cb)
	return (<- out)
}

type Handler struct {
	callbacks map[string]Callback
}

type Callback func(*Request) Response

func (h *Handler) Add(s string, fn Callback) {
	if(h.callbacks == nil) {
		h.callbacks = make(map[string]Callback)
	}
	h.callbacks[s] = fn
}

func (h *Handler) Run() {
	responseCallback := js.Global().Call("_getCallback")
	if len(os.Args) != 2 {
		println("ERROR: subscribe() must be called with one arg")
		return
	}
	jsonRequest := os.Args[1]

	var request = new(Request)
	err := json.Unmarshal([]byte(jsonRequest), request)
	if err != nil {
		return
	}

	var response = BuildResponse()

	for pathStr, pathHandler := range h.callbacks {
		if matched, _ := path.Match(pathStr, request.Pathname); matched {
			response = pathHandler(request)
		}
	}

	responseBytes, err := json.Marshal(response)
	responseStr := string(responseBytes)
	result := responseStr

	responseCallback.Invoke(result)
}

var Router = Handler{}