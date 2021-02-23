<template>
  <el-container>
    <el-container>
      <el-header height="70px">
        <div><b>ReCoS Admin</b></div>
        <div class="settings">
          <el-button round  size="small" icon="el-icon-user" @click="dialogFormVisible = true">Login</el-button>
          <el-button round  size="small" icon="el-icon-setting" @click="dialogFormVisible = true"></el-button>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px">Aside</el-aside>
        <el-main height="700px">Main</el-main>
      </el-container>
      <el-footer height="30px">ReCoS - Remote Command Service (C) 2021 Wifried Klaas</el-footer>
    </el-container>
  </el-container>

  <el-dialog title="Login" v-model="dialogFormVisible">
    <el-form :model="form">
      <el-form-item label="Name" :label-width="formLabelWidth">
        <el-input v-model="form.name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="Password" :label-width="formLabelWidth">
        <el-input v-model="form.password" show-password></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="dialogFormVisible = false"
          >Confirm</el-button
        >
      </span>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: "Page",
  props: {
    msg: String,
  },
  data() {
    return {
      dialogFormVisible: false,
      form: {
        name: "",
        pasword: "",
      },
      formLabelWidth: "120px",
    };
  },
};
</script>

<style scoped>
.settings {
  text-align: right;
}

.el-header {
  background-color: #b3c0d1;
  color: #333;
  text-align: left;
  line-height: 20px;
}
.el-footer {
  background-color: #b3c0d1;
  color: #333;
  text-align: right;
  line-height: 30px;
}

.el-aside {
  background-color: #d3dce6;
  color: #333;
  text-align: center;
  line-height: 300px;
}

.el-main {
  background-color: #e9eef3;
  color: #333;
  text-align: center;
  line-height: 400px;
}
</style>
