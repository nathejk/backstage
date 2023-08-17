<script>
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { email, required } from '@vuelidate/validators'
import axios from 'axios'

export default {
  setup () {
    const route = useRoute()

    const initialState = {
      name: '',
      address: '',
      email: '',
      phone: '',
      team: null,
      days: ['friday', 'saturday', 'sunday'],
      transport: null,
      seatCount: null,
      info: null,
      video: null,
    }

    const load = async() => {
      if (!route.params.id) return initialState
      try {
        const rsp = await axios.get('/api/v1/participant/' + route.params.id)
        if (rsp.status == 200) {
      for (const [key, value] of Object.entries(rsp.data.participant)) {
        state[key] = value
      }
          //return rsp.data.participant
//            state = reactive({
//                ...rsp.data.participant,
//    })
        }
        return initialState
      } catch(error) {
          console.log("error happend", error)
          throw new Error(error.response.data)
      }
    }

    const state = reactive({
      ...initialState,
    })
    onMounted( async() =>{
      load()
    })

    //const values = await load()
    const rules = {
      name: { required },
      address: { required },
      email: { required, email },
      phone: { required },
      team: { required },
      days: { required },
      transport: { required },
      seatCount: { },
      info: { required },
      video: { required },
    }

    const v$ = useVuelidate(rules, state)

    function clear () {
      v$.value.$reset()

      for (const [key, value] of Object.entries(initialState)) {
        state[key] = value
      }
    }

    return { state, load, v$ }
  },
  data() {
    return {
      saved: ref(false),
      teams: ref([
        { slug: 'start', title: 'Primært start (Tommy, Susan, Martin, Dennis)' },
        { slug: 'checkpoint', title: 'Primært checkpoint (Niels Erik, Tobias, Niels, Jakob)' },
        { slug: 'open', title: 'Placer mig' },
        { slug: 'logistik', title: 'Logistik (forbeholdt gøglerlogistik, efter nærmere aftale med Mismis)' },
      ]),
      days: ref([
        { slug: "friday", title: 'Fredag' },
        { slug: "saturday", title: 'Lørdag' },
        { slug: "sunday", title: 'Søndag' },
      ]),
      transports: ref([
        { slug: 'owncar', title: 'Jeg medbringer egen bil' },
        { slug: 'haveseat', title: 'Jeg har en plads i en bil' },
        { slug: 'needseat', title: 'Jeg har brug for plads i en bil' },
      ]),
      seatCounts: ref([
        { slug: '1', title: '1' },
        { slug: '2', title: '2' },
        { slug: '3', title: '3' },
        { slug: '4', title: '4' },
        { slug: '5', title: '5' },
        { slug: '6', title: '6' },
        { slug: '7', title: '7' },
        { slug: '0', title: 'Jeg har ikke ekstra plads' },
      ]),
      yesNo: ref([
        { slug: 'yes', title: 'JA' },
        { slug: 'no', title: 'Nej' },
      ]),
      videos: ref([
        { slug: 'juni', title: 'Ja, I må gerne kontakte mig ifb. optaksvideoer d. 17. juni' },
        { slug: 'august', title: 'Ja, I må gerne kontakte mig ifb. optaksvideoer d. 12.-13. august' },
        { slug: 'both', title: 'Ja, I må gerne kontakte mig ifb. optaksvideoer på begge dage' },
        { slug: 'none', title: 'Nej, jeg ønsker ikke at deltage i optaksvideoerne' },
      ]),
    }
  },
  computed: {
    userId() {
      return this.$route.params.id
    },
    buttonLabel() {
      return this.state.paid ? 'Opdater' : 'Tilmeld'
    },
    paying() {
      return this.saved && !this.state.paid
    },
    receipt() {
      return this.saved && this.state.paid
    },
  },
  methods: {
    async pay() {

        const rsp = await axios.put('/api/v1/participant/' + this.userId, { phone: this.state.mobilepay.replace(/\s/g, '') }) //, { withCredentials: true })
    },
    async save() {
      const result = await this.v$.$validate()
      if (!result) {
          console.log(this.state)
        // notify user form is invalid
        return
      }
      if (this.userId) {
        const data = {
            name: this.state.name,
            address:this.state.address,
            phone:this.state.phone,
            email:this.state.email,
            team:this.state.team,
            days:this.state.days,
            transport:this.state.transport,
            seatCount:this.state.seatCount || "",
            info:this.state.info,
            video:this.state.video,
        }

        const rsp = await axios.patch('/api/v1/participant/' + this.userId, data) //, { withCredentials: true })
          this.saved = true
      } else {
        const rsp = await axios.post('/api/v1/participant', Object.assign({}, this.state)) //, { withCredentials: true })
        if (rsp.status == 201) {
          //this.team = rsp.dat
          //rsp.body.uuid
          this.$router.push({ name: 'deltager', params: { id: rsp.data.participant.uuid }, replace: true })
          this.saved = true
        }
      }

      // perform async actions
      //alert(JSON.stringify(this.state, null, 2))
    }
  }
}
</script>

<template>
  <h1 class="pb-3" >NATHEJK - Gøglertilmelding</h1>
  <v-container class="bg-surface-variant mb-6">
    <v-row justify="center">
      <v-col cols="8">
        <form v-if="paying" @submit.prevent="pay">
          <p class="pb-3">Det koster 50,- kr at deltage som skal indbetales på MobilePay - der kan gå op til 10 minutter før betalingen er registreret, herefter vil du modtage en SMS.</p>
          <v-text-field v-model="state.mobilepay" label="Telefonnummer" prepend-icon="mdi-cellphone-basic"></v-text-field>
          <v-btn class="me-4" type="submit">Send betalings SMS</v-btn>
        </form>
        <div v-else-if="receipt">
          <p class="pb-3">Dine oplysninger er opdateret.</p>
        </div>
        <form v-else @submit.prevent="save">
        <v-text-field v-model="state.name"    label="Navn" density="compact" :error-messages="v$.name.$errors.map(e => e.$message)" @input="v$.name.$touch"
      @blur="v$.name.$touch"></v-text-field>
        <v-text-field v-model="state.address" label="Adresse" density="compact" :error-messages="v$.address.$errors.map(e => e.$message)" @input="v$.address.$touch"
      @blur="v$.address.$touch"></v-text-field>
        <v-text-field v-model="state.email"   label="E-mail" density="compact" :error-messages="v$.email.$errors.map(e => e.$message)"></v-text-field>
        <v-text-field v-model="state.phone"   label="Telefonnummer" density="compact" :error-messages="v$.phone.$errors.map(e => e.$message)"></v-text-field>
        <v-select v-model="state.team" clearable label="Gøglerhold" density="compact" :items="teams" item-title="title" item-value="slug" :error-messages="v$.team.$errors.map(e => e.$message)"></v-select>
        <v-select v-model="state.days" chips multiple label="Deltager" density="compact" :items="days" item-title="title" item-value="slug" :error-messages="v$.days.$errors.map(e => e.$message)"></v-select>
        <v-select v-model="state.transport" clearable label="Transport under Nathejk" density="compact" :items="transports" item-title="title" item-value="slug" :error-messages="v$.transport.$errors.map(e => e.$message)"></v-select>
        <v-select v-model="state.seatCount" clearable label="Hvis du medbringer egen bil, hvor mange ekstra pladser har du" density="compact" :items="seatCounts" item-title="title" item-value="slug" :error-messages="v$.seatCount.$errors.map(e => e.$message)"></v-select>
        <v-select v-model="state.info" clearable label="Må vi sende information på SMS" density="compact" :items="yesNo" item-title="title" item-value="slug" :error-messages="v$.info.$errors.map(e => e.$message)"></v-select>
        <v-select v-model="state.video" clearable label="Har du lyst til at medvirke i årets videoer" density="compact" :items="videos" item-title="title" item-value="slug" :error-messages="v$.video.$errors.map(e => e.$message)"></v-select>
        <p class="py-3">OBS. Ved tilmelding giver man samtidig også fototilladelse.</p>
        <v-btn class="me-4" type="submit">{{ buttonLabel }}</v-btn>
        </form>
      </v-col>
    </v-row>
  </v-container>

</template>
