package template

import (
	"bytes"
	"html/template"
	"git.lewoof.xyz/gitbrowse/config"
)

func CommonHead(c *config.PageConfig) string {
	var headBuffer bytes.Buffer
	t := template.Must(template.New("head").Parse(`
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>{{.Title}}</title>
			<meta name="description" content="{{.Description}}">
			{{range .Styles}}
				<link rel="stylesheet" href="{{.}}">
			{{end}}
			<meta name="theme-color" content="#d8a657">
			<meta property="og:image" content="{{.Thumbnail}}" />
			<link rel="icon" href="{{.Favicon}}">
		</head>
	`))
	t.Execute(&headBuffer, c)
	return headBuffer.String()
}
