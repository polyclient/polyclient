// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xhtml

import "html/template"

// HTMLTemplateData is the data structure for the HTML template output.
type HTMLTemplateData struct {
	Headers []string
	Rows    [][]string
	UseCSS  bool
}

// HTMLTemplate is the HTML template for the HTML exporter.
const HTMLTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Export</title>
    {{- if .UseCSS }}
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            padding: 0;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background-color: #fff;
            border: 1px solid #ddd;
        }
        th, td {
            padding: 8px;
            text-align: left;
            border: 1px solid #ddd;
        }
        th {
            background-color: #f4f4f4;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
    </style>
    {{- end }}
</head>
<body>
    <table>
        {{- if .Headers }}
        <thead>
            <tr>
                {{- range .Headers }}
                <th>{{ . }}</th>
                {{- end }}
            </tr>
        </thead>
        {{- end }}
        <tbody>
            {{- range .Rows }}
            <tr>
                {{- range . }}
                <td>{{ . }}</td>
                {{- end }}
            </tr>
            {{- end }}
        </tbody>
    </table>
</body>
</html>
`

// GetTemplate loads the HTML template for the HTML exporter.
func GetTemplate() *template.Template {
	return template.Must(template.New("htmlExport").Parse(HTMLTemplate))
}
