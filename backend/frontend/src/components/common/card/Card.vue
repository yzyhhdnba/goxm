<template>
  <div class="card">
    <card-list :maindata-item="maindata" :mdname-item="mdname"></card-list>
  </div>
</template>

<script>
import CardList from "./cardList";

export default {
  name: "Card",
  components: {
    CardList,
  },
  props: {
    maindata: {
      type: Array,
      default() {
        return [];
      },
    },
    mdname: {
      type: String,
      default: "",
    },
  },
};
</script>

<style scoped>

</style>