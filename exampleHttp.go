package main

import (
  "fmt"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  // http.ResponseWriter value assembles the HTTP server's response, 
  // by writing to it, we sen ddata to the HTTP client
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
  // a http.Request is a data structure that represents the client HTTP request
  // r.URL.Path is the path component  of the request url
  // the trailing [1:] creates a sub slice dropping the "/"
}

func main() {
  // tells http package to handle all requests to the web root("/") with "handler"
  http.HandleFunc("/", handler)
  // then calls http.ListenAndServe specifying that it should listen on port 8080 on any interface
  http.ListenAndServe(":8080", nil)
}
