<template>
  <div>
    <div>
      <div>
        <div class="item">
          <v-row>
            <v-col>
              <v-card elevation="0">
                <div class="overline">Title</div>

                <div class="body-1 mt-1">{{ thisitem.title }}</div>
              </v-card>
            </v-col>
            <v-col>
              <v-card elevation="0">
                <v-chip-group>
                  <v-chip
                    outlined
                    small
                    class="caption mt-1"
                    v-for="previewtag in thisitem.tags"
                    :key="previewtag"
                  >
                    <v-icon small left> mdi-tag-outline </v-icon>
                    {{ previewtag }}
                  </v-chip>
                  <!--<v-chip 
        class="ma-1" outlined small
      >
<v-rating
  background-color="grey"
  color="black"
  dense
 
  readonly
  length="5"
  size="10"
  :value="thisitem.condition"
></v-rating> </v-chip>-->
                </v-chip-group>
              </v-card>
            </v-col>
          </v-row>
          <v-card elevation="0">
            <div class="overline">Description</div>
        
              <span class="caption mt-1">{{ thisitem.description }}</span>
       
          </v-card>

          <v-divider class="ma-4"></v-divider>
          <div>
       <!--   <v-chip class="ma-1 caption" medium label outlined>
            <v-icon left> mdi-account-outline </v-icon>
            Your Address: {{ thisitem.creator }}
          </v-chip>-->

          <v-chip class="ma-1 caption" label outlined medium>
            <v-icon left> mdi-account-badge-outline </v-icon>
            TPP ID: {{ thisitem.id }}
          </v-chip>

          <v-chip
            v-if="thisitem.shippingcost"
            class="ma-1 caption"
            label
            outlined
            medium
          >
            <v-icon left> mdi-package-variant </v-icon>
            Shipping
          </v-chip>
          <v-chip
            v-if="thisitem.localpickup != ''"
            class="ma-1 caption"
            label
            outlined
            medium
          >
            <v-icon left> mdi-map-marker-outline </v-icon>
            Pickup
          </v-chip>
          <v-chip
            v-if="thisitem.shippingcost > 0"
            class="ma-1 caption"
            label
            outlined
            medium
          >
            <v-icon left> mdi-package-variant-closed </v-icon>
            Shipping cost: {{ thisitem.shippingcost }}<v-icon small right>$vuetify.icons.custom</v-icon>  
          </v-chip>

          <v-chip 
            outlined
            medium
            label
            class="ma-1 caption"
            v-for="country in thisitem.shippingregion"
            :key="country"
          >
            <v-icon small left> mdi-flag-variant-outline </v-icon
            >{{ country }}</v-chip
          >
          <v-chip class="ma-1 caption" label outlined medium>
            <v-icon left> mdi-star-outline </v-icon>
            Condition: {{ thisitem.condition }}/5
          </v-chip>

       
          </div>
          <v-divider class="ma-4 "></v-divider>

          <div class="mt-2 text-center">
            <p  class="font-weight-medium headline"> TPP ID: {{thisitem.id}}  </p><p class="caption"> Tip: Show TPP ID: {{thisitem.id}} on your photos. This creates trust to estimators and buyers, thereby making the item more valueable.</p>
            <v-btn block large outlined @click="addPhoto()" color="primary">
             <span v-if="!img[i]"> <v-icon large left> mdi-plus </v-icon>Add Photo</span><span v-else> <v-icon large left> mdi-refresh </v-icon> Change photo</span>
            </v-btn>
          
          </div>
          <div v-if="img[0] != null"> 
            <div v-for="(image, index) in img " :key="index"> <p> {{uploadValue[index]}}</p>
           <input
              type="file"
              :ref="'input' + index"
              style="display: none"
              @change="previewImage"
              accept="image/*"
            />
            <v-card class="text-center mt-4 elevation-4">
              <v-card-title v-if="img[index] == img[0]">Primary photo</v-card-title>   <v-card-title v-else>Photo {{index + 1}}</v-card-title>
              <v-img class="rounded contain" :src="img[index]" />
  <v-progress-linear
      v-model="uploadValue[index]"
    
    ></v-progress-linear>
              <br />
            </v-card>  <span v-if="img[index]" class="pa-4"><v-btn  block  outlined @click="replacePhoto(index)" color="primary">
            <v-icon  left> mdi-refresh </v-icon> Change photo
            </v-btn></span><span  v-if="img.length - 1 == index " class="pa-4"><v-btn  block  outlined @click="addPhoto()" color="primary">
           <v-icon  left> mdi-plus </v-icon>Add Photo
            </v-btn></span>
            </div>
          </div>
        
      <div class="pt-4 text-right">
        <v-btn
          :disabled="!this.imageData[1]|| !hasAddress"
          color="primary"
          @click="create()"
        >
          Place {{thisitem.title}} <v-icon> mdi-arrow-right-bold</v-icon>
          
        </v-btn>
      </div>
    </div>
  </div>  </div>  </div>
</template>

<script>

import CreateItemForm from "./CreateItemForm.vue";
import { fb, databaseRef } from "./firebase/db";

export default {
  props: ["thisitem"],
  components: { CreateItemForm },
  data() {
    return {
     
      imageData: [],
     i: 0,
      img: [],
   
      //thisitem: {},
      //itemid: "",

      uploadValue: []
  

    };
  },

  computed: {
    hasAddress() {
      return !!this.$store.state.account.address;
    },
   
    /*thisitem() {
      this.itemid = this.$store.state.newitemID;
      return this.$store.getters.getItemByID(this.itemid);
    },*/
  },

  methods: {
    click1(i) {
      this.$refs.input[i].click();
    },

    create() {

      let uploadDate = fb.database.ServerValue.TIMESTAMP
      const post = { photos: this.img
        //_id: this.$store.state.user.uid,
        //itemid: this.thisitem.id,
    , id: { username: this.thisitem.creator, _id: this.$store.state.user.uid, uploadDate: uploadDate }}
       /*databaseRef
        .ref("ItemPhotoGallery/0").set(post) .then((response) => {
          console.log(response);
        })
        .catch((err) => {
          console.log(err);
        });*/

      databaseRef
        .ref("ItemPhotoGallery/" + this.thisitem.id)

        .set(post)
        .then((response) => {
          console.log(response);
        })
        .catch((err) => {
          console.log(err);
        });
      this.$emit("changeStep", "3");
      //this.updateStep();
    },

  
    async previewImage(event) {
       console.log("PREVIEW " + this.i)
      this.uploadValue[this.i] = 0;
      this.img[this.i] = null;

      this.imageData[this.i] = event.target.files[0];
      await this.onUpload();


    },

    addPhoto(){
    
      let refff = this.$refs['input' + this.i]

        console.log(this.$refs)
        console.log(refff)
        refff.click();
     
    },

    newPhoto(){
      

        if(this.img.length > 5){
        this.i = this.img[this.img.length - 1]
       
      }else{
     this.i = this.i + 1
      }
     console.log("NEW INDEX " + this.i)
    },

    replacePhoto(index){
  
        this.i = index
       
   
     console.log("REPLACE INDEX " + this.i)
     this.addPhoto()
    },



    async onUpload(i) {
      this.img[this.i] = null;
      let storageRef = fb
        .storage()
        .ref(`${this.imageData[this.i].name}`)
        .put(this.imageData[this.i]);
      storageRef.on(
        `state_changed`,
        (snapshot) => {
          this.uploadValue[this.i] =
            (snapshot.bytesTransferred / snapshot.totalBytes) * 100;
        },
        (error) => {
          console.log(error.message);
        },
        () => {
          this.uploadValue[this.i] = 100;
          storageRef.snapshot.ref.getDownloadURL().then((url) => {
            this.img[this.i] = url;
  
            console.log("NEW UPLOAD " + this.img[this.i])
            console.log(this.imageData[this.i]);
               this.newPhoto()
          });
        }
      );
      
       
    },
   
  }
}
</script>
