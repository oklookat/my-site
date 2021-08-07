import ElvenPlayerControls from "@/common/plugins/ElvenPlayer/core/ElvenPlayerControls"


export const state = {
    // service vars start //
    audioPlayerControls: undefined,
    addSource: String,
    // service vars end //
    isSourcesTriggered: false,
    sources: [],
    currentSource: '',
}

export const methods = {
    mounted() {
        state.audioPlayerControls = new ElvenPlayerControls()
    },
    play(){
        state.audioPlayerControls.play()
    },
    isPlaying(){
        return !state.audioPlayerControls.isPlaying()
    },
    addToPlaylist(url){
        state.audioPlayerControls.addToPlaylist(url)
        state.isSourcesTriggered = true
        state.isSourcesTriggered = false
    },
    setPlaylist(url) {
        const playlist = [url]
        state.audioPlayerControls.setPlaylist(playlist)
    },
    getPlaylistLength(){
        return state.audioPlayerControls.getPlaylistLength()
    },
}