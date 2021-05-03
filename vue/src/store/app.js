module.exports = {
    types: [
        // this line is used by starport scaffolding
        { type: "estimator", fields: ["estimator", "estimation", "itemid", "interested", "comment", "flag"] },
       // { type: "estimator/delete", fields: ["creator", "itemid"] },

       // { type: "estimator/set", fields: ["creator","itemid", "interested"] },
       // { type: "estimator/flag", fields: ["creator","flag", "itemid", ] },
        //{ type: "item/reveal", fields: ["creator", "itemid",] },
       // { type: "item/transferable", fields: ["seller", "transferbool", "itemid",] },
        //{ type: "item/transfer", fields: ["seller", "transferbool", "itemid",] },
        //{ type: "item/shipping", fields: ["seller", "tracking", "itemid",] },
        { type: "item", fields: ["seller", "title", "description", "shippingcost", "localpickup", "estimationcount", "estimationprice", "buyer", "status", "transferable", "bestestimator", "lowestestimator", "highestestimator", "comments", "tags", "condition", "shippingregion", "depositamount", "note", "discount", "creator", "rating"] },
       // { type: "item/set", fields: ["seller", "shippingcost", "localpickup", "shippingregion",]},
       // { type: "item/delete", fields: ["seller", "id"]},
        { type: "buyer", fields: ["buyer", "deposit", "itemid"] },
    ],
}