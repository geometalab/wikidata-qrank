// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dataLoader *DataLoader

func main() {
	var portFlag = flag.Int("port", 0, "port for serving HTTP requests")
	var dataFlag = flag.String("data", "./cache", "directory with data files for serving")
	flag.Parse()

	port := *portFlag
	if port == 0 {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	var err error
	dataLoader, err = NewDataLoader(*dataFlag)
	if err != nil {
		log.Fatal(err)
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := dataLoader.Reload(); err != nil {
					// Log an error, but keep serving previous data.
					log.Printf("failed to reload data: %q", err)
				}
			}
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/download/qrank.gz", handleDownloadQRank)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
	done <- true
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s",
		`<html>
<head>
<link href='https://tools-static.wmflabs.org/fontcdn/css?family=Roboto+Slab:400,700' rel='stylesheet' type='text/css'/>
<style>
* {
  font-family: 'Roboto Slab', serif;
}
h1 {
  color: #0066ff;
  margin-left: 1em;
  margin-top: 1em;
}
p {
  margin-left: 5em;
}
</style>
</head>
<body><h1>Wikidata QRank</h1>

<p>QRank is ranking <a href="https://www.wikidata.org/">Wikidata entities</a>
by aggregating page views on Wikipedia, Wikispecies, Wikibooks, Wikiquote,
and other Wikimedia projects. Only requests from actual users get counted;
automated requests from known bots and spiders do not influence the ranking.
For background about QRank, for examples that illustrate the ranking, and for
the source code of the program that computes QRank, please <a
href="https://github.com/brawer/wikidata-qrank">see here</a>.</p>

<p>To <b>download</b> the latest QRank data, <a href="/download/qrank.gz">click
here</a>.  The file gets updated periodically; please use
<a href="https://tools.ietf.org/html/rfc7232">HTTP Conditional
Requests</a> when checking for updates.
The QRank data is dedicated to the <b>Public Domain</b> via <a
href="https://creativecommons.org/publicdomain/zero/1.0/">Creative
Commons Zero 1.0</a>. To the extent possible under law, we have waived
all copyright and related or neighboring rights to this work. This
work is published from the United States.</p>

<p><img src="https://mirrors.creativecommons.org/presskit/buttons/88x31/svg/cc-zero.svg"
width="88" height="31" alt="Public Domain" style="float:left"/></p>

</body></html>`)
}

func handleDownloadQRank(w http.ResponseWriter, req *http.Request) {
	stats, qrankFile := dataLoader.Get()

	etag := fmt.Sprintf(`"%s"`, stats.QRankSha256)
	var lastModified time.Time
	if fstat, err := qrankFile.Stat(); err == nil {
		lastModified = fstat.ModTime()
	}

	w.Header().Set("ETag", etag)
	http.ServeContent(w, req, stats.QRankFilename, lastModified, qrankFile)
}