import { reactive } from 'vue'

export const globalData = reactive({
  WorldRanks: null,
  LocalRanks: null,
  PersonalRanks: null,
  isLoginPageDisabled: false,
  username: '',
  isLoggedIn: false
})
