package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

const tpl = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Domain}}{{.Name}} git https://{{.Target}}{{.Name}}.git">
<meta name="go-source" content="{{.Domain}}{{.Name}} https://{{.Target}}{{.Name}} https://{{.Target}}{{.Name}}/tree/master{/dir} https://{{.Target}}{{.Name}}/blob/master{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here; <a href="https://{{.Target}}{{.Name}}">move along</a>.
</body>
</html>
`

func main() {

	domain := flag.String("domain", "localhost.local", "root domain for go")
	target := flag.String("target", "github.com", "target domain")

	var t = template.Must(template.New("goimport").Parse(tpl))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pkg := r.URL.Path
		data := struct {
			Name   string
			Domain string
			Target string
		}{
			Name:   pkg,
			Domain: *domain,
			Target: *target,
		}
		if err := t.Execute(w, data); err != nil {
			log.Println(err)
		}
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
