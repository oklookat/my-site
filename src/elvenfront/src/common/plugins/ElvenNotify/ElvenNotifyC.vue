<template>
  <div class="notify-container">
    <transition-group class="notifications" name="notification-list" tag="ul">
        <li class="notification" :key="notification" v-for="notification in notifications"
              :class="(notification.self.type === 'error') ? 'not-error': ''">
          {{ notification.execute() }}
          <div class="not-message">
            {{ notification.self.message }}
          </div>
          <div class="not-timer-wrap">
            <div class="not-timer" v-bind:style="{ width: `${notification.percents}%` }">
            </div>
          </div>
        </li>
    </transition-group>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";

interface INotification {
  type: string
  message: string
}

export default defineComponent({
  name: "ElvenNotifyC",
  data() {
    return {
      // service vars start //
      SERVICE: 'ELVEN_NOTIFY_C',
      warn: false,
      error: false,
      info: false,
      success: false,
      // service vars end //
      notifications: [],
      notificationsCounter: 0,
    }
  },
  methods: {
    addError(message) {
      const notification: INotification = {
        type: 'error',
        message: message
      }
      this.setNotification(notification)
    },
    setNotification(notification: INotification) {
      const _this = this
      this.notifications.push({
        id: this.notificationsCounter,
        timeoutID: null,
        intervalID: null, // need for calc progress
        timeWhenGone: null, // ms when timeout ends
        percents: 0,
        self: notification, // notification object
        execute: function () {
          if (this.timeoutID) {
            // ---- already initialized
            return
          }
          // ---- init
          const deletedIn = 4000 // ms

          this.timeoutID = setTimeout(() => {
            // delete himself from array
            const index = _this.notifications.findIndex(obj => obj.id === this.id)
            _this.notifications.splice(index, 1)
            _this.notificationsCounter--
          }, deletedIn)

          // calc time to notification deleted
          this.timeWhenGone = new Date().getTime() + deletedIn // set time once item deleted
          this.intervalID = setInterval(() => {
            _this.calcPercents(this, deletedIn)
          }, 200)

          _this.notificationsCounter++
        }
      })
    },
    calcPercents(objContext, deletedIn) {
      const now = new Date().getTime()
      if (now >= objContext.timeWhenGone) {
        // if date when item should be deleted
        clearInterval(objContext.intervalID)
      }
      // get the difference between current date and time when item deleted
      const diff = Math.abs(now - objContext.timeWhenGone)
      // get how much is left as a percentage. deletedIn = 100%
      objContext.percents = Math.round((diff / deletedIn) * 100)
    },

  }
})
</script>

<style scoped>
.notify-container {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  width: 100%;
  margin-top: 8px;
}

.notifications {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 4px;
}

.notification {
  width: 220px;
  height: 52px;
  border-radius: 12px;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 1fr min-content;
}

.not-message {
  margin-top: 2px;
  margin-left: 12px;
}

.not-timer-wrap {
  width: 50%;
  display: flex;
  margin-bottom: 6px;
  justify-self: center;
}

.not-timer {
  transition: width 500ms linear;
  background-color: white;
  border-radius: 6px;
  height: 4px;
}

.notification.not-error {
  background-color: #FF0000;
  color: white;
}


/* animations start */
.notification {
  transition: all 0.4s;
}
.notification-list-enter-from,
.notification-list-leave-to {
  opacity: 0;
  float: left;
  width: 0;
}

.notification-list-leave-active {
  position: absolute;
}
/* animations end */
</style>