<template>
  <div class="app-container">
    <el-row>
      <el-form label-position="top" ref="form" :model="form" style="margin-left: 10px;" size="medium">
        <el-row>
          <el-col :span="6">
            <el-form-item label="Name">
              <el-input v-model="form.name"/>
            </el-form-item>
          </el-col>

          <el-col :span="6" style="margin-left: 30px;">
            <el-form-item label="License Plate Number">
              <el-input v-model="form.licensePlateNumber"/>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="6">
            <el-form-item label="Type">
              <el-select v-model="form.type">
                <el-option label="Car" value="1"/>
                <el-option label="Motorbike" value="0"/>
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="Status" style="margin-left: 30px;">
              <el-select v-model="form.status">
                <el-option label="Enable" value="1"/>
                <el-option label="Disable" value="0"/>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Account" v-if="scope !== 'user'">
          <el-select v-model="form.accountID">
            <el-option
              v-for="(account, index) in childAccounts"
              :key="index"
              :label="account.label"
              :value="account.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button v-if="!isUpdateForm && scope !== 'user'" type="primary" @click="onSubmit">Create</el-button>
          <el-button v-if="isUpdateForm" type="primary" @click="onUpdate">Update</el-button>
          <el-button @click="onCancel">Cancel</el-button>
        </el-form-item>
      </el-form>
    </el-row>

    <el-row>
      <el-col>
        <el-table ref="datatable" table-layout="fixed"
                  :data="tableData.filter(data => !search || data.email.toLowerCase().includes(search.toLowerCase()))"
                  style="width: 100%;">
          <el-table-column align="center" label="Index" type="index" width="70"></el-table-column>
          <el-table-column align="center" label="Name" prop="name"></el-table-column>
          <el-table-column align="center" label="License Plate Number" prop="licensePlateNumber"></el-table-column>
          <el-table-column align="center" label="Type" prop="type"></el-table-column>
          <el-table-column align="center" label="Status" prop="status"></el-table-column>
          <el-table-column align="center" label="Account" prop="email"></el-table-column>
          <el-table-column align="center" label="Created At" prop="createdTime" width="250"></el-table-column>
          <el-table-column align="right" width="300">
            <template v-slot:[slotHeader]="tableScope">
              <el-input v-model="search" size="mini" placeholder="Type to search"/>
            </template>
            <template v-slot="tableScope">
              <el-button size="mini" type="primary" @click="handleEdit(tableScope.$index, tableScope.row)">Edit
              </el-button>
              <el-button size="mini" style="width: 25%;" type="info"
                         @click="toggle(tableScope.$index, tableScope.row)">
                {{ tableScope.row.status === "Enable" ? "Disable" : "Enable" }}
              </el-button>
              <el-button size="mini" v-if="scope !== 'user'" type="danger"
                         @click="handleDelete(tableScope.$index, tableScope.row)">Delete
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination layout="total, sizes, prev, pager, next" :page-size="pageSize" :page-sizes="pageSizes"
                       :total="total"
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
import {accountInfo, getChildAccounts,} from "@/api/account";
import SweetCode from '@/components/SweetCode/index'
import {sanitizeObject} from '@/utils'
import {addDevice, datatable, deleteDevice, getDetail, toggleStatus, updateDevice} from "@/api/device";

export default {
  name: "Device",
  components: {SweetCode},
  data() {
    return {
      childAccounts: [],
      scope: '',
      slotHeader: 'header',
      isUpdateForm: false,
      deviceIdUpdate: '',
      form: {
        name: '',
        licensePlateNumber: '',
        type: '',
        status: '',
        accountID: '',
        email: '',
      },
      tableData: [],
      page: 1,
      pageSize: 10,
      pageSizes: [10, 20, 30, 50, 100],
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
      addDevice(this.form).then(() => {
        this.$message('add device successfully')
        this.clear()
        this.getDatatable()
      }).catch(() => {
        this.$message({
          message: 'unable to device account',
          type: 'error'
        })
      })

    },
    onUpdate() {
      updateDevice(this.deviceIdUpdate, this.form).then(() => {
        this.$message('update device successfully')
        this.clear()
        this.getDatatable()
      }).catch(() => {
        this.$message({
          message: 'unable to device account',
          type: 'error'
        })
      })
    },
    onCancel() {
      this.clear()
    },
    toggle(index, row) {
      let status = 0
      if (row.status === "Disable") {
        status = 1
      }

      if (confirm("Do you really want to update status device?")) {
        toggleStatus(row.id, status).then(resp => {
          this.$message('update status device successfully')
          this.getDatatable()
        }).catch(e => {
          this.$message({
            message: 'unable to update status device: ' + e,
            type: 'error'
          })
        })
      }

    },
    handleEdit(index, row) {
      getDetail(row.id).then(resp => {
        this.form.name = resp.data.name
        this.form.licensePlateNumber = resp.data.licensePlateNumber
        this.form.type = '' + resp.data.type
        this.form.status = '' + resp.data.status
        this.form.accountID = resp.data.accountID

        this.isUpdateForm = true
        this.deviceIdUpdate = row.id
      }).catch(e => {
        this.$message({
          message: 'unable to update device: ' + e,
          type: 'error'
        })
      })
    },
    handleDelete(index, row) {
      if (confirm("Do you really want to delete device?")) {
        deleteDevice(row.id).then(resp => {
          this.$message('delete device successfully')
          this.getDatatable()
        }).catch(e => {
          console.log('error: ', e)
        })
      }
    },
    getDatatable() {
      datatable(this.page, this.pageSize).then(resp => {
        this.tableData = resp.data["Data"]
        this.beautifulDatatable()
        this.total = resp.data["Total"]
      }).catch(e => {
        console.log('error: ', e)
      })
    },
    clear() {
      this.form = {}
      this.isUpdateForm = false
      this.deviceIdUpdate = ''
    },
    beautifulDatatable() {
      for (let i = 0; i < this.tableData.length; i++) {
        if (this.tableData[i].type === 0) {
          this.tableData[i].type = 'Motorbike'
        }
        if (this.tableData[i].type === 1) {
          this.tableData[i].type = 'Car'
        }

        if (this.tableData[i].status === 0) {
          this.tableData[i].status = 'Disable'
        }
        if (this.tableData[i].status === 1) {
          this.tableData[i].status = 'Enable'
        }

        for (let j = 0; j < this.childAccounts.length; j++) {
          if (this.tableData[i].accountID === this.childAccounts[j].value) {
            this.tableData[i].email = this.childAccounts[j].label
          }
        }
      }
    },
    getChildAccounts() {
      getChildAccounts().then(resp => {
        for (let i = 0; i < resp.data.length; i++) {
          this.childAccounts.push({
            label: resp.data[i].email,
            value: resp.data[i].id
          })
        }
      })
    }
  },
  async created() {
    await accountInfo().then(resp => {
      this.scope = resp.data.scope
    })
    await this.getChildAccounts()
    this.getDatatable()
  },
}
</script>

<style scoped>
.app-container {
  background: #f0f2f5;
}

</style>
