<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="sign.png"
          transition="scale-transition"
          width="40"
        />

        <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="samar.png"
          width="100"
        />
      </div>

      <v-spacer></v-spacer>
      <v-switch v-model="$vuetify.theme.dark" label="Dark theme" style="margin-top: 25px"></v-switch>
      <v-spacer></v-spacer>
 
      <v-btn href="https://github.com/dkt64/hires-vector" target="_blank" text>
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <Main />
    </v-content>
  </v-app>
</template>

<script>
import Main from "./components/Main";

export default {
  name: "App",

  components: {
    Main
  },

  data: () => ({
    //
  })
};
</script>
