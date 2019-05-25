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
  highlighter highlighter
  renderer    markdown.Renderer
}

func New (highlighter highlighter) *View {
  view := &View{highlighter: highlighter}
  options := html.RendererOptions{
    RenderNodeHook: view.renderNode,
  }
  view.renderer = html.NewRenderer(options)
  return view
}

func (view *View) Render(w http.ResponseWriter, filename string, bytes []byte) error {
  w.Header().Set("Content-Type", "text/html")
  w.Write([]byte(""+
    "<!doctype html>\n"+
    "<html>\n"+
    "<head>\n"+
    "  <meta charset=\"utf-8\">"+
    "  <title>"+ filename +"</title>\n"+
    "  <style>\n"+
    "    .code-block {\n"+
    "      background: #eee;\n"+
    "      border: 1px solid #aaa;\n"+
    "      border-radius: 4px;\n"+
    "      margin: 1em 0em;\n"+
    "      padding: 0em 1em;\n"+
    "      \n"+
    "      \n"+
    "      \n"+
    "      \n"+
    "      \n"+
    "    }\n"+
    "  </style>\n"+
    "</head>\n"+
    "<body>\n"+
    string(markdown.ToHTML(bytes, nil, view.renderer))+
    "</body>\n"+
    "</html>\n"+
  ""))
  return nil
}

func (view *View) renderNode (w io.Writer, gen ast.Node, entering bool) (ast.WalkStatus, bool) {
  switch node := gen.(type) {
    case *ast.CodeBlock:
      return view.renderCodeBlock(w, node, entering)
  }
  return ast.GoToNext, false
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
