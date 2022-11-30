<template>
  <div class="app-container">
    <el-row v-if="scope === 'user'">
      <el-col :span="24">
        <el-form label-position="left" ref="form" :model="form">
          <el-form-item label="Title">
            <el-input v-model="form.title"/>
          </el-form-item>
          <el-form-item label="Body">
            <el-input type="textarea" v-model="form.body"></el-input>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="onSubmit">Submit</el-button>
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
          <el-table-column align="center" label="Status" prop="status"></el-table-column>
          <el-table-column align="center" label="Title" prop="title"></el-table-column>
          <el-table-column align="right" width="300">
            <template v-slot:[slotHeader]="tableScope">
              <el-input v-model="search" size="mini" placeholder="Type to search"/>
            </template>
            <template v-slot="tableScope">
              <el-button size="mini" type="info" @click="handleView(tableScope.$index, tableScope.row)">Detail
              </el-button>
              <el-button :disabled="tableScope.row.status==='Processed'" v-if="scope !== 'user'" size="mini"
                         type="primary"
                         @click="handleDone(tableScope.$index, tableScope.row)">Done
              </el-button>
              <el-button v-if="scope !== 'user'" size="mini" type="danger"
                         @click="handleDelete(tableScope.$index, tableScope.row)">Delete
              </el-button>
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
import {sanitizeObject} from "@/utils";
import {accountInfo} from "@/api/account";
import {datatable, deleteContact, submitContact, submitDoneContact} from "@/api/contact";
import SweetCode from "@/components/SweetCode";

const statusUnprocessed = 'Unprocessed';
const statusProcessed = 'Processed';

export default {
  name: "Contact",
  components: {SweetCode},
  data() {
    return {
      scope: '',
      slotHeader: 'header',
      form: {
        title: '',
        body: '',
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
      submitContact(this.form).then(() => {
        this.$message('submit contact successfully')
        this.clear()
        this.getDatatable()
      }).catch(() => {
        this.$message({
          message: 'unable to submit contact',
          type: 'error'
        })
      })
    },
    onCancel() {
      this.clear()
      this.getDatatable()
    },
    handleView(index, row) {
      this.dialogVisible = true
      this.titleModel = 'Account: ' + row.email
      this.jsonModel = row
    },
    handleDone(index, row) {
      if (confirm("Are you sure?")) {
        submitDoneContact(row.id).then( () => {
          this.$message.info("successfully")
          this.clear()
          this.getDatatable()
        }).catch(err => {
          console.log('error: ', err)
        })
      }
    },
    handleDelete(index, row) {
      if (confirm("Do you want to delete this contact?")) {
        deleteContact(row.id).then( () => {
          this.$message.info("successfully")
          this.clear()
          this.getDatatable()
        }).catch(err => {
          console.log('error: ', err)
        })
      }
    },
    getAccountInfo() {
      accountInfo().then(resp => {
        this.scope = resp.data.scope
      }).catch(err => {
        console.log('error get account info: ', err)
      })
    },
    getDatatable() {
      datatable(this.page, this.pageSize).then(resp => {
        this.tableData = resp.data["Data"]
        for (let i = 0; i < this.tableData.length; i++) {
          if (this.tableData[i].status === 0) {
            this.tableData[i].status = statusUnprocessed
          }
          if (this.tableData[i].status === 1) {
            this.tableData[i].status = statusProcessed
          }
        }
        this.total = resp.data["Total"]
      }).catch(e => {
        console.log('error: ', e)
      })
    },
    clear() {
      this.form = {}
    }
  },
  async created() {
    await this.getAccountInfo()
    this.getDatatable()
  }
}
</script>

<style scoped>

</style>
