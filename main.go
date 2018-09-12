package main

import (
	"bufio"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

const tpl = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Domain}}{{.Src}} git {{.Target}}">
<meta name="go-source" content="{{.Domain}}{{.Src}} {{.Target}} {{.Target}}/tree/master{/dir} {{.Target}}/blob/master{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here.
</body>
</html>
`

type row struct {
	order    int
	matcher  *regexp.Regexp
	replacer string
}

func main() {
	config := flag.String("conf", "/data/urls.conf", "configuration for the url")
	domain := flag.String("domain", "localhost", "source domain")
	flag.Parse()

	if len(*domain) == 0 {
		log.Fatal("error need a domain")
	}
	_, err := os.Stat(*config)
	if os.IsNotExist(err) {
		log.Fatalf("error file %s does not exist", *config)
	}
	file, err := os.Open(*config)
	if err != nil {
		log.Fatalf("error while opeing the config file: %v", err)
	}
	defer file.Close()

	matcher := []row{}
	scanner := bufio.NewScanner(file)
	counter := 1
	for scanner.Scan() {
		matches := strings.Split(scanner.Text(), "->")
		if len(matches) != 2 {
			log.Fatalf("wrong line nee '->': %s", scanner.Text())
		}
		matcher = append(matcher, row{
			order:    counter,
			matcher:  regexp.MustCompile(strings.TrimSpace(matches[0])),
			replacer: strings.TrimSpace(matches[1]),
		})
		counter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error while scanning file: %v", err)
	}
	log.Printf("loading %v lines", len(matcher))

	sort.Slice(matcher, func(i, j int) bool {
		return matcher[i].order < matcher[j].order
	})

	var t = template.Must(template.New("goimport").Parse(tpl))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		src := r.URL.Path
		var target string
		for _, value := range matcher {
			if value.matcher.MatchString(src) {
				target = value.matcher.ReplaceAllString(src, value.replacer)
				break
			}
		}
		if target == "" {
			http.Error(w, "No repo", http.StatusNotFound)
			return
		}

		data := struct {
			Domain string
			Src    string
			Target string
		}{
			*domain,
			src,
			target,
		}
		if err := t.Execute(w, data); err != nil {
			log.Println(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
