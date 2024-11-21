<template>
    <div v-if="currentLoginPageStatus">
        <div class="login-page">
            <div class="login">
                <div class="login-header">
                    <img src="../assets/wira-logo-full.cc14df8b.png" width="100" height="100"/>
                    <h1>
                        <span class="line">
                            Welcome
                        </span>
                        <span class="line">
                            Players
                        </span>
                    </h1>
                    <div v-if="!twoFacCodeEnabled">
                        <input class="login-header-close" type="button" value="X" @click="toggleLoginPage()"/>
                    </div>
                </div>
                <div class="login-ui">
                    <div class="username-group">
                        <label>Username: </label>
                        <input v-model="username" class="username-input" type="text"/>
                    </div>
                    <div class="password-group">
                        <label>Password: </label>
                        <input v-model="password" class="password-input" type="password"/>
                    </div>
                    <div class="login-group-button">
                        <div class="button-row">
                            <input class="login-submit" type="button" value="Login" @click="userLogin"/>
                            <input class="login-sign-up" type="button" value="Sign Up"/>
                        </div>
                        <input class="login-forgot-password" type="button" value="Forgot Password"/>
                    </div>
                </div>
                <!-- <div v-if="message" :class="messageClass">{{ message }}</div> -->
                <div v-if="currentTwoFacCodeEnabled">
                    <div class="twoFac">
                        <label>Enter your code</label>
                        <div class="twoFacCode">
                            <input v-for="(input, index) in 4" :key="index" type="text" maxlength="1" pattern="\d*" ref="inputs" @input="handleInput(index, $event)" />
                        </div>
                    </div>
                </div>
                <div class="sign-up-ui"></div>
            </div>
        </div>
    </div>
</template>

<script>
import { globalData } from '@/utils/GlobalData'
export default {
  data () {
    return {
      isLoginPageDisabled: false,
      twoFacCodeEnabled: false,
      username: '',
      password: '',
      message: '',
      messageClass: ''
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
    toggleLoginPage () {
      globalData.isLoginPageDisabled = !globalData.isLoginPageDisabled
    },
    async userLogin () {
      try {
        const response = await fetch('http://localhost:8181/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify({
            username: this.username,
            password: this.password
          })
        })

        if (!response.ok) {
          throw new Error('Login failed!')
        }

        const result = await response.json()
        if (result.message.includes('successful')) {
          this.twoFacCodeEnabled = true
        }

        this.message = result.message
        this.messageClass = 'success'
        // globalData.username = this.username
      } catch (error) {
        this.message = error.message || 'Server error'
        this.messageClass = 'error'
      }
    },
    handleInput (index, event) {
      const value = event.target.value
      if (value.length === 1) {
        const nextInput = this.$refs.inputs[index + 1]
        if (nextInput) {
          nextInput.focus()
        }
      }

      const allValueIsFilled = this.$refs.inputs.every(input => input.value.length === 1)
      if (allValueIsFilled) {
        this.twoFacCodeEnabled = false
        globalData.username = this.username
      }
    }
  },
  computed: {
    globalData () {
      return globalData
    },
    currentLoginPageStatus () {
      return globalData.isLoginPageDisabled
    },
    currentUsername () {
      return globalData.username
    },
    currentTwoFacCodeEnabled () {
      return this.twoFacCodeEnabled
    }
  }
}
</script>

<style lang="sass" scoped>
.login-page
    position: fixed
    display: flex
    top: 50%
    left: 50%
    justify-content: right
    align-items: right
    width: 100%
    height: 100%
    transform: translate(-50%, -50%)
    background-color: rgba(0,0,0,0.5)
    .login
        background-color: rgba(0,0,0,0.6)
        padding: 25px
        border-radius: 10px
        .twoFac
            display: flex
            flex-direction: column
            align-items: center
            justify-content: center
            margin-top: 50px
            color: rgb(192, 171, 133)
            .twoFacQR
                width: 200px
                height: 200px
                margin: 20px 0 20px 0
            .twoFacCode
                display: flex
                justify-content: center
                gap: 8px
                input[type=text]
                    width: 50px
                    height: 50px
                    background-color: rgba(192, 171, 133, 0.3)
                    border-width: 0 0 5px 0
                    border-color: rgb(192, 171, 133)
                    outline: none
                    text-align: center
                    color: white
                    font-size: 20px
                    margin-top: 20px
        .login-header
            display: flex
            align-items: center
            margin-bottom: 10vh
            .img
                margin-left: 200px
            h1
                color: rgb(192, 171, 133)
                margin-left: 20px
                display: flex
                flex-direction: column
                text-align: left
            .login-header-close
                position: fixed
                top: 0
                right: 0
                background-color: rgba(0,0,0,0.5)
                color: rgb(192, 171, 133)
                font-size: 17px
                width: 30px
                border-style: none
                margin: 10px
        .login-ui
            .username-group
                display: flex
                justify-content: space-between
                margin: 20px 0 10px 0
                label
                    color: rgb(192, 171, 133)
                .username-input
                    margin-left: 10px
                    border-color: rgb(192, 171, 133)
                    border-width: 0 0 2px 2px
                    background-color: rgba(0,0,0, 0.5)
                    outline: none
                    color: rgb(192, 171, 133)
                    padding: 0 5px 0 5px
                    width: 200px
            .password-group
                display: flex
                justify-content: space-between
                margin-bottom: 10px
                label
                    color: rgb(192, 171, 133)
                .password-input
                    margin-left: 10px
                    border-color: rgb(192, 171, 133)
                    border-width: 0 0 2px 2px
                    background-color: rgba(0,0,0, 0.5)
                    outline: none
                    color: rgb(192, 171, 133)
                    padding: 0 5px 0 5px
                    width: 200px
            .login-group-button
                display: flex
                flex-direction: column
                margin-top: 40px
                .login-forgot-password
                    width: 100%
                    text-align: center
                    border-width: 0 0 2px 0
                    border-color: rgb(192, 171, 133)
                    background-color: rgba(0,0,0,0.5)
                    color: rgb(192, 171, 133)
                    margin-bottom: 20px
                .button-row
                    display: flex
                    justify-content: space-between
                    margin-bottom: 10px
                    margin-top: 3px
                    margin-bottom: 30px
                    .login-submit
                        float: left
                        width: 100%
                        border-width: 0 0 2px 0
                        border-color: rgb(192, 171, 133)
                        background-color: rgba(0,0,0,0.5)
                        color: rgb(192, 171, 133)
                        margin-right: 10px
                    .login-sign-up
                        float: right
                        width: 100%
                        border-width: 0 0 2px 0
                        border-color: rgb(192, 171, 133)
                        background-color: rgba(0,0,0,0.5)
                        color: rgb(192, 171, 133)
</style>
