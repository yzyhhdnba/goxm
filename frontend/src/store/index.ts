import { createStore } from 'vuex'

export default createStore({
  state: {
    videoInfo: {},
    videoId: '',
    role: localStorage.getItem('role') || '',
    userInfo: (() => {
      const raw = localStorage.getItem('userInfo');
      return raw ? JSON.parse(raw) : null;
    })(),
  },
  mutations: {
    changeVideoId(state,video){
      state.videoInfo = video;
    },
    addVideoId(state,videoId){
      state.videoId = videoId;
    },
    changeRole(state, role) {
      state.role = role || '';
    },
    setAuth(state, payload) {
      state.userInfo = payload.userInfo || null;
      state.role = payload.role || '';
    },
    clearAuth(state) {
      state.userInfo = null;
      state.role = '';
      state.videoInfo = {};
      state.videoId = '';
    }
  },
  actions: {
  },
  modules: {
  }
})
