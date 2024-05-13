<template>

  <h1>Nathejk</h1>
    <v-container class="text-left">
      <v-row>
        <v-col cols="8" offset="2">
          <v-card title="Gøgler" subtitle="maks 100 pladser" variant="tonal">
              <template v-slot:title>
        Gøgler
                <!-- v-chip class="ma-2" color="red" text-color="white">Tilmelding lukket</v-chip -->
      </template>
            <v-card-text>

<p>Vær med til at skabe de fedeste oplevelser for spejderne på Nathejk - Nu kan du tilmelde dig som Gøgler på Nathejk. For at deltage som gøgler, skal du være fyldt 17 år. Det koster 50 kr at deltage som gøgler på Nathejk.</p>

<p class="py-3">
<strong>Vær opmærksom på følgende datoer:</strong>
<ul class="ms-5">
    <li>16.-18. august: Gøglerweekend i DDS Bjældehytten (Bjældevej 25, Roskilde), hvor vi bygger rekvisitter, løser andre opgaver og hygger.</li>
    <li>28. august kl. 18: Hjælpermøde i Glostrup Bibliotek, hvor vi laver lidt gøgl for at sætte stemningen for de andre hjælpere og banditter, og
gøglerne får en briefing omkring posterne.</li>
    <li>20.-22. september: Nathejk et sted på Sjælland!</li>
</ul>
      </p>

<a href="https://sky.nathejk.dk/api/public/dl/kPxSoH1M/G%C3%B8gler/G%C3%B8gler%20pa%CC%8A%20Nathejk.pdf"><v-icon icon="mdi-file-pdf-box"></v-icon> Gøgler på Nathejk</a>

            </v-card-text>
            <v-card-actions>
              <v-btn variant="text" to="/tilmeld" :disabled="false">Gå til tilmelding</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
</template>
