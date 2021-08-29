<template>
  <div class="login-container">
    <div class="login-inputs">
      <div class="logo-text">elven</div>
      <input type="text" name="username" placeholder="Имя пользователя" v-model="username" @input="trackInput">
      <input type="password" name="password" placeholder="Пароль" v-model="password" @input="trackInput">
    </div>
    <div class="auth-button" @click="makeLogin">
      Войти
    </div>
    <div class="error" v-if="isError">{{ error }}</div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import AuthAdapter from "@/common/adapters/Main/AuthAdapter";

export default defineComponent({
  name: 'Login',
  data() {
    return {
      isError: false,
      error: '',
      username: '',
      password: '',
    }
  },
  methods: {
    async makeLogin() {
      this.isError = false
      AuthAdapter.login(this.username, this.password)
          .then(() =>{
            this.$router.push({name: 'Index'})
          })
          .catch(error => {
            this.isError = true
            this.error = error
          })
    }
  }
})
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 18px;
}

.login-inputs {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 108px;
  gap: 14px;
}
.logo-text{
  margin-bottom: 16px;
}

input {
  height: 48px;
  width: 264px;
  box-shadow: 0px 0px 19px 0px rgba(34, 60, 80, 0.02);
}

.auth-button {
  box-shadow: 0px 0px 19px 0px rgba(34, 60, 80, 0.02);




  background-color: var(--color-level-1);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  user-select: none;
  cursor: pointer;
  font-size: 1.2rem;
  width: 264px;
  height: 48px;
}

</style>
