module.exports = {
    types: [
        // this line is used by starport scaffolding
        { type: "estimator", fields: ["estimator", "estimation", "itemid", "interested", "comment", "flag"] },
        { type: "estimator/delete", fields: ["creator", "itemid"] },

        { type: "estimator/set", fields: ["creator","itemid", "interested"] },
        { type: "estimator/flag", fields: ["estimator","flag", "itemid", ] },
        { type: "item/reveal", fields: ["creator", "itemid",] },
        { type: "item/transferable", fields: ["creator", "transferbool", "itemid",] },
        { type: "item/transfer", fields: ["creator", "transferbool", "itemid",] },
        { type: "item/shipping", fields: ["creator", "tracking", "itemid",] },
        { type: "item", fields: ["creator", "title", "description", "shippingcost", "localpickup", "estimationcount", "estimationprice", "buyer", "status", "transferable", "bestestimator", "lowestestimator", "highestestimator", "comments", "tags", "condition", "shippingregion", "depositamount"] },
        { type: "item/set", fields: ["creator", "shippingcost", "localpickup", "shippingregion",]},
        { type: "item/delete", fields: ["creator", "id"]},
        { type: "buyer", fields: ["buyer", "deposit", "itemid",] },
    ],
}