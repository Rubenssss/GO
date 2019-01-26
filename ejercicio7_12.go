package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var dbmux sync.Mutex

var listTemp = template.Must(template.New("list").Parse(`
<html>
<body>
{{ range $key, $value := .ItemMap }}
<p>{{$key}}: ${{$value}}</p>
{{ end }}
</body>
</html>
`))

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	
	
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]int



type TemplateData struct {
	ItemMap database
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	dbmux.Lock()
	if err := listTemp.Execute(w, &TemplateData{db}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to execute template: %q\n", err)
	}
	

	dbmux.Unlock()
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	dbmux.Lock()
	price, ok := db[item]
	dbmux.Unlock()
	if ok {
		fmt.Fprintf(w, "$%d\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) 
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "wrong price: %s\n", priceStr)
	} else {
		dbmux.Lock()
		db[item] = price
		dbmux.Unlock()
		fmt.Fprintf(w, "created or updated %s: $%d\n", item, price)
	}
}





