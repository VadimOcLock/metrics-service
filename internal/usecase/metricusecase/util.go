package metricusecase

import (
	"bytes"
	"html/template"
	"sort"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

func buildHTML(metrics []entity.Metric) (string, error) {
	const tpl = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Metrics</title>
		<style>
			table {
			  font-family: arial, sans-serif;
			  border-collapse: collapse;
			  width: 100%;
			}
			td, th {
			  border: 1px solid #dddddd;
			  text-align: left;
			  padding: 8px;
			}
			tr:nth-child(even) {
			  background-color: #dddddd;
			}
		</style>
    </head>
    <body>
        <table>
            <tr>
                <th>Type</th>
                <th>Name</th>
                <th>Value</th>
            </tr>
            {{range .}}
            <tr>
                <td>{{.Type}}</td>
                <td>{{.Name}}</td>
                <td>{{.Value}}</td>
            </tr>
            {{end}}
        </table>
    </body>
    </html>`

	t, err := template.New("metrics").Parse(tpl)
	if err != nil {
		return "", err
	}

	sortMetrics(&metrics)
	var buf bytes.Buffer
	if err = t.Execute(&buf, metrics); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func sortMetrics(metrics *[]entity.Metric) {
	sort.Slice(*metrics, func(i, j int) bool {
		if (*metrics)[i].Type == (*metrics)[j].Type {
			return (*metrics)[i].Name < (*metrics)[j].Name
		}

		return (*metrics)[i].Type < (*metrics)[j].Type
	})
}
