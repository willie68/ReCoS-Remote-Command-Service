<template>
  <Dialog v-model:visible="dialogVisible" :modal="true" :closable="false">
    <template #header>
      <h3>Settings</h3>
    </template>
    This is the settings dialog. Take a look to all tabs with different
    settings.
    <TabView>
      <TabPanel>
        <template #header>
          <i class="pi pi-cog"></i>
          <span> Header I</span>
        </template>
        Content I
      </TabPanel>
      <TabPanel>
        <template #header>
          <i class="pi pi-desktop"></i>
          <span>Open Hardware Monitor</span>
        </template>
        <label for="ohmActive">active: </label>
        <Checkbox id="ohmActive" v-model="ohmActive" :binary="true" />

        <a href="https://openhardwaremonitor.org/" target="_blank"
          >Open Hardware Monitor</a
        >
        updateperiod: 5 url: http://127.0.0.1:12999/data.json
      </TabPanel>
      <TabPanel>
        <template #header>
          <i class="pi pi-volume-up"></i>
          <span>Audioplayer</span>
        </template>
        <div class="p-fluid">
          <div class="p-field">
            <label for="apActive">active: </label>
            <Checkbox id="apActive" v-model="apActive" :binary="true" /><br />
          </div>
          <div class="p-field">
            <Dropdown
              v-model="apSampleRate"
              :options="apSampletrates"
              optionLabel="name"
              placeholder="Select a samplerate"
              optionValue="value"
            /><br />
          </div>
        </div>
        For the audioplayer i need simply the sample rate to work with. <br />
        For convinience you can only switch between 44,1kHz and 48kHz. <br />
        {{ apSampleRate }}
      </TabPanel>
      <TabPanel>
        <template #header>
          <i class="pi pi-sun"></i>
          <span>Philips Hue</span>
        </template>
        <label for="phActive">active: </label>
        <Checkbox id="phActive" v-model="phActive" :binary="true" />

        Philips Hue configuration active: false username:
        IwtURJmST8b44mvZSZ2nl73nZhghVltMvgzlH7UC device: recos#hue_user
        ipaddress: 192.168.178.81 updateperiod: 5
      </TabPanel>
      <TabPanel>
        <template #header>
          <i class="pi pi-home"></i>
          <span>Homematic</span>
        </template>
        <label for="hmActive">active: </label>
        <Checkbox id="hmActive" v-model="hmActive" :binary="true" />

        updateperiod: 5 url: http://192.168.178.80
      </TabPanel>
    </TabView>
    <template #footer>
      <div class="p-pt-4">
        <Button label="Cancel" icon="pi pi-times" @click="cancel" />
        <Button
          label="Save"
          icon="pi pi-check"
          @click="save"
          :disabled="!isSaveOK"
        />
      </div>
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "Settings",
  components: {},
  props: {
    profile: {},
    visible: Boolean,
  },
  emits: ["cancel", "save"],
  data() {
    return {
      dialogVisible: false,
      isFinishOK: false,
      ohmActive: false,
      apActive: false,
      phActive: false,
      hmActive: false,
      apSampletrates: [
        {
          name: "44,1kHz",
          value: 44100,
        },
        {
          name: "48kHz",
          value: 48000,
        },
      ],
      apSampleRate: 44100,
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.activeProfile);
    },
    checkButtons() {
      this.isFinishOK = true;
    },
    updateCommands() {
      let url = this.$store.state.baseURL + "config/commands";
      const myHeaders = new Headers();

      myHeaders.append("Content-Type", "application/json");
      myHeaders.append(
        "Authorization",
        `Basic ${btoa(`admin:${this.$store.state.password}`)}`
      );
      myHeaders.append("X-mcs-profile", this.activeProfile.name);

      fetch(url, {
        method: "GET",
        mode: "cors",
        headers: myHeaders,
      })
        .then((res) => res.json())
        .then((data) => {
          this.commandTypes = [];
          data.forEach((element) => {
            if (element.wizard && element.wizard == true) {
              this.commandTypes.push(element);
            }
          });
        })
        .catch((err) => console.log(err.message));
    },
  },
  mounted() {},
  beforeUnmount() {},
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      if (visible == true) {
        this.step = 0;
      }
      this.checkButtons();
    },
  },
};
</script>

<style>
.w-page {
  width: 600px;
  height: 300px;
  text-align: left;
}
</style>