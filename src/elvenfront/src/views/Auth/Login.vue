<template>
  <div class="login-container">
    <div class="login-inputs">
      <div class="logo-text">elven</div>
      <input type="text" name="username" placeholder="Имя пользователя" v-model="username">
      <input type="password" name="password" placeholder="Пароль" v-model="password">
    </div>
    <button :disabled="!username || !password" class="auth-button" @click="makeLogin">
      Войти
    </button>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import AuthAdapter from "@/common/adapters/Main/AuthAdapter";

export default defineComponent({
  name: 'Login',
  data() {
    return {
      username: '',
      password: '',
    }
  },
  methods: {
    async makeLogin() {
      AuthAdapter.login(this.username, this.password)
          .then(() =>{
            this.$router.push({name: 'Index'})
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
  box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
}

.auth-button {
  border: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
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
