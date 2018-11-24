<template>
  <div class="server-control">
      <button v-on:click="runClenaup" type="button" class="btn btn-primary">Run Cleanup</button>
  </div>
</template>

<script>
import api from '../api'

export default {
  name: 'ServerControl',
  props: {
    msg: String
  },
  methods: {
    runClenaup: function(event) {
      api.startCleanup()
    },
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
