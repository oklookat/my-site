<template>
  <div class="notify-container">
    <transition-group class="notifications" name="notification-list" tag="p">
      <span
        class="notification"
        v-on:click="deleteNotificationByID(notification.id)"
        :key="notification.id"
        v-for="notification in notifications"
        :class="classSetter(notification.self.type)"
      >
        {{ notification.execute() }}
        <div class="not-message">{{ notification.self.message }}</div>
        <div class="not-timer-wrap">
          <div
            class="not-timer"
            v-if="notification.percents"
            v-bind:style="{
              transform: `scaleX(${notification.percents / 100})`
            }"
          ></div>
        </div>
      </span>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { ref, Ref } from '@vue/reactivity'
interface INotificationFull {
  id: number,
  timeoutID: NodeJS.Timeout | null,
  intervalID: ReturnType<typeof setInterval> | null, // need for calc progress
  timeWhenGone: number, // ms when timeout ends
  percents: number,
  self: INotification, // notification object
  executed: boolean,
  execute: () => void
}
interface INotification {
  type: string
  message: string
}

const SERVICE = 'ELVEN_NOTIFY_C'
const notifications: Ref<Array<INotificationFull>> = ref([])
let notificationsCounter = 0
const deletedIn = 5000
// desktop
const maxNotificationsD = 8
// mobile
const maxNotificationsM = 2

function classSetter(errorType) {
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
}

function addNotification(type, message) {
  const notification = {
    type: type,
    message: message
  }
  if (notifications.value.length < 1) {
    notificationsCounter = 0
  }
  let isMaxNotifications = false
  if (window.screen.width > 765) {
    isMaxNotifications = notifications.value.length > maxNotificationsD - 1
  } else {
    isMaxNotifications = notifications.value.length > maxNotificationsM - 1
  }
  if (isMaxNotifications) {
    deleteNotificationByID(notifications[0].id)
  }
  setNotification(notification)
}

function setNotification(notification: INotification) {
  const fullNotification: INotificationFull = {
    id: notificationsCounter,
    timeoutID: null,
    intervalID: null, // need for calc progress
    timeWhenGone: 0, // ms when timeout ends
    percents: 0,
    self: notification, // notification object
    executed: false,
    execute: function () {
      if (this.executed) {
        // already initialized
        return
      }
      // init
      this.timeoutID = setTimeout(() => {
        // delete himself from array after time
        deleteNotification(this)
      }, deletedIn)
      // calc time when notification be deleted
      this.timeWhenGone = new Date().getTime() + deletedIn // set time once item deleted
      this.intervalID = setInterval(() => {
        calcPercents(this, deletedIn)
      }, 100) // interval time = performance. timer transition time = this time + 20ms
      this.executed = true
      notificationsCounter++
    }
  }
  notifications.value.push(fullNotification)
}

function deleteNotification(objContext) {
  const index = notifications.value.findIndex(obj => obj.id === objContext.id)
  if (index > -1) {
    clearTimeouts(index)
    notifications.value.splice(index, 1)
  }
}

function deleteNotificationByID(id) {
  const index = notifications.value.findIndex(obj => obj.id === id)
  if (index > -1) {
    clearTimeouts(index)
    notifications.value.splice(index, 1)
  }
}

function clearTimeouts(index) {
  clearTimeout(notifications.value[index].timeoutID as unknown as number)
  clearInterval(notifications.value[index].intervalID as unknown as number)
}

function calcPercents(objContext, deletedIn) {
  const now = new Date().getTime()
  // if date when item should be deleted
  if (now >= objContext.timeWhenGone) {
    clearInterval(objContext.intervalID)
  }
  // get the difference between current date and time when item deleted
  const diff = Math.abs(now - objContext.timeWhenGone)
  // get how much is left as a percentage. deletedIn = 100%
  objContext.percents = (diff / deletedIn) * 100
}

</script>

<style scoped>
.notify-container {
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