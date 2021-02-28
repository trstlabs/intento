<template>
  <div>
    <div class="sp-container">
      <sp-sign-in />
      <sp-bank-balances />
      <sp-token-send />
      <!-- this line is used by starport scaffolding # 4 -->
		<sp-type-form path="danieljdd.tpp.tpp" type="estimator" :fields="[ ['estimator', 1, 'string'] , ['estimation', 2, 'int64']  , ['itemid', 4, 'string'] , ['deposit', 5, 'string'] , ['interested', 6, 'string'] , ['comment', 7, 'string']  ]" />
		<sp-type-form path="danieljdd.tpp.tpp" type="buyer" :fields="[ ['buyer', 1, 'string'] , ['itemid', 2, 'string'] , ['transferable', 3, 'bool'] , ['deposit', 4, 'string'] ]" />
		<sp-type-form path="danieljdd.tpp.tpp" type="item" :fields="[ ['creator', 1, 'string'] , ['title', 2, 'string'] , ['description', 3, 'string'] , ['shippingcost', 4, 'int64'] , ['localpickup', 5, 'bool'] , ['estimationcount', 6, 'int64'], ['tags', 18, 'repeated string'] , ['condition', 20, 'int64'] , ['shippingregion', 21, 'repeated string'] ]" />
    </div>
  </div>
</template>

<script>
import * as sp from "@tendermint/vue";

export default {
  components: { ...sp },
};
</script>
