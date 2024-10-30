package metricusecase

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

func buildHTML(metrics []entity.Metric) (string, error) {
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
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	SortMetrics(&metrics)
	var buf bytes.Buffer
	if err = t.Execute(&buf, metrics); err != nil {
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	return buf.String(), nil
}

func SortMetrics(metrics *[]entity.Metric) {
	sort.Slice(*metrics, func(i, j int) bool {
		if (*metrics)[i].Type == (*metrics)[j].Type {
			return (*metrics)[i].Name < (*metrics)[j].Name
		}

		return (*metrics)[i].Type < (*metrics)[j].Type
	})
}

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

func buildHTMLNew(metrics []entity.Metrics) (string, error) {
	t, err := template.New("metrics").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	views := make([]metricView, 0)
	for _, m := range metrics {
		vl, err := m.MetricValue()
		if err != nil {
			log.Err(err).Send()
			continue
		}
		views = append(views, metricView{
			Name:  m.ID,
			Type:  m.MType,
			Value: vl,
		})
	}

	sortMetricsNew(&views)
	var buf bytes.Buffer
	if err = t.Execute(&buf, views); err != nil {
		return "", fmt.Errorf("metricusecase.buildHTML: %w", err)
	}

	return buf.String(), nil
}

type metricView struct {
	Name  string
	Type  string
	Value string
}

func sortMetricsNew(metrics *[]metricView) {
	sort.Slice(*metrics, func(i, j int) bool {
		if (*metrics)[i].Type == (*metrics)[j].Type {
			return (*metrics)[i].Name < (*metrics)[j].Name
		}

		return (*metrics)[i].Type < (*metrics)[j].Type
	})
}
