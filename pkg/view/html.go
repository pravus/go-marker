package view

import (
  "text/template"
  "strings"
)

var htmlRoot = template.Must(template.New("root").Parse(strings.TrimSpace(`
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Filename }}</title>
  <style>
    .code {
      background: #eee;
      font-family: monospace;
      margin: 0em 0em;
      padding: 0em 0em;
    }
    .code-block {
      background: #eee;
      border: 1px solid #aaa;
      border-radius: 4px;
      margin: 1em 0em;
      padding: 0em 1em;
    }
  </style>
{{- if .Stylesheet }}
  <link rel="stylesheet" href="{{ .Stylesheet }}">
{{- end }}
</head>
<body>{{ .Body }}</body>
</html>
`)))
