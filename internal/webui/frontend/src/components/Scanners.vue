<template>
  <div class="scanners">
    <b-table striped hover small :items="scanners" :fields="fields">
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

export default class Scanners extends Vue {
  @Prop() private fields!: object;
  @Prop() private scanners!: object[];
  @Prop() private errors!: object[];

  private created() {
    axios.get(`/api/scanners`)
        .then( (response) => {
            this.scanners = response.data;
            this.fields = {
                namespace: {
                    label: 'Namespace',
                    sortable: true,
                },
                label: {
                    label: 'Label',
                    sortable: true,
                },
                schedule: {
                    label: 'Schedule',
                    sortable: true,
                },
            };
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
