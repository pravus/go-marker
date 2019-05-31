package view

import (
  "fmt"
  "io"
  "net/http"

  "github.com/gomarkdown/markdown"
  "github.com/gomarkdown/markdown/ast"
  "github.com/gomarkdown/markdown/html"
)

type highlighter interface {
  Render(io.Writer, string, string) error
}

type View struct {
  stylesheet  string
  highlighter highlighter
  renderer    markdown.Renderer
}

func New (stylesheet string, highlighter highlighter) *View {
  view := &View{stylesheet: stylesheet, highlighter: highlighter}
  options := html.RendererOptions{
    RenderNodeHook: view.renderNode,
  }
  view.renderer = html.NewRenderer(options)
  return view
}

func (view *View) Render(w http.ResponseWriter, filename string, bytes []byte) error {
  w.Header().Set("Content-Type", "text/html")
  htmlRoot.Execute(w, struct {
    Filename   string
    Stylesheet string
    Body       string
  }{
    Filename:   filename,
    Stylesheet: view.stylesheet,
    Body:       string(markdown.ToHTML(bytes, nil, view.renderer)),
  })
  return nil
}

func (view *View) renderNode (w io.Writer, gen ast.Node, entering bool) (ast.WalkStatus, bool) {
  switch node := gen.(type) {
    case *ast.Code:
      return view.renderCode(w, node, entering)
    case *ast.CodeBlock:
      return view.renderCodeBlock(w, node, entering)
  }
  return ast.GoToNext, false
}

func (view *View) renderCode (w io.Writer, node *ast.Code, entering bool) (ast.WalkStatus, bool) {
  rendered := true
  w.Write([]byte("<span class=\"code\">\n"))
  w.Write(node.Literal)
  w.Write([]byte("</span>\n"))
  return ast.GoToNext, rendered
}

func (view *View) renderCodeBlock (w io.Writer, node *ast.CodeBlock, entering bool) (ast.WalkStatus, bool) {
  rendered := false
  if len(node.Info) > 0 {
    rendered = true
    w.Write([]byte("<div class=\"code-block\">\n"))
    err := view.highlighter.Render(w, string(node.Info), string(node.Literal))
    if err != nil {
      rendered = false
      fmt.Printf("view.code-block.highlight: %v\n", err)
      w.Write([]byte(fmt.Sprintf("<!-- error: code-block: %s %v -->", string(node.Info), err)))
    }
    w.Write([]byte("</div>\n"))
  }
  return ast.GoToNext, rendered
}
