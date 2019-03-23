<template>
  <div class="objects">

      <ul id="objectlist">
        <li v-for="object in objects">
          {{ object.namespace }} | {{ object.name }}
          <Schedule :schedule="object.schedule"/>
        </li>
      </ul>

  </div>
</template>

<script lang="ts">
import axios from 'axios';
import { Component, Prop, Vue } from 'vue-property-decorator';
import Schedule from '@/components/Schedule.vue';

@Component({
  components: {
    Schedule,
  },
})

export default class Objects extends Vue {
  @Prop() private objects!: object[];
  @Prop() private errors!: object[];

  private created() {
    axios.get(`/api/objects`)
        .then( (response) => {
            this.objects = response.data;
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
