<template>
  <div class="triggerConfig">
    <div v-for="value,key in settings">
      <div class="row">
        <div class="col-1">{{ key }}:</div><div class="col">{{ value }}</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class TriggerConfig extends Vue {
  @Prop() private settings!: object[];
}

</script>
