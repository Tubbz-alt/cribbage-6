package strategy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func isIn(c model.Card, h []model.Card) bool {
	for _, thisCard := range h {
		if c == thisCard {
			return true
		}
	}
	return false
}

func validateHand(origHand, thisHand []model.Card) bool {
	for _, c := range thisHand {
		if !isIn(c, origHand) {
			return false
		}
	}
	cardCache := make(map[model.Card]int)
	for i, c := range thisHand {
		_, inMap := cardCache[c]
		if inMap {
			return false
		}
		cardCache[c] = i
	}
	return true
}

func factorial(n int) int {
	if n == 1 {
		return n
	}
	return n * factorial(n-1)
}

func nchoosek(n, k int) (int, error) {
	if k > n {
		return 0, errors.New(`k must be less than or equal to n`)
	}
	return factorial(n) / (factorial(k) * factorial(n-k)), nil
}

func generateHand(n int) []model.Card {
	hand := make([]model.Card, n)
	for i := 0; i < n; i++ {
		hand[i] = model.NewCardFromNumber(i)
	}
	return hand
}

func TestChooseFrom(t *testing.T) {
	tests := []struct {
		desc   string
		hand   []model.Card
		nCards int
		expErr bool
	}{{
		desc:   `6 choose 2`,
		hand:   generateHand(6),
		nCards: 2,
		expErr: false,
	}, {
		desc:   `6 choose 3`,
		hand:   generateHand(6),
		nCards: 3,
		expErr: false,
	}, {
		desc:   `6 choose 4`,
		hand:   generateHand(6),
		nCards: 4,
		expErr: false,
	}, {
		desc:   `5 choose 2`,
		hand:   generateHand(5),
		nCards: 2,
		expErr: false,
	}, {
		desc:   `5 choose 3`,
		hand:   generateHand(5),
		nCards: 3,
		expErr: false,
	}, {
		desc:   `5 choose 4`,
		hand:   generateHand(5),
		nCards: 4,
		expErr: false,
	}, {
		desc:   `5 choose 6`,
		hand:   generateHand(5),
		nCards: 6,
		expErr: true,
	}, {
		desc:   `choose zero cards`,
		hand:   generateHand(5),
		nCards: 0,
		expErr: true,
	}}
	for _, tc := range tests {
		all, err := chooseFrom(tc.nCards, tc.hand)
		if tc.expErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			expNum, err := nchoosek(len(tc.hand), tc.nCards)
			require.Nil(t, err)
			assert.Equal(t, expNum, len(all))
			for _, h := range all {
				assert.Equal(t, tc.nCards, len(h))
				ok := validateHand(tc.hand, h)
				assert.True(t, ok)
			}
		}
	}
}
