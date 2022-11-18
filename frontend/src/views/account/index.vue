<template>
  <div class="app-container">
    <el-row>
      <el-col :span="24">
        <el-form label-position="left" ref="form" :model="form">
          <el-row>
            <el-col :span="11">
              <el-form-item label="Name">
                <el-input v-model="form.name"/>
              </el-form-item>
            </el-col>
            <el-col :span="11" style="margin-left: 20px">
              <el-form-item label="Identify ID">
                <el-input v-model="form.identityID"/>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row>
            <el-col :span="11">
              <el-form-item label="Unit">
                <el-input v-model="form.unit"/>
              </el-form-item>
            </el-col>
            <el-col :span="11" style="margin-left: 20px">
              <el-form-item label="Phone Number">
                <el-input v-model="form.phoneNumber"/>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row>
            <el-col :span="11">
              <el-form-item label="Email">
                <el-input v-model="form.email"/>
              </el-form-item>
            </el-col>
            <el-col :span="11" style="margin-left: 20px">
              <el-form-item label="Password">
                <el-input v-model="form.password"/>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row>
            <el-col :span="24">
              <el-form-item>
                <el-select v-model="form.scope" placeholder="Role">
                  <el-option label="Root" value="root"/>
                  <el-option label="Admin" value="admin"/>
                  <el-option label="User" value="user"/>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item>
            <el-button type="primary" @click="onSubmit">Create</el-button>
            <el-button @click="onCancel">Cancel</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>

    <el-row>
      <el-col >
        <el-table table-layout="fixed"
                  :data="tableData.filter(data => !search || data.name.toLowerCase().includes(search.toLowerCase()))"
                  height="500" style="width: 100%;">
          <el-table-column label="Date" prop="date"></el-table-column>
          <el-table-column label="Name" prop="name"></el-table-column>
          <el-table-column align="right">
            <template v-slot:[slotHeader]="scope">
              <el-input v-model="search" size="mini" placeholder="Type to search"/>
            </template>
            <template v-slot="scope">
              <el-button size="mini" @click="handleEdit(scope.$index, scope.row)">Edit</el-button>
              <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">Delete</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>

  </div>
</template>

<script>
export default {
  name: "Account",
  data() {
    return {
      slotHeader: 'header',
      form: {
        name: '',
        identityID: '',
        unit: '',
        phoneNumber: '',
        email: '',
        password: '',
        scope: '',
      },
      tableData: [
        {
          date: '2016-05-03',
          name: 'Tom',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-02',
          name: 'John',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-04',
          name: 'Morgan',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        },
        {
          date: '2016-05-01',
          name: 'Jessy',
          address: 'No. 189, Grove St, Los Angeles'
        }
      ],
      search: '',
    }
  },
  methods: {
    onSubmit() {
      this.$message('submit!')
    },
    onCancel() {
      this.$message({
        message: 'cancel!',
        type: 'warning'
      })
    },
    handleEdit(index, row) {
      console.log(index, row);
    },
    handleDelete(index, row) {
      console.log(index, row);
    }
  }
}
</script>

<style scoped>
.app-container {
  background: #f0f2f5;
}

</style>
