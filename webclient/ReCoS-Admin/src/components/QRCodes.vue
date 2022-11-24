<template>
  <Dialog v-model:visible="dialogVisible" :modal="true" :closable="false">
    <template #header>
      <h3>QR Codes</h3>
    </template>
    This dialog presents you the right QR Code for connecting a smart client.<br/>
    On the left you see all possible ips of your network adapters.<br/>
    If you select one, you can see on the right the associated QR Code.<br/>
    Scan this wih your device and you will start the client.<br/>
    
    <Splitter style="height: 280px">
      <SplitterPanel :size="30">
        <Listbox
          v-model="network"
          :options="networks"
          optionLabel="IP"
          listStyle="max-height:280px"
        >
        </Listbox>
      </SplitterPanel>
      <SplitterPanel :size="70">
          <img :src="network.QRCode"/><br/>
          {{ network.URL }}
      </SplitterPanel>
    </Splitter>
    <template #footer>
      <div class="p-pt-4">
        <Button label="Close" icon="pi pi-times" @click="close" />
      </div>
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "QRCodes",
  components: {},
  props: {
    visible: Boolean,
  },
  emits: ["close"],
  data() {
    return {
      dialogVisible: false,
      networks: [],
      network: {},
    };
  },
  methods: {
    close() {
      this.$emit("close");
    },
    GetNetworks() {
      fetch(this.$baseURL + "config/networks", {
        method: "GET",
      })
        .then((res) => res.json())
        .then((data) => {
          //console.log(data);
          this.networks = data;
          if (this.networks.length > 0) {
              this.network = this.networks[0]
          }
        })
        .catch((err) => {
          console.log(err.message);
        });
    },
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      if (visible) {
        this.GetNetworks();
      }
    },
  },
};
</script>

<style>
</style>