<template>
  <div class="pa-2 mx-lg-auto">
    

  
      <div v-if="showinfo === false">

  <div class="card__empty" v-if="showinfo === false">  

        <v-btn :ripple="false" text @click="getItemToEstimate"><v-icon color="primary" large left>
        mdi-refresh
      </v-icon> all items</v-btn>
      </div>
        <v-skeleton-loader
         class="mx-auto"
      
      
          type="list-item-three-line, image, article"
        ></v-skeleton-loader>
    
      
    </div>

     
   

    <v-card class="pa-2 mx-auto" elevation="2" rounded="lg" v-if="showinfo" >
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

     

      <div elevation="8" v-if="photos.photo">
        <v-carousel 
          delimiter-icon="mdi-minus"
          carousel-controls-bg="primary"
          contain
          
          hide-delimiter-background
          show-arrows-on-hover
        >
          <v-carousel-item v-for="(photo, i) in photos" :key="i" :src="photo" >
           
          </v-carousel-item>
        </v-carousel>
        
      </div>

  
    <div class="pa-2">
    <v-row>
  <v-col >
    <v-card elevation="0" >
    
     <div class="overline">Title</div>
     
  <div class="body-1"> "
           {{item.title }} "
         </div>  </v-card>
  </v-col>
  <v-col >
  <v-card elevation="0">
  <v-chip-group>
    <v-chip outlined small
            v-for="itemtag in item.tags" :key="itemtag"
          > <v-icon small left>
        mdi-tag-outline
      </v-icon>
            {{ itemtag }}
          </v-chip>
        </v-chip-group> 
        
        <v-dialog transition="dialog-bottom-transition"
        max-width="300"> <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"
        >
        <v-chip
      class="ma-1"
      
      outlined
      small

    >
    <v-icon left small>
        mdi-star-outline
      </v-icon>
     {{item.condition}}/5
    </v-chip> </span> </template> <template v-slot:default="dialog">
          <v-card>
            <v-toolbar 
              color="default"
              
            >Condition (provided by seller)</v-toolbar>
            <v-card-text class="text-left">
           
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                Bad 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>Fixable 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                Good 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                As New 
                 
              </div><div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon>
                Perfect 
                 
              </div>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn
                text
                @click="dialog.value = false"
              >Close</v-btn>
            </v-card-actions>
            
          </v-card>
        </template>
    </v-dialog >
    
    </v-card>
  </v-col>
  </v-row>
  <v-card elevation="0" > 
    
     <div class="overline">Description</div>
     
  <div class="body-1"> "
           {{item.description }} "
         </div>  </v-card>

   
  
   </div>


    <div class="pa-2 mx-auto text-center" elevation="8" v-if="lastitem">
      
   
   <v-chip
   v-if="lastitem"
      class="mt-2"
      label
      outlined
      medium
      color="warning"
      
    >
    <v-icon left>
        mdi-alarm
      </v-icon>
      This was the last item, check again later.
    </v-chip>

    </div>
   
    

<v-divider></v-divider>

<div class=" mx-auto">

      

      <v-row>
      <v-col cols="4" >

    <v-dialog transition="dialog-bottom-transition"
        max-width="600"
      >
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"
        ><v-subheader class="text-center" >To me, the item is worth</v-subheader>
        </span>
      </template> <template v-slot:default="dialog">
          <v-card>
            <v-toolbar
              color="default"
              
            >Rules</v-toolbar>
            <v-card-text>
              <div class="text-p pt-4">- Earn ~3% of the item value when you are the best estimator.
                
              </div>
              <div class="text-p pa-2">
                - Your deposit is lost when you are the lowest estimator and the final estimation price is not acceoted by the seller. 
                 
              </div>
              <div class="text-p pa-2">
                
                  - Your deposit is lost when you are the highest estimator and the item is not bought by the buyer that provided prepayment. 
              </div>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn
                text
                @click="dialog.value = false"
              >Close</v-btn>
            </v-card-actions>
            
          </v-card>
        </template>
      
    </v-dialog>
        
      </v-col>
      <v-col cols="8">
        <v-text-field
          label="Amount"
          

          type="number"
          v-model="estimation"

          prefix="$"
          suffix="tokens"
        ></v-text-field>
      </v-col>
    </v-row>
      </div>
        <v-divider></v-divider>
        <v-card elevation="0" > <div class="pa-2">
          <div >
<v-chip-group 
            active-class="primary--text"
            column 
          >
            <v-chip small 
              v-for="(option, text) in options"
              :key="text"
              @click="updateComment(option.attr)"
            >
              {{ option.name }}
            </v-chip>
          </v-chip-group>
</div>
     
  

  

    
<div class="mx-auto">
  <h3 class="text-left"> " </h3>
        <v-text-field rounded dense clearable
          placeholder="leave a comment (optional)"
          
          v-model="comment"
        />
      </div>
      <h3 class="text-right"> " </h3>
     
    </div>
    </v-card> 

      <div >
        <v-btn block elevation="4" color="primary"
          :disabled="!valid || !hasAddress || flight"
          @click="submit(estimation, item.id, interested, comment)"
      ><v-icon left>
        mdi-check
      </v-icon>  
          Estimate item
          <div class="button__label" v-if="flight">
            <div class="button__label__icon">
              <icon-refresh />
            </div>
            Creating estimation...
          </div>
        </v-btn>
        <!-- tag bar
 <v-chip-group 
    
          
          active-class="primary--text"
        >
          <v-chip @click="updateList(tag)"  outlined
            v-for="tag in tags" :key="tag"
          ><v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag }}
          </v-chip>
        </v-chip-group>-->

      


      </div>

</v-card>

<div class="pa-4 mx-auto" v-if="showinfo"> 
 

   <v-row class="text-center"> <v-col class="pa-0">
<v-tooltip bottom :disabled="!interested" v-if="showinfo">
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
           <v-btn
      class="mx-2"
      fab
      dark
      small
      color="pink"
      :outlined="interested == false"
      @click="interested = !interested"
    >
      <v-icon dark>
        mdi-heart
      </v-icon>
    </v-btn>
          </span>
    </template>  <span >Find your liked items in the account section when they are available. </span> 
      </v-tooltip>
</v-col><v-col class="pa-0">
     <!-- <v-tooltip bottom :disabled="flag" v-if="showinfo">
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
           <v-btn
      class="mx-2"
      fab
      dark
      small
      color="red"
      :outlined="flag == false"
      @click="submitFlag(true, item.id)"
    >
      <v-icon dark>
        mdi-alert-octagon
      </v-icon>
    </v-btn>
          </span>
    </template>  <span > When this item is not OK, report it. Thank You.</span> 
      </v-tooltip> -->

      <v-dialog bottom :disabled="flag" v-if="showinfo"
      v-model="dialog"
      persistent
      max-width="290"
    >
    
   
      <template v-slot:activator="{ on, attrs }">
         
        <v-btn
      class="mx-2"
      fab
      dark
      small
      v-bind="attrs"
          v-on="on"
      color="red"
      :outlined="flag == false"
      
    >
      <v-icon dark>
        mdi-alert-octagon
      </v-icon>
    </v-btn>

        
      </template>
      <v-card>
        <v-card-title class="headline">
          Report this item?
        </v-card-title>
        <v-card-text>If this item is not OK, you can report it here. TPP protocol will automatically remove items that are reported often.</v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="red darken-1"
            text
            @click="dialog = false"
          >
            Close
          </v-btn>
          <v-btn
            color="red darken-1"
            text
            @click="submitFlag()"
          >
            Report Item
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

</v-col><v-col class="pa-0">
    <v-btn 
      :disabled="estimation > 1 || !hasAddress || !showinfo"
      outlined
      @click="getNewItemByIndex" color="primary"
    >
      <v-icon dark>
        mdi-arrow-right-bold
      </v-icon>
    </v-btn> 
    </v-col> </v-row>



 <div class="pt-12 mx-lg-auto">
 <v-select append-icon="mdi-tag-outline" dense v-model="selectedFilter" v-on:input="updateList(selectedFilter)" cache-items :items="tags" label="Categories"
  clearable
  rounded
  solo
  persistent-hint
  hint="Specify your expertise"
></v-select>
</div>
 </div>

     
 
  </div>
</template>

<script>
import ToEstimateTagBar from "./ToEstimateTagBar.vue";
import { coins } from "@cosmjs/launchpad";
export default {

  components: { ToEstimateTagBar },
  data() {
    return {
      estimation: "",
      comment: "",
    
      
      options: [ 
      { name: "Great Photos!", attr: "Great Photos!"},
            {name: "Unclear Photos", attr: "I find the photos unclear."},
               { name: "Excellent Description", attr: "I find the description excellent."},
      { name: "Too Vage", attr: "I find the description too vague."},
      { name: "Clear", attr: "I find the item well described, the buyer will know what to expect."},
      { name: "Looks damaged", attr: "The item appears to be damaged."},
      { name: "Repairable", attr: "The item seems damaged, but I think it can be repaired."},  
      { name: "Used", attr: "The item seems used."},
            { name: "As good as new!", attr: "The item appears to look as good as new!"},
      { name: "Dirty", attr: "The item looks dirty to me."},
     



      ],
      
      interested: false,
      flag: false,
      flight: false,
      item: "",
      index: 0,
      showinfo: false,
      lastitem: false,
      loadingitem: true,
      photos: [],
      selectedFilter: "",

      dialog: false,
      conditionInfo: false,
    };
  },
  mounted() {
    
    //console.log(input)
    if (this.$store.state.client != null ) { 
      let input = this.$store.state.account.address
      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setEstimatorItemList", input);
  this.$store.dispatch("setToEstimateList", this.index);
  this.showinfo = true;
this.item = this.items[this.index];
this.loadItemPhotos();
    };
   },


  computed: {
    items() {
      return this.$store.getters.getToEstimateList;
    },
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.estimation.trim().length > 0;
    },
    tags() {
   
      return this.$store.getters.getTagList },

    
  },

  methods: {
    async submit(estimation, itemid, interested, comment) {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
        this.loadingitem = true;
        const type = { type: "estimator" };
        const body = { estimation: estimation, itemid: itemid, interested: interested, deposit: "5tpp",  comment: comment };
        
        await this.$store.dispatch("estimationSubmit", { ...type, body });
        console.log("success!")
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("accountUpdate");
        this.submitRevealEstimation(itemid);
        //this.flight = false;
        //this.loadingitem = false;
      
        
      }
    },

    async submitFlag() {
      if ( !this.flight && this.hasAddress) {
        this.flight = true;
        this.loadingitem = true;
        this.flag = true;
        const type = { type: "estimator/flag" };
        const body = { flag: true, itemid: this.item.id };
        
        await this.$store.dispatch("entitySubmit", { ...type, body });
        

        this.estimation = "";

        this.getNewItemByIndex();
        this.dialog = false;
        this.flag = false;
        
      }
    },


    async getItemToEstimate() {
      if (this.$store.state.client == null) {alert("log in first");};
      let input = this.$store.state.account.address;
      this.$store.dispatch("setEstimatorItemList", input);
      //let index = 0;
      //this.$store.dispatch("setToEstimateList");
      this.item = this.items[this.index];
      this.lastitem = false
      this.loadItemPhotos();
      if (this.showinfo == true) {
        this.showinfo = false;
      } else return (this.showinfo = true);
    },
  

    async getNewItemByIndex() {
      let oldindex = this.index;
      if (oldindex >= 0 && oldindex < this.items.length - 1) {
        this.index = oldindex + 1;
      }

      console.log(oldindex, this.index);
      this.item = this.items[this.index];
      if (oldindex === this.index) {
        this.lastitem = true;
      }
      this.loadItemPhotos();
    },
    loadItemPhotos() {
      this.loadingitem = true;
      const id = this.item.id;
      const db = firebase.database();

      const imageRef = db.ref("ItemPhotoGallery/" + id);
      imageRef.on("value", (snapshot) => {
        const data = snapshot.val();

        if (data != null && data.photo != null) {
          console.log(data.photo);
          this.photos = data;

          this.loadingitem = false;
        }
      });
      this.loadingitem = false;
      this.interested = false;
      this.flight = false;
    },

    async submitRevealEstimation(itemid) {
      if (this.hasAddress) {
        this.estimation = "";
        this.comment = "";
        this.getNewItemByIndex();
         const type = { type: "item/reveal" };
        const body = { itemid: itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        
        
        

        
      }
    },



     updateComment(newComment){
this.comment = newComment;


  },
  updateList(tag) {
      //console.log(this.tag);
      this.$store.dispatch("tagToEstimateList", tag);
      if (!!this.items[0]) {this.item = this.items[0];
      if (!this.items[1]) {
        this.lastitem = true;
      }else {this.lastitem = false};
      this.loadItemPhotos();
      }else{alert("No Items to estimate for:" + tag);this.$store.dispatch("setToEstimateList"); this.getItemToEstimate();}
      },
    
  },
 
 
      
};


</script>

<style scoped>

.short{
  width:100px;
}

button {
  background: none;
  border: none;
  color: #3062C6;
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
.button.button__valid__true:active {
  opacity: 0.65;
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
.button.button__valid__false {
  color: rgba(0, 0, 0, 0.25);
  cursor: not-allowed;
}
.card__empty {
  margin-bottom: 1rem;
  border: 1px rgba(0, 0, 0, 0.1);
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
