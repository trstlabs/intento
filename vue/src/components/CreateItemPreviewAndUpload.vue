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
            <span v-if="!img[0]"><v-btn  block large outlined @click="addPhoto(i)" color="primary">
            <v-icon large left> mdi-plus </v-icon>Add Photo
            </v-btn> 
       </span>
          </div>
          <input
              type="file"
              :ref="'input'+ i"
             style="display: none"
              @change="previewImage"
              accept="image/*"
            />
          <div>  
             <v-progress-circular class="ma-2"
      v-model="uploadValue"
    v-if="uploadValue != 0 && uploadValue != 100"
    ></v-progress-circular>
            <div v-for="(image, index) in imageData" :key="index">
           <input
              type="file"
              :ref="'input'+ index"
            style="display: none"
              @change="previewImage"
              accept="image/*"
            />
            <v-card class="text-center mt-4 elevation-4">
              <v-card-title v-if="img[index] == img[0]">Primary photo</v-card-title>   <v-card-title v-else>Photo {{index + 1}}</v-card-title>
              <v-img class="rounded contain" :src="img[index]"> <template v-slot:placeholder>
        <v-row
          class="fill-height ma-0"
          align="center"
          justify="center"
        >
          <v-progress-circular
            indeterminate
            color="grey lighten-5"
          ></v-progress-circular>
        </v-row>
      </template> </v-img>
  
              <br />
            </v-card>  <span v-if="img[index]" class="pa-4"><v-btn  block  outlined @click="replacePhoto(index)" color="primary">
            <v-icon  left> mdi-refresh </v-icon> Change photo
            </v-btn></span><span  v-if="img.length - 1 == index && index < 11 " class="pa-4"><v-btn  block  outlined @click="addPhoto()" color="primary">
           <v-icon  left> mdi-plus </v-icon>Add Photo
            </v-btn></span>
            </div>
          </div>
        
      <div class="pt-4 text-right">
        <v-btn rounded
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
import AppText from "./AppText.vue";
import CreateItemForm from "./CreateItemForm.vue";
import { fb, databaseRef } from "./firebase/db";

export default {
  props: ["thisitem"],
  components: { AppText, CreateItemForm },
  data() {
    return {
     
      imageData: [],
     i: 0,
      img: [],
   
      //thisitem: {},
      //itemid: "",

      uploadValue: 0,
  

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
      this.uploadValue = 0;
      this.img[this.i] = null;

      this.imageData[this.i] = event.target.files[0];
    this.onUpload();


    },

    addPhoto(){
      let refff = this.$refs['input' + this.i]

        console.log(this.$refs)
        console.log(refff)
        refff.click();
     
    },

    newPhoto(){
      

        if(this.img.length > 10){
        this.i = this.img[this.img.length - 1]
       alert("Maximum amount reached")
      }else{
         console.log("LENG "+this.img.length)
      this.i = this.img.length
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
          this.uploadValue = Math.ceil(
            (snapshot.bytesTransferred / snapshot.totalBytes) * 100);
            console.log(this.uploadValue)
        },
        (error) => {
          console.log(error.message);
        },
        () => {
          //this.uploadValue[this.i] = 100;
          storageRef.snapshot.ref.getDownloadURL().then((url) => {
            this.img[this.i] = url;
  
            console.log("NEW UPLOAD " + this.img[this.i])
            console.log(this.imageData[this.i]);
               this.newPhoto()
          });
        }
      );
      

       
    },

process() {
  const file = this.imageData[this.i]

  if (!file) return;

  const reader = new FileReader();

  reader.readAsDataURL(file);

  reader.onload = function (event) {
    const imgElement = createElement("img");
    imgElement.src = event.target.result;
    //const refff = this.$refs['input' + this.i]
    //console.log(refff)
    //refff.src = event.target.result;

    imgElement.onload = function (e) {
      const canvas = createElement("canvas");
      const MAX_WIDTH = 400;

      const scaleSize = MAX_WIDTH / e.target.width;
      canvas.width = MAX_WIDTH;
      canvas.height = e.target.height * scaleSize;

      const ctx = canvas.getContext("2d");

      ctx.drawImage(e.target, 0, 0, canvas.width, canvas.height);

      this.imageData[this.i] = ctx.canvas.toDataURL(e.target, "image/jpeg");
  this.onUpload();
      // you can send srcEncoded to the server
      //document.querySelector("#output").src = srcEncoded;
    };
  };
}    
   
  }
}
</script>
