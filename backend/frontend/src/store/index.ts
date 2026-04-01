import { createStore } from 'vuex'

export default createStore({
  state: {
    videoInfo:{
     
    },
    videoId:''
  },
  mutations: {
    changeVideoId(state,video){
      state.videoInfo = video;
      console.log(video,state.videoInfo);
      
    },
    addVideoId(state,videoId){
      state.videoId = videoId;
    }
  },
  actions: {
  },
  modules: {
  }
})
