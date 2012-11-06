package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"
)

var htmlRawTemplate = `<!DOCTYPE html>
<html>
	<head> <title>Proxy error</title> </head>
	<body>
		<h1>Proxy Error: {{.H1}}</h1>
		{{.Msg}}
		<br /> <br />
		<hr />
		Generated by <i>cow-proxy</i> at {{.T}}
	</body>
</html>
`

// Do not end with "\r\n" so we can add more header later
var headRawTemplate = "HTTP/1.1 {{.Code}} {{.Reason}}\r\n" +
	"Connection: keep-alive\r\n" +
	"Cache-Control: no-cache\r\n" +
	"Pragma: no-cache\r\n" +
	"Content-Type: text/html\r\n" +
	"Content-Length: {{.Length}}\r\n"

var htmlTmpl, headTmpl *template.Template

func init() {
	var err error
	if headTmpl, err = template.New("errorHead").Parse(headRawTemplate); err != nil {
		fmt.Println("Internal error on generating error head template")
		os.Exit(1)
	}
	if htmlTmpl, err = template.New("errorPage").Parse(htmlRawTemplate); err != nil {
		fmt.Println("Internal error on generating error page template")
		os.Exit(1)
	}
}

func genErrorPage(errMsg, detailedMsg string) (string, error) {
	var err error
	data := struct {
		H1  string
		Msg string
		T   string
	}{
		errMsg,
		detailedMsg,
		time.Now().Format(time.ANSIC),
	}

	buf := new(bytes.Buffer)
	err = htmlTmpl.Execute(buf, data)
	return buf.String(), err
}

func sendErrorPage(w *bufio.Writer, errCode, errReason, errMsg, detailedMsg string) {
	page, err := genErrorPage(errMsg, detailedMsg)
	if err != nil {
		errl.Println("Error generating error page:", err)
		return
	}

	data := struct {
		Code   string
		Reason string
		Length int
	}{
		errCode,
		errReason,
		len(page),
	}
	buf := new(bytes.Buffer)
	if err := headTmpl.Execute(buf, data); err != nil {
		errl.Println("Error generating error page header:", err)
		return
	}

	w.WriteString(buf.String())
	w.WriteString("\r\n")
	w.WriteString(page)
	w.Flush()
}
