package main

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", homeHandler)
}

// The &0x20 thing is a hack that compares the case of the first letter to the
// case of the last letter.
func grandtotal(in map[string]interface{}) int {
	var t int

	for attacker, a := range in {
		if attacker[0]&0x20 == attacker[len(attacker)-1]&0x20 {
			continue
		}

		for victim, v := range a.(map[string]interface{}) {
			if victim[0]&0x20 != victim[len(victim)-1]&0x20 {
				continue
			}

			for _, w := range v.(map[string]interface{}) {
				t += int(w.(float64))
			}
		}
	}

	return t
}

func total(in map[string]interface{}) int {
	var t int

	for _, v := range in {
		if f, ok := v.(float64); ok {
			t += int(f)
		} else {
			t += total(v.(map[string]interface{}))
		}
	}

	return t
}

var homeTemplate = template.Must(template.New("home").Funcs(template.FuncMap{
	"grandtotal": grandtotal,
	"total":      total,
}).Parse(`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	{{with $total := grandtotal .}}<title>{{$total}} zombies. Dead.</title>
	<link href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.3.1/css/bootstrap-combined.min.css" rel="stylesheet">
</head>
<body>
<div class="container">
	<h1>{{$total}} zombies. Dead.</h1>{{end}}

	<table class="table">
	<thead><tr>
		<th></th>
		<th>Nick</th>
		<th>Rochelle</th>
		<th>Coach</th>
		<th>Ellis</th>
		<th>Bill</th>
		<th>Zoey</th>
		<th>Louis</th>
		<th>Francis</th>
	</tr></thead>
	<tr>
		<th>Common Infected</th>
		<td>{{total .Gambler.infected}}</td>
		<td>{{total .Producer.infected}}</td>
		<td>{{total .Coach.infected}}</td>
		<td>{{total .Mechanic.infected}}</td>
		<td>{{total .NamVet.infected}}</td>
		<td>{{total .TeenGirl.infected}}</td>
		<td>{{total .Manager.infected}}</td>
		<td>{{total .Biker.infected}}</td>
	</tr>
	<tr>
		<th>Tanks</th>
		<td>{{total .Gambler.TANK}}</td>
		<td>{{total .Producer.TANK}}</td>
		<td>{{total .Coach.TANK}}</td>
		<td>{{total .Mechanic.TANK}}</td>
		<td>{{total .NamVet.TANK}}</td>
		<td>{{total .TeenGirl.TANK}}</td>
		<td>{{total .Manager.TANK}}</td>
		<td>{{total .Biker.TANK}}</td>
	</tr>
	<tr>
		<th>Boomers</th>
		<td>{{total .Gambler.BOOMER}}</td>
		<td>{{total .Producer.BOOMER}}</td>
		<td>{{total .Coach.BOOMER}}</td>
		<td>{{total .Mechanic.BOOMER}}</td>
		<td>{{total .NamVet.BOOMER}}</td>
		<td>{{total .TeenGirl.BOOMER}}</td>
		<td>{{total .Manager.BOOMER}}</td>
		<td>{{total .Biker.BOOMER}}</td>
	</tr>
	<tr>
		<th>Chargers</th>
		<td>{{total .Gambler.CHARGER}}</td>
		<td>{{total .Producer.CHARGER}}</td>
		<td>{{total .Coach.CHARGER}}</td>
		<td>{{total .Mechanic.CHARGER}}</td>
		<td>{{total .NamVet.CHARGER}}</td>
		<td>{{total .TeenGirl.CHARGER}}</td>
		<td>{{total .Manager.CHARGER}}</td>
		<td>{{total .Biker.CHARGER}}</td>
	</tr>
	<tr>
		<th>Hunters</th>
		<td>{{total .Gambler.HUNTER}}</td>
		<td>{{total .Producer.HUNTER}}</td>
		<td>{{total .Coach.HUNTER}}</td>
		<td>{{total .Mechanic.HUNTER}}</td>
		<td>{{total .NamVet.HUNTER}}</td>
		<td>{{total .TeenGirl.HUNTER}}</td>
		<td>{{total .Manager.HUNTER}}</td>
		<td>{{total .Biker.HUNTER}}</td>
	</tr>
	<tr>
		<th>Jockeys</th>
		<td>{{total .Gambler.JOCKEY}}</td>
		<td>{{total .Producer.JOCKEY}}</td>
		<td>{{total .Coach.JOCKEY}}</td>
		<td>{{total .Mechanic.JOCKEY}}</td>
		<td>{{total .NamVet.JOCKEY}}</td>
		<td>{{total .TeenGirl.JOCKEY}}</td>
		<td>{{total .Manager.JOCKEY}}</td>
		<td>{{total .Biker.JOCKEY}}</td>
	</tr>
	<tr>
		<th>Smokers</th>
		<td>{{total .Gambler.SMOKER}}</td>
		<td>{{total .Producer.SMOKER}}</td>
		<td>{{total .Coach.SMOKER}}</td>
		<td>{{total .Mechanic.SMOKER}}</td>
		<td>{{total .NamVet.SMOKER}}</td>
		<td>{{total .TeenGirl.SMOKER}}</td>
		<td>{{total .Manager.SMOKER}}</td>
		<td>{{total .Biker.SMOKER}}</td>
	</tr>
	<tr>
		<th>Spitters</th>
		<td>{{total .Gambler.SPITTER}}</td>
		<td>{{total .Producer.SPITTER}}</td>
		<td>{{total .Coach.SPITTER}}</td>
		<td>{{total .Mechanic.SPITTER}}</td>
		<td>{{total .NamVet.SPITTER}}</td>
		<td>{{total .TeenGirl.SPITTER}}</td>
		<td>{{total .Manager.SPITTER}}</td>
		<td>{{total .Biker.SPITTER}}</td>
	</tr>
	<tr>
		<th>Witches</th>
		<td>{{total .Gambler.witch}}</td>
		<td>{{total .Producer.witch}}</td>
		<td>{{total .Coach.witch}}</td>
		<td>{{total .Mechanic.witch}}</td>
		<td>{{total .NamVet.witch}}</td>
		<td>{{total .TeenGirl.witch}}</td>
		<td>{{total .Manager.witch}}</td>
		<td>{{total .Biker.witch}}</td>
	</tr>
	</table>
</div>
</body>
</html>`))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	view, err := bucket.View("l4d2", "kills", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = homeTemplate.Execute(w, view.Rows[0].Value)
	if err != nil {
		panic(err)
	}
}
