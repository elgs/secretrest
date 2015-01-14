// handlers
package main

import (
	"fmt"
	"github.com/elgs/gorest"
	"net/http"
	"os"
	"os/exec"
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

	gorest.RegisterHandler("/last",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("key") != gr.SessionKey {
					fmt.Fprintln(w, "Attack!!!")
					return
				}

				command := fmt.Sprint("last -F | grep ppp")
				output, err := exec.Command("bash", "-c", command).CombinedOutput()
				if err != nil {
					fmt.Println("Failed to execute:", err, command)
					fmt.Println(string(output))
				} else {
					fmt.Fprint(w, string(output))
				}
			}
		})

	gorest.RegisterHandler("/lastasof",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("key") != gr.SessionKey {
					fmt.Fprintln(w, "Attack!!!")
					return
				}

				asof := r.FormValue("asof")
				command := fmt.Sprint("lastasof -F -t " + asof + " | grep ppp")
				fmt.Println(command)
				output, err := exec.Command("bash", "-c", command).CombinedOutput()
				if err != nil {
					fmt.Println("Failed to execute:", err, command)
					fmt.Println(string(output))
				} else {
					fmt.Fprint(w, string(output))
				}
			}
		})
}
