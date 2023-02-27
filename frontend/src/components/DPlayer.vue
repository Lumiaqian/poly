<template>
    <div class="p-4">
        <div id="dplayer"></div>
    </div>
</template>
  
<script lang="ts">
import flvjs from 'flv.js';
import Hls from 'hls.js';
import DPlayer from 'dplayer';
import { ref, onBeforeUnmount, defineComponent, onMounted } from 'vue';
import { GetLiveRoom, LogInfo } from '../../wailsjs';

export default defineComponent({
    props: {
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
        screenshot: {
            type: String,
            defalut: ''
        }
    },
    setup(props, ctx) {
        LogInfo(JSON.stringify(ctx))

        const platform = ref('')
        const roomId = ref('')

        const videoType = ref('m3u8')
        platform.value = props.platform as string
        roomId.value = props.roomId as string
        let dp: any;

        const initPlayer = () => GetLiveRoom(platform.value, roomId.value).then((res) => {
            LogInfo(JSON.stringify(res))
            let videoType = res.liveUrl.indexOf('m3u8') > 0 ? 'customHls' : 'customFlv'
            LogInfo(videoType)
            if (dp) {
                dp.destroy()
            }
            dp = new DPlayer({
                container: document.getElementById('dplayer'),
                lang: 'zh-cn',
                autoplay: true,
                live: true,
                screenshot: true,
                airplay: true,
                video: {
                    url: res.liveUrl,
                    type: videoType,
                    // quality: res.quality,
                    // defaultQuality: 0,
                    pic: props.screenshot,
                    customType: {
                        customHls: (video: any, player: any) => {
                            LogInfo(video);
                            LogInfo(player);
                            const hls = new Hls();
                            hls.loadSource(video.src);
                            hls.attachMedia(video);
                            player.events.on("destroy", () => {
                                hls.destroy();
                            });
                        },
                        customFlv: function (video: any, player: any) {
                            LogInfo(video);
                            LogInfo(player);
                            const flvPlayer = flvjs.createPlayer({
                                type: "flv",
                                url: video.src,
                            });
                            flvPlayer.attachMediaElement(video);
                            flvPlayer.load();
                            player.events.on("destroy", () => {
                                flvPlayer.unload();
                                flvPlayer.detachMediaElement();
                                flvPlayer.destroy();
                            });


                        },
                    },
                },
            });
            LogInfo(JSON.stringify(dp))
        })
        onMounted(() => {
            initPlayer();
        });

        onBeforeUnmount(() => {
            LogInfo('onBeforeUnmount')
            if (dp) {
                dp.destroy()
            }
        });


        return {
            videoType,
            dp
        }
    }
})
</script>

<style scoped>
.p-4 {
    width: 100%;
    height: 100%;
    /*pointer-events: none;*/
}
</style>
  