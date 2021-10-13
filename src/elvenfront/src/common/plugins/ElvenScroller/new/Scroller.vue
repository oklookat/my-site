<template>
    <div class="scroller__container">
        <div ref="items" class="scroller__items"></div>
        <div id="scroller__templates">
            <div class="scroller__notTombstone">no tombstone</div>
            <div class="scroller__tombstone">tombstone</div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from '@vue/reactivity'
import { onMounted } from '@vue/runtime-core'
import { ScrollerLogic, ItemsTester } from './ScrollerLogic'
import { IScrollerSource } from './types'
const items = ref(null)


onMounted(() => {
    const source: IScrollerSource = {
        tombstone_: document.querySelector("#scroller__templates > .scroller__tombstone"),
        messageTemplate: document.querySelector("#scroller__templates > .scroller__notTombstone"),
        nextItem_: 0,

        fetch(count: number): Promise<Array<Object>> {
            // Fetch at least 30 or count more objects for display.
            count = Math.max(30, count)
            return new Promise((resolve, reject) => {
                // Assume 50 ms per item.
                setTimeout(() => {
                    let items = []
                    for (let i = 0; i < Math.abs(count); i++) {
                        items[i] = ItemsTester.getRandom()
                        this.nextItem++
                    }
                    resolve(Promise.all(items));
                }, 1000 /* Simulated 1 second round trip time */);
            });
        },

        createTombstone() {
            return this.tombstone_.cloneNode(true) as HTMLElement
        },

        render(item, div) {
            div = div || this.messageTemplate.cloneNode(true);
            div.textContent = item.toString()
            return div
        }

    }
    window.scroller = new ScrollerLogic(items.value, source)
})
</script>

<style scoped>
#scroller__templates {
    display: none;
}

.scroller__container {
    height: 100%;
    width: 100%;
    background-color: red;
}

.scroller__items {
    height: 100%;
    width: 100%;
    background-color: blue;
}

.scroller__tombstone {
    height: 50px;
    width: 150px;
    background-color: salmon;
}

.scroller__notTombstone {
    height: 50px;
    width: 150px;
    background-color: yellow;
}
</style>