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
          <v-chip class="ma-1 caption" medium label outlined>
            <v-icon left> mdi-account-outline </v-icon>
            Your Address: {{ thisitem.creator }}
          </v-chip>

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
            Local pickup
          </v-chip>
          <v-chip
            v-if="thisitem.shippingcost != ''"
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
            <v-btn block large outlined @click="click1" color="primary">
             <span v-if="img1 == null"> <v-icon large left> mdi-plus </v-icon>Add Photos</span><span v-else> <v-icon large left> mdi-refresh </v-icon> Change primary photo</span>
            </v-btn>
            <input
              type="file"
              ref="input1"
              style="display: none"
              @change="previewImage"
              accept="image/*"
            />
          </div>
          <div v-if="img1 != null">
            <v-card class="text-center mt-4 elevation-4">
              <v-card-title>Primary photo</v-card-title>
              <v-img class="rounded contain" :src="img1" />
  <v-progress-linear
      v-model="uploadValue"
    
    ></v-progress-linear>
              <br />
            </v-card>
          </div>
          <div v-if="img1 != null" class="mt-2" >
            <v-btn outlined @click="click2" color="primary">
             <span v-if="img2"> <v-icon large left> mdi-refresh </v-icon> Change photo 2 </span> <span v-else> <v-icon large left> mdi-plus </v-icon>  Additonal</span>
            </v-btn>
            <input
              type="file"
              ref="input2"
              style="display: none"
              @change="previewImage2"
              accept="image/*"
            />
          </div>
          <div v-if="img2 != null">
            <v-card class="text-center mt-4">
             <v-progress-linear
      v-model="uploadValue2"
    
    ></v-progress-linear>

              <v-img class="rounded contain" :src="img2" />

              <br />
            </v-card>
          </div>
          <div class="mt-2">
            <v-btn outlined v-if="imageData2" @click="click3" color="primary">
               <span v-if="img3"><v-icon large left> mdi-refresh </v-icon> Change photo 3 </span> <span v-else>  <v-icon large left> mdi-plus </v-icon> Additonal</span></v-btn>
            <input
              type="file"
              ref="input3"
              style="display: none"
              @change="previewImage3"
              accept="image/*"
            />
          </div>

          <div v-if="img3 != null">
            <v-card class="text-center mt-4">
              <v-card-title>Photo 3</v-card-title>
    <v-progress-linear
      v-model="uploadValue2"
    
    ></v-progress-linear>
              <v-img class="rounded contain" :src="img3" />

              <br />
            </v-card>
          </div>
        </div>
      </div>
      <div class="pt-4 text-right">
        <v-btn
          :disabled="!valid || !!flight || !hasAddress"
          color="primary"
          @click="create()"
        >
          Place item <v-icon> mdi-arrow-right-bold</v-icon>
          <div class="button__label" v-if="flight">
            <div class="button__label__icon">
              <icon-refresh />
            </div>
            Creating item...
          </div>
        </v-btn>
      </div>
    </div>
  </div>
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
     
      imageData: null,
      imageData2: null,
      imageData3: null,
      img1: null,
      img2: null,
      img3: null,
      //thisitem: {},
      //itemid: "",
      flight: false,
      uploadValue: 0,
      uploadValue2: 0,
      uploadValue3: 0,

    };
  },

  computed: {
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      if (this.imageData != null) return true;
    },
    /*thisitem() {
      this.itemid = this.$store.state.newitemID;
      return this.$store.getters.getItemByID(this.itemid);
    },*/
  },

  methods: {
    click1() {
      this.$refs.input1.click();
    },
    click2() {
      this.$refs.input2.click();
    },
    click3() {
      this.$refs.input3.click();
    },

    create() {

      let uploadDate = fb.database.ServerValue.TIMESTAMP
      const post = { photos: {
        photo: this.img1,
        photo2: this.img2,
        photo3: this.img3,
        //_id: this.$store.state.user.uid,
        //itemid: this.thisitem.id,
      }, id: { username: this.thisitem.creator, _id: this.$store.state.user.uid, uploadDate: uploadDate }};
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

    //i am lazy and busy so I double the functions for the other images. Code needs to be improved later ofc.
    previewImage(event) {
      this.uploadValue = 0;
      this.img1 = null;
      this.imageData = event.target.files[0];
      this.onUpload();
    },

    onUpload() {
      this.img1 = null;
      let storageRef = fb
        .storage()
        .ref(`${this.imageData.name}`)
        .put(this.imageData);
      storageRef.on(
        `state_changed`,
        (snapshot) => {
          this.uploadValue =
            (snapshot.bytesTransferred / snapshot.totalBytes) * 100;
        },
        (error) => {
          console.log(error.message);
        },
        () => {
          this.uploadValue = 100;
          storageRef.snapshot.ref.getDownloadURL().then((url) => {
            this.img1 = url;
            console.log(this.img1);
            console.log(this.imageData);
          });
        }
      );
    },
    previewImage2(event) {
      this.uploadValue2 = 0;
      this.img2 = null;
      this.imageData2 = event.target.files[0];
      this.onUpload2();
    },

    onUpload2() {
      this.img2 = null;
      let storageRef = fb
        .storage()
        .ref(`${this.imageData2.name}`)
        .put(this.imageData2);
      storageRef.on(
        `state_changed`,
        (snapshot) => {
          this.uploadValue2 =
            (snapshot.bytesTransferred / snapshot.totalBytes) * 100;
        },
        (error) => {
          console.log(error.message);
        },
        () => {
          this.uploadValue2 = 100;
          storageRef.snapshot.ref.getDownloadURL().then((url) => {
            this.img2 = url;
            console.log(this.img2);
          });
        }
      );
    },
    previewImage3(event) {
      this.uploadValue3 = 0;
      this.img3 = null;
      this.imageData3 = event.target.files[0];
      this.onUpload3();
    },

    onUpload3() {
      this.img3 = null;
      let storageRef = fb
        .storage()
        .ref(`${this.imageData3.name}`)
        .put(this.imageData3);
      storageRef.on(
        `state_changed`,
        (snapshot) => {
          this.uploadValue3 =
            (snapshot.bytesTransferred / snapshot.totalBytes) * 100;
        },
        (error) => {
          console.log(error.message);
        },
        () => {
          this.uploadValue3 = 100;
          storageRef.snapshot.ref.getDownloadURL().then((url) => {
            this.img3 = url;
            console.log(this.img3);
          });
        }
      );
    },
  },
};
</script>
