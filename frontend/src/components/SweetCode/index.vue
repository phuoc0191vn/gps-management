<template>
  <pre v-highlightjs="code" :class="['sweet-code']"><slot></slot></pre>
</template>

<script>
import Vue from 'vue'
import VueHighlightJS from 'vue-highlightjs'
import 'highlight.js/styles/googlecode.css'
Vue.use(VueHighlightJS)
export default {
  props: {
    code: {
      type: String
    }
  }
}
</script>

<style lang="scss">
.sweet-code {
  tab-size: 4;
  -webkit-font-smoothing: auto;
  background-color: transparent;
}
</style>
