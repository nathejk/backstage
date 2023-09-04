<template>

  <h1>Nathejk</h1>
    <v-container class="text-left">
      <v-row>
        <v-col cols="8" offset="2">
          <v-card title="Gøgler" subtitle="maks 100 pladser" variant="tonal">
              <template v-slot:title>
        Gøgler
                <v-chip class="ma-2" color="red" text-color="white">Tilmelding lukket</v-chip>
      </template>
            <v-card-text>

<p>Vær med til at skabe de fedeste oplevelser for spejderne på Nathejk - Nu kan du tilmelde dig som Gøgler på Nathejk. For at deltage som gøgler, skal du være fyldt 17 år (undtagelser kan aftales - send en mail til skuespil@nathejk.dk). Det koster 50 kr at deltage som gøgler på Nathejk.</p>

<p class="py-3">
Husk allerede nu at sætte kryds i kalenderen i weekenden 25.-27. august, hvor der er gøglerweekend. 
Husk også hjælpermødet d. 29. august.
      </p>

<a href="/Gøgler på Nathejk.pdf"><v-icon icon="mdi-file-pdf-box"></v-icon> Gøgler på Nathejk</a>

            </v-card-text>
            <v-card-actions>
              <v-btn variant="text" to="/tilmeld" :disabled="true">Gå til tilmelding</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
</template>
