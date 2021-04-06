package walker

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"

	"github.com/google/uuid"
	. "github.com/heroiclabs/nakama-common/runtime"
)

func (m *Walker) beerLoop(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, mState *MatchState, messages []MatchData) {
	if !mState.started {
		return
	}

	mState.timer--
	if mState.timer == 0 {
		mState.timer = mState.tps

		b := NewBeer()
		mState.beers[b.ID] = b
		dispatcher.BroadcastMessage(2, b.Bytes(), nil, nil, false)
	}

	for _, message := range messages {
		switch message.GetOpCode() {
		case OpCodeClaimBeer:
			b := (&Beer{}).FromBytes(message.GetData())
			delete(mState.beers, b.ID)
			m.incrementScore(ctx, mState, message)
		}
	}
}

type Vector struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
type Beer struct {
	ID       string `json:"id"`
	Position Vector `json:"position"`
}

func RandomPosition() Vector {
	return Vector{
		X: rand.Intn(1240) + 20,
		Y: rand.Intn(680) + 20,
	}
}

func NewBeer() Beer {
	return Beer{
		ID:       uuid.NewString(),
		Position: RandomPosition(),
	}
}

func (b Beer) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(b)
	return buf.Bytes()
}

func (b *Beer) FromBytes(d []byte) *Beer {
	json.Unmarshal(d, b)
	return b
}
