# go-urlredirect

Adding config urls.conf:
```bash
/x/([a-zA-Z0-9-_]) -> https://github.com/myname/$1
```

Running the container:
```bash
docker run --rm -v $(PWD):/data/ -p 8080:8080 blackreloaded/go-urlredirect -domain=localhost:8080
```