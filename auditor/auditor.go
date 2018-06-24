package auditor

import (
	"log"
	"github.com/sumwonyuno/cp-scoring/model"
)

func Audit(state model.State, templates []model.Template) model.Report {
	var report model.Report

	log.Println(templates)

	return report
}