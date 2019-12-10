package ext

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/nevata/session"
	"github.com/nevata/txtcolor"
)

//debug 日志输出文件标识
var debug = "debug_print"

//debugfile 日志文件路径
var debugfile = filepath.Join(ExeDir, debug)

//Log 记录访问日志
func Log(inner session.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			var bodyBytes []byte
			if r.Body != nil {
				// 把request的内容读取出来
				buf, err := ioutil.ReadAll(r.Body)
				if err != nil {
					HandleError(w, err)
					return
				}
				// 把刚刚读出来的再写进去
				bodyBytes = buf
				r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			}
			if len(bodyBytes) > 0 {
				buf.WriteString(fmt.Sprintf("Body = %s\n", string(bodyBytes)))
			}
			buf.WriteString("\n")
			log.Print(buf.String())
		}
		inner.ServeHTTP(nil, w, r)
	})
}
