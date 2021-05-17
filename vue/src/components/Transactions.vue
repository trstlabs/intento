<template>
  <div v-if="depsLoaded">
    <div class="overline ma-2 text-center">Send Transaction</div>
 
      <v-row class="ma-0">
        <v-col cols="11" sm="3"
          ><v-text-field
            label="Amount"
            type="number"
            solo
            rounded
            v-model="amount"
            :rules="[rules.price]"
            append-icon="$vuetify.icons.custom"
          ></v-text-field></v-col
        ><v-col cols="12" sm="9">
          <v-text-field
            v-model="toAddress"
            @click:append-outer="sendMsg"
            solo
            rounded
            append-outer-icon="mdi-send"
            label="Cosmos-address"
          ></v-text-field>
        </v-col>
      </v-row>

      <div class="overline mx-2 text-center">Query account</div>
      <v-row class="ma-0">
        <v-text-field class="mx-4"
          v-model="bankAddress"
          @click:append="queryBalance"
          solo
          rounded
          append-icon="mdi-magnify"
          label="Cosmos-address"
        ></v-text-field>
      </v-row>
  <div class="overline mx-2 text-center">TPP Transactions</div>
      <v-expansion-panels class="caption ma-0">
        <v-expansion-panel v-for="(tx, i) in transactions" :key="i">
          <v-expansion-panel-header
            >TX: {{ tx.response.logs[0].events[0].attributes[0].value }}
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="caption"
                  >Block height: {{ tx.response.height }}</v-list-item-title
                >
                <v-list-item-title class="caption"
                  >Response:<span v-if="(tx.response.code = '0')">
                    Successful<v-icon color="success" small>
                      mdi-checkbox-marked-circle</v-icon
                    ></span
                  >
                  <span v-else>
                    <v-icon color="warning" right small> mdi-close</v-icon
                    >Failed (code {{ tx.response.code }}</span
                  ></v-list-item-title
                >
                <v-list-item-title class="caption"
                  >Timestamp: {{ getFmtTime(tx.response.timestamp) }}</v-list-item-title
                >
                <v-list-item-title class="caption"
                  >Gas used: {{ tx.response.gas_used }}</v-list-item-title
                >
                <v-list-item-title
                  v-if="tx.auth_info.fee.amount[0]"
                  class="caption"
                  >Fee: {{ tx.auth_info.fee.amount[0].amount
                  }}<v-icon right small
                    >$vuetify.icons.custom</v-icon
                  ></v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item v-if="tx.body.messages[0].itemid">
              <v-list-item-content>
                <v-list-item-title
                  ><v-btn
                    outlined
                    block
                    rounded
                    target="_blank"
                    :to="{
                      name: 'BuyItemDetails',
                      params: { id: tx.body.messages[0].itemid },
                    }"
                  >
                    TPP ID: {{ tx.body.messages[0].itemid }}
                  </v-btn></v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item v-else-if="tx.body.messages[0].id">
              <v-list-item-content>
                <v-list-item-title
                  ><v-btn
                    outlined
                    block
                    rounded
                    target="_blank"
                    :to="{
                      name: 'BuyItemDetails',
                      params: { id: tx.body.messages[0].id },
                    }"
                  >
                    TPP ID: {{ tx.body.messages[0].id }}
                  </v-btn></v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item v-if="tx.body.messages[0].estimation">
              <v-list-item-content>
                <v-list-item-title class="caption"
                  >Estimation:
                  {{ tx.body.messages[0].estimation }} TPP</v-list-item-title
                >
                <estimator-item-item-info
                  :itemid="tx.body.messages[0].itemid"
                />
              </v-list-item-content>
            </v-list-item>
            <div v-if="tx.body.messages[0].seller">
              <v-list-item-content>
                <seller-item-item-info
                  class="ma-0 pa-0"
                  :itemid="tx.body.messages[0].id"
                />
              </v-list-item-content>
            </div>
            <div v-if="tx.body.messages[0].buyer">
              <buyer-item-item-info
                class="ma-0 pa-0"
                :itemid="tx.body.messages[0].itemid"
              />
            </div>

            <div
              v-for="(event, eventi) in tx.response.logs[0].events"
              :key="eventi"
            >
              <v-card class="rounded-lg my-2" outlined>
                <span
                  v-for="(attribute, attributei) in event.attributes"
                  :key="attributei"
                >
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title
                        ><span class="caption">{{
                          attribute.key.toUpperCase()
                        }}</span>
                        :
                        <span v-if="attribute.value == moduleAddress"
                          ><v-icon small left>mdi-shield-lock</v-icon>TPP Module
                          Account</span
                        ><span v-else-if="attribute.value == bankAddress"
                          ><v-icon small>mdi-account</v-icon> You</span
                        ><span v-else>
                          {{ attribute.value }}</span
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                </span>
              </v-card>
            </div>
          </v-expansion-panel-content>
        </v-expansion-panel> <div v-if="!transactions[0] ">
      <p class="caption pa-12 text-center">No transactions found for this account address
    </p>
    </div>  
      </v-expansion-panels>


    <!--<v-data-table
    :headers="headers"
    :items="transactions"
    :items-per-page="5"
    class="elevation-1"
  ></v-data-table>
-->

    <sign-tx
      v-if="submitted"
      :key="submitted"
      :fields="fields"
      :value="value"
      :typeUrl="msg"
      @clicked="afterSubmit"
    ></sign-tx>
  </div>
</template>
<script>
import dayjs from 'dayjs'
//import { decode } from 'js-base64'

export default {
  name: "TransferList",

  //	props: { address: String, refresh: Boolean },
  data: function () {
    return {
      toAddress: "",
      amount: 0,
      submitted: false,
      rules: {
        price: (value) => value > 0 || "Must be positive :)",
      },
    };
  },
  computed: {
    depsLoaded() {
      console.log("dep");
      return true;
    },
    sentTransactions() {
      console.log("se");
      return this.$store.getters.getSentTransactions || [];
    },
    receivedTransactions() {
      console.log("rec");
      return this.$store.getters.getReceivedTransactions || [];
    },
    fullBalances() {
      return this.balances.map((x) => {
        this.addMapping(x);
        return x;
      });
    },
    transactions() {
      console.log(this.sentTransactions);
      let sent =
        this.sentTransactions.txs?.map((tx, index) => {
          tx.response = this.sentTransactions.tx_responses[index];
          return tx;
        }) || [];
      let received =
        this.receivedTransactions.txs?.map((tx, index) => {
          tx.response = this.receivedTransactions.tx_responses[index];
          return tx;
        }) || [];
      console.log(
        [...sent, ...received].sort(
          (a, b) => b.response.height - a.response.height
        )
      );

      return [...sent, ...received].sort(
        (a, b) => b.response.height - a.response.height
      );
    },
  },
  beforeCreate() {},

  async created() {
    this.bankAddress = this.$store.state.account.address;
    this.moduleAddress = process.env.VUE_APP_MODULE;

    console.log(this.bankAddress);

    if (this.depsLoaded) {
      console.log("TEXT");
    }
  },
  methods: {
    getFmtTime(time) {
      	const momentTime = dayjs(time)
      return momentTime.format("D MMM, YYYY HH:mm:ss");
    },

    async afterSubmit(value) {
      this.msg = "";
      this.fields = [];
      this.value = {};
      if (value == true) {
        await this.$store.dispatch("setTransactions", this.bankAddress);

        await this.$store.dispatch("bankBalancesGet");
      }
      this.submitted = false;
    },

    async queryBalance() {
      await this.$store.dispatch("setTransactions", this.bankAddress);
    },

    async sendMsg() {
      this.loadingitem = true;
      this.flightre = true;

      this.fields = [
        ["fromAddress", 1, "string", "optional"],
        ["toAddress", 2, "string", "optional"],
        ["amount", 3, "string", "repeated"],
      ];
      this.msg = "/cosmos.bank.v1beta1.MsgSend";

      (this.value = {
        fromAddress: this.$store.state.account.address,
        toAddress: this.toAddress,
        amount: [{ amount: this.amount, denom: "tpp" }],
      }),
        (this.submitted = true);
    },

    /*
async sendMsg() {
 
  this.loadingitem = true;
      this.flightre = true;

 const wallet = this.$store.state.wallet;





const client = await SigningStargateClient.connectWithSigner( process.env.VUE_APP_RPC, wallet);

const fee = {
  amount: [{ amount: '0', denom: 'tpp' }],
  gas: '200000'
};

const msg = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
      amount:  [{ amount: '5', denom: 'tpp' }],
      fromAddress: this.$store.state.account.address,
      toAddress: this.toAddress,
  }
};
const result = await client.signAndBroadcast(this.$store.state.account.address, [msg], fee, "Welcome to the Trust Price Protocol community");
assertIsBroadcastTxSuccess(result);



}*/
  },
};
</script>
