package syntax

import (
  "errors"
  "io"

  "github.com/alecthomas/chroma/formatters/html"
  "github.com/alecthomas/chroma/lexers"
  "github.com/alecthomas/chroma/styles"
)

type Highlighter struct {
  Style string
}

func New(style string) *Highlighter {
  return &Highlighter{Style: style}
}

func (h *Highlighter) Render(w io.Writer, lang, code string) error {
  lexer := lexers.Get(lang)
  if lexer == nil {
    return errors.New("No lexer found")
  }

  style := styles.Get(h.Style)
  if style == nil {
    style = styles.Fallback
  }

  formatter := html.New()

  iterator, err := lexer.Tokenise(nil, code)
  if err != nil {
    return err
  }

  err = formatter.Format(w, style, iterator)
  if err != nil {
    return err
  }

  return nil
}
