package main

import (
    "log"
    "net/http"

    "job-crawler/apps/api/internal/bootstrap"
)

func main() {
    app := bootstrap.NewApp()
    log.Printf("api listening on %s", app.Config.HTTPAddr)
    if err := http.ListenAndServe(app.Config.HTTPAddr, app.Router); err != nil {
        log.Fatal(err)
    }
}
