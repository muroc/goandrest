/*
   Copyright 2014 Maciej Chałapuk

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

/*
  Provides:
   * HTTP Router capable of registering controller
     objects as handlers for calls to RESTful resources,
   * Interfaces are intended to be implemented
     in controllers in client code (that uses this library).
*/
package restful

import (
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

/*
  Routes HTTP requests to controllers registered as RESTful resource handlers.
  Main structure of the library.
*/
type Router struct {
  impl *mux.Router
}
/*
  Constructs new router.
*/
func NewRouter() *Router {
  return &Router{ mux.NewRouter() }
}

/*
  Registers given controller as handler for resource under given path.
  Path should start with "/". Controller should implement at least one
  of interfaces defined in controller.go file (to make all RESTful
  operations available it must implement all of them).
*/
func (router *Router) HandleResource(path string, controller interface{}) {
  api := router.impl.PathPrefix(path).Subrouter()
  registered := false

  if indexer, ok := controller.(Indexer); ok {
    api.Path("/").Methods("GET").HandlerFunc(newIndexerHandler(indexer))
    registered = true
  }
  if shower, ok := controller.(Shower); ok {
    api.Path("/{id}").Methods("GET").HandlerFunc(newShowerHandler(shower))
    registered = true
  }
  if creator, ok := controller.(Creator); ok {
    api.Path("/").Methods("POST").HandlerFunc(newCreatorHandler(creator))
    registered = true
  }
  if updater, ok := controller.(Updater); ok {
    api.Path("/{id}").Methods("PUT").HandlerFunc(newUpdaterHandler(updater))
    registered = true
  }
  if deleter, ok := controller.(Deleter); ok {
    api.Path("/{id}").Methods("DELETE").HandlerFunc(newDeleterHandler(deleter))
    registered = true
  }

  if !registered {
    fmt.Errorf("%#v doesn't implement any REST interfaces", controller)
  }
}

/*
  Called by HTTP server when request is received.
*/
func (router *Router) ServeHTTP(out http.ResponseWriter, in *http.Request) {
  router.impl.ServeHTTP(out, in)
}

