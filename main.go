package main

import (
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

type Server struct {
	middlewares []Middleware //middleware
	routes      []route      //routes
	addr        string       //address
	l           net.Listener //save the listener so it can be closed.
}

type route struct {
	r       string        //route url
	method  string        //method
	handler reflect.Value //handle func
}

//Middleware will be called before each request.
type Middleware interface {
	Handle(*Context)
}

//TODO:Use add a new Middleware that implements Handle(*Context)
func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

//TODO:Next calls next middleware.
func (ctx *Context) Next() {
	ctx.idx += 1
}

//TODO:Invok calls middleware at index of ctx.idx.
func (ctx *Context) Invok() {
	ctx.middlewares[ctx.idx].Handle(ctx)
}

//implements http.Handle
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, r := range s.routes {
		if r.r == req.URL.Path && r.method == req.Method {
			//function handler
			//*context must be the first argument.
			ctx := &Context{req, res, s, make(map[string]string), s.middlewares, 0}
			//call the middlewares
			ctx.Invok()

			var args []reflect.Value
			if requiresContext(r.handler.Type()) {
				args = append(args, reflect.ValueOf(ctx))
			}
			ret := r.handler.Call(args)
			if len(ret) == 0 {
				return
			}
			//if has return value,write to response.
			sval := ret[0]
			var content []byte
			if sval.Kind() == reflect.String {
				content = []byte(sval.String())
			} else if sval.Kind() == reflect.Slice && sval.Type().Elem().Kind() == reflect.Uint8 {
				content = sval.Interface().([]byte)
			}
			ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
			ctx.ResponseWriter.Write(content)
		}
	}
}

//SetHeader sets a response header. If `unique` is true, the current value
//of that header will be overwritten . If false, it will be appended.
func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}

//requiresContext determines whether 'handlerType' contains
//an argument to 'web.Ctx' as its first argument
func requiresContext(handlerType reflect.Type) bool {
	//if the method doesn't take arguments, no
	if handlerType.NumIn() == 0 {
		return false
	}

	//if the first argument is not a pointer, no
	a0 := handlerType.In(0)
	if a0.Kind() != reflect.Ptr {
		return false
	}
	//if the first argument is a context, yes
	if a0.Elem() == contextType {
		return true
	}

	return false
}

//Close closes the server
func (s *Server) Close() {
	if s.l != nil {
		s.l.Close()
	}
}

//Run runs the server
func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.Handle("/", s)
	log.Printf("start serverving...\nPlease visit http://localhost:3000")
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	s.l = l
	err = http.Serve(s.l, mux)
	s.l.Close()
}

//Get adds a handler for the 'GET' http method for server.
func (s *Server) Get(rt string, handler interface{}) {
	s.addRoute(rt, "GET", handler)
}

//Post adds a handler for the 'GET' http method for server.
func (s *Server) Post(rt string, handler interface{}) {
	s.addRoute(rt, "POST", handler)
}
func (s *Server) addRoute(rt string, method string, handler interface{}) {
	switch handler.(type) {
	case reflect.Value:
		fv := handler.(reflect.Value)
		s.routes = append(s.routes, route{r: rt, method: method, handler: fv})
	default:
		fv := reflect.ValueOf(handler)
		s.routes = append(s.routes, route{r: rt, method: method, handler: fv})
	}

}

//Server returns a new Server.
func NewServer() *Server {
	return &Server{addr: "localhost:3000"}
}

//provide context for each request,and bebore handling,run the Middlewares iterately.
type Context struct {
	*http.Request
	http.ResponseWriter
	Server      *Server
	Params      map[string]string
	middlewares []Middleware //middleware
	idx         int          //index of middleware
}

var contextType reflect.Type

func init() {
	contextType = reflect.TypeOf(Context{})
}

//To add a middleware,we just need to implement Handle(*Context) and call the ctx.Next(),if we want to continue the middleware layer.
type ParseForm struct {
}

func (p *ParseForm) Handle(ctx *Context) {
	ctx.ParseForm()
	ctx.Next()
}

func main() {
	println(`Now we have a context for each request.
To decouple the request handling,middleware is very useful,each middleware just deals with part of the handling,
and pass the control to the next,middleware is pluggable.
That means you can just add middlewares you want to use to handle the request.
Suppose you have done the task l6,and have a context,now we need to add a middleware layer to tiny webframework,and make it pluggable.
Edit main.go,and finish the task.Notice something new added.
`)
}
