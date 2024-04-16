// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package report

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

// Report is a template for a report to be produced from monitoring data
// collected by Norma. To render the report, use the Report`s render function.
type Report struct {
	name     string
	template []byte
}

//go:embed single_eval_report.Rmd
var singleEvalReportTemplate []byte

//go:embed multi_eval_report.Rmd
var multiEvalReportTemplate []byte

var (
	// SingleEvalReport is a report template covering metrics collected in a single
	// scenario evaluation.
	SingleEvalReport = Report{
		name:     "single_eval_report",
		template: singleEvalReportTemplate,
	}

	// MultiEvalReport is a report template comparing the results of multiple
	// scenariou evaluations. The input CSV should be the concatenation of the
	// individual measurement CSV files.
	MultiEvalReport = Report{
		name:     "multi_eval_report",
		template: multiEvalReportTemplate,
	}
)

//go:embed render.R
var renderScript []byte

// Render renders this report using the given input data file (in CSV format)
// and places its results into the defined output directory.
func (r *Report) Render(datafile, outputdir string) (string, error) {
	script, err := createTempFile(renderScript, ".R")
	if err != nil {
		return "", err
	}
	defer os.Remove(script)

	template, err := createTempFile(r.template, ".Rmd")
	if err != nil {
		return "", err
	}
	defer os.Remove(template)

	outputfile := r.name + ".html"

	cmd := exec.Command("Rscript", script, template, datafile, outputdir, outputfile)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%v", out.String(), err)
	}

	return outputfile, nil
}

func createTempFile(content []byte, suffix string) (string, error) {
	file, err := os.CreateTemp("", "tmp_*"+suffix)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return "", err
	}
	return file.Name(), nil
}
