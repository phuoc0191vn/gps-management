<template>
  <div class="app-container">
    <div style="width: 100%; height: 500px;">
      <MglMap
        :key="coordinates[0]"
        :accessToken="mapboxAccessToken"
        :mapStyle.sync="mapStyle"
        :center="coordinates"
        :zoom="15">

        <MglMarker v-if="showMarker" :key="showMarker" :coordinates="coordinates" color="blue">
          <MglPopup :showed="true">
            <div>{{ deviceName }}</div>
          </MglPopup>
        </MglMarker>

        <MglGeojsonLayer
          :key="showPolyLine"
          v-if="showPolyLine"
          sourceId="mySource"
          :source="geoJsonSource"
          :layerId="geoJsonLayer.id"
          :layer="geoJsonLayer"/>
      </MglMap>
    </div>

    <el-row style="margin-bottom: 20px;">
      <p style="margin-bottom: 10px;">History Journey:</p>
      <el-col :span="3">
        <el-date-picker
          v-model="date"
          type="date"
          format="yyyy/MM/dd"
          value-format="yyyy-MM-dd"
          placeholder="date"
          style="width: 100%"/>
      </el-col>
      <el-col :span="3" style="margin-left: 10px">
        <el-button type="primary" @click="onShowPolyLine">Show</el-button>
      </el-col>
    </el-row>

    <el-row v-if="scope !== 'user'">
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
              <el-button size="mini" type="info" @click="handleShow(tableScope.$index, tableScope.row)">View</el-button>
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
import {MglGeojsonLayer, MglMap, MglMarker, MglPopup} from "vue-mapbox";
import {getActivityLog, getCurrentLocation} from "@/api/activityLog";
import {sanitizeObject} from "@/utils";
import {accountInfo, datatable} from "@/api/account";
import Device from "@/views/device";

export default {
  components: {Device, MglMap, MglMarker, MglPopup, MglGeojsonLayer},
  data() {
    return {
      scope: '',
      accountID: 'own-account',
      deviceName: '',
      mapboxAccessToken: 'pk.eyJ1IjoicGh1b2NubiIsImEiOiJjbGIwOHc3dDAxMG5vM25tdWZhcTk5M3dmIn0.Y7yzo9ktn1Vuup2d5cG2mQ',
      mapStyle: "mapbox://styles/mapbox/streets-v12",
      coordinates: [1, 1],
      showPolyLine: false,
      showMarker: true,
      geoJsonSource: {
        type: "geojson",
        data: {
          type: 'Feature',
          geometry: {
            type: 'LineString',
            coordinates: [0, 0]
          }
        }
      },
      geoJsonLayer: {
        id: 'testLayer',
        type: 'line',
        layout: {
          'line-join': 'round',
          'line-cap': 'round'
        },
        paint: {
          'line-color': '#ff0000',
          'line-width': 8
        }
      },
      date: '',
      slotHeader: 'header',
      tableData: [],
      page: 1,
      pageSize: 10,
      total: 0,
      search: '',
    };
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
    getCurrentLocation() {
      getCurrentLocation(this.accountID).then(resp => {
        this.deviceName = resp.data.name
        this.coordinates = [resp.data.currentLocation.latitude, resp.data.currentLocation.longitude]
      }).catch(err => {
        console.log('error: ', err)
      })
    },
    getDatatable() {
      datatable(this.page, this.pageSize).then(resp => {
        this.tableData = resp.data["Data"]
        this.total = resp.data["Total"]
      }).catch(e => {
        console.log('error: ', e)
      })
    },
    handleShow(index, row) {
      this.clear()
      this.accountID = row.id
      this.getCurrentLocation()
    },
    onShowPolyLine() {
      this.geoJsonSource.data.geometry.coordinates = []
      getActivityLog(this.accountID, this.date).then(resp => {
        let lineCoordinates = []
        for (let i = 0; i < resp.data.locations.length; i++) {
          lineCoordinates.push([resp.data.locations[i].longitude, resp.data.locations[i].latitude])
        }

        this.geoJsonSource.data.geometry.coordinates = lineCoordinates
        this.showPolyLine = true
        this.showMarker = false
      }).catch(err => {
        console.log('error: ', err)
      })
    },
    clear() {
      this.showPolyLine = false
      this.showMarker = true
      this.coordinates = []
      this.geoJsonSource = {
        type: "geojson",
        data: {
          type: 'Feature',
          geometry: {
            type: 'LineString',
            coordinates: [0, 0]
          }
        }
      }
      this.geoJsonLayer = {
        id: 'testLayer',
        type: 'line',
        layout: {
          'line-join': 'round',
          'line-cap': 'round'
        },
        paint: {
          'line-color': '#ff0000',
          'line-width': 8
        }
      }
    }
  },
  created() {
    this.getCurrentLocation()
    this.getDatatable()
    accountInfo().then(resp => {
      this.scope = resp.data.scope
    })
  }
};
</script>

<style scoped>
.app-container {
  /*background: #f0f2f5;*/
}

</style>
