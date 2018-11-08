import $http from '@/api'

const state = {

}

const getters = {

}

const actions = {
  yovoleSearchTest ({ commit, state, dispatch }, {bkBizId, config}) {
    return $http.get(`collector/yovole/test/${bkBizId}/filecontent`, {}, config)
  }
}

const mutations = {

}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
