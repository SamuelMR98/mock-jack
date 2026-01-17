package game

import (
	"fmt"
	"math/rand"
	"time"
)

type State int

const (
	WaitingDeal State = iota
	PlayerTurn
	DealerTurn
	RoundOver
)

type Suit int
type Rank int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

const (
	Ace Rank = 1
	Two Rank = 2
	Three Rank = 3
	Four Rank = 4
	Five Rank = 5
	Six Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine Rank = 9
	Ten Rank = 10
	Jack Rank = 11
	Queen Rank = 12
	King Rank = 13
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	rank := map[Rank]string{
		Ace:   "A",
		Two:   "2",
		Three: "3",
		Four:  "4",
		Five:  "5",
		Six:   "6",
		Seven: "7",
		Eight: "8",
		Nine:  "9",
		Ten:   "10",
		Jack:  "J",
		Queen: "Q",
		King:  "K",
	}[c.Rank]

	suit := map[Suit]string{
		Clubs:    "♣",
		Diamonds: "♦",
		Hearts:   "♥",
		Spades:   "♠",
	}[c.Suit]

	return rank + suit
}

type Hand struct {
	Cards []Card
}

func (h *Hand) Clear() { h.Cards = h.Cards[:0] }
func (h *Hand) Add(c Card) { h.Cards = append(h.Cards, c) }
func (h Hand) String() string {
	if len(h.Cards) == 0 { return "<empty>" }
	out := ""
	for i, c := range h.Cards {
		if i > 0 { out += " " }
		out += c.String()
	}
	return out
}

// Value computes the blackjack value of the hand and whether it is a soft hand.
func (h Hand) Value() (best int, isSoft bool) {
	total := 0
	aces := 0
	for _, card := range h.Cards {
		switch {
			case card.Rank == Ace:
				aces++
				total += 11
			case card.Rank >= Ten:
				total += 10
			default:
				total += int(card.Rank)
		}
	}

	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}

	// If there are remaining aces counted as 11, it's a soft hand.
	isSoft = false
	hard := 0
	aceCount := 0
	for _, card := range h.Cards {
		if card.Rank == Ace {
			aceCount++
			hard += 1
		} else if card.Rank >= Ten {
			hard += 10
		} else {
			hard += int(card.Rank)
		}
	}
	if aceCount > 0 && hard+10 <= 21 {
		isSoft = true
	}

	return total, isSoft
}

type Deck struct {
	cards []Card
	rng *rand.Rand
	shoe int
}

func NewDeck(shoe int) *Deck {
	if shoe < 1 {
		shoe = 1
	}
	d := &Deck{
		shoe: shoe,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.reset()
	return d
}

func (d *Deck) reset() {
	d.cards = d.cards[:0]
	for s := Clubs; s <= Spades; s++ {
		for r := Ace; r <= King; r++ {
			for i := 0; i < d.shoe; i++ {
				d.cards = append(d.cards, Card{Suit: s, Rank: r})
			}
		}
	}
	d.shuffle()
}

func (d *Deck) shuffle() {
	// Fisher-Yates shuffle
	n := len(d.cards)
	for i := n - 1; i > 0; i-- {
		j := d.rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Draw() Card {
	if len(d.cards) == 0 {
		d.reset()
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card
}

type Game struct {
	Deck   *Deck
	Player Hand
	Dealer Hand
	State  State
	Result string
}

func NewGame(shoe int) *Game {
	return &Game{
		Deck:  NewDeck(shoe),
		State: WaitingDeal,
	}
}

func (g *Game) Deal() {
	g.Player.Clear()
	g.Dealer.Clear()
	g.Result = ""
	g.State = PlayerTurn

	g.Player.Add(g.Deck.Draw())
	g.Dealer.Add(g.Deck.Draw())
	g.Player.Add(g.Deck.Draw())
	g.Dealer.Add(g.Deck.Draw())

	// Check for immediate blackjack
	playerValue, _ := g.Player.Value()
	dealerValue, _ := g.Dealer.Value()
	if playerValue == 21 || dealerValue == 21 {
		g.finishRound()
	}
}

func (g *Game) PlayerHit() {
	if g.State != PlayerTurn {
		return
	}
	g.Player.Add(g.Deck.Draw())
	playerValue, _ := g.Player.Value()
	if playerValue > 21 {
		g.finishRound()
	}
}

func (g *Game) PlayerStand() {
	if g.State != PlayerTurn {
		return
	}
	g.State = DealerTurn
	for {
		dealerValue, _ := g.Dealer.Value()
		// Dealer hits on soft 17
		// TODO: Implement rule variations if needed
		if dealerValue < 17 || (dealerValue == 17 && g.isDealerSoft()) {
			g.Dealer.Add(g.Deck.Draw())
		} else {
			break
		}
	}
	g.finishRound()
}

func (g *Game) isDealerSoft() bool {
	_, isSoft := g.Dealer.Value()
	return isSoft
}

func (g *Game) finishRound() {
	g.State = RoundOver
	playerValue, _ := g.Player.Value()
	dealerValue, _ := g.Dealer.Value()

	switch {
		case playerValue > 21:
			g.Result = fmt.Sprintf("Player busts (%d). Dealer wins.", playerValue)
		case dealerValue > 21:
			g.Result = fmt.Sprintf("Dealer busts (%d). Player wins!", dealerValue)
		case playerValue > dealerValue:
			g.Result = fmt.Sprintf("Player wins! (%d vs %d)", playerValue, dealerValue)
		case playerValue < dealerValue:
			g.Result = fmt.Sprintf("Dealer wins. (%d vs %d)", dealerValue, playerValue)
		default:
			g.Result = fmt.Sprintf("Push. (%d vs %d)", playerValue, dealerValue)
	}
}

