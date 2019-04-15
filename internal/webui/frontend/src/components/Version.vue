<template>
  <div class="text">
      <div class="row">
          <div class="col-1">Version:</div><div class="col-6">{{ version.version }}</div>
      </div>
      <div class="row">
          <div class="col-1">Date:</div><div class="col-6">{{ version.date }}</div>
      </div>
      <div class="row">
          <div class="col-1">Build:</div><div class="col-6">{{ version.build }}</div>
      </div>
  </div>
</template>

<script lang="ts">
import axios from 'axios';
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Version extends Vue {
  @Prop() private version!: object[];
  @Prop() private error!: object;

  private created() {
    axios.get(`/api/version`)
        .then( (response) => {
            this.version = response.data;
        })
        .catch( (e) => {
            this.error = e;
            this.$root.$emit('bv::show::modal', 'failed', '#btnShow');
        });
  }
}
</script>
