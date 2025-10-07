package main

import (
	"fmt"
	"strconv"
	"log"
	"net/http"
)

puissant := [[ , , , , , ],
			[ , , , , , ],
			[ , , , , , ],
			[ , , , , , ],
			[ , , , , , ],
			[ , , , , , ],
			[ , , , , ,'x']]

func main() {
    var colone string
    fmt.Print("Quelle colonne : ")
    _, err := fmt.Scanln(&colone)
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }

    coloneInt, err := strconv.Atoi(colone)
    if err != nil {
        fmt.Println("Erreur de conversion :", err)
        return
    }

    if [coloneInt] [6] = None {

    }

    http.HandleFunc("/", func(w http. ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello from Go web server")) 
    })
        log.Println("Serveur lancé sur http://localhost:8080") 
        log.Fatal(http. ListenAndServe (":8080", nil))
}
