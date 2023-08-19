package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// --- 0. Main function ---
func main() {
	game := &Game{}
	// Firstly let's ask for the number of players (between 1-5)
	playersCount := 1
	for {
		fmt.Println("How many players (1-5) ? ")
		var prompt string
		_, err := fmt.Scanln(&prompt)
		if err != nil {
			panic(err)
		}
		count, err := strconv.Atoi(prompt)
		if err == nil && count >= 1 && count <= 5 {
			playersCount = count
			break
		}
		fmt.Println("Invalid number; please select between 1 up to 5 players")
	}
	// Sets up the game
	game.Configure(playersCount)
	// Start the game!
	game.Begin(playersCount)
}

// --- 1. Types and Consts ---

// Action represents the actions a player can take
type Action string

// List of the available actions
const (
	HIT         Action = "hit"
	STAND       Action = "stand"
	SCORE_LIMIT int    = 21
)

// Lists of possible ranks and suits for the cards
var ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var suits = []string{"Hearts", "Spades", "Diamonds", "Clubs"}

// Card is made of a Rank and a Suit
type Card struct {
	Rank string
	Suit string
}

// Deck struct represents a collection of cards
type Deck struct {
	Cards []Card
}

// Player struct contains all information to describe player's status
type Player struct {
	Hand   []Card
	Bet    int64
	Wallet float64
}

// Game struct holds all other structs that are necessary to run the game
type Game struct {
	Deck    Deck
	Dealer  Player
	Players []Player
}

// --- 2. Methods to handle Card type ---

// Display displays the card as a combination of its rank and suit
func (c *Card) Display() string {
	return c.Rank + "/" + c.Suit
}

// Score gets card score based on the rank
func (c *Card) Score() int {
	switch c.Rank {
	case "A":
		return 11
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	default:
		return 10
	}
}

// --- 3. Methods to handle Deck type ---

// Draw removes a card from the top of the deck
// TODO: handle the case the deck is empty
func (d *Deck) Draw() (*Card, error) {
	if len(d.Cards) == 0 {
		return nil, errors.New("Deck is empty")
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card, nil
}

// Initialize initializes the deck
func (d *Deck) Initialize() {
	d.Cards = nil
	for _, suit := range suits {
		for _, rank := range ranks {
			d.Cards = append(d.Cards, Card{Suit: suit, Rank: rank})
		}
	}
}

// Shuffle randomizes the order of the cards in the deck
func (d *Deck) Shuffle() {
	maxIdx := len(d.Cards) - 1
	for idx := 0; idx < maxIdx; idx++ {
		currCard := d.Cards[idx]
		randIdx := rand.Intn(maxIdx-idx) + idx + 1
		d.Cards[idx] = d.Cards[randIdx]
		d.Cards[randIdx] = currCard
	}
}

// --- 4. Methods to handle Player type ---

// Action asks for the user's action
func (p *Player) Action() (action Action) {
userInput:
	for {
		fmt.Println("Do you want to hit 'h' or stand 's'? ")
		var prompt string
		_, err := fmt.Scanln(&prompt)
		// TODO handle the case the cmd line buffer fails for any reason
		if err != nil {
			panic(err)
		}
		switch prompt {
		case "h", "H", "hit", "Hit", "HIT":
			action = HIT
			break userInput
		case "s", "S", "stand", "Stand", "STAND":
			action = STAND
			break userInput
		default:
			fmt.Println("Invalid option; please type H to hit or S to stand")
		}
	}
	return action
}

// Add adds a new card to the player's hand
func (p *Player) Add(card *Card) {
	p.Hand = append(p.Hand, *card)
}

// Balance updates user's wallet based on the round's results
func (p *Player) Balance(playerIdx int, earnings float64, isPush bool) {
	if isPush {
		fmt.Println("Player #", playerIdx+1, "won back $", p.Bet, "this round")
	} else {
		fmt.Println("Player #", playerIdx+1, "won $", earnings-float64(p.Bet), "this round")
	}
	p.Wallet += earnings
}

// Bid asks the player to input their bid
func (p *Player) Bid() {
	var bid int64
	for {
		walletStr := fmt.Sprintf("%.2f", p.Wallet)
		fmt.Println()
		fmt.Println("Your Wallet: $" + walletStr + ". How much would you like to bet? ")
		var prompt string
		_, err := fmt.Scanln(&prompt)
		if err != nil {
			panic(err)
		}
		amount, err := strconv.ParseInt(prompt, 10, 64)
		if err == nil && amount > 0 && amount <= int64(p.Wallet) {
			bid = amount
			break
		}
		fmt.Println("Invalid bet; please input a new value")
	}
	p.Bet = bid
	p.Wallet -= float64(bid)
}

// InitialHand draws the initial cards for each player
func (p *Player) InitialHand(d *Deck, numberOfCards int) error {
	p.Hand = nil
	for idx := 0; idx < numberOfCards; idx++ {
		card, err := d.Draw()
		if err != nil {
			return err
		}
		p.Add(card)
	}
	return nil
}

// IsBlackjack checks if player's hand is a Blackjack
func (p *Player) IsBlackjack() bool {
	pScore := p.Score()
	if pScore == 21 && len(p.Hand) == 2 {
		return true
	}
	return false
}

// Score scores the hand total of a player
func (p *Player) Score() int {
	score := 0
	aces := 0
	for i := 0; i < len(p.Hand); i++ {
		score += p.Hand[i].Score()
		if p.Hand[i].Rank == "A" {
			aces++
		}
	}
	for score > SCORE_LIMIT && aces > 0 {
		score -= 10
		aces--
	}
	return score
}

// ShowHand displays all cards in player's or dealer's hand
func (p *Player) ShowHand(playerIdx int, isDealer bool) {
	if isDealer {
		fmt.Println("Dealer Hand: ")
	} else {
		fmt.Println("Player #", playerIdx+1, "Hand: ")
	}
	for _, card := range p.Hand {
		fmt.Println(card.Display())
	}
	if isDealer {
		fmt.Println("Dealer Total Score: ", p.Score())
	} else {
		fmt.Println("Total Score: ", p.Score())
	}
	fmt.Println()
}

// ShowTable displays all cards in the table - dealer and current user
func (p *Player) ShowTable(playerIdx int, dealerCard Card) {
	// display player's index
	fmt.Println()
	fmt.Println("Player #", playerIdx+1)
	fmt.Println()
	// display dealer's top card
	fmt.Println("Dealer Hand: ")
	fmt.Println(dealerCard.Display())
	fmt.Println()
	// display current player's cards
	fmt.Println("Your Hand: ")
	for _, card := range p.Hand {
		fmt.Println(card.Display())
	}
	fmt.Println()
}

// --- 5. Methods to handle Game type ---

// Begin starts a new game play
func (g *Game) Begin(numOfPlayers int) {
	// Shuffle the deck
	g.Deck.Shuffle()
	// Ask for the initial bids
	for idx := range g.Players {
		player := &g.Players[idx]
		fmt.Println()
		fmt.Println()
		fmt.Println("Getting bid for player", idx+1)
		player.Bid()
	}
	// Draw 2 cards to each player
	for idx := range g.Players {
		player := &g.Players[idx]
		err := player.InitialHand(&g.Deck, 2)
		if err != nil {
			panic(err)
		}
	}
	// Draw 2 cards to the dealer
	err := g.Dealer.InitialHand(&g.Deck, 2)
	if err != nil {
		panic(err)
	}
	// Check if dealer has a Blackjack
	if g.Dealer.IsBlackjack() {
		fmt.Println("Dealer has a Blackjack!")
		g.Dealer.ShowHand(-1, true)
		// Check what other players won this round
		g.Payout()
		// Display round's results
		g.Results()
		if g.Continue() {
			g.Begin(numOfPlayers)
		}
	}
	// Game continues, check each player...
	for idx := range g.Players {
		currPlayer := &g.Players[idx]
		if currPlayer.IsBlackjack() {
			continue
		}
		g.Round(idx)
	}
	// Dealer's turn
	g.Round(-1)
	// Check earnings on this round and update the wallets
	g.Payout()
	// Display round's results
	g.Results()
	if g.Continue() {
		g.Begin(numOfPlayers)
	}
}

// Configure prepares the game with its initial configuration
func (g *Game) Configure(numOfPlayers int) {
	// Creates the players - each of them starts with $1000 in the wallet
	for idx := 0; idx < numOfPlayers; idx++ {
		player := &Player{}
		player.Wallet += 1000
		g.Players = append(g.Players, *player)
	}
	fmt.Println(numOfPlayers, "player(s) have been added with 1000 in the wallet each of them. Let's start the game!")
	// Creates the dealer
	g.Dealer = Player{}
	// Creates the deck and shuffles it
	g.Deck = Deck{}
	g.Deck.Initialize()
}

// Continue asks if the user wants to still play the game
func (g *Game) Continue() (keepPlaying bool) {
	fmt.Println()
stillPlay:
	for {
		fmt.Println("Do you want to start a new round? 'y' or 'n'")
		var prompt string
		_, err := fmt.Scanln(&prompt)
		// TODO handle the case the cmd line buffer fails for any reason
		if err != nil {
			panic(err)
		}
		switch prompt {
		case "y", "Y", "yes", "Yes", "YES":
			keepPlaying = true
			break stillPlay
		case "n", "N", "no", "No", "NO":
			fmt.Println("Thanks for playing!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option; please type Y to continue or N to exit")
		}
	}
	return
}

// Payout calculates and distributes the cash to each player when the round is over
// the payout calculation may vary based on multiple different rules
func (g *Game) Payout() {
	dScore := g.Dealer.Score()
	for idx := range g.Players {
		player := &g.Players[idx]
		pBet := float64(player.Bet)
		player.ShowHand(idx, false)
		// Player wins with Blackjack: pays 3 to 2 the bet
		if player.IsBlackjack() && !g.Dealer.IsBlackjack() {
			fmt.Println("Player #", idx+1, "has a Blackjack!")
			fmt.Println("Player #", idx+1, "wins!")
			player.Balance(idx, pBet*2.5, false)
			continue
		}
		pScore := player.Score()
		// Player wins but with no Blackjack: player earns double the bet
		if pScore <= SCORE_LIMIT && (dScore > SCORE_LIMIT || pScore > dScore) {
			fmt.Println("Player #", idx+1, "wins!")
			player.Balance(idx, pBet*2, false)
			continue
		}
		// Player lost, dealer won; no changes in their wallet
		if (pScore > SCORE_LIMIT && dScore <= SCORE_LIMIT) ||
			(pScore <= SCORE_LIMIT && pScore < dScore) {
			fmt.Println("Dealer defeats Player #", idx+1, "!")
			continue
		}
		// Otherwise, both dealer and player went over 21,
		// or both have a Blackjack;
		// in this case we return the bet to the wallet
		if player.IsBlackjack() {
			fmt.Println("Player #", idx+1, "has a Blackjack!")
		}
		fmt.Println("It's a push...")
		player.Balance(idx, pBet, true)
	}
}

// Results displays the results of the round
func (g *Game) Results() {
	for idx := range g.Players {
		player := &g.Players[idx]
		walletStr := fmt.Sprintf("%.2f", player.Wallet)
		fmt.Println("Player #", idx+1, "Wallet: $", walletStr)
	}
}

// Round performs the current round for the dealer or the player
// playerIdx < 0 means it's dealer's turn
func (g *Game) Round(playerIdx int) {
	if playerIdx < 0 {
		// Dealer stands for any score greater than 16
		stand := false
		for !stand {
			if g.Dealer.Score() > 16 {
				stand = true
			} else {
				card, err := g.Deck.Draw()
				if err != nil {
					panic(err)
				}
				g.Dealer.Add(card)
			}
		}
		g.Dealer.ShowHand(-1, true)
		return
	}
	// Player's round
	currPlayer := &g.Players[playerIdx]
	stand := false
	for !stand {
		currPlayer.ShowTable(playerIdx, g.Dealer.Hand[1])
		action := currPlayer.Action()
		switch action {
		case HIT:
			card, err := g.Deck.Draw()
			if err != nil {
				panic(err)
			}
			currPlayer.Add(card)
			if currPlayer.Score() > SCORE_LIMIT {
				fmt.Println("Player #", playerIdx+1, "is over 21!")
				currPlayer.ShowHand(playerIdx, false)
				stand = true
			}
		case STAND:
			stand = true
		}
	}
}
