<template>
  <div>
    <div class="sp-container">
      <sp-sign-in />
      <sp-bank-balances />
      <sp-token-send />
      <!-- this line is used by starport scaffolding # 4 -->
		<sp-type-form path="danieljdd.tpp.tpp" type="estimator" :fields="[ ['creator', 1, 'string'] , ['estimation', 2, 'string'] , ['estimatorestimationhash', 3, 'string'] , ['itemid', 4, 'string'] , ['deposit', 5, 'string'] , ['interested', 6, 'string'] , ['comment', 7, 'string'] , ['flag', 8, 'string'] ]" />
		<sp-type-form path="danieljdd.tpp.tpp" type="buyer" :fields="[ ['creator', 1, 'string'] , ['itemid', 2, 'string'] , ['transferable', 3, 'string'] , ['deposit', 4, 'string'] ]" />
		<sp-type-form path="danieljdd.tpp.tpp" type="item" :fields="[ ['creator', 1, 'string'] , ['title', 2, 'string'] , ['description', 3, 'string'] , ['shippingcost', 4, 'string'] , ['localpickup', 5, 'string'] , ['estimationcounthash', 6, 'string'] , ['bestestimator', 7, 'string'] , ['lowestestimator', 8, 'string'] , ['highestestimator', 9, 'string'] , ['estimationprice', 10, 'string'] , ['estimatorlist', 11, 'string'] , ['estimatorestimationhashlist', 12, 'string'] , ['transferable', 13, 'string'] , ['buyer', 14, 'string'] , ['tracking', 15, 'string'] , ['status', 16, 'string'] , ['comments', 17, 'string'] , ['tags', 18, 'string'] , ['flags', 19, 'string'] , ['condition', 20, 'string'] , ['shippingregion', 21, 'string'] ]" />
    </div>
  </div>
</template>

<script>
import * as sp from "@tendermint/vue";

export default {
  components: { ...sp },
};
</script>
