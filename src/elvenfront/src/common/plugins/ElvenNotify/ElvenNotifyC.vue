<template>
  <div class="notify-container">
    <transition-group class="notifications" name="notification-list" tag="p">
      <span class="notification" v-on:click="deleteNotificationByID(notification.id)"
            :key="notification" v-for="notification in notifications"
            :class="classSetter(notification.self.type)">
        {{ notification.execute() }}
        <div class="not-message">
          {{ notification.self.message }}
        </div>
        <div class="not-timer-wrap">
          <div class="not-timer" v-if="notification.percents"
               v-bind:style="{
               transform: `scaleX(${notification.percents / 100})`}">
          </div>
        </div>
      </span>
    </transition-group>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";

export default defineComponent({
  name: "ElvenNotifyC",
  data() {
    return {
      SERVICE: 'ELVEN_NOTIFY_C',
      notifications: [],
      notificationsCounter: 0,
      deletedIn: 5000,
      maxNotificationsD: 8, // desktop
      maxNotificationsM: 2, // mobile
    }
  },
  methods: {
    classSetter(errorType) {
      switch (errorType) {
        case 'error':
          return 'not-error'
        case 'warn':
          return 'not-warn'
        case 'info':
          return 'not-info'
        case 'success':
          return 'not-success'
        default:
          return 'not-unknown'
      }
    },
    addNotification(type, message) {
      const notification = {
        type: type,
        message: message
      }
      if (this.notifications.length < 1) {
        this.notificationsCounter = 0
      }
      let isMaxNotifications = false
      if(window.screen.width > 765){
        isMaxNotifications = this.notifications.length > this.maxNotificationsD - 1
      } else {
        isMaxNotifications = this.notifications.length > this.maxNotificationsM - 1
      }
      if (isMaxNotifications) {
        this.deleteNotificationByID(this.notifications[0].id)
      }
      this.setNotification(notification)
    },
    setNotification(notification) {
      const _this = this
      const fullNotification = {
        id: this.notificationsCounter,
        timeoutID: null,
        intervalID: null, // need for calc progress
        timeWhenGone: null, // ms when timeout ends
        percents: null,
        self: notification, // notification object
        executed: false,
        execute: function () {
          if (this.executed) {
            // ---- already initialized
            return
          }
          // ---- init
          this.timeoutID = setTimeout(() => {
            // delete himself from array
            _this.deleteNotification(this)
          }, _this.deletedIn)
          // calc time to notification deleted
          this.timeWhenGone = new Date().getTime() + _this.deletedIn // set time once item deleted
          this.intervalID = setInterval(() => {
            _this.calcPercents(this, _this.deletedIn)
          }, 100) // time = performance and timer transition time = this time + 20ms
          _this.notificationsCounter++
          this.executed = true
        }
      }
      this.notifications.push(fullNotification)
    },
    deleteNotification(objContext) {
      const index = this.notifications.findIndex(obj => obj.id === objContext.id)
      if (index > -1) {
        this.clearTimeouts(index)
        this.notifications.splice(index, 1)
      }
    },
    deleteNotificationByID(id) {
      const index = this.notifications.findIndex(obj => obj.id === id)
      if (index > -1) {
        this.clearTimeouts(index)
        this.notifications.splice(index, 1)
      }
    },
    clearTimeouts(index) {
      clearTimeout(this.notifications[index].timeoutID)
      clearInterval(this.notifications[index].intervalID)
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
      objContext.percents = (diff / deletedIn) * 100
    },

  }
})
</script>

<style scoped>
.notify-container {
  /*pointer-events: none;*/
  position: fixed;
  bottom: 0;
  width: 100%;
  margin-bottom: 8px;
  overflow: hidden;
}


.notifications {
  height: max-content;
  width: 100%;
  display: flex;
  flex-direction: column-reverse;
  gap: 8px;
  box-sizing: border-box;
  position: relative;
}

.notification {
  cursor: pointer;
  align-self: center;
  width: 240px;
  min-height: 52px;
  border-radius: 6px;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 1fr min-content;
}

.not-message {
  margin-top: 4px;
  margin-left: 12px;
  margin-right: 8px;
}

.not-timer-wrap {
  width: 50%;
  justify-self: center;
  margin-bottom: 12px;
  margin-top: 12px;
}

.not-timer {
  border-radius: 4px;
  transition: transform 120ms linear;
  background-color: rgba(255, 255, 255, 0.8);
  height: 4px;
}

/* NOT TYPES STYLING START */
.notification {
  color: black;
  /*border: 1px solid rgba(255, 255, 255, 0.325);*/
  backdrop-filter: blur(15px) saturate(180%);
  border-radius: 6px;

}

.not-error {
  background-color: rgba(255, 0, 0, 0.35);
  border: 1px solid rgba(255, 0, 0, 0.25);
}

.notification.not-error > .not-timer-wrap > .not-timer {
  background-color: rgba(190, 13, 20, 1);
}

.not-warn {
  background-color: rgba(255, 240, 0, 0.75);
  border: 1px solid rgba(191, 179, 0, 0.35);
}

.notification.not-warn > .not-timer-wrap > .not-timer {
  background-color: rgba(190, 180, 0, 1);
}

.not-info {
  background-color: rgba(100, 200, 255, 0.75);
  border: 1px solid rgba(75, 150, 191, 0.35);
}

.notification.not-info > .not-timer-wrap > .not-timer {
  background-color: rgba(75, 150, 190, 1);
}

.not-success {
  background-color: rgba(140, 255, 50, 0.75);
  border: 1px solid rgba(105, 191, 38, 0.35);
}

.notification.not-success > .not-timer-wrap > .not-timer {
  background-color: rgba(105, 190, 140, 1);
}

/* NOT TYPES STYLING START */


/* ANIMATIONS START */
.notification {
  transition: all 0.4s;
}

.notification-list-enter-from,
.notification-list-leave-to {
  opacity: 0;
}

.notification-list-leave-active {

}

/* ANIMATIONS END */

/* ADAPTIVE START */
@media screen and (min-width: 765px) {
  .notify-container {
    margin-right: 12px;
    height: min-content;
    width: 224px;
    right: 0;
    bottom: 0;
  }

  .notifications {
    height: max-content;
    width: max-content;
    flex-direction: column;
  }

  .notification {
    position: relative;
    width: 214px;
    min-height: 52px;
  }

  .not-timer-wrap {
    margin-bottom: 8px;
    margin-top: 8px;
  }
}
/* ADAPTIVE END */
</style>