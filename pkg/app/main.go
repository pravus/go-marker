package app

import (
  "fmt"
  "html/template"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/go-chi/chi"

  "gitlab.com/jhord/marker/pkg/syntax"
  "gitlab.com/jhord/marker/pkg/view"
)

type App struct {
  bind       string
  stylesheet string
  style      string
  renderer   renderer
}

type renderer interface {
  Render(http.ResponseWriter, string, []byte) error
}

func New(bind, stylesheet, style string) *App {
  return &App{
    bind:       bind,
    stylesheet: stylesheet,
    style:      style,
    renderer:   view.New(stylesheet, syntax.New(style)),
  }
}

func (app *App) Run() {
  r := chi.NewRouter()
  r.Get("/",           app.index())
  r.Get("/{basename}", app.root())

  fmt.Printf("http.bind: %s\n", app.bind)
  err := http.ListenAndServe(app.bind, r)
  if err != nil {
    fmt.Printf("http.error: %v\n", err)
  }
}

func (app *App) index() http.HandlerFunc {
  index := template.Must(template.New("index").Parse(htmlIndex))
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    list, err := ioutil.ReadDir(".")
    if err != nil {
      log.Printf("readdir: error=%v", err)
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      return
    }
    names := []string{}
    for _, info := range list {
      name := info.Name()
      if strings.HasSuffix(name, ".md") {
        names = append(names, name[0:len(name) - 3])
      }
    }
    err = index.Execute(w, struct {
      Names []string
    }{
      Names: names,
    })
    if err != nil {
      log.Printf("template: error=%v", err)
    }
  })
}

func (app *App) root() http.HandlerFunc {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    basename := chi.URLParam(r, "basename")

    if basename == "favicon.ico" {
      w.WriteHeader(http.StatusNotFound)
      return
    }

    if basename == app.stylesheet {
      app.renderStylesheet(w, r)
      return
    }

    app.renderMarkdown(w, r, basename + ".md")
  })
}

func (app *App) renderMarkdown(w http.ResponseWriter, r *http.Request, filename string) {
  file, err := os.Open(filename)
  if err != nil {
    fmt.Printf("markdown.open: %s: %v\n", filename, err)
    w.WriteHeader(http.StatusNotFound)
    return
  }
  defer file.Close()

  all, err := ioutil.ReadAll(file)
  if err != nil {
    fmt.Printf("markdown.read: %s: %v\n", filename, err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  err = app.renderer.Render(w, filename, all)
  if err != nil {
    fmt.Printf("markdown.render: %s: %v\n", filename, err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (app *App) renderStylesheet(w http.ResponseWriter, r *http.Request) {
  file, err := os.Open(app.stylesheet)
  if err != nil {
    fmt.Printf("stylesheet.open: %s: %v\n", app.stylesheet, err)
    w.WriteHeader(http.StatusNotFound)
    return
  }
  defer file.Close()

  w.Header().Set("Content-Type", "text/css")
  _, err = io.Copy(w, file)
  if err != nil {
    fmt.Printf("stylesheet.copy: %s: %v\n", app.stylesheet, err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}
