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
              <el-form-item v-if="!isUpdateForm" label="Email">
                <el-input v-model="form.email"/>
              </el-form-item>
            </el-col>
            <el-col :span="11" style="margin-left: 20px">
              <el-form-item v-if="!isUpdateForm" label="Password">
                <el-input v-model="form.password"/>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row>
            <el-col :span="24">
              <el-form-item v-if="!isUpdateForm">
                <el-select v-model="form.scope" placeholder="Role">
                  <el-option label="Root" value="root"/>
                  <el-option label="Admin" value="admin"/>
                  <el-option label="User" value="user"/>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item>
            <el-button v-if="!isUpdateForm && scope !== 'user'" type="primary" @click="onSubmit">Create</el-button>
            <el-button v-if="isUpdateForm" type="primary" @click="onUpdate">Update</el-button>
            <el-button @click="onCancel">Cancel</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <el-table ref="datatable" table-layout="fixed"
                  :data="tableData.filter(data => !search || data.email.toLowerCase().includes(search.toLowerCase()))"
                  style="width: 100%;">
          <el-table-column align="center" label="Index" type="index" width="70"></el-table-column>
          <el-table-column align="center" label="Email" prop="email"></el-table-column>
          <el-table-column align="center" label="Scope" prop="scope"></el-table-column>
          <el-table-column align="center" label="Created By" prop="createdBy"></el-table-column>
          <el-table-column align="center" label="Created Time" prop="createdTime"></el-table-column>
          <el-table-column align="center" label="Updated Time" prop="updatedTime"></el-table-column>
          <el-table-column align="right" width="300">
            <template v-slot:[slotHeader]="tableScope">
              <el-input v-model="search" size="mini" placeholder="Type to search"/>
            </template>
            <template v-slot="tableScope">
              <el-button size="mini" type="info" @click="handleView(tableScope.$index, tableScope.row)">View</el-button>
              <el-button size="mini" type="primary" @click="handleEdit(tableScope.$index, tableScope.row)">Edit</el-button>
              <el-button v-if="scope !== 'user'" size="mini" type="warning" @click="handleReset(tableScope.$index, tableScope.row)">Reset</el-button>
              <el-button v-if="scope !== 'user'" size="mini" type="danger" @click="handleDelete(tableScope.$index, tableScope.row)">Delete</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination layout="total, sizes, prev, pager, next" :page-size="pageSize" :total="total"
                       @size-change="changePerPage" @current-change="setPage">
        </el-pagination>
      </el-col>
    </el-row>

    <el-dialog id="eModal" :title="this.titleModel" :visible.sync="dialogVisible">
      <sweet-code :code="SanitizeObject(this.jsonModel)">
        <code class="json"></code>
      </sweet-code>
    </el-dialog>
  </div>
</template>

<script>
import {addUser, datatable, deleteUser, resetUser, updateUser, userDetail, userInfo} from "@/api/user";
import SweetCode from '@/components/SweetCode/index'
import {sanitizeObject} from '@/utils'

export default {
  name: "Account",
  components: {SweetCode},
  data() {
    return {
      scope: '',
      slotHeader: 'header',
      isUpdateForm: false,
      userIDUpdate: '',
      form: {
        name: '',
        identityID: '',
        unit: '',
        phoneNumber: '',
        email: '',
        password: '',
        scope: '',
      },
      tableData: [],
      page: 1,
      pageSize: 10,
      total: 0,
      search: '',
      dialogVisible: false,
      titleModel: '',
      jsonModel: '',
    }
  },
  methods: {
    SanitizeObject(obj) {
      return sanitizeObject(obj)
    },
    changePerPage(val) {
      this.pageSize = val
      this.getDatatable()
    },
    setPage(val) {
      this.page = val
      this.getDatatable()
    },
    onSubmit() {
      addUser(this.form).then(() => {
        this.$message('add account successfully')
        this.clear()
        this.getDatatable()
      }).catch(() => {
        this.$message({
          message: 'unable to add account',
          type: 'error'
        })
      })

    },
    onUpdate() {
      this.form.email = ''
      this.form.password = ''

      updateUser(this.userIDUpdate, this.form).then(() => {
        this.$message('update account successfully')
        this.clear()
      }).catch(() => {
        this.$message({
          message: 'unable to update account',
          type: 'error'
        })
      })
    },
    onCancel() {
      this.clear()
    },
    handleView(index, row) {
      userDetail(row.id).then(resp => {
        this.dialogVisible = true
        this.jsonModel = resp.data
        this.titleModel = 'Account: ' + resp.data.name
      }).catch(e => {
        console.log('error: ', e)
      })

    },
    handleEdit(index, row) {
      userDetail(row.id).then(resp => {
        this.isUpdateForm = true
        this.form.name = resp.data.name
        this.form.identityID = resp.data.identityID
        this.form.unit = resp.data.unit
        this.form.phoneNumber = resp.data.phoneNumber
        this.form.scope = resp.data.scope

        this.userIDUpdate = resp.data.userID
      }).catch(e => {
        console.log('error: ', e)
      })
    },
    handleReset(index, row) {
      if (confirm("Do you really want to reset account?")) {
        resetUser(row.id).then(resp => {
          this.$message('reset account successfully')
        }).catch(e => {
          console.log('error: ', e)
        })
      }
    },
    handleDelete(index, row) {
      if (confirm("Do you really want to delete account?")) {
        deleteUser(row.id).then(resp => {
          this.$message('delete account successfully')
          this.getDatatable()
        }).catch(e => {
          console.log('error: ', e)
        })
      }
    },
    getDatatable() {
      datatable(this.page, this.pageSize).then(resp => {
        this.tableData = resp.data["Data"]
        this.total = resp.data["Total"]
      }).catch(e => {
        console.log('error: ', e)
      })
    },
    clear() {
      this.form = {}
      this.isUpdateForm = false
      this.userIDUpdate = ''
    }
  },
  created() {
    this.getDatatable()
    userInfo().then(resp => {
      this.scope = resp.data.scope
    })
  },
}
</script>

<style scoped>
.app-container {
  background: #f0f2f5;
}

</style>
