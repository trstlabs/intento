<template>
  <div>
    
    <div>
      <div>
        <div class="item">
  <v-row>
  <v-col >
    <v-card elevation="0" >
    
     <div class="overline">Title</div>
     
  <div class="body-1"> "
           {{thisitem.title }} "
         </div>  </v-card>
  </v-col>
  <v-col >
  <v-card elevation="0">
  <v-chip-group>
    <v-chip outlined small
            v-for="itemtag in thisitem.tags" :key="itemtag"
          > <v-icon small left>
        mdi-tag-outline
      </v-icon>
            {{ itemtag }}
          </v-chip><!--<v-chip 
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
        </v-chip-group> </v-card>
  </v-col>
  </v-row>
           <v-card elevation="0" >  <div class=" overline">Description</div> <v-card-text>
    
     
  <div class="body-1 "> "
           {{ thisitem.description }} "
         </div> </v-card-text> </v-card>

<v-divider class="ma-4"></v-divider>
         <v-chip
      class="ma-1"
      medium label outlined

    >
    
    <v-icon left>
        mdi-account-outline
      </v-icon>
      Your Address: {{ thisitem.creator }}
    </v-chip>

          
          
          <v-chip
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-account-badge-outline
      </v-icon>
      Identifier: {{ thisitem.id }}
    </v-chip>

   
           <v-chip v-if="thisitem.shippingcost"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant
      </v-icon>
      Shipping option available
    </v-chip>
<v-chip v-if="thisitem.localpickup"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
         mdi-map-marker-outline
      </v-icon>
      Local pickup option available
    </v-chip>
    <v-chip v-if="thisitem.shippingcost"
      class="ma-1"
      label
      outlined
      medium

    >
    <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      Shipping cost: ${{thisitem.shippingcost}} TPP
    </v-chip>
    <v-chip
      class="ma-1"
      label
      outlined
      medium

    >
    <v-icon left>
        mdi-star-outline
      </v-icon>
      Condition: {{thisitem.condition}}/5
    </v-chip>

  <!--  <v-chip label color="grey"
        class="ma-1" disabled
      medium>
<v-rating
  background-color="#eee"
  color="white"
  dense
 
  readonly
  length="5"
  size="12"
  :value="thisitem.condition"
></v-rating> </v-chip>-->
    <v-row class="text-center">

</v-row>
    
   <v-divider class="ma-4"></v-divider>
          
          
          <div class="mt-2">
            <v-btn  outlined @click="click1"> <v-icon
          
          left
        >
          mdi-image
        </v-icon>Primary  </v-btn>
            <input
              type="file"
              ref="input1"
              style="display: none"
              @change="previewImage"
              accept="image/*"
            />
          </div>
          <div v-if="img1 != null">
            <v-card class="text-center mt-4">
            <v-card-title>Primary photo</v-card-title>
            <v-img class=" rounded contain"   :src="img1" />

            <br /> </v-card>
          </div>
          <div class="mt-2">
            <v-btn outlined @click="click2"> <v-icon
          
          left
        >
          mdi-image
        </v-icon>Additonal </v-btn>
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
            <v-card-title>Photo 2</v-card-title>
            
            <v-img class=" rounded contain" :src="img2" />

            <br /> </v-card>
          </div>
          <div class="mt-2">
            <v-btn outlined v-if="imageData2" @click="click3"> <v-icon
          
          left
        >
          mdi-image
        </v-icon>
              Additonal
            </v-btn>
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
           
            <v-img class=" rounded contain" :src="img3" />

            <br /> </v-card>
          </div>
        </div>

        <v-btn
          :class="[`button__valid__${!!valid && !flight && hasAddress}`]"
          @click="create()"
        >
          Place item
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

export default {
  props: ["thisitem"],
  components: { AppText, CreateItemForm },
  data() {
    return {
      fields: {
        title: "",
        description: "",
        shippingcost: "0",
        localpickup: false,
        estimationcount: "5",
      },
      imageData: null,
      imageData2: null,
      imageData3: null,
      img1: null,
      img2: null,
      img3: null,
      //thisitem: {},
      //itemid: "",
      flight: false,
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
      const post = {
        photo: this.img1,
        photo2: this.img2,
        photo3: this.img3,
        //itemid: this.thisitem.id,
      };
      console.log(firebase);
      firebase
        .database()
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
      console.log(firebase);
      this.uploadValue = 0;
      this.img1 = null;
      this.imageData = event.target.files[0];
      this.onUpload();
    },

    onUpload() {
      this.img1 = null;
      const storageRef = firebase
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
      console.log(firebase);
      this.uploadValue2 = 0;
      this.img2 = null;
      this.imageData2 = event.target.files[0];
      this.onUpload2();
    },

    onUpload2() {
      this.img2 = null;
      const storageRef = firebase
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
      console.log(firebase);
      this.uploadValue3 = 0;
      this.img3 = null;
      this.imageData3 = event.target.files[0];
      this.onUpload3();
    },

    onUpload3() {
      this.img3 = null;
      const storageRef = firebase
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
        },

        
      );

    },
 
  },

};
</script>


<style scoped>
button {
  background: none;
  border: none;
  color: rgba(0, 125, 255);
  padding: 0;
  font-size: inherit;
  font-weight: 800;
  font-family: inherit;
  text-transform: uppercase;
  margin-top: 0.5rem;
  cursor: pointer;
  transition: opacity 0.1s;
  letter-spacing: 0.03em;
  transition: color 0.25s;
  display: inline-flex;
  align-items: center;
}
.item {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow: hidden;
}
.item__field {
  display: grid;
  line-height: 1.5;
  grid-template-columns: 15% 1fr;
  grid-template-rows: 1fr;
  word-break: break-all;
}
.item__field__key {
  color: rgba(0, 0, 0, 0.25);
  word-break: keep-all;
  overflow: hidden;
}
button:focus {
  opacity: 0.85;
  outline: none;
}
.v-btn.button__valid__true:active {
  opacity: 0.65;
  color: rgba(0, 125, 255);
}
.button__label {
  display: inline-flex;
  align-items: center;
}
.button__label__icon {
  height: 1em;
  width: 1em;
  margin: 0 0.5em 0 0.5em;
  fill: rgba(0, 0, 0, 0.25);
  animation: rotate linear 4s infinite;
}
.v-btn.button__valid__false {
  color: rgba(0, 0, 0, 0.25);
  cursor: not-allowed;
}
.card__empty {
  margin-bottom: 1rem;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 8px;
  color: rgba(0, 0, 0, 0.25);
  text-align: center;
  min-height: 8rem;
}
@keyframes rotate {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(-360deg);
  }
}
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>