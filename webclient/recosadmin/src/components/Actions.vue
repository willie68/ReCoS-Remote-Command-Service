<template>
  <Toolbar class="p-mt-0 p-pt-0 p-pb-0">
    <template #left>
    <p class="p-ml-2">Pages:</p>
    <Dropdown
      class="p-ml-1 dropdownwidth"
      v-model="activePage"
      :options="profile.pages"
      optionLabel="name"
      placeholder="Select a Page"
    />
    <Button icon="pi pi-plus" @click="newPage" class="p-ml-1" />
    <Button icon="pi pi-trash" class="p-button-warning" />

    <p class="p-ml-2">Actions:</p>
    <Dropdown
      class="p-ml-1 dropdownwidth"
      v-model="selectedAction"
      :options="profile.actions"
      optionLabel="name"
      placeholder="Select an action"
    />
    <SplitButton
      v-tooltip="'Edit'"
      icon="pi pi-pencil"
      :model="actionMenuItems"
      class="p-button-warning"
    ></SplitButton>
    </template>
  </Toolbar>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <div class="p-pb-2 p-pt-2"><b>Descriptions</b></div>
      <Fieldset
        :legend="'Profile: ' + profile.name"
        :toggleable="true"
        class="p-pt-2"
      >
        {{ profile.description }}
      </Fieldset>
      <Fieldset :legend="'Page: ' + activePage.name" :toggleable="true">
        {{ activePage.description }}
      </Fieldset>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <PageSettings :page="activePage" :profile="profile"></PageSettings>
    </SplitterPanel>
  </Splitter>
</template>

<script>
export default {
    name: "Actions",
    props: {
        profile: {},
    },
    data() {
        return {}
    },
}
</script>