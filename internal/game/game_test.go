package game_test

import (
	"fmt"
	"testing"

	g "github.com/matheusverissimo/unlucky-hat/internal/game"
)

func TestItemsShufflingInsideHat(t *testing.T) {
	hat := g.Hat{Items: make([]*g.Item, 0)}
	fireball := g.Item{
		Name: "Fireball",
		Use: func(g *g.Game, target *g.Player) {
			fmt.Printf("%s lost 1 HP!", target.Name)
			target.Health--
		},
	}

	nothing := g.Item{
		Name: "Nothing",
		Use: func(g *g.Game, target *g.Player) {
			fmt.Println("Puff! There was nothing on the hat this turn.")
		},
	}

	possibleItems := []g.Item{fireball, nothing}
	maxCountOfSameItem := 3
	for i := 0; i < 200; i++ {
		hat.ShuffleItemsInside(possibleItems, 6, maxCountOfSameItem)
		m := getItemsCounterMap(&hat)
		for itemName, count := range m {
			if count > maxCountOfSameItem {
				t.Errorf("There were %d %s when the max should be %d!", count, itemName, maxCountOfSameItem)
			}
		}
	}
}

func getItemsCounterMap(h *g.Hat) map[string]int {
	m := make(map[string]int)

	for _, i := range h.Items {
		if count, ok := m[i.Name]; ok {
			m[i.Name] = count + 1
		} else {
			m[i.Name] = 1
		}
	}

	return m
}
