<template>
  <div class="home">
    <Scanners/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Scanners from '@/components/Scanners.vue'; // @ is an alias to /src

@Component({
  components: {
    Scanners,
  },
})
export default class Home extends Vue {}
</script>
