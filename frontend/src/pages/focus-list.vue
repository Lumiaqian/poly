<template>
  <div class="q-pa-md">
    <q-table grid card-container-class="cardContainerClass" title="关注列表" :rows="rows" row-key="roomId" :filter="filter"
      hide-header v-model:pagination="pagination" :rows-per-page-options="rowsPerPageOptions">
      <template v-slot:item="props">
        <div class="q-pa-xs col-xs-12 col-sm-6 col-md-3">
          <router-link :to="{ path: '/index/room', query: { platform: props.row.platform, roomId: props.row.roomId } }">
            <q-card>
              <q-img :src="props.row.screenshot">
                <div class="absolute-bottom text-h8 text-left">
                  {{ props.row.gameFullName }}
                </div>
                <div class="absolute-bottom text-h8 text-right">
                  <q-icon name="person" />{{ props.row.count }}
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
                    <img :src="props.row.avatar">
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ props.row.platformName }} · {{ props.row.anchor }}</q-item-label>
                  <q-item-label caption>{{ props.row.roomName }}</q-item-label>
                </q-item-section>
                <q-item-label>
                  <q-chip v-if="props.row.isLive" dense color="red" text-color="white" label="直播中" />
                  <q-chip v-else dense color="orange" text-color="white" label="未开播" />
                </q-item-label>
              </q-item>
              <q-separator />
            </q-card>
          </router-link>
        </div>
      </template>
    </q-table>
  </div>
</template>

<script lang="ts">
import { useQuasar } from 'quasar'
import { ref, computed, watch, reactive, onMounted } from 'vue'
import { GetFocus, LogInfo } from '../../wailsjs';



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

    const initRoomInfo = () => GetFocus().then((res) => {
      rows.length = 0
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

        rows.push(roomInfo)
      })
      LogInfo('roomInfoSize: ' + rows.length)
    })
    const $q = useQuasar()

    function getItemsPerPage() {
      if ($q.screen.lt.sm) {
        return 4
      }
      if ($q.screen.lt.md) {
        return 8
      }
      return 12
    }

    const filter = ref('')
    const pagination = ref({
      page: 1,
      rowsPerPage: getItemsPerPage()
    })

    watch(() => $q.screen.name, () => {
      pagination.value.rowsPerPage = getItemsPerPage()
    })

    onMounted(() => {
      initRoomInfo();
    });

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
      filter,
      pagination,

      cardContainerClass: computed(() => {
        return $q.screen.gt.xs
          ? 'grid-masonry grid-masonry--' + ($q.screen.gt.sm ? '3' : '2')
          : null
      }),

      rowsPerPageOptions: computed(() => {
        return $q.screen.gt.xs
          ? $q.screen.gt.sm ? [4, 8, 12] : [4, 8]
          : [12]
      })
    }
  }
}
</script>

<style lang="sass">
.grid-masonry
  flex-direction: column
  height: 700px

  &--2
    > div
      &:nth-child(2n + 1)
        order: 1
      &:nth-child(2n)
        order: 2

    &:before
      content: ''
      flex: 1 0 100% !important
      width: 0 !important
      order: 1
  &--3
    > div
      &:nth-child(3n + 1)
        order: 1
      &:nth-child(3n + 2)
        order: 2
      &:nth-child(3n)
        order: 3

    &:before,
    &:after
      content: ''
      flex: 1 0 100% !important
      width: 0 !important
      order: 2
.router-link-active
  text-decoration: none
a 
  text-decoration: none
</style>
