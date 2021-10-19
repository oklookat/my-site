<template>
  <div class="notify__container">
    <transition-group class="notify__notifications" name="notification__list" tag="div">
      <div
        class="notify__notification"
        v-on:click="deleteNotificationByID(notification.id)"
        :key="notification.id"
        v-for="notification in notifications"
      >
        {{ notification.execute() }}
        <div class="notification__message">{{ notification.self.message }}</div>
        <div class="notification__timer-wrapper">
          <div
            class="notification__timer"
            v-if="notification.percents"
            v-bind:style="{
              transform: `scaleX(${notification.percents / 100})`
            }"
          ></div>
        </div>
      </div>
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
  message: string
}

// used by plugin
const SERVICE = 'ELVEN_NOTIFY_C'

const notifications: Ref<Array<INotificationFull>> = ref([])
let notificationsCounter = 0
const deletedIn = 5000
// desktop
const maxNotificationsD = 8
// mobile
const maxNotificationsM = 2

// used by plugin.
function addNotification(message) {
  const notification = {
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
.notify__container {
  position: fixed;
  bottom: 0;
  width: 100%;
  margin-bottom: 8px;
  overflow: hidden;
}

.notify__notifications {
  height: max-content;
  width: 100%;
  display: flex;
  flex-direction: column-reverse;
  gap: 8px;
  box-sizing: border-box;
  position: relative;
}

.notify__notification {
  cursor: pointer;
  align-self: center;
  width: 240px;
  min-height: 52px;
  border-radius: 6px;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 1fr min-content;
}

.notification__message {
  margin-top: 4px;
  margin-left: 12px;
  margin-right: 8px;
}

.notification__timer-wrapper {
  width: 50%;
  justify-self: center;
  margin-bottom: 12px;
  margin-top: 12px;
}

.notification__timer {
  border-radius: 4px;
  transition: transform 120ms linear;
  background-color: rgba(255, 255, 255, 0.8);
  height: 4px;
}

.notify__notification {
  color: black;
  /*border: 1px solid rgba(255, 255, 255, 0.325);*/
  backdrop-filter: blur(15px) saturate(180%);
  border-radius: 6px;
}

.notify__notification > .notification__timer-wrapper > .notification__timer {
  background-color: rgb(190, 190, 190);
}

/* ANIMATIONS START */
.notify__notification {
  transition: all 0.4s;
}

@media (prefers-color-scheme: light) {
  .notify__notification {
    color: #fff;
    background-color: rgb(130, 130, 130);
    border: 1px solid rgb(120, 120, 120);
  }
}
@media (prefers-color-scheme: dark) {
  .notify__notification {
    color: #fff;
    background-color: rgb(50, 50, 50);
    border: 1px solid rgb(60, 60, 60);
  }
}

.notification__list-enter-from,
.notification__list-leave-to {
  opacity: 0;
}
/* ANIMATIONS END */

/* ADAPTIVE START */
@media screen and (min-width: 765px) {
  .notify__container {
    margin-right: 12px;
    height: min-content;
    width: 224px;
    right: 0;
    bottom: 0;
  }

  .notify__notifications {
    height: max-content;
    width: max-content;
    flex-direction: column;
  }

  .notify__notification {
    position: relative;
    width: 214px;
    min-height: 52px;
  }

  .notification__timer-wrapper {
    margin-bottom: 8px;
    margin-top: 8px;
  }
}
/* ADAPTIVE END */
</style>