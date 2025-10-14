package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

const None = ' '

var grid [7][6]rune
var currentPlayer rune = '🟡'

func init() {
	resetGrid()
}

func resetGrid() {
	for i := 0; i < 7; i++ {
		for j := 0; j < 6; j++ {
			grid[i][j] = None
		}
	}
}

type Move struct {
	Column int `json:"column"`
}

type GameState struct {
	Grid    [7][6]rune
	Message string
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/move", handleMove)

	log.Println("✅ Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/index.html"))
	tpl.Execute(w, nil)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var move Move
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := playMove(move.Column)
	resp := GameState{Grid: grid, Message: message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func playMove(column int) string {
	if column < 0 || column >= 7 {
		return "Colonne invalide."
	}

	for i := 5; i >= 0; i-- {
		if grid[column][i] == None {
			grid[column][i] = currentPlayer
			if checkVictory() {
				return string(currentPlayer) + " a gagné 🎉 !"
			}
			switchPlayer()
			return "Coup joué par " + string(grid[column][i])
		}
	}
	return "Colonne pleine !"
}

func switchPlayer() {
	if currentPlayer == '🟡' {
		currentPlayer = '🔴'
	} else {
		currentPlayer = '🟡'
	}
}

func checkVictory() bool {
	for y := 0; y < 6; y++ {
		for x := 0; x < 4; x++ {
			if grid[x][y] != None &&
				grid[x][y] == grid[x+1][y] &&
				grid[x][y] == grid[x+2][y] &&
				grid[x][y] == grid[x+3][y] {
				return true
			}
		}
	}
	return false
}
