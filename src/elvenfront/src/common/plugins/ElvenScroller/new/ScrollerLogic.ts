/**
 * Copyright 2015 Google Inc. All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * 
 * changes made (https://github.com/oklookat): rewritten to TypeScript and modern JS
 */


import { IScrollerSource, TAnchorItem } from "./types"

export class ScrollerLogic {
    // Number of items to instantiate beyond current view in the scroll direction.
    private readonly RUNWAY_ITEMS = 50
    // Number of items to instantiate beyond current view in the opposite direction.
    private readonly RUNWAY_ITEMS_OPPOSITE = 10
    // The number of pixels of additional length to allow scrolling to.
    private readonly SCROLL_RUNWAY = 2000
    // The animation interval (in ms) for fading in content from tombstones.
    private readonly ANIMATION_DURATION_MS = 200
    // main vars
    private anchorItem: TAnchorItem = { index: 0, offset: 0 }
    private firstAttachedItem_: number = 0
    private lastAttachedItem_: number = 0
    private anchorScrollTop: number = 0
    private tombstoneSize_: number = 0
    private tombstoneWidth_: number = 0
    private tombstones_: Array<HTMLElement> = []
    private scroller_: HTMLElement
    private source_: IScrollerSource
    private items_: Array<any> = []
    private loadedItems_: number = 0
    private requestInProgress_: boolean = false
    private scrollRunway_: HTMLDivElement
    private scrollRunwayEnd_: number = 0

    /**
    * Construct an infinite scroller.
    * @param {HTMLElement} scroller The scrollable element to use as the infinite
    *     scroll region.
    * @param {IScrollerSource} source A provider of the content to be
    *     displayed in the infinite scroll region.
    */
    constructor(scrollable: HTMLElement, source: IScrollerSource) {
        this.scroller_ = scrollable
        this.source_ = source
        this.scroller_.addEventListener('scroll', this.onScroll_.bind(this))
        window.addEventListener('resize', this.onResize_.bind(this))
        // Create an element to force the scroller to allow scrolling to a certain
        // point.
        this.scrollRunway_ = document.createElement('div')
        // Internet explorer seems to require some text in this div in order to
        // ensure that it can be scrolled to.
        this.scrollRunway_.textContent = ' '
        this.scrollRunwayEnd_ = 0
        this.scrollRunway_.style.position = 'absolute'
        this.scrollRunway_.style.height = '1px'
        this.scrollRunway_.style.width = '1px'
        this.scrollRunway_.style.transition = 'transform 0.2s'
        this.scroller_.appendChild(this.scrollRunway_)
    }

    /**
    * Called when the browser window resizes to adapt to new scroller bounds and
    * layout sizes of items within the scroller.
    */
    private onResize_() {
        // TODO: If we already have tombstones attached to the document, it would
        // probably be more efficient to use one of them rather than create a new
        // one to measure.
        const tombstone = this.source_.createTombstone()
        tombstone.style.position = 'absolute'
        this.scroller_.appendChild(tombstone)
        tombstone.classList.remove('invisible')
        this.tombstoneSize_ = tombstone.offsetHeight
        this.tombstoneWidth_ = tombstone.offsetWidth
        this.scroller_.removeChild(tombstone)

        // Reset the cached size of items in the scroller as they may no longer be
        // correct after the item content undergoes layout.
        for (let i = 0; i < this.items_.length; i++) {
            this.items_[i].height = this.items_[i].width = 0
        }
        this.onScroll_()
    }

    /**
    * Called when the scroller scrolls. This determines the newly anchored item
    * and offset and then updates the visible elements, requesting more items
    * from the source if we've scrolled past the end of the currently available
    * content.
    */
    private onScroll_() {
        const delta = this.scroller_.scrollTop - this.anchorScrollTop
        // Special case, if we get to very top, always scroll to top.
        if (this.scroller_.scrollTop == 0) {
            this.anchorItem = { index: 0, offset: 0 }
        } else {
            this.anchorItem = this.calculateAnchoredItem(this.anchorItem, delta)
        }
        this.anchorScrollTop = this.scroller_.scrollTop
        const lastScreenItem = this.calculateAnchoredItem(this.anchorItem, this.scroller_.offsetHeight)
        if (delta < 0) {
            this.fill(this.anchorItem.index - this.RUNWAY_ITEMS, lastScreenItem.index + this.RUNWAY_ITEMS_OPPOSITE)
        } else {
            this.fill(this.anchorItem.index - this.RUNWAY_ITEMS_OPPOSITE, lastScreenItem.index + this.RUNWAY_ITEMS)
        }
    }

    /**
    * Calculates the item that should be anchored after scrolling by delta from
    * the initial anchored item.
    * @param {TAnchorItem} initialAnchor The initial position
    *     to scroll from before calculating the new anchor position.
    * @param {number} delta The offset from the initial item to scroll by.
    * @return {{index: number, offset: number}} Returns the new item and offset
    *     scroll should be anchored to.
    */
    private calculateAnchoredItem(initialAnchor: TAnchorItem, delta: number): TAnchorItem {
        if (delta == 0) {
            return initialAnchor
        }
        delta += initialAnchor.offset
        let i = initialAnchor.index
        let tombstones = 0
        if (delta < 0) {
            while (delta < 0 && i > 0 && this.items_[i - 1].height) {
                delta += this.items_[i - 1].height
                i--
            }
            tombstones = Math.max(-i, Math.ceil(Math.min(delta, 0) / this.tombstoneSize_))
        } else {
            while (delta > 0 && i < this.items_.length && this.items_[i].height && this.items_[i].height < delta) {
                delta -= this.items_[i].height
                i++
            }
            if (i >= this.items_.length || !this.items_[i].height)
                tombstones = Math.floor(Math.max(delta, 0) / this.tombstoneSize_)
        }
        i += tombstones
        delta -= tombstones * this.tombstoneSize_
        return {
            index: i,
            offset: delta,
        }
    }

    /**
    * Sets the range of items which should be attached and attaches those items.
    * @param {number} start The first item which should be attached.
    * @param {number} end One past the last item which should be attached.
    */
    private fill(start: number, end: number) {
        this.firstAttachedItem_ = Math.max(0, start)
        this.lastAttachedItem_ = end
        this.attachContent()
    }

    /**
    * Creates or returns an existing tombstone ready to be reused.
    * @return {Element} A tombstone element ready to be used.
    */
    private getTombstone(): HTMLElement {
        const tombstone = this.tombstones_.pop()
        if (tombstone) {
            tombstone.classList.remove('invisible')
            tombstone.style.opacity = '1'
            tombstone.style.transform = ''
            tombstone.style.transition = ''
            return tombstone
        }
        return this.source_.createTombstone()
    }

    /**
    * Attaches content to the scroller and updates the scroll position if
    * necessary.
    */
    private attachContent() {
        // Collect nodes which will no longer be rendered for reuse.
        // TODO: Limit this based on the change in visible items rather than looping
        // over all items.
        let i: number
        const unusedNodes: Array<HTMLElement> = []
        for (i = 0; i < this.items_.length; i++) {
            // Skip the items which should be visible.
            if (i == this.firstAttachedItem_) {
                i = this.lastAttachedItem_ - 1
                continue
            }
            if (this.items_[i].node) {
                if (this.items_[i].node.classList.contains('tombstone')) {
                    this.tombstones_.push(this.items_[i].node)
                    this.tombstones_[this.tombstones_.length - 1].classList.add('invisible')
                } else {
                    unusedNodes.push(this.items_[i].node)
                }
            }
            this.items_[i].node = null
        }

        const tombstoneAnimations = {}
        // Create DOM nodes.
        for (i = this.firstAttachedItem_; i < this.lastAttachedItem_; i++) {
            while (this.items_.length <= i)
                this.addItem_()
            if (this.items_[i].node) {
                // if it's a tombstone but we have data, replace it.
                if (this.items_[i].node.classList.contains('tombstone') &&
                    this.items_[i].data) {
                    // TODO: Probably best to move items on top of tombstones and fade them in instead.
                    if (this.ANIMATION_DURATION_MS) {
                        this.items_[i].node.style.zIndex = 1
                        tombstoneAnimations[i] = [this.items_[i].node, this.items_[i].top - this.anchorScrollTop]
                    } else {
                        this.items_[i].node.classList.add('invisible')
                        this.tombstones_.push(this.items_[i].node)
                    }
                    this.items_[i].node = null
                } else {
                    continue
                }
            }
            const node = this.items_[i].data ? this.source_.render(this.items_[i].data, unusedNodes.pop()) : this.getTombstone()
            // Maybe don't do this if it's already attached?
            node.style.position = 'absolute'
            this.items_[i].top = -1
            this.scroller_.appendChild(node)
            this.items_[i].node = node
        }

        // Remove all unused nodes
        while (unusedNodes.length) {
            this.scroller_.removeChild(unusedNodes.pop())
        }

        // Get the height of all nodes which haven't been measured yet.
        for (i = this.firstAttachedItem_; i < this.lastAttachedItem_; i++) {
            // Only cache the height if we have the real contents, not a placeholder.
            if (this.items_[i].data && !this.items_[i].height) {
                this.items_[i].height = this.items_[i].node.offsetHeight
                this.items_[i].width = this.items_[i].node.offsetWidth
            }
        }

        // Fix scroll position in case we have realized the heights of elements
        // that we didn't used to know.
        // TODO: We should only need to do this when a height of an item becomes
        // known above.
        this.anchorScrollTop = 0
        for (i = 0; i < this.anchorItem.index; i++) {
            this.anchorScrollTop += this.items_[i].height || this.tombstoneSize_
        }
        this.anchorScrollTop += this.anchorItem.offset

        // Position all nodes.
        let curPos = this.anchorScrollTop - this.anchorItem.offset
        i = this.anchorItem.index
        while (i > this.firstAttachedItem_) {
            curPos -= this.items_[i - 1].height || this.tombstoneSize_
            i--
        }
        while (i < this.firstAttachedItem_) {
            curPos += this.items_[i].height || this.tombstoneSize_
            i++
        }
        // Set up initial positions for animations.
        for (let i in tombstoneAnimations) {
            const anim = tombstoneAnimations[i]
            this.items_[i].node.style.transform = `translateY(${(this.anchorScrollTop + anim[1])}px scale(${(this.tombstoneWidth_ / this.items_[i].width)}, ${(this.tombstoneSize_ / this.items_[i].height)})`
            // Call offsetTop on the nodes to be animated to force them to apply current transforms.
            this.items_[i].node.offsetTop
            anim[0].offsetTop
            this.items_[i].node.style.transition = `transform ${this.ANIMATION_DURATION_MS}ms`
        }
        for (i = this.firstAttachedItem_; i < this.lastAttachedItem_; i++) {
            const anim = tombstoneAnimations[i]
            if (anim) {
                anim[0].style.transition = `transform ${this.ANIMATION_DURATION_MS}ms, opacity ${this.ANIMATION_DURATION_MS}ms`
                anim[0].style.transform = `translateY(${curPos}px) scale(${(this.items_[i].width / this.tombstoneWidth_)}, ${(this.items_[i].height / this.tombstoneSize_)})`
                anim[0].style.opacity = 0
            }
            if (curPos != this.items_[i].top) {
                if (!anim)
                    this.items_[i].node.style.transition = ''
                this.items_[i].node.style.transform = `translateY(${curPos}px)`
            }
            this.items_[i].top = curPos
            curPos += this.items_[i].height || this.tombstoneSize_
        }

        this.scrollRunwayEnd_ = Math.max(this.scrollRunwayEnd_, curPos + this.SCROLL_RUNWAY)
        this.scrollRunway_.style.transform = `translate(0, ${this.scrollRunwayEnd_}px)`
        this.scroller_.scrollTop = this.anchorScrollTop

        if (this.ANIMATION_DURATION_MS) {
            // TODO: Should probably use transition end, but there are a lot of animations we could be listening to.
            setTimeout(() => {
                for (let i in tombstoneAnimations) {
                    const anim = tombstoneAnimations[i]
                    anim[0].classList.add('invisible')
                    this.tombstones_.push(anim[0])
                    // Tombstone can be recycled now.
                }
            }, this.ANIMATION_DURATION_MS)
        }

        this.maybeRequestContent()
    }

    /**
     * Requests additional content if we don't have enough currently.
    */
    private maybeRequestContent() {
        // Don't issue another request if one is already in progress as we don't
        // know where to start the next request yet.
        if (this.requestInProgress_) {
            return
        }
        const itemsNeeded = this.lastAttachedItem_ - this.loadedItems_
        if (itemsNeeded <= 0) {
            return
        }
        this.requestInProgress_ = true
        this.source_.fetch(itemsNeeded).then(this.addContent.bind(this))
    }

    /**
    * Adds an item to the items list.
    */
    private addItem_() {
        this.items_.push({
            'data': null,
            'node': null,
            'height': 0,
            'width': 0,
            'top': 0,
        })
    }

    /**
    * Adds the given array of items to the items list and then calls
    * attachContent to update the displayed content.
    * @param {Array<Object>} items The array of items to be added to the infinite
    *     scroller list.
    */
    private addContent(items: Array<Object>) {
        this.requestInProgress_ = false
        for (let i = 0; i < items.length; i++) {
            if (this.items_.length <= this.loadedItems_)
                this.addItem_()
            this.items_[this.loadedItems_++].data = items[i]
        }
        this.attachContent()
    }
}

export class ItemsTester {
    private static readonly domains = ['.com', '.ru', '.net', '.io', '.us']
    private static readonly sites = ['google', 'facebook', 'apple', 'amazon', 'youtube']

    public static getRandom(): string {
        const randomDomain = this.domains[Math.floor(Math.random()*this.domains.length)]
        const randomSite = this.sites[Math.floor(Math.random()*this.sites.length)]
        return `${randomSite}${randomDomain}`
    }
}