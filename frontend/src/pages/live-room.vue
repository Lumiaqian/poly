<template>
    <div class="room-main-container">
        <div class="room">
            <div class="room-video">
                <DPlayer :platform="platform" :room-id="roomId" :is-live="isLive()" :screenshot="roomInfo.screenshot"
                    class="room-video-play" />
            </div>
            <div class="room-info">
                <q-avatar square class="room-info-head">
                    <img :src=roomInfo.avatar>
                </q-avatar>
                <div class="room-info-after-head">
                    <div class="room-info-after-head-name">
                        <q-chip v-if="isLive()" dense color="red" text-color="white" label="直播中" />
                        <q-chip v-else dense color="orange" text-color="white" label="未开播" />
                        {{ roomInfo.roomName }}
                    </div>
                    <div class="room-info-after-head-anchor">
                        {{ platformName }} · {{ roomInfo.gameFullName }} · {{ roomInfo.anchor }}
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
   
<script setup lang="ts">
import { ref, onMounted, reactive,onBeforeMount } from 'vue'
import DPlayer from 'components/DPlayer.vue';
import { GetLiveRoomInfo, LogInfo } from '../../wailsjs';
import { useRoute } from 'vue-router'; 
const route = useRoute()
const platform = ref('')
const roomId = ref('')
const live = ref(true)
const platformName = ref('')

const roomInfo = reactive({
    platform: '',
    roomId: '',
    roomName: '',
    anchor: '',
    avatar: '',
    onLineCount: 0,
    screenshot: '',
    gameFullName: '',
})


const initRoomInfo = () => GetLiveRoomInfo(platform.value, roomId.value).then((res) => {
    LogInfo('roomInfo: ' + JSON.stringify(res))
    roomInfo.platform = res.platform;
    roomInfo.roomId = res.roomId;
    platform.value = res.platform;
    roomId.value = res.roomId;
    if (res.liveStatus == 2) {
        live.value = true;
    } else {
        live.value = false;
    }
    roomInfo.roomName = res.roomName;
    roomInfo.anchor = res.anchor;
    roomInfo.avatar = res.avatar;
    roomInfo.gameFullName = res.gameFullName;
    roomInfo.onLineCount = res.onLineCount;
    roomInfo.screenshot = res.screenshot;
    platformName.value = res.platformName;
    
})

onBeforeMount(() => {
    platform.value = String(route.query.platform);
    roomId.value = String(route.query.roomId);
})

onMounted(() => {
    initRoomInfo();
});

function isLive(): boolean {
    LogInfo('onMounted avatar: ' + roomInfo.avatar);
    return live.value;
}


</script>
 
<style scoped>
.room-main-container {
    height: 100%;
    width: 100%;
}

.room {
    position: relative;
    width: 100%;
    height: 100%;
}

.room-video {
    position: relative;
    width: 100%;
    height: 84%;
    background-color: black;
    top: 0px;
    left: 0px;
    bottom: 80px;
}

.room-video-play {
    width: 100%;
    height: 100%;
}

.room-right {
    width: 22%;
    height: 100%;
    position: fixed;
    top: 0px;
    right: 0px;
    border: 1px solid #c8c8c9;
}

.room-info-head {
    float: left;
    margin-top: 9px;
    margin: 8px;
    width: 60px;
    height: 60px;
    box-shadow: #2b2b2b 0px 0px 5px 1px;
    border-radius: 10px;
}

.room-info-after-head {
    float: left;
    margin: 10px;
    margin-top: 8px;
}

.room-info-after-head-name {
    font-weight: bold;
    font-size: 20px;
}

.room-info-after-head-anchor {
    margin-top: 10px;
    font-weight: bold;
    font-size: 15px;
}
</style>
   