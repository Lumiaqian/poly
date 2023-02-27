<template>
    <div ref="artRef" id="player" class="artplayer-app"></div>
</template>
  
<script setup lang="ts">
import flvjs from 'flv.js';
import Hls from 'hls.js';
import Artplayer from 'artplayer';
import { ref, onBeforeUnmount, onMounted, nextTick } from 'vue';
import { GetLiveRoom, LogInfo } from '../../wailsjs';



const props = defineProps({
    platform: {
        type: String,
        defalut: ''
    },
    roomId: {
        type: String,
        defalut: ''
    },
    isLive: {
        type: Boolean,
        defalut: true
    },
})

const artRef = ref('')

const platform = ref('');
const roomId = ref('');


platform.value = props.platform as string;
roomId.value = props.roomId as string;
let art: any;
interface Quality {
    default: boolean;
    html: string;
    url: string;
};
let qualities: Array<Quality>;
const initPlayer = () => GetLiveRoom(platform.value, roomId.value).then((res) => {
    LogInfo(JSON.stringify(res))
    let videoType = res.liveUrl.indexOf('m3u8') > 0 ? 'customHls' : 'customFlv'
    LogInfo(videoType)

    qualities = new Array<Quality>;
    res.quality.forEach((val, idx) => {
        let quality: Quality = {
            default: false,
            html: val.name,
            url: val.url,
        };
        if (idx === 0) {
            quality.default = true;
        }
        qualities.push(quality);
    })
    LogInfo('qualities ' + JSON.stringify(qualities));
    art = new Artplayer({
        container: artRef.value,
        lang: 'zh-cn',
        autoplay: true,
        isLive: true,
        url: res.liveUrl,
        type: videoType,
        autoSize: true, //固定视频比例
        pip: true,  //画中画
        fullscreen: true, //全屏按钮
        aspectRatio: true,  // 长宽比
        setting: true, // 设置按钮
        fullscreenWeb: true,  //网页全屏按钮
        volume: 1, //默认音量
        flip: true, //翻转
        screenshot: true,//截图
        mutex: false, //假如页面里同时存在多个播放器，是否只能让一个播放器播放

        quality: qualities,
        airplay: true,
        customType: {
            customHls: (video: any, url: any, art: any,) => {
                LogInfo('播放customHls')
                if (Hls.isSupported()) {
                    const hls = new Hls();
                    hls.loadSource(url);
                    hls.attachMedia(video);
                    art.hls = hls;
                    art.once('url', () => hls.destroy());
                    art.once('destroy', () => hls.destroy());
                } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
                    video.src = url;
                } else {
                    art.notice.show = 'Unsupported playback format: m3u8';
                }
            },
            customFlv: function (video: any, url: any, art: any,) {
                if (flvjs.isSupported()) {
                    LogInfo('播放flv')
                    const flv = flvjs.createPlayer({ type: 'flv', url });
                    flv.attachMediaElement(video);
                    flv.load();
                    // optional
                    art.flv = flv;
                    art.on('selector', (item: any) => {
                        LogInfo('切换清晰度')
                        LogInfo(JSON.stringify(item))
                        // flv.destroy()
                    });
                    art.once('url', () => flv.destroy());
                    art.once('destroy', () => flv.destroy());
                    // art.once('switch', () => flv.destroy());
                } else {
                    art.notice.show = 'Unsupported playback format: flv';
                }
            },
        }
    });
})

onMounted(() => {
    initPlayer()
    nextTick(() => {
        console.log('')
    });
});
onBeforeUnmount(() => {
    if (art) {
        art.destroy(false);
    }
});
</script>

<style scoped>
.artplayer-app{
  width: 100%;
  height: 100%;
  pointer-events: none;
}
</style>

  