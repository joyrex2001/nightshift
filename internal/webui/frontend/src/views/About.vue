<template>
  <div id="about">
    <Version/>
    <div class="text">
      For more information visit the <a href="https://github.com/joyrex2001/nightshift">github</a> page.
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Version from '@/components/Version.vue';

@Component({
  components: {
    Version,
  },
})
export default class About extends Vue {}
</script>

<style scoped>
#about {
  border-top: 1px solid #dee2e6;
}
.text {
  margin: 10px 10px;
}
</style>
