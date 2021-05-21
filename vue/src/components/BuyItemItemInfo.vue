<template>
  <div class="mx-auto">
    <v-card  rounded="xl" :raised="showinfo"  v-click-outside="clickOutside ">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>
<div class="pa-2 mx-auto" >
      <router-link style="text-decoration: none; color: inherit;"  :to="{ name: 'BuyItemDetails', params: { id: itemid } }" >
        <v-row >
          <v-col class="pa-2" cols="12" md="7">
            <p :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              v-if="thisitem.creator != thisitem.seller"
              class="text-capitalize subtitle-2 pa-2 text-left"
            >
              <v-icon left> mdi-repeat </v-icon>{{ thisitem.title }}
            </p>

            <p v-else class="text-capitalize subtitle-2 pa-2 text-left">
              <v-icon color="primary" left>mdi-check-all </v-icon>
              {{ thisitem.title }}
            </p>

            <span class="ma-1">
              <p class="ma-1 caption font-weight-light" v-if="thisitem.description.length < 80">
                {{ thisitem.description }}
              </p>
              <p class="ma-1 caption font-weight-light" v-else>
                {{ thisitem.description.substring(0, 78) + ".." }}
              </p>
            </span>
          </v-col>

          <v-col class="pa-2" cols="12" md="5">
            <div v-if="imageurl">
              <v-img class="rounded-lg contain ma-2 mb-0" :aspect-ratio="4/3" :src="imageurl"></v-img>
            </div>
          </v-col>
        </v-row>
      </router-link>     </div>
      <v-card-actions>
           <v-btn plain class="mr-2" icon @click="(showinfo = !showinfo), getItemPhotos()">
          <v-icon>{{
            showinfo ? "mdi-chevron-up" : "mdi-chevron-down"
          }}</v-icon>
        </v-btn>
        <div>
          <!--<v-btn class="rounded-pill ml-6 mr-6"  small
             color="primary"
            :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
            plain
          >
            Details
          </v-btn>-->
        </div>
        <div v-if="thisitem.creator != thisitem.seller">
       
          
           <span><router-link style="text-decoration: none; color: inherit;" :to="{ name: 'BuyItemDetails', params: { id: itemid } }">
   <v-chip style="cursor: pointer;" class="mr-2 pr-0"  v-if="
               thisitem.shippingcost > 0 &&
              thisitem.localpickup == '' &&
              thisitem.discount == 0
            " small dark color="primary lighten-1"  
              >
              <v-hover v-slot="{ hover }" close-delay="100" open-delay="30" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else>{{
                 Number(thisitem.estimationprice) + Number(thisitem.shippingcost)
              }}<v-icon small right>$vuetify.icons.custom</v-icon> </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary"
            
          ><v-icon small left> mdi-repeat </v-icon>
            <v-icon small left> mdi-plus </v-icon
            ><v-icon small left> mdi-package-variant-closed </v-icon>
          </v-chip>
          </v-chip
            ></router-link>
          </span>

         


 <span><router-link style="text-decoration: none; color: inherit;" :to="{ name: 'BuyItemDetails', params: { id: itemid } }">
   <v-chip  v-if="
              thisitem.shippingcost > 0 &&
              thisitem.localpickup == '' &&
              thisitem.discount > 0
            " small dark color="primary lighten-1"  style="cursor: pointer;" class="mr-2 pr-0" 
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else>{{
                Number(thisitem.estimationprice) +
                Number(thisitem.shippingcost) -
                Number(thisitem.discount)
              }}<v-icon small right>$vuetify.icons.custom</v-icon> </span>
              </span>
</v-hover>


                  <v-chip label
           
            class="pl-0 caption"
            color="primary"
            
          ><v-icon small  right> mdi-repeat </v-icon
            ><v-icon small  right> mdi-plus </v-icon
            ><v-icon small  right> mdi-package-variant-closed </v-icon>
            <v-icon small right> mdi-plus </v-icon
            ><v-icon small right> mdi-brightness-percent</v-icon>
          </v-chip>
          </v-chip
            > </router-link>
          </span>


          <span><router-link style="text-decoration: none; color: inherit;" :to="{ name: 'BuyItemDetails', params: { id: itemid } }">
   <v-chip  v-if="thisitem.discount > 0 && thisitem.localpickup != ''"
             small dark color="primary lighten-1"  style="cursor: pointer;" class="mr-2 pr-0"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else>{{
                thisitem.estimationprice - thisitem.discount
              }}<v-icon small right>$vuetify.icons.custom</v-icon> </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary "
            
          ><v-icon small right> mdi-repeat </v-icon>
            <v-icon small right> mdi-plus </v-icon
            ><v-icon small right> mdi-brightness-percent</v-icon>
          </v-chip>
          </v-chip
            ></router-link>
          </span>

        </div>
        <div v-else>

          
         <span v-if="thisitem.localpickup == ''">
            <!--<v-chip class="ma-1 caption"  color="primary lighten-1" small>
              <v-chip label dark color="primary">
                <v-icon small right>$vuetify.icons.custom</v-icon>{{
                  Number(thisitem.estimationprice) +
                  Number(thisitem.shippingcost)
                }}</v-chip
              >
              <v-icon right> mdi-check-all </v-icon
              ><v-icon small right> mdi-plus </v-icon
              ><v-icon small right> mdi-package-variant-closed </v-icon>
            </v-chip>-->
            <router-link
              style="text-decoration: none; color: inherit"
              :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
            >
              <v-chip style="cursor: pointer;" class="mr-2 pr-0" small dark color="primary" >
                <v-hover v-slot="{ hover }" close-delay="300" open-delay="60">
                  <span>
                    <span class="pr-2" v-if="hover"> Buy Now </span
                    ><span class="pr-2" v-else
                      >{{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost)
                      }}<v-icon small right>$vuetify.icons.custom</v-icon
                      >
                    </span>
                  </span>
                </v-hover>


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-1"
            
          > <v-icon right> mdi-check-all </v-icon
              ><v-icon small right> mdi-plus </v-icon
              ><v-icon small right> mdi-package-variant-closed </v-icon>
          </v-chip>
          </v-chip
            >
     </router-link
            >
            
            </span
          >

          <span v-else>
       <router-link style="text-decoration: none; color: inherit;" :to="{ name: 'BuyItemDetails', params: { id: itemid } }">
            <v-chip style="cursor: pointer;" class="mr-2 pr-0" 
             small dark color="primary lighten-1" 
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-3 caption" v-else>{{
                thisitem.estimationprice
              }}<v-icon small right>$vuetify.icons.custom</v-icon></span>
              </span>
</v-hover>


                  <v-chip label
           
            class="pl-0 caption"
            color="primary"
            
          > <v-icon right> mdi-check-all </v-icon>
          
          </v-chip>
          </v-chip
            > </router-link>
            
            </span
          >
  
        </div> <v-chip color="primary lighten-2" small v-if="thisitem.discount > 0" class="mx-2  font-weight-medium" >
                
                  {{ Math.floor(thisitem.discount/thisitem.estimationprice * 100)}}%   <v-icon small right> mdi-brightness-percent </v-icon>
                
                </v-chip>

                <v-chip  color="primary" small v-if="thisitem.condition " class="mx-2 font-weight-medium" >
               {{thisitem.condition}} <!--<v-rating
                            :value="Number(thisitem.condition)"
                            readonly
                            color="white"
                            background-color="primary lighten-1"
                            x-small
                            dense
                          ></v-rating> --><v-icon small right>mdi-star</v-icon>
                </v-chip>

        <v-spacer></v-spacer>
       
      </v-card-actions>

      <v-expand-transition >
        <div>
          <div class="pa-2 mx-auto" elevation="8" v-if="showinfo">
            <div>
              <div v-if="photos[0]">
                <v-divider></v-divider>
                <v-carousel
                  cycle
                  height="400"
                  hide-delimiter-background
                  show-arrows-on-hover
                >
                  <v-carousel-item
                    v-for="(photo, i) in photos"
                    :key="i"
                    :src="photo"
                  >
                  </v-carousel-item>
                </v-carousel>
              </div>

              <span>
                 <div class="pa-2 overline text-center">Description</div>
                <v-card-text>
                  <div class="caption" >{{ thisitem.description }}</div>
                </v-card-text>
              </span>
              <v-chip :to="{ name: 'SearchTag', params: { tag: itemtags } }"
                outlined
                medium
                
                class="ma-1 caption font-weight-light text-capitalize"
                v-for="itemtags in thisitem.tags"
                :key="itemtags"
              >
                <v-icon small left> mdi-tag-outline </v-icon
                >{{ itemtags }}</v-chip
              >
              <v-chip :to="{ name: 'SearchRegion', params: { region: selected } }"
                outlined
                
                
                class="ma-1 caption font-weight-light "
                v-for="selected in thisitem.shippingregion"
                :key="selected"
              >
                <v-icon small left> mdi-flag-variant-outline </v-icon
                >{{ selected.toUpperCase() }}</v-chip
              >
              <v-chip class="ma-1 caption font-weight-light "   :to="{ name: 'BuyItemDetails', params: { id: itemid } }" outlined >
                <v-icon left> mdi-account-badge-outline </v-icon>
                TPP ID: {{ thisitem.id }}
              </v-chip>

              <v-chip class="ma-1 caption font-weight-light "  outlined >
                <v-icon small left> mdi-star </v-icon>
                Condition: {{ thisitem.condition }}/5
              </v-chip>

              <v-chip
                v-if="thisitem.localpickup != ''"
                class="ma-1 caption font-weight-light "
                
                outlined
                
                ><v-icon left> mdi-map-marker-outline </v-icon>
                Pickup available</v-chip
              >

              <v-chip
                v-if="thisitem.shippingcost > 0"
                class="ma-1 caption font-weight-light "
                
                outlined
                
              >
                <v-icon left> mdi-package-variant-closed </v-icon>
                Shipping: {{ thisitem.shippingcost}} <v-icon right small>$vuetify.icons.custom</v-icon> 
              </v-chip>

              <v-chip
                v-if="thisitem.estimationprice > 0"
                class="ma-1 caption font-weight-light "
                
                outlined
                
              >
                <v-icon left> mdi-check-all </v-icon>
                Price: {{ thisitem.estimationprice}} <v-icon right small>$vuetify.icons.custom</v-icon> 
              </v-chip>
<v-btn class="rounded-pill ml-6 mr-6"  small
             color="primary"
            :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
            plain
          >
            Details
          </v-btn>
            </div>
          </div>
        </div>
      </v-expand-transition>
    </v-card>
  </div>
</template>

<script>
import { databaseRef } from "./firebase/db";
import ItemListBuy from "./ItemListBuy.vue";



export default {
  props: ["itemid"],
  components: { ItemListBuy },
  data() {
    return {
      //itemid: this.item.id,
      //make sure deposit is number+token before sending tx

      showinfo: false,
      imageurl: "",
     
      photos: [],
    };
  },

   beforeCreate(){
    this.loadingitem = true;
  },

  mounted() {

    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null ) {
        //console.log(data[0]);
        this.imageurl = data[0];
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
          if (data != null ) {
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

