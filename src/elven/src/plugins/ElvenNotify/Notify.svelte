<script lang="ts">
  import type {
    ElvenNotify,
    Notification,
    NotificationFull,
  } from "@/plugins/ElvenNotify/types";
  import { quintOut } from "svelte/easing";
  import { crossfade } from "svelte/transition";

  // when notification will be deleted (in ms)
  const deletedIn = 5000;
  // max notifications on desktop
  const maxNotificationsD = 8;
  // max notifications on mobile
  const maxNotificationsM = 2;
  // this array displayed in component
  let notifications: NotificationFull[] = [];
  // used for set notification id
  let notificationsCounter = 0;

  /** plugin controls */
  class Plugin implements ElvenNotify {
    public add(notification: Notification) {
      add(notification);
    }
  }
  window.$elvenNotify = new Plugin();

  /** add user notification, then create full notification*/
  function add(n: Notification) {
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
    set(n);
  }

  /** create full notification and push to array */
  function set(notification: Notification) {
    const full: NotificationFull = {
      id: notificationsCounter++,
      self: notification, // notification object
      percents: 0,
      timeWhenGone: 0, // ms when timeout ends
      executed: false,
      timeoutID: null,
      intervalID: null, // need for calc progress
    };
    notifications = [...notifications, full];
  }

  function deleteNotification(context: NotificationFull) {
    const index = notifications.findIndex((obj) => obj.id === context.id);
    if (index > -1) {
      clearTimeouts(index);
      notifications.splice(index, 1);
      notifications = notifications;
    }
  }

  function clearTimeouts(index: number) {
    clearTimeout(notifications[index].timeoutID as unknown as number);
    clearInterval(notifications[index].intervalID as unknown as number);
  }

  function calcPercents(ctx: NotificationFull, deletedIn: number) {
    const now = new Date().getTime();
    // if date when item should be deleted
    if (now >= ctx.timeWhenGone) {
      clearInterval(ctx.intervalID);
    }
    // get the difference between current date and time when item deleted
    const diff = Math.abs(now - ctx.timeWhenGone);
    // get how much is left as a percentage. deletedIn = 100%
    ctx.percents = (diff / deletedIn) * 100;
  }

  /** init timers */
  function execute(n: NotificationFull): NotificationFull {
    // already initialized
    if (n.executed) {
      return n;
    }
    // init
    n.timeoutID = setTimeout(() => {
      // delete from array after time
      deleteNotification(n);
    }, deletedIn);
    // calc time when notification be deleted
    n.timeWhenGone = new Date().getTime() + deletedIn; // set time once item deleted
    n.intervalID = setInterval(() => {
      calcPercents(n, deletedIn);
      // make reactive
      notifications[notificationsCounter + 1] =
        notifications[notificationsCounter + 1];
    }, 100); // interval time = performance. timer transition time = this time + 20ms
    n.executed = true;
    return n;
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

<div class="notify">
  <div class="notify__notifications">
    {#each notifications.map(execute) as notification (notification.id)}
      <div
        class="notification"
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

<style lang="scss">
  @mixin _desktop {
    @media (min-width: 765px) {
      @content;
    }
  }

  .notify {
    width: 100%;
    z-index: 9999;
    bottom: 0;
    margin-bottom: 8px;
    position: fixed;
    overflow: hidden;
    &__notifications {
      width: 100%;
      gap: 8px;
      height: max-content;
      display: flex;
      flex-direction: column-reverse;
      box-sizing: border-box;
      position: relative;
    }
    @include _desktop() {
      margin-right: 12px;
      height: min-content;
      width: 224px;
      right: 0;
      bottom: 0;
      &__notifications {
        height: max-content;
        width: max-content;
        flex-direction: column;
      }
    }
  }

  .notification {
    cursor: pointer;
    border-radius: 6px;
    color: white;
    backdrop-filter: blur(15px) saturate(180%);
    width: 240px;
    min-height: 52px;
    border-radius: 6px;
    align-self: center;
    display: grid;
    grid-template-columns: 1fr;
    grid-template-rows: 1fr min-content;
    &__message {
      margin-top: 4px;
      margin-left: 12px;
      margin-right: 8px;
    }
    &__timer-wrapper {
      width: 50%;
      margin-bottom: 12px;
      margin-top: 12px;
      justify-self: center;
    }
    &__timer {
      border-radius: 4px;
      height: 4px;
      transition: transform 120ms linear;
      background-color: rgb(190, 190, 190);
    }
    @include _desktop() {
      position: relative;
      width: 214px;
      min-height: 52px;
      &__timer-wrapper {
        margin-bottom: 8px;
        margin-top: 8px;
      }
    }
    @media (prefers-color-scheme: light) {
      background-color: rgb(130, 130, 130);
      border: 1px solid rgb(120, 120, 120);
    }
    @media (prefers-color-scheme: dark) {
      background-color: rgb(50, 50, 50);
      border: 1px solid rgb(60, 60, 60);
    }
  }
</style>
