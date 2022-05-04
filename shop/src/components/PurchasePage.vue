<template>
  <div class="text-center">
    <v-dialog v-model="dialog" width="500">
      <template v-slot:activator="{ on, attrs }">
        <v-btn color="success" :disabled="item.quantity<=0" dark v-bind="attrs" v-on="on"> Buy </v-btn>
      </template>

      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          {{ item.name }}
        </v-card-title>

        <v-card-text>
          <v-col>
            <img
              src="https://www.bmc-switzerland.com/media/catalog/product/b/m/bmc-22-10503-009-bmc-fourstroke-01-three-mountain-bike-black-01.png"
            />
            <v-container>
              <v-slider
                v-model="purchaseQuantity"
                min="1"
                :max="item.quantity"
                thumb-label
                Label="Quantity"
              ></v-slider>
              <v-row>
                <v-spacer></v-spacer>
                <h3>Quantity: {{ purchaseQuantity }}</h3>
                <v-spacer></v-spacer>
              </v-row>
            </v-container>
          </v-col>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-btn color="red" text @click="dialog = false"> Cancel </v-btn>
          <v-spacer></v-spacer>
          <h2>Total: C${{ purchaseQuantity * item.price }}</h2>
          <v-spacer></v-spacer>
          <v-btn color="success" text @click="purchase"> Purchase </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  name: "PurchasePage",
  data() {
    return {
      dialog: false,
      purchaseQuantity: 1,
    };
  },
  props: {
    item: Object,
  },
  methods: {
    purchase() {
        this.$emit('purchase', this.purchaseQuantity, this.item)
        this.dialog = false
    },
  },
};
</script>

<style scoped>
img {
  border-radius: 8px;
  max-width: 80%;
  max-height: 80%;
  margin-left: auto;
  margin-right: auto;
  display: block;
}

.theme--light.v-btn.v-btn--disabled:not(.v-btn--flat):not(.v-btn--text):not(.v-btn-outlined) {
  color: fuchsia !important;
}
</style>