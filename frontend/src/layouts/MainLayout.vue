<template>
  <q-layout view="lHh Lpr lFf" class="bg-grey-9 text-black">
    <q-header elevated class="bg-grey-6">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          @click="toggleLeftDrawer"
          aria-label="Menu"
          icon="menu"
        />

        <q-toolbar-title>
          Poly
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      class="bg-grey-8"
    >
      <q-list dark>
        <q-item-label header>Author:LumiaQian</q-item-label>
        <q-item clickable rel="noopener" to="/index">
          <q-item-section avatar>
            <q-icon name="home" />
          </q-item-section>
          <q-item-section>
            <q-item-label>首页推荐</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable to="/index/focus" >
          <q-item-section avatar>
            <q-icon name="favorite" />
          </q-item-section>
          <q-item-section>
            <q-item-label>关注</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable  >
          <q-item-section avatar>
            <q-icon name="list" />
          </q-item-section>
          <q-item-section>
            <q-item-label>平台</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable >
          <q-item-section avatar>
            <q-icon name="view_list" />
          </q-item-section>
          <q-item-section>
            <q-item-label>分类</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable to="/art-room" >
          <q-item-section avatar>
            <q-icon name="tv" />
          </q-item-section>
          <q-item-section>
            <q-item-label>电视</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'MyLayout',

  setup () {
    const leftDrawerOpen = ref(false)

    function toggleLeftDrawer () {
      leftDrawerOpen.value = !leftDrawerOpen.value
    }

    return {
      leftDrawerOpen,
      toggleLeftDrawer
    }
  }
}
</script>
