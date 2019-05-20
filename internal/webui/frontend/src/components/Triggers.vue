<template>
  <div class="triggers">
    <b-table class="noselect" striped hover bordered small :items="triggers" :fields="fields">
      <template slot="id" slot-scope="data">
         <div class="nowrap">{{ data.value }}</div>
      </template>
      <template slot="settings" slot-scope="data">
         <triggerConfig :settings="data.value"/>
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
import TriggerConfig from '@/components/TriggerConfig.vue';

@Component({
  components: {
    TriggerConfig,
  },
})

export default class Triggers extends Vue {
  @Prop() private fields!: object;
  @Prop() private triggers!: object[];
  @Prop() private error!: object;

  private created() {
    this.fields = {
        id: {
            label: 'Id',
            sortable: true,
        },
        type: {
            label: 'Type',
            sortable: true,
        },
        settings: {
            label: 'Settings',
            sortable: true,
        },
    };
    axios.get(`/api/triggers`)
        .then( (response) => {
            this.triggers = response.data;
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
.nowrap {
    white-space: nowrap;
}
</style>
