<template>
  <v-container>
    <v-row>
      <v-col :key="1"></v-col>
      <v-col :key="2" :cols="12">
        <v-card>
          <v-card-title>
            Bikes
            <v-spacer></v-spacer>
            <v-text-field
              v-model="search"
              append-icon="mdi-magnify"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
          </v-card-title>
          <v-data-table :headers="headers" :items="inventory" :search="search">
            <template slot="item.state" slot-scope="props">
              <v-icon v-if="props.item.state == 'IN_STOCK'" color="green" large
                >mdi-checkbox-marked-circle</v-icon
              >
              <v-icon v-else color="red" large>mdi-minus-circle</v-icon>
            </template>
            <template slot="item.purchase" slot-scope="props">
              <PurchasePage :item="props.item" @purchase="purchase" />
            </template>
          </v-data-table>
        </v-card>
      </v-col>
      <v-col :key="3"></v-col>
    </v-row>
  </v-container>
</template>

<script>
import PurchasePage from "./PurchasePage";
export default {
  name: "InventoryPage",
  props: {
    msg: String,
  },
  data() {
    return {
      catalogue: null,
      inventory: [],
      search: "",
      headers: [
        {
          text: "Name",
          align: "start",
          sortable: true,
          value: "name",
        },
        { text: "Quantity", value: "quantity" },
        { text: "In Stock", value: "state" },
        { text: "Price CAD", value: "price" },
        { text: "", value: "purchase" },
      ],
    };
  },
  components: {
    PurchasePage,
  },
  mounted() {
    this.$http
      .get("http://localhost:8081/v1/catalogue")
      .then((result) => {
        this.catalogue = result.data.objects;
        var idList = this.catalogue.map((object) => object.id);
        this.$http
          .post("http://localhost:8081/v1/inventory", idList)
          .then((result) => {
            this.inventory = result.data.counts.map((item) => {
              var catalogue = this.catalogue.find((obj) => {
                return obj.id === item.catalog_object_id;
              });
              return {
                id: catalogue.id,
                state: item.quantity>0?item.state:"out",
                quantity: item.quantity,
                name: catalogue.item_variation_data.name,
                price: catalogue.item_variation_data.price_money.amount,
              };
            });
          });
      })
      .catch((error) => {
        console.log(error.response);
      });
  },
  methods: {
    purchase(quantity, item) {
      const req = {
        quantity: quantity.toString(),
        id: item.id
      }
      this.$http
        .post("http://localhost:8081/v1/purchase",req)
        .then(() => {
          item.quantity = item.quantity - quantity;
          if (item.quantity <= 0) {
            item.quantity = 0;
            item.state = "out";
          }
        })
        .catch((error) => {
          console.log(error.response);
        });
    },
  },
};
</script>

<style scoped>
</style>
