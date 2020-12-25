package main

import (
	"log"
	"net/http"
	"text/template"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	tpl *template.Template
	sshPublicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDB50LyYeb6V6JQ/flA0IupNFXeDDdvqDuWWiYHhA/2FHfnGnUQXVOaNZWaoO48xja5eilNI1pTCvIvie1yUzY2PvPr6LVtBObb5wtHfmXiIK+CIou62S2nYVEQvQiEccFJyatvaEVhH2KhNvd85oF9/ryPy+yhWR7QYGCM39y4o7dt6bNVyEVtCf+XGEOwEcD74mlQLDnaHAtIvpKU2KzD1bReerGkMAvEud8Mb3zTITqlwDsLTRTOFVZzqD9yPM8Jg11gnSisojXO8ihLL+9PaYBpCLT6NztyHH5flgfWHofEbk4XIwftqAcfDO+/noH3PuRX45ByTvDr//dQ05HD ansible"
)

type SeedTemplate struct {
	Hostname string
	SSHPublicKey string
}

func createSeed(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	templateName := vars["arch"]
	if templateName == "" {
		templateName = "rasp.tmpl"
	} else {
		templateName = templateName + ".tmpl"
	}

	query := r.URL.Query()
	hostname := query.Get("hostname")
	if hostname == "" {
		http.Error(w, "hostname have to be provided", 400)
		return
	}
	data := SeedTemplate{
		Hostname:hostname,
		SSHPublicKey:sshPublicKey,
	}

	if err := tpl.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main()  {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	rtr := mux.NewRouter()
	rtr.Handle("/metrics", promhttp.Handler())
	rtr.HandleFunc("/seed/{arch}", createSeed).Methods("GET")
	rtr.HandleFunc("/seed", createSeed).Methods("GET")
	http.Handle("/", rtr)
	tpl, err = template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	zap.L().Info("Server started. Loaded templates", zap.String("templates", tpl.ParseName))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

