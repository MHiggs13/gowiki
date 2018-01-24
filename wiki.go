package main

import (
  "fmt"
  "io/ioutil"
)

// start by defining the datastructures
// a wiki consists of a series of interconnected pages, each of which has a title and a body
type Page struct {
  Title string
  Body  []byte    // a byte slice, type expect by th eio libraries
}

func (p *Page) save() error {
  // save method to allow for persistent storage of a page        
  // Signature reads: this is a method named save that takes as its receive p, a pointer to Page,
  // it takes no parameters and returns a value of type error
  filename := p.Title + ".txt"
  // WriteFile returns an error value so our method must, if all goes well nil will be returned
  // the octal integer literal 0600 indicates that the file should be created with read-write 
  // permissions for the current user only
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  // constructs the file name from the title param
  // readds the file's contents intoa  new variable body
  // returns a pointer to a Page literal constructed with the correct title and body values
  filename := title + ".txt"
  // _ throws away the error part of the return from ReadFile
  //body, _ := ioutil.ReadFile(filename)
  // instead we can do
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  // callers can now check the second param, if it is nil  then a page is loaded successfully
  return &Page{Title:title, Body: body}, nil
}

func main() {
  p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  p1.save()
  p2, _ := loadPage("TestPage")
  fmt.Println(string(p2.Body))
}
