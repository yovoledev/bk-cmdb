<template>
    <div class="yovole-layout clearfix">
        <cmdb-business-selector class="business-selector" v-model="business"></cmdb-business-selector>
        <bk-button class="yovole-table-btn-refresh" type="primary" style="display: none;" :disabled="$loading()" @click="handleRefresh">
            {{$t("HostResourcePool['刷新查询']")}}
        </bk-button>
        <cmdb-table class="yovole-table" ref="table" 
            :loading="$loading('post_searchBusiness_list')" 
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="157">
        </cmdb-table>
    </div>
</template>

<script>
    import {mapActions, mapGetters} from 'vuex'
    import cmdbBusinessSelector from '@/components/ui/selector/business'

    export default {
        components: {
            cmdbBusinessSelector
        },
        data () {
            return {
                business: '',
                businessResolver: null,
                table: {
                    header: [{id: 'bk_biz_id', name: 'bk_biz_id'},
                        {id: 'bk_host_id', name: 'bk_host_id'},
                        {id: 'bk_host_innerip', name: 'bk_host_innerip'},
                        {id: 'file_content', name: 'file_content'}],
                    list: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    }
                }
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
                await Promise.all([this.getBusiness()])
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
                this.yovoleSearchTest({
                    bkBizId: this.business,
                    config: {
                        requestId: 'yovoletest',
                        cancelPrevious: true
                    }
                }).then(data => {
                    this.table.pagination.count = data.length
                    this.table.list = data
                    return data
                }).catch(e => {
                    this.table.list = []
                    this.table.pagination.count = 0
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

    .yovole-table {
       margin: 20px 20px;
    }
</style>