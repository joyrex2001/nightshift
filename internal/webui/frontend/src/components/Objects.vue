<template>
  <div class="objects">

    <b-navbar v-show="selected.length" fixed="bottom" type="dark" variant="light">
      <b-nav-form>
        <b-form-input class="mr-sm-2" type="number" v-model.number="replicas" placeholder="Replicas" />
        <b-button size="sm" class="my-2 my-sm-0" type="button" v-on:click="scale">Scale now</b-button>
      </b-nav-form>
    </b-navbar>

    <b-table
        striped hover bordered small
        :select-mode="range" selectable @row-selected="rowSelected"
        :items="objects" :fields="fields">
      <template slot="schedule" slot-scope="data">
        <schedule :schedule="data.value"/>
      </template>
    </b-table>

    <b-modal ok-only title="Error" id="invalid">
      <div class="d-block">Invalid number of replicas!</div>
    </b-modal>

    <b-modal ok-only title="Scaling" id="scaling">
      <div class="d-block">NOT YET IMPLEMENTED</div>
      <div class="d-block">Scaling to {{ replicas }} replicas.</div>
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

export default class Objects extends Vue {
  @Prop() private fields!: object;
  @Prop() private objects!: object[];
  @Prop() private errors!: object[];
  @Prop() private replicas!: number;

  @Prop() private selected!: object[];
  private rowSelected(items: object[]) {
        this.selected = items;
  }

  private scale() {
      if (typeof(this.replicas) === 'undefined' || this.replicas < 0) {
          this.$root.$emit('bv::show::modal', 'invalid', '#btnShow');
          return;
      }
      this.$root.$emit('bv::show::modal', 'scaling', '#btnShow');
  }

  private created() {
    this.selected = [];
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

<style>
tr:focus {
    outline: none;
}
th:focus {
    outline: none;
}
.b-row-selected:focus  {
    border: 1px solid black;
}
</style>
