<template>
  <div class="pa-2 mx-lg-auto">
   
     <v-dialog v-if=(!fields.title)
      v-model="dialog"
      width="500"
    >
      <template v-slot:activator="{ on, attrs }">
       
    <h2  v-bind="attrs"
          v-on="on" class="headline pa-4 text-center">Place a new item</h2>
      </template>

      <v-card class="text-center">
        <v-card-title class="h2 lighten-2 ">
          Info
        </v-card-title>

        <v-card-text>
          After placing the item, an estimation will be made.
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="dialog = false"
          >
            Let's go
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-else
      v-model="dialog"
      width="500"
    >
      <template v-slot:activator="{ on, attrs }">
       <h2  v-bind="attrs"
          v-on="on"  class="headline  pa-4 text-center">Place '{{fields.title}}'</h2>
    
      </template>

      <v-card class="text-center">
        <v-card-title class="headline lighten-2 ">
          Info
        </v-card-title>

        <v-card-text>
          After placing the item, an estimation will be made.
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="dialog = false"
          >
            Let's go
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
   
    
   

     

    <v-stepper class="elevation-0" v-model="e1">
      <v-stepper-header>
        <v-stepper-step :complete="e1 > 1" step="1"> Data </v-stepper-step>

        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 2" step="2">
          Pictures
        </v-stepper-step>
        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 3" step="3"> Done! </v-stepper-step>
      </v-stepper-header>
      <v-stepper-items>
        <div>
          <v-stepper-content step="1" >
            <div class="ma-5">
             
              
              <v-text-field class="ma-1" prepend-icon="mdi-format-title"
                :rules="rules.titleRules"
                label="Title"
                v-model="fields.title" required
                
              />
              
              <v-textarea class="ma-1" prepend-icon="mdi-text"
              :rules="rules.descriptionRules"
            v-model="fields.description"
            label="Description"
            auto-grow
            
            
          > </v-textarea>

               <v-combobox 
                 prepend-icon="mdi-tag-outline"
                 hint="At least one and at most 5 category tags"
                :persistent-hint="selectedTags == 0 || selectedTags == 5"
                 :search-input.sync="search"
          v-model="selectedTags"
          :items="taglist"
         :rules="rules.tagRules"
          label="Categories"
          deletable-chips
          multiple
          chips
          
        > <template v-slot:no-data>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title>
              No results matching "<strong>{{ search }}</strong>". Press <kbd>enter</kbd> to create a new one
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template></v-combobox> 
      
      
      
              <div>
        

<v-row><v-btn class="pa-2"
        
        text
        icon
        @click="fields.estimationcount = 3"
      >
        <v-icon > mdi-check</v-icon></v-btn>
              <v-slider class="pa-2 "
                hint="Lower for faster results, higher for better accuracy"
                thumb-label
                
        
       
                :persistent-hint="fields.estimationcount != 3"
                label="Accuracy"
                :thumb-size="90"
                max="12"
                :rules="rules.estimationcountRules"
                placeholder="Estimation count"
                v-model="fields.estimationcount"
              ><template v-slot:thumb-label="item">
            {{ item.value }} Estimations
          </template></v-slider>
      </v-row>
      <v-row> <v-btn class="pa-2"
        
        text
        icon
        @click="fields.condition = 0"
      >
        <v-icon >{{fields.condition === 0 ? 'mdi-star-outline' : 'mdi-star'}}</v-icon>
      </v-btn>
                <v-slider class="pa-2"  label="Condition"
                  :hint="'Condition is '+ conditionLabel() + ', explain condition in description'"
        v-model="fields.condition"
       
        :max="4"
     
        :persistent-hint="fields.condition != 0"
        
       
        :thumb-size="24"
          thumb-label
      ><template v-slot:thumb-label="{ value }">
            {{ satisfactionEmojis[value] }}
          </template> </v-slider></v-row> 

           <v-row>    
           <v-btn class="pa-2 mt-2"
        
        text
        icon
        @click="fields.shippingcost = 0"
      >
        <v-icon > {{fields.shippingcost === 0 ? 'mdi-package-variant' : 'mdi-package-variant-closed'}} </v-icon>
      </v-btn>

                <v-slider class="pa-2 mt-2"
                  hint="Set to 0 tokens no for shipping"
                  
                  thumb-label
                  label="Shipping cost"
                  suffix="tokens"
                  :persistent-hint="fields.shippingcost != 0"
                  
                  placeholder="Shipping cost"
                  :thumb-size="70"
                  v-model="fields.shippingcost"
                  
                ><template v-slot:thumb-label="item">
            {{ item.value }} tokens
          </template> </v-slider>
</v-row>  <v-row  > <v-col>  <v-row>    
           <v-btn class="pa-2"
        
        text
        icon
        @click="fields.localpickup = !fields.localpickup"
      >
        <v-icon > {{fields.localpickup ? 'mdi-map-marker' : 'mdi-map-marker-off'}} </v-icon>
      </v-btn>


      <v-switch class="ml-2" 
      v-model="fields.localpickup"
      inset
      label="Local pickup"
      :persistent-hint="fields.shippingcost !=0 && fields.localpickup == true && selectedCountries.length > 1"
      hint="Specify local pickup location in description"
    ></v-switch>  

    

               
                </v-row></v-col><v-col> <v-select
                 prepend-icon="mdi-earth"
                 hint="At least one"
                 :persistent-hint="selectedCountries == 0"
                
          v-model="selectedCountries"
          :items="countryCodes"
         :rules="rules.shippingRules"
          label="Location"
          deletable-chips
          multiple
          chips
          
        > </v-select> 
</v-col></v-row>

               <!-- <v-row v-if="fields.shippingcost == 0 ">    
           <v-btn class="pa-2"
        
        text
        icon
        @click="fields.localpickup = !fields.localpickup"
      >
        <v-icon > {{fields.localpickup ? 'mdi-map-marker' : 'mdi-map-marker-off'}} </v-icon>
      </v-btn>


      <v-switch class="ml-2" 
      v-model="fields.localpickup"
      inset
      label="Local pickup"
      
      
    ></v-switch>  

    

               
                </v-row>-->
                

     <!--  <v-combobox
                 hint="Maximum of 5 tags"
                 
          v-model="selectedTags"
          :items="taglist"
          label="Tags"
          deletable-chips
          multiple
          chips
          
        > </v-combobox>-->

              </div>
              <div class="text-center pt-6">
              <v-btn color="primary" :disabled="!valid || !!flight || !hasAddress"
              
                @click="submit()"
              > Next  
                <v-icon > mdi-arrow-right-bold</v-icon>
                <div class="v-btn__label" v-if="flight">
                  <div class="v-btn__label__icon">
                    <icon-refresh />
                  </div>
                  Creating item ID...
                </div>
              </v-btn>
              </div>
            </div>
          </v-stepper-content>
        </div>

        <v-stepper-content step="2">
          <div v-if="showpreview">
            <create-item-preview-and-upload
              :thisitem="thisitem"
              v-on:changeStep="updateStepCount($event)"
            />
          </div>
        </v-stepper-content>
        <v-stepper-content step="3">
          <v-alert rounded-lg type="success">
            Submitted, the item will be estimated
          </v-alert>
          <p>
            You can always find your item in your
            <router-link to="/account">account</router-link>. Your item will be available to buy after you make it transferable.
          </p>
          
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </div>
</template>


<script>
//import AppText from "./AppText.vue";
import CreateItemPreviewAndUpload from "./CreateItemPreviewAndUpload.vue";
export default {
  components: { CreateItemPreviewAndUpload },
  data() {
    return {
      fields: {
        title: "",
        description: "",
        shippingcost: "0",
        localpickup: false,
        estimationcount: "3",
        //openbox: false,
        condition: "0",

      },
      
      rules: {
        titleRules: [
          (v) => !!v || "Title is required",
          (v) =>
            (v && v.length <= 80) || "Title must be less than 80 characters",
        ],
        descriptionRules: [
          (v) => !!v || "Description is required",
          (v) =>
            (v && v.length > 4) || "Description must be more than 4 characters",
        ],
        estimationcountRules: [
          (v) => !!v || "Estimation count is required",
          (v) =>
            (v && v > 2) || "Estimation count must be more than 2 estimators",
          (v) =>
            (v && v < 12) || "Estimation count must be less than 12 estimators",
        ],
        tagRules: [
          (v) => !!v.length == 1 || "Category tag is required",
          (v) =>
            (v && v.length < 6) || "Category tags must be less than 6",
        ],
        
         shippingRules:  [ 
          (v) => !!v.length == 1 || "A country is required",
          
        
        ], 
        
      },
      itemid: "",
      selectedTags: [],
      selectedCountries: [],
      thisitem: {},
      flight: false,
      showpreview: false,
      e1: 1,
      search: null,
      dialog: false,
      
      satisfactionEmojis: ['😭', '🙁', '🙂', '😊', '😄'],
      countryCodes:["NL", "BE", "UK", "DE", "US","CA"]
    };
    
    
  },
  watch: {
      selectedTags (val) {
        if (val.length > 5) {
          this.$nextTick(() => this.selectedTags.pop())
        }
      },
    },

  mounted() {
    
    //console.log(input)
    
      this.$store.dispatch("setSortedTagList");
  
   },


  

  computed: {
    taglist(){
      //this.$store.dispatch("setTagList");
      return this.$store.getters.getTagList },
      //return ["sadfd","dasf"]; },
   
   

    hasAddress() {
     
      return !!this.$store.state.account.address || alert("Sign in first");
    },
    
    valid() {
      if (
        this.fields.title.trim().length > 3 &&
      this.fields.description.trim().length > 4 && this.selectedTags.length > 0 && this.selectedCountries.length > 0 
      )
       {
        return true;
    };
  }, },

  methods: {
    async submit() {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
        const type = { type: "item" };
       
       const fields = [
          ["creator", 1,'string', "optional"],
           [ "title", 2,'string', "optional"] ,                                                    
          ["description",3,'string', "optional"],
         ["shippingcost",4,'int64', "optional"],
          ["localpickup",5,'bool', "optional"],
          ["estimationcount",6,'int64', "optional" ],
         ["tags", 7,'string', "repeated"],
           ["condition",  8, 'int64', "optional"],
           ["shippingregion",  9,'string', "repeated"],
           ["depositamount", 10, "int64", "optional" ]
      
        ];
 //const body = [this.$store.state.account.address,"dsaf", "asdf", 33, 1, "sdfsdf", "asdf", 4, "sfda"]
        const body = {
          creator: this.$store.state.account.address,
            title: this.fields.title,                                                    
          description: this.fields.description,
          shippingcost: this.fields.shippingcost,
           localpickup: this.fields.localpickup,
        estimationcount: this.fields.estimationcount,
         tags: this.selectedTags,
           condition: this.fields.condition,
          shippingregion: this.selectedCountries,
          depositamount: this.fields.estimationcount,
        };
        
        
      
        await this.$store.dispatch("itemSubmit", { ...type,fields, body });
        //const payload = { ...type, body }
        //await this.$store.dispatch("entityFetch", payload);
        //await this.$store.dispatch("accountUpdate");
       
        


        this.flight = false;
        //this.fields.title = "";
       // this.fields.description = "";
        //this.fields.shippingcost = "";
       // this.fields.localpickup = false;
        //this.fields.estimationcount = "";
        this.itemid = await this.$store.state.newitemID;
        //console.log()
        console.log(this.itemid);
        this.thisitem = await this.$store.getters.getItemByID(this.itemid);
        this.e1 = 2;
        this.showpreview = true;
        //alert("Submitted, find the item in the account section");
      }
    },
    updateStepCount(e1) {
      this.e1 = e1;
    },

    conditionLabel(){

      if (this.fields.condition === 0){ return "'bad'"; };
      if (this.fields.condition === 1){ return "'fixable'"; };
      if (this.fields.condition === 2){ return "'decent'"; };
      if (this.fields.condition === 3){ return "'as new'"; };
      if (this.fields.condition === 4){ return "'perfect'"; };

      
    }

   
  },
};
</script>

