<template>
  <div class="main">
    <h1>Server Manager</h1>
    <ServerControl/>
  </div>
</template>

<script>
// @ is an alias to /src
import ServerControl from '@/components/ServerControl.vue'

export default {
  name: 'main',
  components: {
    ServerControl
  }
}
</script>
