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

package restful

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "fmt"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func NewIndexerHandler(controller Indexer) HandlerFunc {
  return func (out http.ResponseWriter, req *http.Request) {
    elems := controller.Index()
    if elems == nil {
      panic(fmt.Errorf("Indexer.All on %#v returned nil\n", controller))
    }

    encoder := json.NewEncoder(out)
    if err := encoder.Encode(elems); err != nil {
      panic(err)
    }
  }
}

func NewShowerHandler(controller Shower) HandlerFunc {
  return func (out http.ResponseWriter, in *http.Request) {
    id := getResourceId(in)

    elem := controller.Show(id)

    bytes, _ := json.Marshal(map[string]interface{} { id : elem })
    if _, err := out.Write(bytes); err != nil {
      panic(err)
    }
  }
}

func NewCreatorHandler(controller Creator) HandlerFunc {
  return func (out http.ResponseWriter, in *http.Request) {
    decoder := json.NewDecoder(in.Body)
    elem := controller.New()
    if elem == nil {
      panic(fmt.Errorf("Creator.New on %#v returned nil", controller))
    }
    if err := decoder.Decode(&elem); err != nil {
      panic(err)
    }

    id := controller.Create(elem)

    bytes, _ := json.Marshal(map[string]interface{} { id : elem })
    if _, err := out.Write(bytes); err != nil {
      panic(err)
    }
  }
}

func NewUpdaterHandler(cotroller Updater) HandlerFunc {
  return func (out http.ResponseWriter, in *http.Request) {
    id := getResourceId(in)

    decoder := json.NewDecoder(in.Body)
    elem := cotroller.New()
    if elem == nil {
      panic(fmt.Errorf("Creator.New on %#v returned nil", cotroller))
    }
    if err := decoder.Decode(&elem); err != nil {
      panic(err)
    }

    cotroller.Update(id, elem)

    bytes, _ := json.Marshal(map[string]interface{} { id : elem })
    if _, err := out.Write(bytes); err != nil {
      panic(err)
    }
  }
}

func NewDeleterHandler(deleter Deleter) HandlerFunc {
  return func (out http.ResponseWriter, in *http.Request) {
    id := getResourceId(in)

    ok := deleter.Delete(id)

    bytes, _ := json.Marshal(map[string]bool { "ok": ok })
    if _, err := out.Write(bytes); err != nil {
      panic(err)
    }
  }
}

func getResourceId(in *http.Request) string {
  id, ok := mux.Vars(in)["id"]
  if !ok {
    panic(fmt.Errorf("resource id not found in request data"))
  }
  return id
}

