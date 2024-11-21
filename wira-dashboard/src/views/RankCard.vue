<template>
    <div class="rank-card">
        <div class="rank-card-container">
          <img :src="charIcon" width="35" height="35" alt=""/>
          <input type="text" class="rank-username" :value="username" maxlength="13" disabled/>
          <p class="rank-class">{{ c_class }}</p>
          <p class="rank-score">{{ score }}</p>
        </div>
    </div>
</template>

<script>
export default {
  props: {
    username: {
      type: String,
      required: true
    },
    c_class: {
      type: String,
      required: true
    },
    score: {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      charIcon: ''
    }
  },
  methods: {
    randomImge () {
      const imgSource = [require('../assets/adventurer.png'), require('../assets/magician.png'), require('../assets/swordsman.png'), require('../assets/wizard.png')]
      this.charIcon = imgSource[Math.floor(Math.random() * imgSource.length)]
    }
  },
  created () {
    this.randomImge()
  }
}
</script>

<style lang="sass" scoped>
    .rank-card
      float: left
      height: 27px
      margin-top: 17px
      background: rgb(30, 32, 31)
      width: 100%
      .rank-card-container
        display: flex
        align-items: center
        height: 27px
        width: 100%
        border-style: outset
        border-width: 2px 2px 2px 0
        border-color: rgb(192, 171, 133)
        background: rgba(0, 0, 0, 0.75)
        border-radius: 0 7px 10px 10px
        img
          border-style: outset
          border-radius: 100px
          border-width: 3px
          border-color: rgb(192, 171, 133)
          background: rgb(0, 0, 0)
        p
          padding: 1px 2px 1px 2px
        .rank-username
          white-space: nowrap
          text-align: left
          width: 13ch
          overflow: hidden
          text-overflow: ellipsis
          background: rgba(0,0,0,0.3)
          border: none
          color: white
          font-size: 15px
          margin-left: 5px
        .rank-class
          text-align: right
          width: 75px
          font-size: 15px
          margin-right: 15px
        .rank-score
          text-align: right
          width: 35px
          font-size: 15px
</style>
