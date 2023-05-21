# Zebra

![Version](https://img.shields.io/badge/Version-Prototype-red)
[![GoDoc](https://godoc.org/github.com/up1-io/zebra?status.svg)](https://godoc.org/github.com/up1-io/zebra)

Zebra is a minimalist web framework for Go that focuses on simplicity, performance, and ease of use.

> Note: Zebra is currently in the prototype stage. It is not yet ready for production use.

## Features

- Directory-based routing
- Middleware support
- Static file serving

## Getting Started

### Installation

To install Zebra, use `go get`:

```bash
go get github.com/up1-io/zebra
```

### Quick Start

Here's a simple example to demonstrate how to create a basic Zebra application:

```go
package main

import (
	"github.com/up1-io/zebra"
)

func main() {
	app, err := zebra.New()
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.ListenAndServe(":8080"))
}
```

Folder structure:

```
.
├── main.go
└── pages
    └── _index.gohtml
    └── _layout.gohtml
    └── _404.gohtml
```

`pages/_index.gohtml`:

```html
{{ define "content" }}
    <h1>Hello, World!</h1>
{{ end }}
```

`pages/_layout.gohtml`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Zebra</title>
</head>
<body>
    {{ template "content" . }}
</body>
</html>
```