package main

import (
  "flag"

  "gitlab.com/jhord/marker/pkg/app"
  "gitlab.com/jhord/marker/pkg/view"
  "gitlab.com/jhord/marker/pkg/syntax"
)

func main() {
  _bind   := flag.String("http",   "127.0.0.1:8000", "sets the http service listen address and port")
  _syntax := flag.String("syntax", "pygments",       "sets the syntax highlighter style")
  flag.Parse()

  app.New(view.New(syntax.New(*_syntax))).Run(*_bind)
}
