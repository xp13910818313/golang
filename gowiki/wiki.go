package main

import (
  "regexp"
  "html/template"
  "net/http"
  "github.com/pvandorp/golang/gowiki/data"
)

var templates = template.Must(template.ParseFiles("views/edit.html", "views/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := data.LoadPage(title)

  if err != nil {
    http.Redirect(w, r, "/edit/" + title, http.StatusFound)
    return
  }

  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := data.LoadPage(title)

  if err != nil {
    p = &data.Page{Title: title}
  }

  renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *data.Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &data.Page{Title: title, Body: []byte(body)}
  p.Save()
  http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))

  http.ListenAndServe(":8080", nil)
}
