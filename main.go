package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Attack struct {
	Name     string
	Damage   int
	MaxUses  int
	UsesLeft int
}

type Pokemon struct {
	Name     string
	BaseName string
	HP       int
	MaxHP    int
	Attacks  []Attack
	Level    int
}

type Player struct {
	Pokemon Pokemon
	Potions map[string]int
	Gold    int
}

type Enemy struct {
	Name   string
	HP     int
	MaxHP  int
	MinDmg int
	MaxDmg int
}

func chooseStarter() Player {
	fmt.Println("Choisis ton starter :")
	fmt.Println("1. Salamèche ")
	fmt.Println("2. Carapuce ")
	fmt.Println("3. Bulbizarre ")

	var choice int
	fmt.Print("Ton choix : ")
	fmt.Scan(&choice)

	var poke Pokemon
	switch choice {
	case 1:
		poke = Pokemon{
			Name: "Salamèche", BaseName: "Salamèche", HP: 50, MaxHP: 50, Level: 1,
			Attacks: []Attack{
				{"Griffe", 10, -1, -1},
				{"Flammèche", 20, 3, 3},
			},
		}
	case 2:
		poke = Pokemon{
			Name: "Carapuce", BaseName: "Carapuce", HP: 55, MaxHP: 55, Level: 1,
			Attacks: []Attack{
				{"Charge", 9, -1, -1},
				{"Pistolet à O", 18, 3, 3},
			},
		}
	case 3:
		poke = Pokemon{
			Name: "Bulbizarre", BaseName: "Bulbizarre", HP: 60, MaxHP: 60, Level: 1,
			Attacks: []Attack{
				{"Fouet Lianes", 8, -1, -1},
				{"Tranch'Herbe", 17, 3, 3},
			},
		}
	default:
		poke = Pokemon{
			Name: "Salamèche", BaseName: "Salamèche", HP: 50, MaxHP: 50, Level: 1,
			Attacks: []Attack{
				{"Griffe", 10, -1, -1},
				{"Flammèche", 20, 3, 3},
			},
		}
	}

	fmt.Printf("Tu as choisi %s ! Veux-tu changer son nom ? (o/n) : ", poke.Name)
	var ans string
	fmt.Scan(&ans)
	if strings.ToLower(ans) == "o" {
		fmt.Print("Entre un nouveau nom : ")
		fmt.Scan(&poke.Name)
	}

	return Player{Pokemon: poke, Potions: map[string]int{"Petite Potion(+10)": 1, "Moyenne Potion(+25)": 1}, Gold: 20}
}

func showMap(posX, posY int, combats [4][2]int, done []bool, center [2]int) {
	grid := [5][5]string{}
	for i := range grid {
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}
	grid[posY][posX] = "X"
	for i, c := range combats {
		if !done[i] {
			grid[c[1]][c[0]] = "C"
		}
	}
	grid[center[1]][center[0]] = "P"

	fmt.Println("=== Carte ===")
	for _, row := range grid {
		for _, cell := range row {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
	fmt.Println("X=Joueur | C=Combat | P=Centre Pokémon")
	fmt.Println("1=haut 2=bas 3=droite 4=gauche")
}

func showFinalMap(_, posY int) {
	grid := [5][1]string{}
	for i := 0; i < 5; i++ {
		grid[i][0] = "."
	}
	grid[posY][0] = "X"
	grid[0][0] = "B"
	fmt.Println("=== Chemin vers le Boss ===")
	for _, row := range grid {
		fmt.Println(row[0])
	}
	fmt.Println("X=Joueur | B=Maître de la Ligue")
	fmt.Println("1=haut 2=bas")
}

func startWorld(player *Player) {
	posX, posY := 2, 2
	combats := [4][2]int{{0, 0}, {4, 0}, {0, 4}, {4, 4}}
	center := [2]int{2, 0}
	done := make([]bool, 4)
	names := []string{"Pikachu", "Racaillou", "Ponyta", "Coconfort"}
	wins := 0

	for {
		showMap(posX, posY, combats, done, center)
		var move int
		fmt.Print("Déplace-toi : ")
		fmt.Scan(&move)

		switch move {
		case 1:
			if posY > 0 {
				posY--
			}
		case 2:
			if posY < 4 {
				posY++
			}
		case 3:
			if posX < 4 {
				posX++
			}
		case 4:
			if posX > 0 {
				posX--
			}
		}

		// combats
		for i, c := range combats {
			if posX == c[0] && posY == c[1] && !done[i] {
				enemy := Enemy{names[i], 30, 30, 5, 12}
				if fight(player, enemy) {
					done[i] = true
					wins++
					player.Pokemon.Level += 2
					fmt.Printf("Ton Pokémon passe au niveau %d !\n", player.Pokemon.Level)
				}
			}
		}

		// centre Pokémon
		if posX == center[0] && posY == center[1] {
			centerMenu(player)
		}

		// après 4 combats -> évolution + combat final
		if wins == 4 {
			if player.Pokemon.Level >= 9 {
				evolvePokemon(&player.Pokemon)
			}
			startFinalPath(player)
			return
		}
	}
}

func centerMenu(player *Player) {
	for {
		fmt.Println("=== Centre Pokémon ===")
		fmt.Println("1. Soigner")
		fmt.Println("2. Acheter potions")
		fmt.Println("3. Quitter")
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			player.Pokemon.HP = player.Pokemon.MaxHP
			fmt.Println("Ton Pokémon est soigné !")
		case 2:
			fmt.Println("1. Petite Potion (+10 PV) - 5 Gold")
			fmt.Println("2. Moyenne Potion (+25 PV) - 10 Gold")
			var buy int
			fmt.Scan(&buy)
			if buy == 1 && player.Gold >= 5 {
				player.Potions["Petite Potion(+10)"]++
				player.Gold -= 5
				fmt.Println("Tu as acheté une Petite Potion !")
			} else if buy == 2 && player.Gold >= 10 {
				player.Potions["Moyenne Potion(+25)"]++
				player.Gold -= 10
				fmt.Println("Tu as acheté une Moyenne Potion !")
			} else {
				fmt.Println("Pas assez d'or !")
			}
		case 3:
			return
		}
	}
}

func startFinalPath(player *Player) {
	posX, posY := 0, 4
	for {
		showFinalMap(posX, posY)
		var move int
		fmt.Print("Déplace-toi : ")
		fmt.Scan(&move)
		if move == 1 && posY > 0 {
			posY--
		} else if move == 2 && posY < 4 {
			posY++
		}

		if posY == 0 {
			fmt.Println("⚔️ Combat final contre le Maître Pokémon !")
			enemy := Enemy{"Rayquaza", 150, 150, 12, 20}
			fight(player, enemy)
			fmt.Println("Félicitations, tu as gagné la Ligue Pokémon !")
			return
		}
	}
}

func fight(player *Player, enemy Enemy) bool {
	fmt.Printf("⚔️ %s sauvage apparaît ! (%d PV)\n", enemy.Name, enemy.HP)
	rand.Seed(time.Now().UnixNano())

	for player.Pokemon.HP > 0 && enemy.HP > 0 {
		fmt.Printf("\n%s (%d/%d PV, Niveau:%d)\n", player.Pokemon.Name, player.Pokemon.HP, player.Pokemon.MaxHP, player.Pokemon.Level)
		fmt.Printf("%s (%d/%d PV)\n", enemy.Name, enemy.HP, enemy.MaxHP)
		fmt.Println("1. Attaquer")
		fmt.Println("2. Potion")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Choisis une attaque :")
			for i, atk := range player.Pokemon.Attacks {
				if atk.MaxUses == -1 {
					fmt.Printf("%d. %s (%d dégâts, illimité)\n", i+1, atk.Name, atk.Damage)
				} else {
					fmt.Printf("%d. %s (%d dégâts, %d restant)\n", i+1, atk.Name, atk.Damage, atk.UsesLeft)
				}
			}
			var atkChoice int
			fmt.Scan(&atkChoice)
			if atkChoice > 0 && atkChoice <= len(player.Pokemon.Attacks) {
				atk := &player.Pokemon.Attacks[atkChoice-1]
				if atk.MaxUses == -1 || atk.UsesLeft > 0 {
					enemy.HP -= atk.Damage
					if atk.MaxUses != -1 {
						atk.UsesLeft--
					}
					fmt.Printf("%s utilise %s ! %s perd %d PV.\n", player.Pokemon.Name, atk.Name, enemy.Name, atk.Damage)
				}
			}
		case 2:
			fmt.Printf("1. Potion +10 (%d)\n", player.Potions["Petite Potion(+10)"])
			fmt.Printf("2. Potion +25 (%d)\n", player.Potions["Moyenne Potion(+25)"])
			var potChoice int
			fmt.Scan(&potChoice)
			if potChoice == 1 && player.Potions["Petite Potion(+10)"] > 0 {
				player.Pokemon.HP += 10
				if player.Pokemon.HP > player.Pokemon.MaxHP {
					player.Pokemon.HP = player.Pokemon.MaxHP
				}
				player.Potions["Petite Potion(+10)"]--
				fmt.Println("Petite Potion(+10) utilisée !")
			} else if potChoice == 2 && player.Potions["Moyenne Potion(+25)"] > 0 {
				player.Pokemon.HP += 25
				if player.Pokemon.HP > player.Pokemon.MaxHP {
					player.Pokemon.HP = player.Pokemon.MaxHP
				}
				player.Potions["Moyenne Potion(+25)"]--
				fmt.Println("Moyenne Potion(+25) utilisée !")
			}
		}

		if enemy.HP <= 0 {
			fmt.Printf("%s est vaincu !\n", enemy.Name)
			player.Gold += 10
			return true
		}

		dmg := rand.Intn(enemy.MaxDmg-enemy.MinDmg+1) + enemy.MinDmg
		player.Pokemon.HP -= dmg
		fmt.Printf("%s attaque et inflige %d dégâts !\n", enemy.Name, dmg)
		if player.Pokemon.HP <= 0 {
			fmt.Printf("%s est K.O. !\n", player.Pokemon.Name)
			return false
		}
	}
	return false
}

func evolvePokemon(p *Pokemon) {
	var newBase string
	var newAttacks []Attack

	switch p.BaseName {
	case "Salamèche":
		newBase = "Dracaufeu"
		newAttacks = []Attack{
			{"Lance-Flammes", 25, -1, -1},
			{"Déflagration", 40, 3, 3},
		}
	case "Carapuce":
		newBase = "Tortank"
		newAttacks = []Attack{
			{"Hydrocanon", 23, -1, -1},
			{"Surf", 38, 3, 3},
		}
	case "Bulbizarre":
		newBase = "Florizarre"
		newAttacks = []Attack{
			{"Fouet Lianes+", 22, -1, -1},
			{"Canon Graine", 37, 3, 3},
		}
	}

	// STATS APRES EVOLUTION
	p.BaseName = newBase
	p.Attacks = newAttacks
	p.MaxHP = 130
	p.HP = 130

	// Si le joueur n'avait pas renommé, on met le nouveau vrai nom
	if p.Name == "Salamèche" || p.Name == "Carapuce" || p.Name == "Bulbizarre" {
		p.Name = newBase
	}

	fmt.Printf(" %s est arrivé au niveau 9 donc il évolue en %s !\n", p.Name, p.BaseName)

	var ans string
	fmt.Printf("Souhaitez-vous changer son nom ? (o/n) : ")
	fmt.Scan(&ans)
	if strings.ToLower(ans) == "o" {
		fmt.Print("Entre un nouveau nom : ")
		fmt.Scan(&p.Name)
	}

	fmt.Printf("Nouvelles stats : %s (%d PV)\n", p.Name, p.MaxHP)
	fmt.Println("Attaques :")
	for _, atk := range p.Attacks {
		if atk.MaxUses == -1 {
			fmt.Printf("- %s : %d dégâts (illimité)\n", atk.Name, atk.Damage)
		} else {
			fmt.Printf("- %s : %d dégâts (%d utilisations)\n", atk.Name, atk.Damage, atk.MaxUses)
		}
	}
}

func main() {
	fmt.Println("=== Mini Pokémon CLI ===")
	player := chooseStarter()
	startWorld(&player)
	fmt.Println("Merci d'avoir joué ")
}
