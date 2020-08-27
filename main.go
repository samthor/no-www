package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	codeNukeServiceWorker = `if(typeof ServiceWorkerGlobalScope!=="undefined"&&self instanceof ServiceWorkerGlobalScope){self.addEventListener("install",event=>{event.waitUntil(self.skipWaiting())});self.addEventListener("activate",event=>{const p=(async()=>{await self.clients.claim();const existingClients=await clients.matchAll({includeUncontrolled:true,type:"window"});try{existingClients.forEach(client=>client.navigate(client.url))}catch(e){}})();event.waitUntil(p)})}`
)

func main() {
	var err error

	http.HandleFunc("/", httpIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Running on :%v...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	log.Fatalf("Stop: %v", err)
}

func httpIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if maybeServiceWorker(r) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, codeNukeServiceWorker)
		return
	}

	var redir bool

	host := r.Host
	for strings.HasPrefix(host, "www.") {
		host = host[4:]
		redir = true
	}

	if !redir {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	r.URL.Host = host
	r.URL.Scheme = "https"
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1hr
	http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
}

func maybeServiceWorker(r *http.Request) bool {
	if r.Header.Get("Service-Worker") != "" {
		return true
	}
	if r.URL.Path == "/sw.js" || r.URL.Path == "/service-worker.js" {
		return true
	}
	return false
}
