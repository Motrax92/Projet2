package main

import (
	"fmt"
	"strconv"
	"log"
	"net/http"
)

const None = ' '
var puissant = [7][6]rune{
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
	{None, None, None, None, None, None},
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello from Go web server"))
    })
        log.Println("Serveur lancé sur http://localhost:8080") 
        log.Fatal(http.ListenAndServe(":8080", nil))
}


func init() {
    // Initialise le plateau vide
    for i := 0; i < 7; i++ {
        for j := 0; j < 6; j++ {
            puissant[i][j] = None
        }
    }
}

func demanderCoup() int {
    var colone string
    fmt.Print("Quelle colonne (0-6) : ")
    _, err := fmt.Scanln(&colone)
    if err != nil {
        fmt.Println("Erreur :", err)
        return -1
    }

    coloneInt, err := strconv.Atoi(colone)
    if err != nil {
        fmt.Println("Erreur de conversion :", err)
        return -1
    }

    if coloneInt < 0 || coloneInt > 6 {
        fmt.Println("Colonne invalide")
        return -1
    }

    return coloneInt
}

func jouerCoup(coloneInt int) {
    // Trouver la première case vide depuis le bas
    for i := 5; i >= 0; i-- {
        if puissant[coloneInt][i] == None {
            puissant[coloneInt][i] = 'X'
            fmt.Println("Coup joué dans la colonne", coloneInt)
            afficherPlateau()
            return
        }
    }
    fmt.Println("Colonne pleine !")
}

func afficherPlateau() {
    for j := 0; j < 6; j++ {
        for i := 0; i < 7; i++ {
            fmt.Printf("[%c]", puissant[i][j])
        }
        fmt.Println()
    }
}

func verifierVictoire() bool {
    // Exemple : vérifie 4 en ligne horizontalement
    for y := 0; y < 6; y++ {
        for x := 0; x < 4; x++ {
            if puissant[x][y] != None &&
                puissant[x][y] == puissant[x+1][y] &&
                puissant[x][y] == puissant[x+2][y] &&
                puissant[x][y] == puissant[x+3][y] {
                fmt.Println("Victoire détectée sur la ligne", y)
                return true
            }
        }
    }
    return false
}

func chips() {
    for {
        col := demanderCoup()
        if col == -1 {
            continue
        }
        jouerCoup(col)
        if verifierVictoire() {
            fmt.Println("Partie terminée !")
            break
        }
    }
}
