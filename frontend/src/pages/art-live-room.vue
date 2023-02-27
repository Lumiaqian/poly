<template>
    <div class="room-main-container">
        <div class="room-left">
            <div class="room-left-video">
                <ArtPlayer class="room-left-video-play" :platform="platform" :room-id="roomId" :is-live="live"
                    />
            </div>
            <div class="room-left-info">
                <q-avatar square class="room-left-info-head">
                    <img :src=roomInfo.avatar>
                </q-avatar>
                <div class="room-left-info-after-head">
                    <div class="room-left-info-after-head-name">
                        <div :class="isLive() ? 'info-isLive' : 'info-notLive'" style="font-size: small">{{ isLive() ? "直播中"
                            : "未开播" }}
                        </div>
                        {{ roomInfo.roomName }}
                    </div>
                    <div class="room-left-info-after-head-anchor">
                        {{ platform }} · {{ roomInfo.gameFullName }} · {{ roomInfo.anchor }}
                    </div>
                </div>
            </div>
        </div>
        <div class="room-right">
            <div class="room-right-top">
                直播聊天
            </div>
        </div>
    </div>
</template>
   
<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import ArtPlayer from 'components/ArtPlayer.vue';
import { GetLiveRoomInfo, LogInfo } from '../../wailsjs';
const platform = ref('huya')
const roomId = ref('222523')
const live = ref(true)

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

const style = ref({
    //width: '100%',
    height: '650px',
    //margin: '60px auto 0',
})
const initRoomInfo = () => GetLiveRoomInfo(platform.value, roomId.value).then((res) => {
    LogInfo('roomInfo: ' + JSON.stringify(res))
    roomInfo.platform = res.platform
    roomInfo.roomId = res.roomId
    platform.value = res.platform
    roomId.value = res.roomId
    if (res.liveStatus == 2) {
        live.value = true
    } else {
        live.value = false
    }
    roomInfo.roomName = res.roomName
    roomInfo.anchor = res.anchor
    roomInfo.avatar = res.avatar
    roomInfo.gameFullName = res.gameFullName
    roomInfo.onLineCount = res.onLineCount
    roomInfo.screenshot = res.screenshot
    
})

onMounted(() => {
    initRoomInfo();
});

function isLive(): boolean {
    LogInfo('onMounted avatar: ' + roomInfo.avatar)
    return live.value
}
</script>
 
<style scoped>
.room-main-container {
    position: relative;
    height: 100%;
    width: 100%;
}

.room-left {
    position: relative;
    width: 78%;
    height: 100%;
}

.room-left-video {
    position: absolute;
    width: 100%;
    height: 84%;
    background-color: black;
    top: 0px;
    left: 0px;
    bottom: 80px;
}

.room-left-video-play {
    width: 100%;
    height: 100%;
}

.room-right {
    width: 22%;
    height: 100%;
    position: fixed;
    top: 0px;
    right: 0px;
    border-left: 1px solid #c8c8c9;
}

.room-left-info-head {
    float: left;
    margin-top: 9px;
    margin-left: 8px;
    width: 60px;
    height: 60px;
    box-shadow: #2b2b2b 0px 0px 5px 1px;
    border-radius: 10px;
}

.room-left-info-after-head {
    float: left;
    margin-left: 10px;
    margin-top: 8px;
}

.room-left-info-after-head-name {
    font-weight: bold;
    font-size: 20px;
}

.room-left-info-after-head-anchor {
    margin-top: 10px;
    font-weight: bold;
    font-size: 15px;
}

.info-isLive {
    margin-top: 6px;
    margin-right: 5px;
    float: left;
    height: 18px;
    width: 45px;
    background-color: #c10f0f;
    border-radius: 10px;
    font-size: 5px;
    font-weight: 600;
    text-align: center;
    color: #F3F6F8;
}

.info-notLive{
  margin-top: 6px;
  margin-right: 5px;
  float: left;
  height: 18px;
  width: 45px;
  background-color: #979797;
  border-radius: 10px;
  font-size: 5px;
  font-weight: 600;
  text-align: center;
  color: #F3F6F8;
}

.head-avatar {
    border-radius: 10px;
    height: 100%;
    width: 100%;
}
</style>
   