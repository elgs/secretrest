// handlers
package main

import (
	"fmt"
	"github.com/elgs/gorest"
	"net/http"
	"os"
	"strings"
)

func init() {
	gorest.RegisterHandler("/shutdown",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.RemoteAddr, "127.0.0.1:") && r.FormValue("key") == gr.SessionKey {
					defer func() {
						os.Exit(0)
					}()
				} else {
					fmt.Fprintln(w, "Attack!!!")
				}
			}
		})
	gorest.RegisterHandler("/ping",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "pong")
			}
		})
}
