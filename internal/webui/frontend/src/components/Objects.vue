<template>
  <div class="objects">
    <b-table striped hover small :items="objects" :fields="fields">
      <template slot="schedule" slot-scope="data">
        <Schedule :schedule="data.value"/>
      </template>
    </b-table>
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
  @Prop() private fields!: object;
  @Prop() private objects!: object[];
  @Prop() private errors!: object[];

  private created() {
    this.fields = {
      namespace: {
          label: 'Namespace',
          sortable: true,
      },
      name: {
          label: 'Name',
          sortable: true,
      },
      schedule: {
          label: 'Schedule',
          sortable: true,
      },
    };
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
