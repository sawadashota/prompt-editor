Prompt Editor
===

[![CircleCI](https://circleci.com/gh/sawadashota/prompt-editor/tree/master.svg?style=svg)](https://circleci.com/gh/sawadashota/prompt-editor/tree/master)
[![GoDoc](https://godoc.org/github.com/sawadashota/prompt-editor?status.svg)](https://godoc.org/github.com/sawadashota/prompt-editor)
[![codecov](https://codecov.io/gh/sawadashota/prompt-editor/branch/master/graph/badge.svg)](https://codecov.io/gh/sawadashota/prompt-editor)
[![Go Report Card](https://goreportcard.com/badge/github.com/sawadashota/prompt-editor)](https://goreportcard.com/report/github.com/sawadashota/prompt-editor)
[![GolangCI](https://golangci.com/badges/github.com/sawadashota/prompt-editor.svg)](https://golangci.com)

Open editor in console app.

Installation
---

```
$ go get -u github.com/sawadashota/prompt-editor
```

Usage
---

```go
type User struct { 
	Name string `json:"name"`
}

func main() {
	user := &User{"Alice"}

	// second argument is default value option
	editor, err := prompteditor.New("vi", prompteditor.MarshalJSON(user))
	if err != nil { 
		// error handling
	}

	if err := editor.Open(); err != nil {
		// error handling
	}

	var edited User
	if err := editor.UnmarshalEdited(&edited); err != nil {
		// error handling
	}

	fmt.Println(user.Name)
	// print edited user name
}
```