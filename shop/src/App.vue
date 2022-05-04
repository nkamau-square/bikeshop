<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <h1>Ngumba's Bike Shop</h1>
      </div>
    </v-app-bar>

    <v-main>
      <InventoryPage/>
    </v-main>
  </v-app>
</template>

<script>
import InventoryPage from './components/InventoryPage';

export default {
  name: 'App',

  components: {
    InventoryPage,
  },

  data: () => ({
    //
  }),
};
</script>
