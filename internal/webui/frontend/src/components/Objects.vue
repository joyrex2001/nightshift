<template>
  <div class="objects">

    <b-navbar v-show="selected.length" fixed="bottom" type="dark" variant="light">
      <b-nav-form>
        <b-form-input class="mr-sm-2" type="number" v-model.number="replicas" placeholder="Replicas" />
        <b-button size="sm" class="my-2 my-sm-0" type="button" v-on:click="showScaleDialog">Scale now</b-button>
        &nbsp;
        <b-button size="sm" class="my-2 my-sm-0" type="button" v-on:click="showRestoreDialog">Restore state</b-button>
      </b-nav-form>
    </b-navbar>

    <b-table class="noselect"
        striped hover bordered small
        select-mode="range" selectable @row-selected="rowSelected"
        :items="objects" :fields="fields">
      <template slot="schedule" slot-scope="data">
        <schedule :schedule="data.value"/>
      </template>
    </b-table>

    <b-modal ok-only title="Error" id="invalid">
      <div class="d-block">Invalid number of replicas!</div>
    </b-modal>

    <b-modal ok-only title="Error" id="failed">
      <div class="d-block">{{ this.error }}</div>
    </b-modal>

    <b-modal@ok="reload" ok-only title="Success" id="success">
      <div class="d-block">
          Resource scaling has been succesfully scheduled.
      </div>
    </b-modal>

    <b-modal @ok="scale" title="Scale selected resources" id="scaling">
      <div class="d-block">
          The selection will be scaled to {{ replicas }} replicas.
          Are you sure?
      </div>
    </b-modal>

    <b-modal @ok="restore" title="Restore selected resources" id="restoring">
      <div class="d-block">
          The selection will be restored to the previously known number of replicas.
          Are you sure?
      </div>
    </b-modal>

    <div>&nbsp;</div>
    <div>&nbsp;</div>

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
  @Prop() private error!: object;
  @Prop() private replicas!: number;

  @Prop() private selected!: object[];
  private rowSelected(items: object[]) {
        this.selected = items;
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
      replicas: {
          label: 'Current',
          sortable: true,
      },
    };
    axios.get(`/api/objects`)
        .then( (response) => {
            this.objects = response.data;
        })
        .catch( (e) => {
            this.error = e;
            this.$root.$emit('bv::show::modal', 'failed', '#btnShow');
        });
  }

  private showScaleDialog() {
      if (typeof(this.replicas) === 'undefined' || this.replicas < 0) {
          this.$root.$emit('bv::show::modal', 'invalid', '#btnShow');
          return;
      }
      this.$root.$emit('bv::show::modal', 'scaling', '#btnShow');
  }

  private scale(evt: object) {
      axios.post(`/api/objects/scale/${this.replicas}`, this.selected)
          .then( (response) => {
              this.$root.$emit('bv::show::modal', 'success', '#btnShow');
          })
          .catch( (e) => {
              this.error = e;
              this.$root.$emit('bv::show::modal', 'failed', '#btnShow');
          });
  }

  private showRestoreDialog() {
      this.$root.$emit('bv::show::modal', 'restoring', '#btnShow');
  }

  private restore(evt: object) {
      axios.post(`/api/objects/restore`, this.selected)
          .then( (response) => {
              this.$root.$emit('bv::show::modal', 'success', '#btnShow');
          })
          .catch( (e) => {
              this.error = e;
              this.$root.$emit('bv::show::modal', 'failed', '#btnShow');
          });
  }

  private reload(evt: object) {
      this.$router.go(0);
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
.noselect {
  -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome and Opera */
}
.b-row-selected:focus  {
    border: 1px solid black;
}
</style>
