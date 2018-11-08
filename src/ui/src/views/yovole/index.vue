<template>
  <div class="yovole-layout clearfix">
    <cmdb-business-selector class="business-selector" v-model="business"></cmdb-business-selector>
    <bk-button class="yovole-table-btn-refresh" type="primary" style="display: none;"
               :disabled="$loading()"
               @click="handleRefresh">
      {{$t("HostResourcePool['刷新查询']")}}
    </bk-button>
    <div class="yovoleinfo">{{ info }}</div>
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import cmdbBusinessSelector from '@/components/ui/selector/business'

  export default {
    components: {
        cmdbBusinessSelector
    },
    data () {
      return {
        business: '',
        businessResolver: null,
        info: null
      }
    },
    computed: {
      ...mapGetters(['supplierAccount'])
    },
    watch: {
      business (business) {
        this.handleRefresh()
      }
    },
    async created () {
      try {
        await Promise.all([
          this.getBusiness()
        ])
        await this.handleRefresh()
      } catch (e) {
        console.log(e)
      }
    },
    methods: {
      ...mapActions('yovole', ['yovoleSearchTest']),

      getBusiness () {
        return new Promise((resolve, reject) => {
          this.businessResolver = () => {
            this.businessResolver = null
            resolve()
          }
        })
      },
      async handleRefresh () {
        console.log('call refresh yovoletest method')
        this.yovoleSearchTest({
          bkBizId: this.business,
          config: {
            requestId: 'yovoletest',
            cancelPrevious: true
          }
        }).then(data => {
          if (data) {
            console.dir(data)
            this.info = JSON.stringify(data)
          }
          return data
        }).catch(e => {
          this.info = null
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .yovole-layout {
    padding: 0;
    height: 100%;
  }

  .business-selector {
    display: block;
    width: auto;
    margin: 20px 20px 13px;
  }

  .yovoleinfo {
    margin: 20px 20px;
  }
</style>
