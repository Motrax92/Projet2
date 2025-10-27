package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

const (
	cols = 7
	rows = 6
	None = ' ' // case vide
)

var grid [cols][rows]rune
var currentPlayer rune = '1' // 1 commence

func init() { resetGrid() }

func resetGrid() {
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			grid[x][y] = None
		}
	}
	currentPlayer = '1'
}

type Move struct {
	Column int `json:"column"`
}

// On renvoie des STRINGS (et pas des runes) -> le front affiche directement 🟡 / 🔴
type GameState struct {
	Grid    [][]string `json:"Grid"`
	Message string     `json:"Message"`
}

func main() {
	// ====== Fichiers statiques ======
	// /style/css.css -> lit dans le dossier "style"
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	// /static/... -> lit dans le dossier "static" (vidéos, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ====== Pages ======
	http.HandleFunc("/", serveHome)        // page du jeu -> static/index.html
	http.HandleFunc("/regle", serveRegle)  // page règles -> static/regle.html
	http.HandleFunc("/regle/", serveRegle) // idem si /regle/ (avec / final)
	http.HandleFunc("/contact", serveContact) // page contact -> static/contact.html
	http.HandleFunc("/contact/", serveContact)  // idem si /contact (sans .html)
	// ====== API jeu ======
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/reset", handleReset) // reset sans JS -> redirection

	log.Println("✅ Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/index.html"))
	_ = tpl.Execute(w, nil)
}

func serveRegle(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/regle.html"))
	_ = tpl.Execute(w, nil)
}
func serveContact(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/contact.html"))
	_ = tpl.Execute(w, nil)
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
	resp := GameState{Grid: exportGrid(), Message: message}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}

// ✅ Reset SANS JavaScript : POST /reset -> reset -> redirection vers "/"
func handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	resetGrid()
	// Redirige vers la home pour réafficher une grille vide
	http.Redirect(w, r, "/", http.StatusSeeOther) // 303 See Other
}

func playMove(column int) string {
	if column < 0 || column >= cols {
		return "Colonne invalide."
	}
	// déposer le jeton (gravité)
	for y := rows - 1; y >= 0; y-- {
		if grid[column][y] == None {
			grid[column][y] = currentPlayer

			// victoire ?
			if checkVictory(column, y, currentPlayer) {
				return string(currentPlayer) + " a gagné 🎉 !"
			}

			// match nul ?
			if isDraw() {
				return "🤝 Match nul !"
			}

			// tour suivant
			switchPlayer()
			return "Coup joué par " + string(grid[column][y])
		}
	}
	return "Colonne pleine !"
}

func switchPlayer() {
	if currentPlayer == '1' {
		currentPlayer = '2'
	} else {
		currentPlayer = '1'
	}
}

// Convertit la grille [7][6]rune -> [][]string avec "" / "🟡" / "🔴"
func exportGrid() [][]string {
	out := make([][]string, cols)
	for x := 0; x < cols; x++ {
		out[x] = make([]string, rows)
		for y := 0; y < rows; y++ {
			if grid[x][y] == None {
				out[x][y] = ""
			} else {
				out[x][y] = string(grid[x][y])
			}
		}
	}
	return out
}

func isDraw() bool {
	// s'il reste une case vide en haut de n'importe quelle colonne, pas nul
	for x := 0; x < cols; x++ {
		if grid[x][0] == None {
			return false
		}
	}
	return true
}

func checkVictory(x, y int, p rune) bool {
	// 4 directions : →, ↓, ↘, ↗  (on compte dans les 2 sens)
	dirs := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		count := 1
		count += countDir(x, y, d[0], d[1], p)
		count += countDir(x, y, -d[0], -d[1], p)
		if count >= 4 {
			return true
		}
	}
	return false
}

func countDir(x, y, dx, dy int, p rune) int {
	n := 0
	for {
		x += dx
		y += dy
		if x < 0 || x >= cols || y < 0 || y >= rows {
			break
		}
		if grid[x][y] != p {
			break
		}
		n++
	}
	return n
}
