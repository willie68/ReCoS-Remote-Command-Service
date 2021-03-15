<template>
  <Dialog v-model:visible="dialogVisible" :modal="true" :closable="false">
    <template #header>
      <h3>Action Wizard for profile: {{ this.profile.name }}</h3>
    </template>
    <Step0 v-if="this.step == 0" :profile="profile"></Step0>
    <Step1
      v-if="this.step == 1"
      :profile="profile"
      :value="newAction"
      :commandTypes="commandTypes"
      v-on:value="updateAction($event)"
      v-on:next="checkNextState(1, $event)"
    ></Step1>
    <Step2
      v-if="this.step == 2"
      :profile="profile"
      :value="newAction"
      :commandTypes="commandTypes"
      v-on:value="updateAction($event)"
      v-on:next="checkNextState(2, $event)"
    ></Step2>
    <Step3
      v-if="this.step == 3"
      :profile="profile"
      :value="newAction"
      :commandTypes="commandTypes"
      v-on:value="updateAction($event)"
      v-on:next="checkNextState(2, $event)"
    ></Step3>
    <template #footer>
      <div class="p-pt-4">
        <Button label="Cancel" icon="pi pi-times" @click="cancel" />
        <Button
          label="Back"
          icon="pi pi-angle-left"
          @click="back"
          :disabled="!isBackOK"
        />
        <Button
          label="Next"
          icon="pi pi-angle-right"
          autofocus
          @click="next"
          :disabled="!isNextOK"
        />
        <Button
          label="Finish"
          icon="pi pi-check"
          @click="save"
          :disabled="!isFinishOK"
        />
      </div>
    </template>
  </Dialog>
</template>

<script>
import Step0 from "./Step0.vue";
import Step1 from "./Step1.vue";
import Step2 from "./Step2.vue";
import Step3 from "./Step3.vue";

export default {
  name: "ActionWizard",
  components: {
    Step0,
    Step1,
    Step2,
    Step3,
  },
  props: {
    profile: {},
    visible: Boolean,
  },
  emits: ["cancel", "save"],
  data() {
    return {
      dialogVisible: false,
      isFinishOK: false,
      isNextOK: true,
      isBackOK: false,
      activeProfile: { name: "" },
      step: 0,
      maxStep: 3,
      newAction: {},
      commandTypes: [],
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.activeProfile);
    },
    updateAction(data) {
      this.newAction = data;
    },
    back() {
      this.step--;
      if (this.step < 0) {
        this.step = 0;
      }
      this.checkButtons();
      this.$forceUpdate();
    },
    next() {
      this.step++;
      if (this.step > this.maxStep) {
        this.step = this.maxStep;
      }
      this.checkButtons();
      this.$forceUpdate();
    },
    checkButtons() {
      this.isNextOK = this.step < this.maxStep;
      this.isBackOK = this.step > 0;
      this.isFinishOK = this.step == this.maxStep;
    },
    checkName(name) {
      if (name == "") {
        this.isNameOK = false;
        return;
      }
      this.isNameOK = !this.profiles
        .map((elem) => elem.toLowerCase())
        .includes(name.toLowerCase());
    },
    checkNextState(actualStep, next) {
      if (actualStep == this.step) {
        console.log("ActionWizard: next event:", next);
        if (next) {
          this.isNextOK = true;
        } else {
          this.isNextOK = false;
          this.isFinishOK = false;
        }
      }
    },
  },
  mounted() {
    let url = this.$store.state.baseURL + "config/commands";
    fetch(url)
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
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      if (visible == true) {
        this.step = 0;
      }
      this.checkButtons();
    },
    profile(profile) {
      this.activeProfile = profile;
    },
    newAction(newAction) {
      console.log("ActionWizard: newAction: ", JSON.stringify(newAction));
    },
  },
};
</script>

<style>
.w-page {
  width: 500px;
  height: 200px;
  text-align: left;
}
</style>