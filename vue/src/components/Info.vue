<template>
  <div class="pa-2 mx-auto">
    <v-card elevation="2" rounded="lg" >
      <p class="text-capitalize font-weight-medium h2 pa-8 text-center"
            >
             Frequently Asked Questions
            </p> 

      <div class="pa-2 mx-auto font-weight-thin">
         <p
              
              class="font-weight-medium pa-2 text-left"
            >
              Question: What do the buttons mean?
            </p>

            <p
              
              class=" text-left"
            >Answer: </p> 
             <p>If an item has the icon "<v-icon small > mdi-check-all </v-icon>", it's price has been determined by experts. Thus, as a buyer you can be sure that the item holds a certain value. This price can not be changed ( Before posting the item to the experts, the seller has to decide what additional shipping costs they want to charge). If there are additional costs (for shipping), the "<v-icon small > mdi-package-variant-closed </v-icon>" icon is desplayed. The cheapest option is always displayed at the Shop section of the app. 

               If there are no additonal costs => <span >
            <v-chip
              class="ma-1 pl-0 caption"
              
              color="primary lighten-1"
              small
            >
              <v-chip label dark color="primary">
                ${{ "111111" }}</v-chip
              >
              <v-icon right> mdi-check-all </v-icon>
            </v-chip></span
          >
            . With additional costs: <span >
            <v-chip
              class="ma-1 pl-0 caption"
              
              color="primary lighten-1"
              small
            >
              <v-chip label dark color="primary">
                ${{ "111111" }}</v-chip
              >
              <v-icon right> mdi-check-all </v-icon>
              <v-icon small right> mdi-plus </v-icon
              ><v-icon small right> mdi-package-variant-closed </v-icon>
            </v-chip></span
          >.


             As an incentive, buyers of 'estimated' items receive a ±5% TPP cashback reward. Indicated by a <v-icon> mdi-hand </v-icon>  with a <v-icon> mdi-heart </v-icon>. Like so: <v-chip
            class="ma-1 pl-0 caption"
            
            dark
            color="green lighten-2"
            small
          ><v-chip dark label color="green">
            ${{ ("111111"* 0.05).toFixed(0)
            }}</v-chip>
            <v-icon small right> mdi-hand-heart </v-icon>
          </v-chip>
            </p>
            <p>If an item has the icon "
              <v-icon small > mdi-repeat </v-icon>", it has been reposted. It's price has been deterimined previously by experts and the item was bought by a buyer before. The seller is not the creator of the item. The value of the item can have appreciated or depreciated in the meantime, therefore it is up to the seller to give a discount "<v-icon small > mdi-label-percent</v-icon>" or add cost "<v-icon small  > mdi-package-variant-closed </v-icon>". This can also be a combination of the two, for example:
            </p><span class="pa-2 font-weight-regular text-left" >
   <v-chip small dark color="primary" class="pr-0" 
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2 " v-if="hover" > Buy Now </span><span class="pr-2" v-else>$11111</span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-2"
            
          ><v-icon small  right> mdi-repeat </v-icon
            ><v-icon small  right> mdi-plus </v-icon
            ><v-icon small  right> mdi-package-variant-closed </v-icon>
            <v-icon small right> mdi-minus </v-icon
            ><v-icon small right> mdi-label-percent</v-icon>
          </v-chip>
          </v-chip
            >

          </span>
            
            
            
             </div>
           
<v-divider class="ma-4" />
           
   <div class="pa-2 mx-auto font-weight-thin">
         <p
              
              class="font-weight-medium pa-2 text-left"
            >
              Question: How do buy an item?
            </p><p
              
              class="pa-2 font-weight-light text-left"
            >Answer: </p>  <p> Provide prepayment using the buttons at the detailed page of an item. If you provided prepayment for an item that has shipping cost, the item seller will complete the item transfer (once they indicate the item has shipped). Else, If you choose to provide prepayment for a 'pickup' item, you complete the item transfer (under Account) when you are meeting up with the seller or at anytime to transfer the digital ownership of the item.
<v-divider class="ma-4" />
                  <p
              
              class="font-weight-medium pa-2 text-left"
            >
              Question: Does selling cost anything?
            </p><p
              
              class="pa-2 font-weight-light text-left"
            >Answer: </p> 
            <p> Selling any item is free at any time. No costs never nada. </p>
<v-divider class="ma-4" />
            <p
              
              class="font-weight-medium pa-2 text-left"
            >
              Question: How do I contact a buyer/seller?
            </p>
<p
              
              class="pa-2 font-weight-light text-left"
            >Answer: </p> 
            <p> Send a message to the seller at the details page of the item. If you are a seller, message the buyer at the item actions under Account. Alternatively, in messages you can add any address signed up to the marketplace.
<v-divider class="ma-4" />
  <p
              
              class="font-weight-medium pa-2 text-left"
            >
              Question: What are the transfer options?
            </p> 

            <p
              
              class="pa-2 font-weight-light text-left"
            >Answer: </p> <p> There are 2 carefully designed options for you to choose, Pickup and Shipping. These processes differ in terms of responsibility. For pickup, the buyer is fully in control of the prepayment up until the (physical) transaction, meaning that the buyer controls when the payment is payed out to the seller. For shipping, the seller will have control of the payment. When the seller indicates the item has been shipped, the tokens will be payed out to the seller. Up until that, the buyer may withdrawl at any moment.
       </p><p> You can access to these functions at the item actions under Account.</p>
            </div>
          <v-row>
           
       

        
        </v-row>
    </v-card>
      <v-card-actions>
       
    
        <div>
         
           

        </div>

     
     
      </v-card-actions>

  

  </div>
</template>

<script>
import { databaseRef } from "./firebase/db";


export default {
 /// props: ["itemid"],
  //components: { ItemListBuy, WalletCoins },
  data() {
    return {
      //itemid: this.item.id,
      //make sure deposit is number+token before sending tx
      amount: "",
      flight: false,
      flightLP: false,
      flightSP: false,
      showinfo: false,
      imageurl: "",
      loadingitem: true,
      photos: [],
    };
  },

  mounted() {
    this.loadingitem = true;
    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null && data.photo != null) {
        //console.log(data.photo);
        this.imageurl = data.photo;
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },
  computed: {
    thisitem() {
      return this.$store.getters.getItemByID(this.itemid);
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
    commentlist() {
      return this.thisitem.comments.filter((i) => i != "") || [];
    },
  },

  methods: {
    getItemPhotos() {
      if (this.showinfo && this.imageurl != "") {
        this.loadingitem = true;
        const id = this.itemid;

        const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
        imageRef.on("value", (snapshot) => {
          const data = snapshot.val();
          if (data != null && data.photo != null) {
            this.photos = data;
            this.loadingitem = false;
          }
        });
        this.loadingitem = false;
      }
    },
    clickOutside() {
      if ((this.showinfo = true)) {
        this.showinfo = false;
      }
    },
  },
};
</script>



<!---
shows item id from buy list
<div id="item-list-buy">
      {{ itemid }}
    </div>
    ---->