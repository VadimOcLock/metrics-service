package metricusecase

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

func buildHTML(metrics []entity.Metrics) (string, error) {
	const tpl = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>MetricsData</title>
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
                <th>Delta</th>
            </tr>
            {{range .}}
            <tr>
                <td>{{.MType}}</td>
                <td>{{.ID}}</td>
                <td>{{.Value}}</td>
                <td>{{.Delta}}</td>
            </tr>
            {{end}}
        </table>
    </body>
    </html>`

	t, err := template.New("metrics").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	SortMetrics(&metrics)
	var buf bytes.Buffer
	if err = t.Execute(&buf, metrics); err != nil {
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	return buf.String(), nil
}

func SortMetrics(metrics *[]entity.Metrics) {
	sort.Slice(*metrics, func(i, j int) bool {
		if (*metrics)[i].MType == (*metrics)[j].MType {
			return (*metrics)[i].ID < (*metrics)[j].ID
		}

		return (*metrics)[i].MType < (*metrics)[j].MType
	})
}
