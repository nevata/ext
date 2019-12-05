package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/nevata/session"
	"github.com/nevata/txtcolor"
)

var debug = "debug_print"
var debugfile string = filepath.Join(ExeDir, debug)

//Log 记录访问日志
func Log(inner session.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//解析json
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			//NZ_RPC不标准，发送json数据时，没有带这个标记过来
			contentType := r.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") || contentType == "" {
				params := make(map[string]string)
				if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
					if err != io.EOF {
						HandleError(w, err)
						return
					}
				}

				if r.Form == nil {
					r.Form = make(url.Values)
				}

				for k, v := range params {
					r.Form.Set(k, v)
				}
			}
		}

		//输出log
		if FileExist(debugfile) {
			buf := bytes.NewBuffer(nil)
			buf.WriteString(fmt.Sprintln(txtcolor.Blue(r.Method), r.URL, r.Proto))
			//log.Println(txtcolor.Blue(r.Method), r.URL, r.Proto)
			for k, v := range r.Header {
				//log.Printf("Header[%q] = %q\n", k, v)
				buf.WriteString(fmt.Sprintf("Header[%q] = %q\n", k, v))
			}

			//log.Printf("Host = %q\n", r.Host)
			//log.Printf("RemoteAddr = %q\n", r.RemoteAddr)
			buf.WriteString(fmt.Sprintf("Host = %q\n", r.Host))
			buf.WriteString(fmt.Sprintf("RemoteAddr = %q\n", r.RemoteAddr))
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}

			for k, v := range r.Form {
				//log.Printf("Form[%q] = %q\n", k, v)
				buf.WriteString(fmt.Sprintf("Form[%q] = %q\n", k, v))
			}
			buf.WriteString("\n")
			log.Print(buf.String())
		}
		inner.ServeHTTP(nil, w, r)
	})
}
