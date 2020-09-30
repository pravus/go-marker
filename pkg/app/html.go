package app

import (
  "strings"
)

var htmlIndex = strings.TrimSpace(`
<html>
<head>
  <title>Marker</title>
</head>
<body>
{{ range $index, $name := .Names }}
<div><a href="/{{ $name }}">{{ $name }}</a></div>
{{ end }}
</body>
</html>
`)
