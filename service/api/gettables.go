package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	// "io/ioutil"
	// "mime"
	// "time"
)

func (rt *_router) getTables(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	err := rt.db.GetTables()
	if err != nil {
		// fmt.Println(err)
	}
}
