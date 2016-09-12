# GoFileBrowser

Go directory watcher that creates a static `index.html` file of the watched directory when changes are made.

Example config:

``` 
static:
  serve: /static
  path: ./static

locations:
  test:
    title: "My Static Files"
    watch: /home/user/go/src/github.com/nananas/gofilebrowser/test
    recursive: Yes
```

A small test server is included in the `./test` directory.