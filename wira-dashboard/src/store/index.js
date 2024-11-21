import { createStore } from 'vuex'

export default createStore({
  state: {
    WorldRanks: null,
    LocalRanks: null,
    PersonalRanks: null
  },
  mutations: {
    setWorldRanks (state, data) {
      state.WorldRanks = data
    },
    setLocalRanks (state, data) {
      state.LocalRanks = data
    },
    setPersonalRanks (state, data) {
      state.PersonalRanks = data
    }
  },
  actions: {
    fetchWorldRanks ({ commit }, counter) {
      fetch(`http://localhost:8181/api/rankings?count=10&page=${counter}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          commit('setWorldRanks', data)
        })
        .catch(error => {
          console.error('Error fetching WorldRanks:', error)
        })
    },
    fetchLocalRanks ({ commit }, counter) {
      fetch(`http://localhost:8181/api/rankings?count=10&page=${counter}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          commit('setLocalRanks', data)
        })
        .catch(error => {
          console.error('Error fetching LocalRanks:', error)
        })
    },
    fetchPersonalRanks ({ commit }, username) {
      fetch(`http://localhost:8181/api/rankings?username=${username}&count=8`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          commit('setPersonalRanks', data)
        })
        .catch(error => {
          console.error('Error fetching PersonalRanks:', error)
        })
    }
  },
  getters: {
    WorldRanks: state => state.WorldRanks,
    LocalRanks: state => state.LocalRanks,
    PersonalRanks: state => state.PersonalRanks
  }
})
