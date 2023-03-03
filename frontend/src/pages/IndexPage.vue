<template>
  <div class="q-pa-md">
    <q-banner inline-actions rounded class="bg-grey text-black">
      首页推荐
    </q-banner>
    <q-infinite-scroll @load="onLoad" :offset="250">
      <div v-for="(item, index) in list" :key="index" class="caption q-pa-md row items-center">
        <div v-for="value in item" :key="value.roomId" class="caption-clo">
          <div class="q-pa-md col-xs-12 col-sm-6 col-md-3">
            <router-link :to="{ path: '/index/room', query: { platform: value.platform, roomId: value.roomId } }">
              <q-card>
                <q-img :src="value.screenshot">
                  <div class="absolute-bottom text-h8 text-left">
                    {{ value.gameFullName }}
                  </div>
                  <div class="absolute-bottom text-h8 text-right">
                    <q-icon name="person" />{{ value.count }}
                  </div>
                  <template v-slot:error>
                    <div class="absolute-full flex flex-center bg-negative text-white">
                      Cannot load image
                    </div>
                  </template>
                </q-img>
                <q-item>
                  <q-item-section avatar>
                    <q-avatar>
                      <img :src="value.avatar">
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ value.platformName }} · {{ value.anchor }}</q-item-label>
                    <q-item-label caption>{{ value.roomName }}</q-item-label>
                  </q-item-section>
                  <q-item-label>
                    <q-chip v-if="value.isLive" dense color="red" text-color="white" label="直播中" />
                    <q-chip v-else dense color="orange" text-color="white" label="未开播" />
                  </q-item-label>
                </q-item>
                <q-separator />
              </q-card>
            </router-link>
          </div>
        </div>
      </div>
      <template v-slot:loading>
        <div class="row justify-center q-my-md">
          <q-spinner-dots color="primary" size="40px" />
        </div>
      </template>
    </q-infinite-scroll>
  </div>
</template>

<script lang="ts">
import { ref, reactive } from 'vue'
import { LogInfo, GetRecommend } from '../../wailsjs';



interface RoomInfo {
  platform: string
  roomId: string
  roomName: string
  anchor: string
  avatar: string
  onLineCount: number
  screenshot: string
  gameFullName: string
  isLive: boolean
  count: string
  platformName: string
}

export default {


  setup() {

    const rows = reactive([] as RoomInfo[])
    const list = reactive([[]] as RoomInfo[][])
    const rowSize = ref(4)
    const page = ref(1)
    const pageSize = ref(16)
    rows.length = 0
    list.length = 0


    function handleOnline(online: number): string {
      let num = online.toString().trim()
      if (num.length > 4) {
        let numCut = num.substring(0, num.length - 4)
        let afterPoint = num.substring(num.length - 4, num.length - 3)
        return numCut + '.' + afterPoint + '万'
      } else {
        return num + '人'
      }
    }

    return {
      rows,
      list,
      rowSize,
      page,
      pageSize,
      onLoad(index: number, done: () => void) {
        LogInfo('index :' + index)
        setTimeout(() => {
          GetRecommend(index, pageSize.value).then((res) => {
            let roomInfoList = reactive([] as RoomInfo[])
            res.forEach((item) => {
              let roomInfo: RoomInfo = {
                platform: item.platform,
                roomId: item.roomId,
                roomName: item.roomName,
                anchor: item.anchor,
                avatar: item.avatar,
                onLineCount: item.onLineCount,
                screenshot: item.screenshot,
                gameFullName: item.gameFullName,
                isLive: false,
                count: handleOnline(item.onLineCount),
                platformName: item.platformName
              }
              if (item.liveStatus == 2) {
                roomInfo.isLive = true
              } else {
                roomInfo.isLive = false
              }
              roomInfoList.push(roomInfo)
            })
            const total = Math.ceil(roomInfoList.length / rowSize.value)
            for (let i = 0; i < total; i++) {
              let row = roomInfoList.slice(i * rowSize.value, (i + 1) * rowSize.value)
              list.push(row)
            }
            LogInfo('onLoad roomInfoList: ' + JSON.stringify(list))
            LogInfo('onLoad roomInfoList length: ' + list.length)
          })
          done()
        }, 1000)

      }

    }
  }
}
</script>

<style lang="sass">
.router-link-active
  text-decoration: none
a 
  text-decoration: none
</style>
