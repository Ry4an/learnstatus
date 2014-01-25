package learnstatus

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
	"time"
)

type Hit struct {
	Date time.Time
}

func init() {
	http.HandleFunc("/", check)
}

func check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")

	// create new record
	c := appengine.NewContext(r)
	hit := Hit{
		Date: time.Now(),
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Hit", nil), &hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// see how long it look the last 10 to come in
	q := datastore.NewQuery("Hit").Order("-Date").Limit(10)
	hits := make([]Hit, 0, 10)
	if _, err := q.GetAll(c, &hits); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var oldest = hits[len(hits)-1]
	var duration = hits[0].Date.Sub(oldest.Date)
	var up = duration/time.Minute > 5
	if err := checkTemplate.Execute(w, up); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var checkTemplate = template.Must(template.New("check").Parse(checkTemplateStr))

const checkTemplateStr = `
<html>
<head>
<title>{{if .}}Up{{else}}Down{{end}}</title>
</head><body>{{if .}}Up{{else}}Down{{end}}</body></html>
`
