package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
  "regexp"
  "errors"
)

// start by defining the datastructures
// a wiki consists of a series of interconnected pages, each of which has a title and a body
type Page struct {
  Title string
  Body  []byte    // a byte slice, type expect by th eio libraries
}

// loaded once, to avoid inefficiently load on demand everytime a page is rendered
// template.Must is a convenience wrapper that panics when passed a non nil error value
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// valid path regex for a Title
// MustCompile - panics where as Compile would return an error as a second param
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  // validates path and extracts the page title
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("Invalid Page Title")
  }
  return m[2], nil
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
  // @@@@@@@ problem here
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
  print(filename, "\n")
  print(title, "\n")
  print(body, "\n")
  // callers can now check the second param, if it is nil  then a page is loaded successfully
  return &Page{Title:title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  // handles URLs prefixed with "/view/"
  // function loads the page data
  p, err := loadPage(title)   // generally bad practice to ignore the error here
  if err != nil {
    // if no page found create a page and let them edit it
    http.Redirect(w,  r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  // byte(body) converts body to a byte to go into the page struct
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string))http.HandlerFunc {
  // the returned function is a closure, it encloses the values defined outside of it      
  return func(w http.ResponseWriter, r *http.Request) {
    // here we will extract the page title from the request
    // and call the provided handler 'fn'
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  // html/templates pacakge helps with security, e.g. it auto escapes > to &gt
  // template.ParseFiles reads the contents of edit.html and returns a *template.Template
  // t.execute, executes the template by writing the generated HTML to the http.ResponseWriter
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil  {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  // p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  // p1.save()
  // p2, _ := loadPage("TestPage")
  // fmt.Println(string(p2.Body))

  // lets the http know to respond to view, edit, save and etc
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.ListenAndServe(":8080", nil)
}
