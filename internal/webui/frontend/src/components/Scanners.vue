<template>
  <div class="scanners">
    <b-table striped hover bordered small :items="scanners" :fields="fields">
      <template slot="schedule" slot-scope="data">
         <schedule :schedule="data.value"/>
      </template>
    </b-table>

    <b-modal ok-only title="Error" id="failed">
      <div class="d-block">{{ this.error }}</div>
    </b-modal>
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
  @Prop() private error!: object;

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
            this.error = e;
            this.$root.$emit('bv::show::modal', 'failed', '#btnShow');
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
