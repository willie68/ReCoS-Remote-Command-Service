<template>
  <Dialog
    v-model:visible="dialogVisible"
    :modal="true"
    :closable="false"
  >
    <template #header>
      <h3>Action Wizard for profile: {{ this.profile.name }}</h3>
    </template>
    <Step0 v-if="this.step == 0" :profile="profile"></Step0>
    <Step1
      v-if="this.step == 1"
      :profile="profile"
      :value="newAction"
      v-on:next="checkNextState(1, $event)"
    ></Step1>
    <Step2
      v-if="this.step == 2"
      :profile="profile"
      :value="newAction"
    ></Step2>
    <template #footer>
      <Button label="Cancel" icon="pi pi-times" @click="cancel" />
      <Button
        label="Back"
        icon="pi pi-angle-left"
        autofocus
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
        autofocus
        @click="save"
        :disabled="!isFinishOK"
      />
    </template>
  </Dialog>
</template>

<script>
import Step0 from "./Step0.vue";
import Step1 from "./Step1.vue";
import Step2 from "./Step2.vue";

export default {
  name: "ActionWizard",
  components: {
    Step0,
    Step1,
    Step2,
  },
  props: {
    profile: {},
    visible: Boolean,
  },
  data() {
    return {
      dialogVisible: false,
      isFinishOK: false,
      isNextOK: true,
      isBackOK: false,
      activeProfile: { name: "" },
      step: 0,
      maxStep: 2,
      newAction: {},
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.activeProfile);
    },
    back() {
      this.step--;
      if (this.step < 0) {
        this.step = 0;
      }
      this.checkButtons();
    },
    next() {
      this.step++;
      if (this.step > this.maxStep) {
        this.step = this.maxStep;
      }
      this.checkButtons();
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
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
    },
    profile(profile) {
      this.activeProfile = profile;
    },
    newAction(newAction) {
      console.log("ActionWizard: newAction: ", JSON.stringify(newAction))
    }
  },
};
</script>

<style>
.w-page {
  width: 500px;
  height: 200px;
}
</style>