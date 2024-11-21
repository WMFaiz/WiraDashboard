<template>
  <div class="dashboard-card">
    <div class="card-title">
      <div class="title-content">
        <img src="../assets/wira-logo-full.cc14df8b.png" width="100" height="100" alt="" srcset="">
        <div class="inputGroup">
          <div v-if="currentUsername === '' && !currentLoginStatus">
            <input class="loginButton" type="button" value="Login" @click="toggleLoginPage"/>
            <div v-if="currentLoginPageStatus">
              <Login />
            </div>
          </div>
          <div v-else>
            <input class="loginButton" type="button" value="Logout" @click="logout"/>
          </div>
          <h1>Players Ranks</h1>
          <!-- <input class="searchbar" type="search" placeholder="Search..."/> -->
        </div>
      </div>
    </div>
    <div class="card-container-main">
      <RankCardContainer title="World Ranks" :ranks="WorldRanks"/>
      <RankCardContainer title="Regional Ranks (MY)" :ranks="LocalRanks"/>
      <RankCardContainer title="Personal Ranks" :ranks="PersonalRanks"/>
    </div>
  </div>
</template>

<script>
import RankCardContainer from './RankCardContainer.vue'
import Login from './Login.vue'
import { globalData } from '@/utils/GlobalData'
export default {
  components: {
    RankCardContainer,
    Login
  },
  data () {
    return {
      username: '',
      isLoginPageDisabled: false,
      WorldRanks: [],
      LocalRanks: [],
      PersonalRanks: [],
      isLoggedIn: false
    }
  },
  watch: {
    isLoginPageDisabled: {
      handle: function (newValue) {
        globalData.isLoginPageDisabled = newValue
      },
      deep: true
    },
    username: {
      handle: function (newValue) {
        globalData.username = newValue
      },
      deep: true
    }
  },
  methods: {
    fetchWorldRanksData () {
      fetch('http://localhost:8181/api/rankings?count=10&page1', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          globalData.WorldRanks = data
        })
        .catch(error => {
          console.error('fetchWorldRanksData -->', error)
        })
    },
    fetchLocalRanksData () {
      fetch('http://localhost:8181/api/rankings?count=10&page=10', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          globalData.LocalRanks = data
        })
        .catch(error => {
          console.error('fetchWorldRanksData -->', error)
        })
    },
    fetchPersonalRanksData () {
      fetch('http://localhost:8181/api/rankings?username=' + this.currentUsername + '&count=8', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
        .then(response => response.json())
        .then(data => {
          globalData.PersonalRanks = data
        })
        .catch(error => {
          console.error('fetchWorldRanksData -->', error)
        })
    },
    toggleLoginPage () {
      globalData.isLoginPageDisabled = !globalData.isLoginPageDisabled
    },
    logout () {
      fetch('http://localhost:8181/api/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include'
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          } else {
            throw new Error('Logout failed')
          }
        })
        .then(() => {
          globalData.username = ''
          globalData.isLoggedIn = false
        })
        .catch(error => {
          console.error('Logout error -->', error)
        })
    },
    checkSession () {
      fetch('http://localhost:8181/api/check-session', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include'
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          } else {
            throw new Error('Session expired or not found')
          }
        })
        .then(() => {
          globalData.isLoggedIn = true
        })
        .catch(() => {
          globalData.isLoggedIn = false
        })
    }
  },
  computed: {
    currentLoginPageStatus () {
      return globalData.isLoginPageDisabled
    },
    currentUsername () {
      return globalData.username
    },
    currentPersonalRanks () {
      return globalData.PersonalRanks || []
    },
    currentLoginStatus () {
      return globalData.isLoggedIn
    }
  },
  mounted () {
    this.checkSession()
    this.fetchWorldRanksData()
    this.fetchLocalRanksData()
    this.fetchPersonalRanksData()
    if (this.currentUsername) {
      this.fetchPersonalRanksData()
    }
  }
}
</script>

<style lang="sass" scoped>
.dashboard-card
  color: white
  font-family: "Poppins", sans-serif
  .card-title
    background: rgba(30, 32, 31, 0.5)
    .title-content
      display: flex
      align-items: center
      img
        margin: 0 20px 0 10px
      .inputGroup
        width: 100%
        padding: 0 0 20px 0
        h1
          color: rgb(192, 171, 133)
          text-align: center
          justify-content: center
          align-items: center
        .loginButton
          float: right
          padding: 10px
          width: 75px
          font-weight: bold
          letter-spacing: 1px
          background: rgba(0, 0, 0, 0.5)
          border-radius: 0 0 10px 10px
          border-width: 2px 0 0 0
          border-color: rgb(192, 171, 133)
          color: rgb(192, 171, 133)
        .searchbar
          height: 10px
          background: rgba(0, 0, 0, 0.5)
          padding: 15px 20px 15px 20px
          color: rgb(192, 171, 133)
          font-size: 15px
          letter-spacing: 2px
          outline: none
          width: 200px
          border-color: rgb(192, 171, 133)
          border-width: 2px 0 0 2px
          border-radius: 0 10px 0 0
  .card-container-main
    display: flex
    justify-content: center
    flex-wrap: wrap
    margin: 0
    height: auto
</style>
