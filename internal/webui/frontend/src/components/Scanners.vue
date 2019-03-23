<template>
  <div class="scanners">

      <ul id="scannerlist">
        <li v-for="scanner in scanners">
          {{ scanner.namespace }} |
          {{ scanner.label }}
        </li>
      </ul>

  </div>
</template>

<script lang="ts">
import axios from 'axios';
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Scanners extends Vue {
  @Prop() private scanners!: object[];
  @Prop() private errors!: object[];

  private created() {
    axios.get(`/api/scanners`)
        .then( (response) => {
            this.scanners = response.data;
        })
        .catch( (e) => {
            this.errors.push(e);
        });
  }
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
