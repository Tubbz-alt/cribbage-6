import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentGameID: '',
  currentGame: {},
  currentAction: {
    numShuffles: 0,
    selectedCards: [],
    percCut: 0.5,
  },
  loading: true,
};

const gameSlice = createSlice({
  name: 'game',
  initialState,
  reducers: {
    goToGame: {
      reducer: (state, action) => {
        state.loading = true;
        state.currentGameID = action.payload.id;
      },
      prepare: (id, history) => {
        return { payload: { id, history } };
      },
    },
    gameRetrieved(state, action) {
      state.loading = false;
      state.currentGame = action.payload.data;
      state.currentAction = initialState.currentAction;
      switch (state.currentGame.phase) {
        case `Deal`:
          // TODO leave numShuffles
          break;
      }
    },
    exitGame: {
      reducer: state => {
        state.loading = false;
        state.currentGameID = '';
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    refreshGame: {
      reducer: (state, action) => {
        if (state.currentGameID !== action.payload.id) {
          throw `bad game id: expected "${state.currentGameID}", got "${action.payload.id}"`;
        }
      },
      prepare: gameID => {
        return { payload: { id: gameID } };
      },
    },
    shuffleDeck(state) {
      isNaN(state.currentAction.numShuffles)
        ? (state.currentAction.numShuffles = 1)
        : (state.currentAction.numShuffles =
            state.currentAction.numShuffles + 1);
    },
    selectCard: {
      reducer: (state, action) => {
        // Nothing here?
        const card = action.payload.card;
        if (!card) {
          return;
        }

        const currentIndex = state.currentAction.selectedCards.findIndex(
          c => c.name === card.name,
        );
        const newSelected = [...state.currentAction.selectedCards];

        if (currentIndex === -1) {
          newSelected.push(card);
        } else {
          newSelected.splice(currentIndex, 1);
        }
        state.currentAction.selectedCards = newSelected;
      },
      prepare: (card, history) => {
        return { payload: { card, history } };
      },
    },
    dealCards: {
      reducer: (state, action) => {
        // Nothing here?
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    buildCrib: {
      reducer: (state, action) => {
        // Nothing here?
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    cutDeck: {
      reducer: (state, action) => {
        // Nothing here?
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    pegCard: {
      reducer: (state, action) => {
        // Nothing here?
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
  },
});

export const { actions, reducer, name: sliceKey } = gameSlice;
