<script>
import { onMounted, reactive, ref } from 'vue'
import { VDataTable } from 'vuetify/labs/VDataTable'
import axios from 'axios'

export default {
  setup () {

    const state = reactive({
        headers: [],
        items: [],
    })
    const load = async() => {
      try {
        const rsp = await axios.get('/api/v1/participant')
        if (rsp.status == 200) {
          state.items = rsp.data.participants
          state.headers = Object.keys(rsp.data.participants[0])
        }
          console.log(state)
      } catch(error) {
          console.log("error happend", error)
          throw new Error(error.response.data)
      }
    }

    //onMounted( async() =>{
      load()
    //})

    return { state, load }
  },
components: {
    VDataTable,
  },
  data: () => ({
      headers: [
        { title: 'Navn', align: 'start', key: 'name' },
        { title: 'Telefon', align: 'start', key: 'phone' },
        { title: 'E-mail', key: 'email' },
        { title: 'Betalt', align: 'end', key: 'paid' },
      ],
  }),
  computed: {
  },
  methods: {
      rowClick(item, row) {
          console.log(item, row)
      }
  }
}
</script>

<template>
    <h1>Tilmeldte</h1>
  <v-data-table
    :headers="headers"
    :items="state.items"
    class="hover"
    _density="compact"
    show-select
    item-key="name"
                @click:row="rowClick"

  ></v-data-table>
</template>
