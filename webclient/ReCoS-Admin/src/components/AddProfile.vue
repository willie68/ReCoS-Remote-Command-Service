<template>
  <Dialog v-model:visible="dialogProfileVisible">
    <template #header>
      <h3>Add profile</h3>
    </template>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="template" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0"
          >template</label
        >
        <div class="p-col-24 p-md-10">
          <Dropdown
            v-model="selectedTemplate"
            :options="templates"
            optionLabel="label"
            optionGroupLabel="label"
            optionGroupChildren="items"
            placeholder="select a template"
            :filter="true"
            filterPlaceholder="Find a template"
            class="p-ml-2"
            @change="changeTemplate"
          >
            <template #optiongroup="slotProps">
              <div class="p-d-flex p-ai-center">
                <div><b>{{ slotProps.option.label }}</b></div>
              </div>
            </template>
          </Dropdown>
        </div>
      </div>
      <div class="p-field p-grid">
        <label for="name" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0">name</label>
        <div class="p-col-24 p-md-10">
          <InputText
            id="name"
            type="text"
            v-model="addProfile.name"
            class="p-ml-2"
            :disabled="edit"
            :class="{ 'p-invalid': !isNameOK }"
            autofocus
          />
        </div>
      </div>
      <div class="p-field p-grid">
        <label for="description" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0"
          >description</label
        >
        <div class="p-col-24 p-md-10">
          <InputText
            id="description"
            type="text"
            v-model="addProfile.description"
            class="p-ml-2"
          />
        </div>
      </div>
    </div>
    <template #footer>
      <Button
        label="Import"
        icon="pi pi-cloud-upload"
        class="p-button-text"
        @click="this.importAction()"
      />
      <Button
        label="Cancel"
        icon="pi pi-times"
        class="p-button-text"
        @click="cancel"
      />
      <Button
        label="Save"
        icon="pi pi-check"
        autofocus
        @click="save"
        :disabled="!isNameOK"
      />
    </template>
  </Dialog>
  <Upload
    :visible="dialogUploadVisible"
    filetype=".profile"
    @cancel="dialogUploadVisible = false"
    @save="doImport"
  />
</template>

<script>
import Upload from "./Upload.vue";
export default {
  name: "AddProfile",
  components: {
    Upload,
  },
  props: {
    profile: {},
    visible: Boolean,
    edit: Boolean,
    profiles: {},
  },
  computed: {
    templates: {
      get: function () {
        let temps = Array();
        for (let i = 0; i < this.temps.length; i++) {
          let temp = this.temps[i];
          let cat = temp.group;
          if (!cat || cat == "") {
            cat = "unknown";
          }

          let found = false;
          for (let x = 0; x < temps.length; x++) {
            if (temps[x].label == cat) {
              temps[x].items.push(temp);
              found = true;
            }
          }
          if (!found) {
            let myTemp = {
              label: cat,
              items: Array(),
            };
            myTemp.items.push(temp);
            temps.push(myTemp);
          }
        }
        return temps;
      },
    },
  },
  data() {
    return {
      dialogProfileVisible: false,
      dialogUploadVisible: false,
      addProfile: { name: "", description: "" },
      isNameOK: true,
      selectedTemplate: "",
      temps: [],
    };
  },
  created() {
    this.updateTemplates();
  },
  methods: {
    importAction() {
      this.dialogUploadVisible = true;
    },
    doImport(event) {
      let newProfile = event;
      console.log("new profile: " + JSON.stringify(newProfile));
      this.dialogUploadVisible = false;
      this.addProfile = newProfile;
      //this.emitter.emit("insertAction", newAction);
    },
    cancel() {
      this.$emit("cancel");
    },
    save() {
      if (this.selectedTemplate) {
        let url = this.$baseURL + "config/templates/"+this.selectedTemplate.name;
        let that = this;
        fetch(url)
          .then((res) => res.json())
          .then((data) => {
            //console.log(data);
            that.addProfile.pages = data.pages;
            that.addProfile.actions = data.actions;
            this.$emit("save", that.addProfile);
          })
          .catch((err) => console.log(err.message));
      } else {
        this.$emit("save", this.addProfile);
      }
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
    updateTemplates() {
      let url = this.$baseURL + "config/templates";
      let that = this;
      fetch(url)
        .then((res) => res.json())
        .then((data) => {
          //console.log(data);
          that.temps = data;
        })
        .catch((err) => console.log(err.message));
    },
    changeTemplate() {
      this.addProfile.name = this.selectedTemplate.name
      this.addProfile.description = this.selectedTemplate.description
      console.log(this.selectedTemplate.name)
    },
  },
  watch: {
    visible(visible) {
      this.dialogProfileVisible = visible;
      this.checkName(this.addProfile.name);
    },
    profile(profile) {
      this.addProfile = profile;
    },
    addProfile: {
      deep: true,
      handler(profile) {
        this.checkName(profile.name);
      },
    },
  },
};
</script>

<style>
.p-dialog {
  width: 30% !important;
}
</style>