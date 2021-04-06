package walker

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	. "github.com/heroiclabs/nakama-common/runtime"
)

// MatchState represents a match's current authoratative state
type MatchState struct {
	tps          int
	users        map[string]MatchmakerEntry
	presences    map[string]Presence
	scores       map[string]int
	publishScore bool
	beers        map[string]Beer
	started      bool
	timer        int
}

// MatchInit initializes the match
func (m *Walker) MatchInit(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	tickRate := 30
	logger.WithField("params", params).Info("MatchInit")
	users := make(map[string]MatchmakerEntry)
	for _, e := range params["users"].([]MatchmakerEntry) {
		users[e.GetPresence().GetSessionId()] = e
	}

	return &MatchState{
		tps:       tickRate,
		users:     users,
		presences: make(map[string]Presence),
		scores:    make(map[string]int),
		beers:     make(map[string]Beer),
		timer:     tickRate,
	}, tickRate, "walker-" + ctx.Value(RUNTIME_CTX_MATCH_ID).(string)
}

// MatchLoop handles the main game loop
func (m *Walker) MatchLoop(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, state interface{}, messages []MatchData) interface{} {
	mState, _ := state.(*MatchState)

	m.beerLoop(ctx, logger, db, nk, dispatcher, tick, mState, messages)

	// Proxy Messages
	for _, message := range messages {
		dispatcher.BroadcastMessage(message.GetOpCode(), message.GetData(), nil, message, false)
	}

	// Publish Score
	if mState.publishScore {
		m.publishScore(ctx, mState, dispatcher)
		mState.publishScore = false
	}
	return mState
}

func (m *Walker) incrementScore(ctx context.Context, mState *MatchState, presence Presence) {
	u := mState.users[presence.GetSessionId()]
	mState.scores[u.GetProperties()["username"].(string)]++
	mState.publishScore = true
}

func (m *Walker) publishScore(ctx context.Context, mState *MatchState, dispatcher MatchDispatcher) {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(mState.scores)
	dispatcher.BroadcastMessageDeferred(OpCodeScores, buf.Bytes(), nil, nil, false)
}

// MatchJoinAttempt handles join attempts
func (m *Walker) MatchJoinAttempt(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, state interface{}, presence Presence, metadata map[string]string) (interface{}, bool, string) {
	return state, true, ""
}

// MatchJoin handles joins
func (m *Walker) MatchJoin(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, state interface{}, presences []Presence) interface{} {
	mState, _ := state.(*MatchState)
	mState.started = true
	for _, p := range presences {
		u := mState.users[p.GetSessionId()]
		mState.presences[p.GetSessionId()] = p
		mState.scores[u.GetProperties()["username"].(string)] = 0
	}
	return mState
}

// MatchLeave handles leaves
func (m *Walker) MatchLeave(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, state interface{}, presences []Presence) interface{} {
	mState, _ := state.(*MatchState)
	for _, p := range presences {
		delete(mState.presences, p.GetSessionId())
	}
	if len(mState.presences) == 0 {
		// terminate match
		return nil
	}
	return mState
}

// MatchTerminate handles deconstruction of the match
func (m *Walker) MatchTerminate(ctx context.Context, logger Logger, db *sql.DB, nk NakamaModule, dispatcher MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	message := "Server shutting down in " + strconv.Itoa(graceSeconds) + " seconds."
	dispatcher.BroadcastMessage(2, []byte(message), nil, nil, false)
	return state
}
