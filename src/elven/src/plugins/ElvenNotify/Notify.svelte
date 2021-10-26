<script lang="ts">
  import { quintOut } from "svelte/easing";
  import { crossfade } from "svelte/transition";

  interface INotificationFull {
    id: number;
    timeoutID: NodeJS.Timeout | null;
    intervalID: ReturnType<typeof setInterval> | null; // need for calc progress
    timeWhenGone: number; // ms when timeout ends
    percents: number;
    self: INotification; // notification object
    executed: boolean;
    execute: () => void;
  }
  interface INotification {
    message: string;
  }

  //// settings
  // when notification will be deleted (in ms)
  const deletedIn = 5000;
  // max notifications on desktop
  const maxNotificationsD = 8;
  // max notifications on mobile
  const maxNotificationsM = 2;

  //// main
  let notifications: Array<INotificationFull> = [];
  let notificationsCounter = 0;

  // used by plugin
  export let addNot: string = null;
  $: watchNotification(addNot);
  function watchNotification(value) {
    addNotification(value);
    addNot = null;
  }

  function addNotification(message) {
    if (!message) {
      return;
    }
    const notification = {
      message: message,
    };
    // clear counter if no notifications
    if (notifications.length < 1) {
      notificationsCounter = 0;
    }
    // adaptive
    let isMaxNotifications = false;
    if (window.screen.width > 765) {
      isMaxNotifications = notifications.length > maxNotificationsD - 1;
    } else {
      isMaxNotifications = notifications.length > maxNotificationsM - 1;
    }
    if (isMaxNotifications) {
      deleteNotification(notifications[0]);
    }
    // add notification
    setNotification(notification);
  }

  function setNotification(notification: INotification) {
    const fullNotification: INotificationFull = {
      id: notificationsCounter++,
      timeoutID: null,
      intervalID: null, // need for calc progress
      timeWhenGone: 0, // ms when timeout ends
      percents: 0,
      self: notification, // notification object
      executed: false,
      execute: function () {
        if (this.executed) {
          // already initialized
          return;
        }
        // init
        this.timeoutID = setTimeout(() => {
          // delete himself from array after time
          deleteNotification(this);
        }, deletedIn);
        // calc time when notification be deleted
        this.timeWhenGone = new Date().getTime() + deletedIn; // set time once item deleted
        this.intervalID = setInterval(() => {
          calcPercents(this, deletedIn);
          // make reactive
          notifications[notificationsCounter + 1] =
            notifications[notificationsCounter + 1];
        }, 100); // interval time = performance. timer transition time = this time + 20ms
        this.executed = true;
      },
    };
    notifications = [...notifications, fullNotification];
  }

  function deleteNotification(context: INotificationFull) {
    const index = notifications.findIndex((obj) => obj.id === context.id);
    if (index > -1) {
      clearTimeouts(index);
      notifications.splice(index, 1);
      notifications = notifications;
    }
  }

  function clearTimeouts(index) {
    clearTimeout(notifications[index].timeoutID as unknown as number);
    clearInterval(notifications[index].intervalID as unknown as number);
  }

  function calcPercents(objContext, deletedIn) {
    const now = new Date().getTime();
    // if date when item should be deleted
    if (now >= objContext.timeWhenGone) {
      clearInterval(objContext.intervalID);
    }
    // get the difference between current date and time when item deleted
    const diff = Math.abs(now - objContext.timeWhenGone);
    // get how much is left as a percentage. deletedIn = 100%
    objContext.percents = (diff / deletedIn) * 100;
  }

  function execute(notification: INotificationFull): INotificationFull {
    notification.execute();
    return notification;
  }

  // animation
  const [send, receive] = crossfade({
    duration: (d) => Math.sqrt(d * 200),

    fallback(node, params) {
      const style = getComputedStyle(node);
      const transform = style.transform === "none" ? "" : style.transform;

      return {
        duration: 400,
        easing: quintOut,
        css: (t) => `
          transform: ${transform} scale(${t});
					opacity: ${t}
				`,
      };
    },
  });
</script>

<div class="notify__container">
  <div class="notify__notifications">
    {#each notifications.map(execute) as notification (notification.id)}
      <div
        class="notify__notification"
        on:click={() => deleteNotification(notification)}
        in:receive={{ key: notification.id }}
        out:send={{ key: notification.id }}
      >
        <div class="notification__message">{notification.self.message}</div>
        <div class="notification__timer-wrapper">
          {#if notification.percents}
            <div
              class="notification__timer"
              style="transform: scaleX({notification.percents / 100})"
            />
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .notify__container {
    width: 100%;
    z-index: 9999;
    bottom: 0;
    margin-bottom: 8px;
    position: fixed;
    overflow: hidden;
  }

  .notify__notifications {
    width: 100%;
    gap: 8px;
    height: max-content;
    display: flex;
    flex-direction: column-reverse;
    box-sizing: border-box;
    position: relative;
  }

  .notify__notification {
    width: 240px;
    min-height: 52px;
    border-radius: 6px;
    cursor: pointer;
    align-self: center;
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
    margin-bottom: 12px;
    margin-top: 12px;
    justify-self: center;
  }

  .notification__timer {
    border-radius: 4px;
    height: 4px;
    transition: transform 120ms linear;
    background-color: rgba(255, 255, 255, 0.8);
  }

  .notify__notification {
    border-radius: 6px;
    color: black;
    backdrop-filter: blur(15px) saturate(180%);
  }

  .notify__notification > .notification__timer-wrapper > .notification__timer {
    background-color: rgb(190, 190, 190);
  }

  /* -------- MEDIA START -------- */
  /* ---- theming start ---- */
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
  /* ---- theming end ---- */
  /* ---- adaptive start ---- */
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
  /* ---- adaptive end ---- */
  /* -------- MEDIA END -------- */
</style>
