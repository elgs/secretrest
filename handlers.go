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

	gorest.RegisterHandler("/ppp_lastasof",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("key") != gr.SessionKey {
					fmt.Fprintln(w, "Attack!!!")
					return
				}

				asof := r.FormValue("asof")
				if asof == "" {
					asof = "19010101000000"
				}
				current := `join -1 2 -2 1 <(last -F | grep 'ppp.*still' | sort -k2) <(cat /proc/net/dev | grep ppp | awk -F ':[ \t]*|[ \t]+' '{print $2,$3,$11}' | sort)`
				history := `lastasof -F -t ` + asof + ` | egrep -v "crash|down|still|gone" | grep ppp`
				command := fmt.Sprint(history + " ; " + current)
				output, err := exec.Command("bash", "-c", command).CombinedOutput()
				if err != nil {
					fmt.Println("Failed to execute:", err, command)
					fmt.Println(string(output))
				} else {
					fmt.Fprint(w, string(output))
				}
			}
		})

	gorest.RegisterHandler("/kill_ppp",
		func(dbo gorest.DataOperator, gr *gorest.Gorest) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("key") != gr.SessionKey {
					fmt.Fprintln(w, "Attack!!!")
					return
				}

				ppp := r.FormValue("ppp")
				if ppp == "" {
					return
				}
				command := fmt.Sprint("kill -TERM `cat /var/run/" + ppp + ".pid`")
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
