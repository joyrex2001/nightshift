<template>
  <div class="scanners">
    <b-table striped hover bordered small :items="scanners" :fields="fields">
      <template slot="schedule" slot-scope="data">
         <schedule :schedule="data.value"/>
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
    this.fields = {
        namespace: {
            label: 'Namespace',
            sortable: true,
        },
        label: {
            label: 'Label selector',
            sortable: true,
        },
        schedule: {
            label: 'Schedule',
            sortable: true,
        },
    };
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

<style>
tr:focus {
    outline: none;
}
th:focus {
    outline: none;
}
</style>
