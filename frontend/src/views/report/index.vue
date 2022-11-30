<template>
  <div class="app-container">
    <el-row>
      <el-col :span="24">
        <el-form label-position="left" ref="form" :model="form">
          <!--          <el-row style="margin-bottom: 10px; font-weight: bold;">Filter: </el-row>-->
          <el-row>
            <el-form-item label="Filter Date">
              <el-col :span="3">
                <el-date-picker
                  v-model="form.startTime"
                  type="date"
                  format="yyyy/MM/dd"
                  value-format="yyyy-MM-dd"
                  placeholder="from"
                  style="width: 100%"/>
              </el-col>
              <el-col :span="3" style="margin-left: 10px;">
                <el-date-picker
                  v-model="form.endTime"
                  type="date"
                  format="yyyy/MM/dd"
                  value-format="yyyy-MM-dd"
                  placeholder="to"
                  style="width: 100%"/>
              </el-col>
            </el-form-item>
          </el-row>

          <el-row>
            <el-form-item label="Device" v-if="scope !== 'user'">
              <el-select v-model="form.accountID" style="margin-left: 20px;">
                <el-option
                  v-for="(account, index) in childAccounts"
                  :key="index"
                  :label="account.label"
                  :value="account.value"
                />
              </el-select>
            </el-form-item>
          </el-row>

          <el-form-item>
            <el-button type="primary" @click="onSubmit">Generate</el-button>
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
          <el-table-column align="center" label="Name" prop="name"></el-table-column>
          <el-table-column align="center" label="Status" prop="status"></el-table-column>
          <el-table-column align="center" label="Created Time" prop="createdTime"></el-table-column>
          <el-table-column align="right" width="300">
            <template v-slot:[slotHeader]="tableScope">
              <el-input v-model="search" size="mini" placeholder="Type to search"/>
            </template>
            <template v-slot="tableScope">
              <el-button size="mini" type="info" @click="handleDownload(tableScope.$index, tableScope.row)">Download
              </el-button>
              <el-button size="mini" type="danger" @click="handleDelete(tableScope.$index, tableScope.row)">Delete
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination layout="total, sizes, prev, pager, next" :page-size="pageSize" :total="total"
                       @size-change="changePerPage" @current-change="setPage">
        </el-pagination>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import {datatable, deleteReport, downloadReport, generateReport} from "@/api/report";
import {getEnableDevice} from "@/api/device";
import {accountInfo, getChildAccounts} from "@/api/account";
import {sanitizeObject} from "@/utils";

const statusUnprocessed = 'Unprocessed';
const statusProcessing = 'Processing';
const statusProcessed = 'Processed';

export default {
  name: "Report",
  data() {
    return {
      deviceID: '',
      scope: '',
      childAccounts: [],
      form: {
        startTime: '',
        endTime: '',
        accountID: '',
      },
      slotHeader: 'header',
      tableData: [],
      page: 1,
      pageSize: 10,
      total: 0,
      search: '',
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
      let id = this.deviceID
      if (this.scope !== 'user') {
        id = this.form.accountID
      }
      generateReport(id, this.form.startTime, this.form.endTime).then(resp => {
        this.$message('your report is generating')
        this.getDatatable()
      }).catch(err => {
        console.log('error: ', err)
      })
    },
    onCancel() {
      this.form = {}
    },
    handleDownload(index, row) {
      downloadReport(row.id).then(resp => {
        let blob = new Blob([resp], {type: 'text/csv'}), url = window.URL.createObjectURL(blob)
        window.open(url)
      }).catch(err => {
        console.log('error: ', err)
      })
    },
    handleDelete(index, row) {
      if (confirm("Do you really want to delete report?")) {
        deleteReport(row.id).then(resp => {
          this.$message.info("successfully")
          this.getDatatable()
        }).catch(err => {
          console.log('error: ', err)
        })
      }
    },
    getDatatable() {
      datatable(this.page, this.pageSize).then(resp => {
        this.tableData = resp.data["Data"]
        for (let i = 0; i < this.tableData.length; i++) {
          if (this.tableData[i].status === 0) {
            this.tableData[i].status = statusUnprocessed
          }
          if (this.tableData[i].status === 1) {
            this.tableData[i].status = statusProcessing
          }
          if (this.tableData[i].status === 2) {
            this.tableData[i].status = statusProcessed
          }
        }
        this.total = resp.data["Total"]
      }).catch(e => {
        console.log('error: ', e)
      })
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
    },
    getEnableDevice() {
      getEnableDevice().then(resp => {
        this.deviceID = resp.data.id
      }).catch(err => {
        console.log('error get device: ', err)
      })
    },
    getAccountInfo() {
      accountInfo().then(resp => {
        this.scope = resp.data.scope
      }).catch(err => {
        console.log('error get account info: ', err)
      })
    }
  },
  async created() {
    await this.getChildAccounts()
    await this.getEnableDevice()
    await this.getAccountInfo()
    await this.getDatatable()
  }
}
</script>

<style scoped>

</style>
