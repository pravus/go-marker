package app

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"

  "github.com/go-chi/chi"
)

type App struct {
  view view
}

type view interface {
  Render(http.ResponseWriter, string, []byte) error
}

func New(view view) *App {
  return &App{view: view}
}

func (app *App) Run(bind string) {
  r := chi.NewRouter()
  r.Get("/{basename}", app.root())

  fmt.Printf("http.bind: %s\n", bind)
  err := http.ListenAndServe(bind, r)
  if err != nil {
    fmt.Printf("http.error: %v\n", err)
  }
}

func (app *App) root() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    basename := chi.URLParam(r, "basename")

    if basename == "favicon.ico" {
      w.WriteHeader(http.StatusNotFound)
      return
    }

    filename := basename +".md"
    file, err := os.Open(filename)
    if err != nil {
      fmt.Printf("root.open: %v\n", err)
      w.WriteHeader(http.StatusNotFound)
      return
    }
    defer file.Close()

    all, err := ioutil.ReadAll(file)
    if err != nil {
      fmt.Printf("root.open: %v\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    err = app.view.Render(w, filename, all)
    if err != nil {
      fmt.Printf("root.render: %v\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
}
