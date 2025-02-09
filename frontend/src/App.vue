<template>
  <div style="padding: 20px">
    <el-table :data="pingResults" style="width: 100%">
      <el-table-column prop="ip" label="IP Address"></el-table-column>
      <el-table-column prop="ping_time" label="Ping Time (ms)"></el-table-column>
      <el-table-column prop="last_success_at" label="Last Successful Ping"></el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
import axios from 'axios';

interface PingResult {
  ip: string;
  ping_time: number;
  last_success_at: string;
}

export default defineComponent({
  name: 'App',
  setup() {
    const pingResults = ref<PingResult[]>([]);

    const fetchData = async () => {
      try {
        const response = await axios.get('/api/stats');
        pingResults.value = response.data;
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    onMounted(() => {
      fetchData();
      setInterval(fetchData, 5000);
    });

    return { pingResults };
  }
});
</script>

<style scoped>
</style>
