<template>
  <div class="Triggers">
      <Triggers/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Triggers from '@/components/Triggers.vue';

@Component({
  components: {
    Triggers,
  },
})

export default class TriggersOverview extends Vue {}
</script>
