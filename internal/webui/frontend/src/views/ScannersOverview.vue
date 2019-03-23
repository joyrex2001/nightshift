<template>
  <div class="Scanners">
      <Scanners/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Scanners from '@/components/Scanners.vue';

@Component({
  components: {
    Scanners,
  },
})

export default class ScannersOverview extends Vue {}
</script>
