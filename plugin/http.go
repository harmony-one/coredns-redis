package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Plugin) reload(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	zone := q.Get("zone")
	w.Header().Set("Content-Type", "application/json")
	if zone == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "require zone"})
		return
	}
	loaded, err := p.loadCacheForZone(zone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "loaded": loaded})

}

func (p *Plugin) serveHttp() {
	go func() {
		http.HandleFunc("/reload", p.reload)
		fmt.Printf("Serving HTTP server at port %v\n", p.Redis.ReloadHttpPort)
		err := http.ListenAndServe(fmt.Sprintf(":%v", p.Redis.ReloadHttpPort), nil)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Printf("Started HTTP server")
}
