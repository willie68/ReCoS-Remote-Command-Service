<template>
  <div class="layout-content">
    <div class="content-section implementation">
      <Toolbar>
        <template #left>
          <h1>ReCoS Admin</h1>
        </template>

        <template #right>
          <span class="p-input-icon-right">
            Kennwort
            <InputText ref="pwd" v-model="pwd" :type="pwdType"></InputText>
            <i class="pi pi-eye-slash" @click="togglePwdView()" />
            <i
              v-if="!showPwd"
              class="pi pi-eye-slash"
              @click="togglePwdView()"
            />
            <i v-if="showPwd" class="pi pi-eye" @click="togglePwdView()" />
          </span>
          <SplitButton
            label="Save"
            icon="pi pi-check"
            :model="items"
            class="p-button-warning"
          ></SplitButton>
          <Button icon="pi pi-cog" class="p-mr-1" />
        </template>
      </Toolbar>
      <Splitter style="height: 400px">
        <SplitterPanel :size="20">
          <ProfileAccordion></ProfileAccordion>
        </SplitterPanel>
        <SplitterPanel :size="80"> Panel 2 </SplitterPanel>
      </Splitter>
    </div>
  </div>
</template>

<script>
import ProfileAccordion from "./ProfileAccordion.vue";

export default {
  name: "Page",
  components: {
    ProfileAccordion,
  },
  props: {
    msg: String,
  },
  data() {
    return {
      showPwd: false,
      pwdType: "password",
      items: [
        {
          label: "Update",
          icon: "pi pi-refresh",
        },
        {
          label: "Delete",
          icon: "pi pi-times",
        },
        {
          label: "Vue Website",
          icon: "pi pi-external-link",
          command: () => {
            window.location.href = "https://vuejs.org/";
          },
        },
        {
          label: "Upload",
          icon: "pi pi-upload",
          command: () => {
            window.location.hash = "/fileupload";
          },
        },
      ],
      profileitems: [
        {
          label: "Options",
          items: [
            {
              label: "Update",
              icon: "pi pi-refresh",
              command: () => {
                this.$toast.add({
                  severity: "success",
                  summary: "Updated",
                  detail: "Data Updated",
                  life: 3000,
                });
              },
            },
            {
              label: "Delete",
              icon: "pi pi-times",
              command: () => {
                this.$toast.add({
                  severity: "warn",
                  summary: "Delete",
                  detail: "Data Deleted",
                  life: 3000,
                });
              },
            },
          ],
        },
        {
          label: "Navigate",
          items: [
            {
              label: "Vue Website",
              icon: "pi pi-external-link",
              url: "https://vuejs.org/",
            },
            {
              label: "Router",
              icon: "pi pi-upload",
              url: "https://vuejs.org/",
            },
          ],
        },
      ],
    };
  },
  methods: {
    togglePwdView() {
      this.showPwd = !this.showPwd;
      if (this.showPwd) {
        this.pwdType = "text";
      } else {
        this.pwdType = "password";
      }
    },
    toggle(event) {
      this.$refs.menu.toggle(event);
    },
  },
};
</script>

<style>
.p-button {
  margin-bottom: 0.5rem;
}
</style>