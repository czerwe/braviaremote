package main

import (
	// "bytes"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"gobravia"
	"html/template"
	"net/http"
	"os"
)

type Options struct {
	Listenport int    `long:"listenport" env:"LISTENPORT" default:"4043" description:"Listening port"`
	broadcast  string `long:"broadcast" short:"b" required:"True" env:"SUBNET" description:"subnet to send wol signal example: 192.168.0.255" default:"10.0.0.255"`
	Host       string `long:"BRAVIAIP" short:"i" required:"True" env:"BRAVIAIP" description:"Hostname or IP of Bravia TV"`
	Pin        string `long:"pin" short:"p" required:"False" env:"PIN" description:"access pin set in TV" default:"0000"`
	Mac        string `long:"mac" short:"m" required:"False" env:"MAC" description:"MAC" default:"FC:F1:52:72:52:5F"`
	Loglevel   string `long:"loglevel" env:"LOGLEVEL" default:"info" description:"loglevel" choice:"warn" choice:"info" choice:"debug"`
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var (
	opts    Options
	level   log.Level
	version string = "0.3.0"
	bravia  *gobravia.BraviaTV
)

func main() {
	// version = "v1.0.0"

	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	switch opts.Loglevel {
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "debug":
		level = log.DebugLevel
	default:
		level = log.DebugLevel
	}

	log.SetLevel(level)

	log.Info("Initiating handlers")

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"requestURI":  r.RequestURI,
			"remoteAddr ": r.RemoteAddr,
		}).Error("failed to find page")
		// fmt.Println(r.RequestURI)
		http.ServeFile(w, r, "public/index.html")
	})

	r.HandleFunc("/commandlist", Commandlist)
	r.HandleFunc("/remote", Remote)
	r.HandleFunc("/keypress/{key:[A-Za-z0-9]*}", Keypress)
	r.HandleFunc("/macro/{key:[A-Za-z0-9]*}", Macro)

	// r.HandleFunc("/temperature/{sensortype:[A-Za-z\\-]+}.{serial:[A-Z0-9]+}:{channel:[0-9]}.{type:[A-Z]+}/{value}", KlimaHandler)
	// r.HandleFunc("/query/{measurement:[A-Za-z]+}/{room:[A-Za-z]+}", ccuprocessing.QueryHandler)

	log.WithFields(log.Fields{
		"port": opts.Listenport,
		"host": "0.0.0.0",
	}).Info("starting listender")

	log.WithFields(log.Fields{
		"host":  opts.Host,
		"pin":   opts.Pin,
		"mac":   opts.Mac,
		"bcast": opts.broadcast,
	}).Info("bravia data")

	bravia = gobravia.GetBravia(opts.Host, opts.Pin, opts.Mac)
	bravia.GetCommands()

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", opts.Listenport), r)

}

type CommandStruct struct {
	Commands map[string]string
}

func Commandlist(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/codelist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := CommandStruct{Commands: bravia.Commands}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Remote(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/remote.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, "nothing"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Keypress(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	log.WithFields(log.Fields{
		"key":         vars["key"],
		"tvconnected": bravia.Connected,
	}).Info("key pressed")

	var ok = false

	if !bravia.Connected {
		ok = bravia.GetCommands()
	} else {
		ok = true
	}

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}

	tmpl, err := template.ParseFiles("templates/codelist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if command, ok := bravia.Commands[vars["key"]]; ok {
		log.WithFields(log.Fields{
			"key": vars["key"],
		}).Debug("key found")

		bravia.SendCommand(command)

	} else {
		if vars["key"] == "poweron" {
			log.WithFields(log.Fields{
				"key": vars["key"],
			}).Debug("wakeup signal")

			bravia.Poweron(opts.broadcast)

		} else {
			log.WithFields(log.Fields{
				"key": vars["key"],
			}).Warning("key does not exists")
		}
	}

	b := make(map[string]string)
	b["hollei"] = "hoschi"

	if err := tmpl.Execute(w, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Macro(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	log.WithFields(log.Fields{
		"macro": vars["key"],
	}).Info("macro called")
}
