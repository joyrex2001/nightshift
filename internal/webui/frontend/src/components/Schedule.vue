<template>
  <div class="schedules">
    <div class="row" v-for="sched in schedule">
      {{ sched.Description }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Schedule extends Vue {
  @Prop() private schedule!: object[];
}

</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}
li {
  margin: 0 10px;
}
</style>
