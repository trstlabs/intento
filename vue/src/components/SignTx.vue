<template>
  <div class="pa-2 mx-lg-auto">
    <v-dialog persistent v-model="dialog" width="500">
    
      <v-card class="text-center rounded-lg">
        <v-card-title class="h2 lighten-2">Send Transaction </v-card-title>

        <v-card-text>
          Set the amount of gas or the fee you are willing to pay.
        </v-card-text>
        <v-row class="ma-2">
          <v-col>
            <p class="caption">Gas maximum:</p>

            <v-text-field
              label="Amount"
              type="number"
              v-model="gas"
              :rules="[rules.price]"
              append-icon="mdi-fuel"
            ></v-text-field>
          </v-col>
         
          <v-col>
            <p class="caption">Fee maximum:</p>

            <v-text-field
              label="Amount"
              type="number"
              v-model="amount"
              :rules="[rules.price]"
              append-icon="$vuetify.icons.custom"
            ></v-text-field>
          </v-col>
         </v-row
        >

        <v-divider></v-divider>
        <v-alert dense class="ma-2" type="warning" v-if="error" >The transaction result is: {{error.rawLog}}</v-alert>
          <v-alert dense class="ma-2" type="success" v-else-if="sent" >The transaction is sent</v-alert>
        <v-card-actions>
          <v-btn text @click="$emit('clicked', false)"> Discard </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="submit()"><span v-if="sent">re</span>Send </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

    
<script>

import {
  SigningStargateClient,
  assertIsBroadcastTxSuccess,isBroadcastTxSuccess,
} from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

export default {
   props: {
 
  fields: Array,
  value: Object,
  msg: String,
  typeUrl: String,

},
   
 //  props: ["type", "fields", "value", "msg"],
  data() {
    return {
      amount: 0,
       dialog: true,
      sent: false,
      error: "",
      gas: 200000,
       rules: {
      
          price: value => value > -1 || 'Must be positive',
       
          },
    };
  },

methods: {
  async submit() {

      console.log("submitting")
      const wallet = this.$store.state.wallet;
      //const type2 = type.charAt(0).toUpperCase() + type.slice(1);

      console.log(this.msg)
      
      if(this.typeUrl){ var typeUrl = this.typeUrl
      
      var client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
      )
      
      
      }else{
      var typeUrl = `/${process.env.VUE_APP_PATH}.`+ this.msg;

  
      let MsgCreate = new Type(`${this.msg}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      this.fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });
   var client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
      { registry }
     
     
      );

          }
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
       value: this.value
      };

      console.log(msg);
      const fee = {
        amount: [{ amount: this.amount.toString(), denom: "tpp" }],
        gas: this.gas.toString(),
      };
      const result = await client.signAndBroadcast(
        this.$store.state.account.address,
        [msg],
        fee
      ); 
    

     if (isBroadcastTxSuccess(result)) {
            this.$emit('clicked', true)
     }else{
         this.error = result
     
     }

     this.sent = true
      assertIsBroadcastTxSuccess(result);
      console.log("success!");

      //passs whether it is success or falure back to any parent component
 
      
    },

}
};
</script>