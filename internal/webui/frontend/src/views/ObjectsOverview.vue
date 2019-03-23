<template>
  <div class="Objects">
      <Objects/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Objects from '@/components/Objects.vue';

@Component({
  components: {
    Objects,
  },
})
export default class ObjectsOverview extends Vue {}
</script>
