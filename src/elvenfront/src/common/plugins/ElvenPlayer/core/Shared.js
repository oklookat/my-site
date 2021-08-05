import ElvenPlayerControls from "@/common/plugins/ElvenPlayer/core/ElvenPlayerControls"

export const state = {
    // service vars start //
    SERVICE: 'ELVEN_PLAYER_C',
    audioPlayer: undefined,
    addSource: String,
    // service vars end //
    sources: [],
    currentSource: '',
}

export const methods = {
    mounted(){
        state.audioPlayer = new ElvenPlayerControls()
    },
    addSource(source){
        state.audioPlayer.addSource(source)
    },
}