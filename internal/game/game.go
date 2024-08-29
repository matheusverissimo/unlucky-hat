package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	Name string
	Use  func(g *Game, target *Player)
}

type Player struct {
	Name string
	Health int
	GetShotTarget func() bool
}

type Hat struct {
	Items []*Item
}

func (h *Hat) GetItemOnTop() *Item {
	item := h.Items[len(h.Items) - 1]
	h.Items = h.Items[:(len(h.Items) - 1)]
	return item
}

func (h *Hat) ShuffleItemsInside(possibleItems []Item, totalItems int, maximumNumberOfSameItem int) {
	items := make([]*Item, 0)
	aux := make(map[string]int)

	for i := 0; i < totalItems; i++ {
		idx := rand.Intn(len(possibleItems))
		chosenItem := &possibleItems[idx]
		if count, ok := aux[chosenItem.Name]; ok && count == maximumNumberOfSameItem {
			if _, notRemoved := aux[chosenItem.Name]; !notRemoved {
				possibleItems = RemoveFromSlice(possibleItems, idx)
			}
			i--
			continue
		}
		items = append(items, chosenItem)
		if count, ok := aux[chosenItem.Name]; ok {
			aux[chosenItem.Name] = count + 1
		} else {
			aux[chosenItem.Name] = 1
		}
	}

	h.Items = items
}

func (h *Hat) IsEmpty() bool {
	return len(h.Items) == 0
}

type Game struct {
	Hat         *Hat
	TurnCounter int
	PlayerTurn  int
	Players     []Player
}

func (g *Game) IsFinished() bool {
	for _, p := range g.Players {
		if p.Health <= 0 {
			return true
		}
	}
	return false
}


func (g *Game) PrintHatItems() {
	aux := make(map[string]int)
	for _, i := range g.Hat.Items {
		if count, ok := aux[i.Name]; ok {
			aux[i.Name] = count + 1
		} else {
			aux[i.Name] = 1
		}
	}
	str := "Inside the magical hat are"
	idx := 0
	separator := ","
	for k, v := range aux {
		if idx == len(aux) - 2 {
			separator = " and"
		}
		if idx == len(aux) - 1 {
			separator = ""
		}
		word := k
		if v > 1 {
			word += "s"
		}
		str += fmt.Sprintf(" %d %s%s", v, word, separator)
		idx++
	}
	str += "."
	fmt.Println(str)
}

func (g *Game) PrintPlayersStatuses() {
	for _, v := range g.Players {
		fmt.Printf("* %s: %d HP\n", v.Name, v.Health)
	}
}

func (g *Game) PrintPlayerTurn() {
	fmt.Printf("The hat is handed to %s!\n", g.Players[g.PlayerTurn].Name)
}

func NewGame() *Game {
	return &Game{
		Hat: &Hat{
			Items: make([]*Item, 0),
		},
		Players: []Player{
			{Name: "Player", Health: 5, GetShotTarget: func() bool {
				fmt.Println("Are you shoting him or yourself (h/m)?")
				var choice string
				fmt.Scan(&choice)
				if choice == "h" {
					return true
				}
				return false
			}},
			{Name: "Bot", Health: 5, GetShotTarget: func () bool {
				var r = rand.Float32()
				fmt.Printf("Hmmm... Im thinking...\n")
				time.Sleep(time.Second * 3)
				if r > 0.5 {
					fmt.Printf("Imma cast it in myself!\n")
					return false
				}
				fmt.Printf("You're gonna feel it!\n")
				return true
			}},
		},
		PlayerTurn:  0,
		TurnCounter: 0,
	}
}

func (g *Game) Run() {
	fireball := Item{
		Name: "Fireball",
		Use: func(g *Game, target *Player) {
			fmt.Printf("%s lost 1 HP!\n", target.Name)
			target.Health--
		},
	}

	nothing := Item{
		Name: "Nothing",
		Use: func(g *Game, target *Player) {
			fmt.Println("Puff! There was nothing on the hat this turn.")
		},
	}

	possibleItems := []Item{fireball, nothing}

	g.Hat.ShuffleItemsInside(possibleItems, 6, 4)

	g.PrintHatItems()
	for !g.IsFinished() {
		if g.Hat.IsEmpty() {
			g.Hat.ShuffleItemsInside(possibleItems, 6, 4)
			fmt.Println("The hat got empty, items are getting shuffled again.")
			g.PrintHatItems()
		}
		fmt.Printf("== ROUND %d ==\n", g.TurnCounter + 1)
		g.PrintPlayerTurn()
		g.PrintPlayersStatuses()
		playerShootOther := g.Players[g.PlayerTurn].GetShotTarget()
		item := g.Hat.GetItemOnTop()
		if playerShootOther {
			item.Use(g, &g.Players[(g.PlayerTurn + 1) & 1])
			g.PlayerTurn = (g.PlayerTurn + 1) & 1
		} else {
			item.Use(g, &g.Players[(g.PlayerTurn)])
		}
		g.TurnCounter++
	}
}


func RemoveFromSlice[T any](slice []T, index int)[]T{
	return append(slice[:index], slice[index+1:]...)
}
