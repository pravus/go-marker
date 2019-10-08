package main

import (
  "flag"

  "gitlab.com/jhord/marker/pkg/app"
)

func main() {
  _bind       := flag.String("http",       "127.0.0.1:8000", "sets the http service listen address and port")
  _stylesheet := flag.String("stylesheet", "",               "sets a stylesheet to be used in the rendered HTML")
  _syntax     := flag.String("syntax",     "pygments",       "sets the syntax highlighter style")

  flag.Parse()

  app.New(*_bind, *_stylesheet, *_syntax).Run()
}
