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
      diet: null,
      tshirt: null,
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
      diet: { required },
      //tshirt: { required },
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
        { slug: 'open', title: 'Ingen holdning' },
        { slug: 'night', title: 'Vil gerne stå på checkpoints om natten' },
        { slug: 'day', title: 'Vil ikke stå på checkpoints om natten' },
        { slug: 'logistik', title: 'Logistik (kun efter aftale)' },
      ]),
      days: ref([
        { slug: "friday", title: 'Fredag' },
        { slug: "saturday", title: 'Lørdag' },
        { slug: "sunday", title: 'Søndag' },
      ]),
      transports: ref([
        { slug: 'owncar', title: 'Jeg medbringer egen bil' },
        { slug: 'needseat', title: 'Jeg har brug for plads i en bil' },
      ]),
      seatCounts: ref([
        { slug: '1', title: '1' },
        { slug: '2', title: '2' },
        { slug: '3', title: '3' },
        { slug: '4', title: '4' },
        { slug: '5', title: '5' },
        { slug: '0', title: 'Jeg har ikke ekstra plads' },
      ]),
      diet: ref([
        { slug: 'all', title: 'Ingen særlige præferencer' },
        { slug: 'vegetarian', title: 'Jeg vil gerne have vegetarmad' },
      ]),
      tshirt: ref([
        { slug: 'ingen', title: 'Ingen t-shirt' },
        { slug: 'xs', title: 'X-Small' },
        { slug: 's', title: 'Small' },
        { slug: 'm', title: 'Medium' },
        { slug: 'l', title: 'Large' },
        { slug: 'xl', title: 'X-Large' },
        { slug: '2xl', title: 'XX-Large' },
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
            diet:this.state.diet,
            //tshirt:this.state.tshirt,
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
      <ul class="text-left ms-5">
          <li>Du skal være fyldt 17 år for at deltage som gøgler på Nathejk.</li>
          <li>Det koster 50 kr., som betales med MobilePay ved tilmeldingen.</li>
          <li>Ved tilmeldingen giver du automatisk fototilladelse.</li>
          <li>Ved tilmeldingen giver du tilladelse til, at vi må kontakte dig på e-mail, telefon og sms i forbindelse med Nathejk.</li>
          <li>Efter tilmeldingen bliver du inviteret til gøglernes Facebook-gruppe, hvor vi planlægger og koordinerer vores indsats.</li>
          <li>Du skal kunne holde på en hemmelighed, da spejdere og banditter ikke må få kendskab til løbsområdet, indholdet på posterne eller gemmeting før tid.</li>
          <li>Du skal som udgangspunkt selv anskaffe dit kostume, så du kan udfylde din rolle på posterne. Du får mere information når vi nærmer os Nathejk samt masser af hjælp fra de andre gøglere.</li>
      </ul>
      <p class="text-left py-3">Det er ikke et krav at deltage i Gøglerweekenden og Hjælpermødet, men det er en god idé, hvis du vil lære de andre gøglere at kende og forstå, præcis hvad din rolle er på posterne. Det er selvfølgelig et krav at kunne komme til Nathejk, men hvis du ikke kan være med hele tiden fra fredag aften til søndag formiddag, kan du i tilmeldingen angive hvilke dage, du er med.</p>
      <p class="text-left py-3">Har du spørgsmål er du velkommen til at kontakte gøglerchefen Niels Jakob på njl@nathejk.dk</p>

        <form v-if="paying" @submit.prevent="pay">
          <p class="pb-3">Det koster 50,- kr at deltage som skal indbetales på MobilePay - betalinger registreres manuelt.</p>
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
        <v-select v-model="state.diet" clearable label="Har du specielle præferencer angående maden på Nathejk" density="compact" :items="diet" item-title="title" item-value="slug" :error-messages="v$.diet.$errors.map(e => e.$message)"></v-select>
        <v-select v-if="false" v-model="state.tshirt" clearable label="Ønsker du at købe en års t-shirt (DKK 175,-) udleveres på Nathejk" density="compact" :items="tshirt" item-title="title" item-value="slug" :error-messages="v$.tshirt.$errors.map(e => e.$message)"></v-select>
        <v-btn class="me-4" type="submit" :disabled="true">{{ buttonLabel }}</v-btn>
        </form>
      </v-col>
    </v-row>
  </v-container>

</template>
