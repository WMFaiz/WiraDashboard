<template>
    <div class="card-container">
        <div class="card-container-header">
          <h3>{{ title }}</h3>
          <div class="search-bar">
            <input class="search-input" type="search" placeholder="Search..."/>
            <input class="search-button" type="button" value="Search"/>
            <input class="search-button-descending" type="button" value="v" @click="sortDescending"/>
            <input class="search-button-ascending" type="button" value="^" @click="sortAscending"/>
          </div>
        </div>
        <div class="pages">
          <input type="button" class="card-container-back" value="<" :disabled="currentPage === 1" @click="goToPreviousPage">
          <p>{{ currentPage }} / {{ totalPages }}</p>
          <input type="button" class="card-container-next" value=">" :disabled="currentPage === totalPages" @click="goToNextPage">
        </div>
        <div class="rank-container-scroller">
          <div v-if="title.includes('World')">
            <div v-if="currentRanks && currentRanks.length">
              <div v-for="(rank, index) in currentRanks" :key="index + rank.username + rank.c_class">
                <RankCard :username="rank.username" :c_class="CharacterClassConverter(rank.c_class)" :score="rank.score"/>
              </div>
            </div>
          </div>
          <div v-else-if="title.includes('Regional')">
            <div v-if="currentRanks && currentRanks.length">
              <div v-for="(rank, index) in currentRanks" :key="index + rank.username + rank.c_class">
                <RankCard :username="rank.username" :c_class="CharacterClassConverter(rank.c_class)" :score="rank.score"/>
              </div>
            </div>
          </div>
          <div v-else-if="title.includes('Personal')">
            <div v-if="currentRanks && currentRanks.length && currentUsername">
              <div v-for="(rank, index) in currentRanks" :key="index + rank.username + rank.c_class">
                <RankCard :username="rank.username" :c_class="CharacterClassConverter(rank.c_class)" :score="rank.score"/>
              </div>
            </div>
            <div class="rank-container-non-data" v-if="!currentUsername">
              <p>Require Login</p>
            </div>
          </div>
        </div>
    </div>
</template>

<script>
import RankCard from './RankCard.vue'
import { globalData } from '@/utils/GlobalData'
export default {
  props: {
    title: {
      type: String,
      required: true
    },
    ranks: {
      type: Array,
      required: true
    }
  },
  components: {
    RankCard
  },
  data () {
    return {
      sortedRanks: [],
      currentPage: 1,
      itemsPerPage: 10
    }
  },
  watch: {
    ranks: {
      handler (newRanks) {
        this.sortedRanks = [...newRanks]
      },
      immediate: true
    }
  },
  methods: {
    sortAscending () {
      this.sortedRanks = [...this.sortedRanks].sort((a, b) => a.reward_score - b.reward_score)
      this.resetPagination()
    },
    sortDescending () {
      this.sortedRanks = [...this.sortedRanks].sort((a, b) => b.reward_score - a.reward_score)
      this.resetPagination()
    },
    goToPreviousPage () {
      if (this.currentPage > this.totalPages) {
        this.currentPage--
        fetch('http://localhost:8181/api/rankings?count=10&page=' + this.currentPage, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        })
          .then(response => response.json())
          .then(data => {
            if (this.title.includes('World')) {
              globalData.WorldRanks = [...data]
            } else if (this.title.includes('Regional')) {
              globalData.LocalRanks = [...data]
            } else if (this.title.includes('Personal')) {
              globalData.PersonalRanks = [...data]
            }
          })
          .catch(error => {
            console.error('Error -->', error)
          })
      }
    },
    goToNextPage () {
      if (this.currentPage < this.totalPages) {
        this.currentPage++
        fetch('http://localhost:8181/api/rankings?count=10&page=' + this.currentPage, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        })
          .then(response => response.json())
          .then(data => {
            if (this.title.includes('World')) {
              globalData.WorldRanks = [...data]
            } else if (this.title.includes('Regional')) {
              globalData.LocalRanks = [...data]
            } else if (this.title.includes('Personal')) {
              globalData.PersonalRanks = [...data]
            }
          })
          .catch(error => {
            console.error('Error -->', error)
          })
      }
    },
    resetPagination () {
      this.currentPage = 1
    },
    CharacterClassConverter (cClass) {
      const classMap = {
        1: 'Warrior',
        2: 'Barbarian',
        3: 'Wizard',
        4: 'Sorcerer',
        5: 'Assassin',
        6: 'Archer',
        7: 'Cleric',
        8: 'Monk'
      }
      return classMap[cClass] || 'Unknown Class'
    }
  },
  computed: {
    globalData () {
      return globalData
    },
    totalPages () {
      if (this.title.includes('World')) {
        return 100
      } else if (this.title.includes('Regional')) {
        return 50
      } else if (this.title.includes('Personal')) {
        return 1
      } else {
        return 0
      }
    },
    currentRanks () {
      if (this.title.includes('World')) {
        return globalData.WorldRanks || []
      } else if (this.title.includes('Regional')) {
        return globalData.LocalRanks || []
      } else if (this.title.includes('Personal')) {
        return globalData.PersonalRanks || []
      }
      return []
    },
    currentUsername () {
      return globalData.username
    },
    paginatedRanks () {
      const start = (this.currentPage - 1) * this.itemsPerPage
      const end = start + this.itemsPerPage
      return this.sortedRanks.slice(start, end)
    },
    uniqueRanks () {
      const seenKeys = new Set()
      return this.sortedRanks.filter(rank => {
        const key = `${rank.username}-${rank.class_id}`
        if (seenKeys.has(key)) {
          return false
        }
        seenKeys.add(key)
        return true
      })
    }
  }
}
</script>

<style lang="sass" scoped>
.card-container-main
  display: flex
  justify-content: center
  flex-wrap: wrap
  margin: 0
  height: auto
  .card-container
    margin: 0 10px 0 10px
    padding: 0 0 20px 0
    .card-container-header
      display: block
      justify-content: center
      text-align: center
      border-style: outset
      border-width: 0 0 2px 0
      border-color: rgb(192, 171, 133)
      background: rgba(0, 0, 0, 0.3)
      padding-top: 1px
      margin-top: 20px
      color: rgb(192, 171, 133)
      .search-bar
        .search-input
          background: rgba(0, 0, 0, 0.5)
          padding: 5px 20px 5px 20px
          margin: 5px
          color: white
          font-size: 15px
          letter-spacing: 2px
          width: 200px
          outline: none
          border-color: rgb(192, 171, 133)
          border-width: 2px 0 0 2px
          border-radius: 0 10px 0 0
        input[type=button]
          background: rgba(0,0,0, 0.5)
          border-width: 1px 0 0 0
          border-color: rgb(192, 171, 133)
          border-radius: 0 0 10px 10px
          color: rgb(192, 171, 133)
          margin-right: 5px
          padding: 5px
          font-size: 15px
    .rank-container-scroller
      padding: 0 5px 5px 5px
      flex-grow: 1
      background: rgba(0, 0, 0, 0.3)
      box-sizing: border-box
      height: 450px
      max-height: 450px
      overflow: hidden
      .rank-container-non-data
        text-align: center
        justify-content: center
        align-items: center
        margin-top: 50%
        font-size: 20px
        color:  rgb(192, 171, 133)
    .pages
      display: flex
      justify-content: center
      align-items: center
      background: rgba(0,0,0, 0.3)
      height: 25px
      width: 100%
      p
        width: 100%
        text-align: center
        font-size: 15px
        color: rgb(192, 171, 133)
      input[disabled]
        background-color: brown
      .card-container-back
        background: rgba(0,0,0, 0.5)
        border-width: 0 0 0 3px
        border-color: rgba(192, 171, 133, 0.55)
        color: rgb(192, 171, 133)
        width: 40px
        height: 25px
        font-size: 15px
      .card-container-next
        background: rgba(0,0,0, 0.5)
        border-width: 0 3px 0 0
        border-color: rgb(192, 171, 133)
        color: rgb(192, 171, 133)
        width: 40px
        height: 25px
        font-size: 15px
</style>
