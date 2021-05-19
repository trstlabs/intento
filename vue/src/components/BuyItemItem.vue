<template>
  <div class="pa-2 mx-lg-auto">
    <v-progress-linear indeterminate :active="loadingitem"></v-progress-linear>
    <v-card v-if="!loadingitem" class="pa-2 ma-auto" elevation="0" rounded="lg">
      <v-row class="ma-0 pa-2 mb-2 text-center"
        ><v-col cols="2">
          <v-tooltip bottom v-if="thisitem.creator != thisitem.seller">
            <template v-slot:activator="{ on, attrs }">
              <span v-bind="attrs" v-on="on"
                ><v-icon color="primary lighten-2" icon plain>
                  mdi-repeat
                </v-icon></span
              >
            </template>
            <span>This item is reposted</span>
          </v-tooltip>

          <v-tooltip bottom v-else>
            <template v-slot:activator="{ on, attrs }">
              <span v-bind="attrs" v-on="on"
                ><v-icon color="primary lighten-2" icon plain>
                  mdi-check-all
                </v-icon></span
              >
            </template>
            <!-- <span v-if="thisitem.estimatorlist[0]">This item is estimated by {{thisitem.estimatorlist.length}} estimators</span>-<span v-else>Priced through estimations</span>--><span
              >Priced through estimations</span
            >
          </v-tooltip>
        </v-col>
        <v-col cols="8">
          <p class="display-1 font-weight-thin">
            {{ thisitem.title }}
          </p> </v-col
        ><v-col cols="2">
          <span class="d-flex d-sm-none">
            <v-btn @click="shareItem" icon plain color="primary">
              <v-icon>mdi-share-variant</v-icon>
            </v-btn>
          </span>
          <v-speed-dial
            v-model="dialShare"
            direction="bottom"
            open-on-hover
            class="d-none d-sm-flex justify-center"
          >
            <template v-slot:activator>
              <v-btn icon plain color="primary">
                <v-icon v-if="dialShare">mdi-close</v-icon>
                <v-icon v-else>mdi-share-variant</v-icon>
              </v-btn>
            </template>
            <v-btn
              dark
              fab
              color="blue"
              small
              :href="`https://twitter.com/share?url=${pageUrl}&text=${encodeURI(
                'Check out this ' + thisitem.title
              )}&via=TrustPrice&hashtags=${thisitem.tags}`"
              target="_blank"
            >
              <v-icon>mdi-twitter</v-icon>
            </v-btn>
            <v-btn
              dark
              fab
              color="blue accent-5"
              small
              :href="`https://www.facebook.com/sharer/sharer.php?u=${pageUrl}`"
              target="_blank"
            >
              <v-icon>mdi-facebook</v-icon>
            </v-btn>
            <v-btn
              dark
              fab
              color="green"
              small
              :href="`https://wa.me/?text=${encodeURI(
                'Check out this ' + thisitem.title
              )}%20${pageUrl}`"
              target="_blank"
            >
              <v-icon>mdi-whatsapp</v-icon>
            </v-btn>
            <v-btn
              dark
              fab
              color="tertiary"
              small
              :href="`mailto:?subject=I found something you might like&amp;body=Check out this ${thisitem.title} at ${pageUrl}`"
              target="_blank"
            >
              <v-icon>mdi-email</v-icon>
            </v-btn>
            <v-btn
              dark
              fab
              color="blue"
              small
              :href="`https://t.me/share/url?url=${pageUrl}&text=${encodeURI(
                'Check out this ' + thisitem.title
              )}`"
              target="_blank"
            >
              <v-icon>mdi-telegram</v-icon>
            </v-btn>
          </v-speed-dial>
        </v-col>
      </v-row>

      <div>
        <div class="pa-2 mx-auto" elevation="8">
          <div>
            <div v-if="photos[0]">
              <v-carousel
                v-if="magnify == false"
                style="height: 100%"
                delimiter-icon="mdi-minus"
                carousel-controls-bg="primary"
                height="300"
                hide-delimiter-background
                show-arrows-on-hover
              >
                <v-carousel-item
                  max-height="300"
                  contain
                  v-for="(photo, i) in photos"
                  :key="i"
                  :src="photo"
                >
                </v-carousel-item>
              </v-carousel>
              <v-carousel
                v-if="magnify == true"
                delimiter-icon="mdi-minus"
                carousel-controls-bg="primary"
                contain
                hide-delimiter-background
                show-arrows-on-hover
              >
                <v-carousel-item
                  contain
                  v-for="(photo, i) in photos"
                  :key="i"
                  :src="photo"
                >
                </v-carousel-item>
              </v-carousel>
            </div>
            <v-row class="ml-4 mt-1 mb-1">
              <span v-for="(photo, index) in photos" :key="index">
                <img
                  class="ma-1"
                  @click="show(photo)"
                  height="56"
                  :src="photo" /></span
              ><v-spacer /><v-btn
                x-small
                class="mr-4"
                color="primary"
                icon
                @click="magnify = !magnify"
              >
                <v-icon> mdi-crop-free</v-icon>
              </v-btn></v-row
            >
            <v-dialog v-model="fullscreen">
              <v-card>
                <v-card-title class="grey lighten-2">
                  {{ thisitem.title }} <v-spacer></v-spacer>
                  <v-btn color="primary" icon @click="fullscreen = false">
                    <v-icon> mdi-close</v-icon>
                  </v-btn>
                </v-card-title>
                <v-img :src="showphoto" />
              </v-card>
            </v-dialog>
            <div>
              <div class="pa-2 overline text-center">Description</div>
              <v-card-text>
                <div class="body-1">{{ thisitem.description }}</div>
              </v-card-text>
            </div>

            <span v-if="thisitem.note">
              <v-divider class="mx-4 pa-2" />
              <div class="pl-4 overline text-center">
                <span v-if="thisitem.rating > 0">Rating</span
                ><span v-else> Reseller's Note </span>
              </div>
              <v-card-text>
                <div v-if="thisitem.rating > 0" class="body-1">
                  <v-dialog
                    transition="dialog-bottom-transition"
                    max-width="300"
                  >
                    <template v-slot:activator="{ on, attrs }">
                      <span v-bind="attrs" v-on="on">
                        <v-rating
                          :value="Number(thisitem.rating)"
                          readonly
                          color="primary darken-1"
                          background-color="primary lighten-1"
                          small
                          dense
                        ></v-rating>
                      </span>
                    </template>
                    <template v-slot:default="dialog">
                      <v-card>
                        <v-toolbar color="default">Rating</v-toolbar>
                        <v-card-text class="text-left">
                          <div class="text-p pa-2">
                            <v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon>
                            Scam/Fake
                          </div>
                          <div class="text-p pa-2">
                            <v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon>Bad
                          </div>
                          <div class="text-p pa-2">
                            <v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon>
                            Ok
                          </div>
                          <div class="text-p pa-2">
                            <v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star-outline </v-icon>
                            Good
                          </div>
                          <div class="text-p pa-2">
                            <v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon
                            ><v-icon left small> mdi-star </v-icon>
                            Great
                          </div>
                        </v-card-text>
                        <v-card-actions class="justify-end">
                          <v-btn text @click="dialog.value = false"
                            >Close</v-btn
                          >
                        </v-card-actions>
                      </v-card>
                    </template>
                  </v-dialog>
                  <v-icon left class="mb-n8"> mdi-account-edit</v-icon>
                  <v-chip
                    class="ma-2 rounded-0 rounded-br-xl rounded-t-xl"
                    color="primary lighten-2"
                  >
                    {{ thisitem.note }}
                  </v-chip>
                </div>
                <div v-else class="body-1 text-right">
                  <v-chip
                    class="ma-2 rounded-0 rounded-bl-xl rounded-t-xl"
                    color="primary darken-1"
                  >
                    {{ thisitem.note }} </v-chip
                  ><v-icon class="mb-n8" right> mdi-account</v-icon>
                </div>
              </v-card-text></span
            >
            <v-divider class="mx-4 pa-2" />
            <div
              class="text-center pa-2"
              v-if="thisitem.transferable && thisitem.status == '' && !thisitem.buyer"
            >
              <v-row class="mx-4">
                <v-btn x-small icon @click="iteminfo = !iteminfo">
                  <v-icon>mdi-information-outline</v-icon>
                </v-btn>
                <span v-if="this.$store.state.account.address">
                  <wallet-coins
                /></span>

                <div
                  v-if="iteminfo"
                  class="text-center caption font-weight-light pa-2"
                >
                  You can buy {{ thisitem.title }}
                  <span v-if="thisitem.shippingcost > 0"
                    >and ship the item if you live in one of the following
                    locations:
                    <span
                      v-for="(loc, i) in thisitem.shippingregion"
                      :key="loc"
                      class="font-weight-medium"
                    >
                      <span v-if="i + 1 == thisitem.shippingregion.length">{{
                        loc
                      }}</span
                      ><span v-else> {{ loc }}, </span>
                    </span>
                    <span v-if="!thisitem.shippingregion[0]">
                      all locations </span
                    >. Additional cost (e.g. shipping) is
                    {{ thisitem.shippingcost }}
                    <v-icon small>$vuetify.icons.custom</v-icon>.</span
                  >
                  <span v-if="thisitem.localpickup != ''">
                    and you can arrange a pickup by sending a message to
                    <a @click="createRoom">{{ thisitem.seller }}. </a>
                  </span>
                  <span v-if="thisitem.discount > 0">
                    Reseller gives a discount of {{ thisitem.discount }}
                    <v-icon small>$vuetify.icons.custom</v-icon> on the original
                    selling price of {{ thisitem.estimationprice }}
                    <v-icon small>$vuetify.icons.custom</v-icon>.</span
                  >
                  <!-- <span v-if="thisitem.creator == thisitem.seller"
                    >If you buy the item you will receive a cashback reward of
                    {{ (thisitem.estimationprice * 0.05).toFixed(0) }}
                    <v-icon small >$vuetify.icons.custom</v-icon>.
                  </span>-->
                  With TPP you can withdrawl your payment at any time, up until
                  the item transaction and no transaction costs are applied.
                </div>
              </v-row>

              <v-img @click="iteminfo = !iteminfo" src="img/design/buy.png">
              </v-img>
              <v-row v-if="thisitem.creator == thisitem.seller">
                <v-col>
                  <v-btn
                    rounded
                    block
                    color="primary lighten-1"
                    :disabled="thisitem.localpickup == ''"
                    @click="
                      submit(thisitem.estimationprice),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                      <v-icon left> mdi-check-all </v-icon> Buy for
                      {{ thisitem.estimationprice
                      }}<v-icon small right>$vuetify.icons.custom</v-icon>
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn> </v-col
                ><v-col>
                  <v-btn
                    block
                    rounded
                    color="primary lighten-1"
                    :disabled="thisitem.shippingcost == 0"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP">
                      <v-icon left> mdi-check-all</v-icon
                      ><v-icon left> mdi-plus</v-icon
                      ><v-icon left> mdi-package-variant-closed </v-icon>

                      Buy for
                      {{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost)
                      }}
                      <v-icon small>$vuetify.icons.custom</v-icon>
                    </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
              <v-row v-else>
                <v-col>
                  <v-btn
                    rounded
                    v-if="thisitem.localpickup != '' && thisitem.discount > 0"
                    block
                    color="primary lighten-1"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) -
                          Number(thisitem.discount)
                      ),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                      Buy for
                      {{
                        Number(thisitem.estimationprice) -
                        Number(thisitem.discount)
                      }}<v-icon small right>$vuetify.icons.custom</v-icon>
                      <v-icon right> mdi-repeat </v-icon>
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn>
                  <v-btn
                    block
                    rounded
                    color="primary lighten-1"
                    v-if="thisitem.localpickup != '' && thisitem.discount == 0"
                    @click="
                      submit(thisitem.estimationprice),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                      Buy for {{ thisitem.estimationprice
                      }}<v-icon small right>$vuetify.icons.custom</v-icon>
                      <v-icon right> mdi-repeat </v-icon>
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn> </v-col
                ><v-col>
                  <v-btn
                    block
                    rounded
                    color="primary lighten-1"
                    v-if="thisitem.shippingcost > 0 && thisitem.discount == 0"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP">
                      Buy for
                      {{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost)
                      }}<v-icon small right>$vuetify.icons.custom</v-icon>
                      <v-icon right> mdi-repeat </v-icon>
                      <v-icon right> mdi-plus </v-icon
                      ><v-icon right> mdi-package-variant-closed </v-icon>
                    </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn>
                  <v-btn
                    block
                    rounded
                    color="primary lighten-1"
                    v-if="thisitem.shippingcost > 0 && thisitem.discount > 0"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost) -
                          Number(thisitem.discount)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP">
                      Buy for
                      {{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost) -
                        Number(thisitem.discount)
                      }}<v-icon small right>$vuetify.icons.custom</v-icon>
                      <v-icon right> mdi-repeat </v-icon>
                      <v-icon right> mdi-plus </v-icon
                      ><v-icon right> mdi-package-variant-closed </v-icon>
                    </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Awaiting transaction...
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
              <v-divider class="mx-4 mt-4" />
            </div>
            <div class="overline mb-2 text-center">Information</div>

            <v-dialog transition="dialog-bottom-transition" max-width="300">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-chip
                    style="cursor: pointer"
                    class="ma-1 font-weight-light"
                    outlined
                    medium
                    >Condition:
                    <v-rating
                      :value="Number(thisitem.condition)"
                      readonly
                      color="primary darken-1"
                      background-color="primary lighten-1"
                      small
                      dense
                    ></v-rating>
                  </v-chip>
                </span>
              </template>
              <template v-slot:default="dialog">
                <v-card>
                  <v-toolbar color="default"
                    >Condition (provided by seller)</v-toolbar
                  >
                  <v-card-text class="text-left">
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      Bad
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>Fixable
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      Good
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      As New
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon>
                      Perfect
                    </div>
                  </v-card-text>
                  <v-card-actions class="justify-end">
                    <v-btn text @click="dialog.value = false">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </template>
            </v-dialog>

            <v-chip
              v-if="thisitem.localpickup"
              class="ma-1 font-weight-light"
              target="_blank"
              :href="
                'https://www.google.com/maps/search/?api=1&query=' +
                thisitem.localpickup
              "
              outlined
              ><v-icon left> mdi-map-marker-outline </v-icon> Pickup
              Location</v-chip
            >

            <v-chip
              :to="{ name: 'SearchRegion', params: { region: country } }"
              outlined
              class="ma-1 font-weight-light text-uppercase"
              v-for="country in thisitem.shippingregion"
              :key="country"
            >
              <v-icon small left> mdi-flag-variant-outline </v-icon
              >{{ country }}</v-chip
            >

            <v-chip
              :to="{ name: 'SearchTag', params: { tag: tag } }"
              outlined
              class="ma-1 font-weight-light text-capitalize"
              v-for="tag in thisitem.tags"
              :key="tag"
            >
              <v-icon small left> mdi-tag-outline </v-icon>{{ tag }}</v-chip
            >
            <v-card class="ma-1 rounded-t-xl" outlined>
              <v-list dense disabled>
                <v-subheader>About</v-subheader>
                <v-list-item-group>
                  <v-list-item>
                    <v-list-item-icon>
                      <v-icon>mdi-account-badge-outline </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>TPP ID: </v-col>
                          <v-col>{{ thisitem.id }}</v-col></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.creator != thisitem.seller">
                    <v-list-item-icon>
                      <v-icon> mdi-account-outline</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col
                            >Original Seller: {{ thisitem.creator }}</v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-icon>
                      <v-icon
                        @click="createRoom"
                        :disabled="!this.$store.state.account.address"
                      >
                        mdi-account</v-icon
                      >
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Seller: {{ thisitem.seller }}</v-col></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.shippingcost > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-package-variant-closed </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Shipping Cost: </v-col>
                          <v-col
                            >{{ thisitem.shippingcost
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.seller != thisitem.creator">
                    <v-list-item-icon>
                      <v-icon> mdi-repeat</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Original Price: </v-col>
                          <v-col
                            >{{ thisitem.estimationprice
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-else-if="thisitem.estimationprice > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-check-all </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Estimation Price: </v-col>
                          <v-col
                            >{{ thisitem.estimationprice
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.discount > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-brightness-percent</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Discount: </v-col>
                          <v-col
                            >{{ thisitem.discount
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-card>
            <v-row
              v-if="
                this.$store.state.account.address ==
                (thisitem.seller || thisitem.buyer)
              "
              class="justify-center my-4"
            >
              <v-divider class="ma-4" />
              <div>
                <v-btn outlined rounded to="/account=placeditems">Go To Actions</v-btn>
              </div>
              <v-divider class="ma-4"
            /></v-row>
            <div class="overline mt-4 text-center">Comments</div>
            <div v-if="thisitem.comments">
              <div
                class="font-weight-light"
                v-for="(single, i) in allcomments"
                v-bind:key="i"
              >
                <v-icon class="mb-n8"> mdi-account-check</v-icon>
                <v-chip
                  class="ma-2 rounded-0 rounded-br-xl rounded-t-xl text-no-wrap primary lighten-2"
                >
                  {{ single }}
                </v-chip>
              </div>
            </div>
            <div v-if="allcomments.length == 0">
              <p class="caption text-center">No comments to show right now</p>
            </div>

            <v-divider class="ma-4" />
          </div>
        </div>
      </div>
      <v-row class="pa-2 mx-auto">
        <v-btn text rounded @click="sellerInfo">
          <v-icon v-if="!info" left> mdi-plus</v-icon
          ><v-icon v-else left> mdi-close</v-icon>Advanced
        </v-btn>

        <v-spacer />
        <v-btn
          rounded
          :disabled="!this.$store.state.account.address"
          text
          @click="createRoom"
          ><v-icon left> mdi-message-text</v-icon> Message Seller</v-btn
        >
      </v-row>
      <div class="pa-2 mx-auto caption" v-if="info">
        <span>
          <p class="text-center my-4">   </p>
            This seller had {{ sold }} before.<div v-if="ItemDate"> Buyable since: {{ItemDate}}</div><div v-if="ItemReadyDate">Time of last estimation: {{ItemReadyDate}}</div>

           <!--<p  Of which _ have been transfered by shipping and _ by local pickup.</p>-->
        </span>
        <v-card-title v-if="SellerItems[0]" class="overline justify-center">
          All Seller items
        </v-card-title>
        <div v-for="item in SellerItems" v-bind:key="item.id">
          <v-card
            outlined
            class="py-2 my-2 rounded-xl"
            :to="{ name: 'BuyItemDetails', params: { id: item.id } }"
          >
            <v-row class="text-left caption ma-2"
              ><span class="font-weight-medium"> {{ item.title }}</span>
              <v-spacer /><v-spacer />
              <v-chip class="mx-1" v-if="item.status != ''" small>
                {{ item.status }}
              </v-chip>
              <v-chip
                class="mx-1"
                small
                v-if="item.transferable && item.status != ''"
              >
                Sold for: {{ item.estimationprice
                }}<v-icon small right>$vuetify.icons.custom</v-icon> </v-chip
              ><v-chip
                small
                class="mx-1"
                v-if="item.buyer && !item.transferable"
              >
                Sold</v-chip
              >
              <v-chip
                class="mx-1"
                small
                v-else-if="!item.transferable && item.status == ''"
                >Not on sale yet</v-chip
              ><v-chip
                class="mx-1"
                small
                v-else-if="
                  item.transferable &&
                  item.status == '' &&
                  item.creator == thisitem.seller
                "
                color="primary lighten-2"
                >{{ item.estimationprice
                }}<v-icon small right>mdi-check-all</v-icon></v-chip
              ><v-chip
                class="mx-1"
                small
                v-else-if="item.creator != thisitem.seller && item.status == ''"
                color="primary lighten-2"
                ><v-icon small>mdi-repeat</v-icon></v-chip
              >

              <v-chip class="mx-1" small v-if="item.rating > 0"
                >Rating<v-rating
                  :value="Number(item.rating)"
                  readonly
                  color="primary darken-1"
                  background-color="primary lighten-1"
                  x-small
                  dense
                ></v-rating>
              </v-chip>
            </v-row>
          </v-card>
        </div>
      </div>
      <v-img src="img/design/transfer.png"></v-img>
    </v-card><sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
  </div>
</template>
<script>
import BuyItemDetails from "../views/BuyItemDetails.vue";
import { usersRef, roomsRef, databaseRef } from "./firebase/db.js";
import dayjs from 'dayjs'

export default {
  components: { BuyItemDetails },
  props: ["itemid"],

  data() {
    return {
      amount: "",
      iteminfo: false,
 
      flightLP: false,
      flightSP: false,
      info: false,
      imageurl: "",
      loadingitem: false,
      photos: [],
      dialog: false,
      fullscreen: false,
      magnify: false,

       fields: [],
      value: {},
      msg: "",
      submitted: false,
        dialShare: false,
      showphoto: null,
    };
  },

  beforeCreate() {
    this.loadingitem = true;
  },
  mounted() {
    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null) {
        //console.log(data[0]);
        this.photos = data;
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },
  computed: {
    thisitem() {
      if (this.$store.state.data.item) {
        return this.$store.getters.getItemByID(this.itemid) || {};
      } else {
        return {};
      }
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
    allcomments() {
      if (this.thisitem) {
        return this.thisitem.comments.filter((i) => i != "") || [];
      }
    },
    SellerItems() {
      return this.$store.getters.getBuySellerList || [];
    },
    pageUrl() {
      return process.env.VUE_APP_URL + "/itemid=" + this.thisitem.id;
    },
    ItemDate() {
      let event = this.$store.getters.getEvent("ItemTransferable") || []
       return (this.getFmtTime(event.tx_responses[0].timestamp));
    },
     ItemReadyDate() {
      let event = this.$store.getters.getEvent("ItemReadyForReveal") || []
       return (this.getFmtTime(event.tx_responses[0].timestamp));
    }
  },

  methods: {

    getFmtTime(time) {
      	const momentTime = dayjs(time)
      return momentTime.format("D MMM, YYYY HH:mm:ss");
    },
    async submit(deposit) {
      if (!this.hasAddress) {
        alert("Sign in first");
      //  this.$router.push("/");
        window.location.reload();
      }

      if (this.hasAddress) {
        // this.flightLP = true;
        this.loadingitem = true;
       this.msg = "MsgCreateBuyer"
        const body = { deposit: deposit, itemid: this.thisitem.id };
        this.fields = [
          ["buyer", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];

        this.value = {
          buyer: this.$store.state.account.address,
          ...body,
        },
         
  this.submitted = true
        this.loadingitem = false;
      }

    },

    async getThisItem() {
      await submit();
      return thisitem();
    },

    shareItem() {
      if (navigator.share) {
        const shareData = {
          title: this.thisitem.title,
          text: "Check out this " + this.thisitem.title,
          url: this.pageUrl,
        };

        navigator.share(shareData);
      } else {
        alert("web share not supported");
      }
    },

    /*async getItem() {

         const url = `${process.env.VUE_APP_API}/${process.env.VUE_APP_PATH.replace(/\./g, '/')}/${"item/"+ this.$route.params.id}`;
      const body = (await axios.get(url)).data.Item

console.log(url)
console.log(body)
       
     this.thisitem = body
    },
*/

    shareItem() {
      if (navigator.share) {
        const shareData = {
          title: this.thisitem.title,
          text: "Check out this " + this.thisitem.title,
          url: this.pageUrl,
        };

        navigator.share(shareData);
      } else {
        alert("web share not supported");
      }
    },

  /*  async paySubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreateBuyer`;
      let MsgCreate = new Type(`MsgCreateBuyer`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      console.log(fields);
      fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          buyer: this.$store.state.account.address,
          ...body,
        },
      };

      console.log(msg);
      const fee = {
        amount: [{ amount: "0", denom: "tpp" }],
        gas: "200000",
      };

      const result = await client.signAndBroadcast(
        this.$store.state.account.address,
        [msg],
        fee
      );
      if (!result.data) {
        alert("TX failed");
        window.location.reload();
      }
      assertIsBroadcastTxSuccess(result);
      alert("Transaction sent");
*/
          async afterSubmit(value){
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true){
        const type = { type: "buyer" };
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("bankBalancesGet");
         this.$store.dispatch("updateItem", this.thisitem.id)//.then(result => this.newitem = result)
         this.$router.push("/account=boughtitems")}

          this.submitted = false
                this.flightLP = false;
      this.flightSP = false;
    },

    show(photo) {
      this.showphoto = photo;

      this.fullscreen = true;
    },

  
    sellerInfo() {
      this.$store.dispatch("setBuySellerItemList", this.thisitem.seller);
    this.$store.dispatch("setEvent", {type: "ItemTransferable", attribute: "Itemid", value: this.thisitem.id});
  this.$store.dispatch("setEvent", {type: "ItemReadyForReveal", attribute: "Itemid", value: this.thisitem.id});
      let rs = this.SellerItems.filter((i) => i.buyer != "");
      this.sold = "no buyers";
      if (rs[1]) {
        this.sold = rs.length + " buyers";
      }else if (rs[0]) {
        this.sold = rs.length + " buyer";
      }
      this.info = !this.info;
    },
    async createRoom() {
      if (this.$store.state.user) {
        const user = await usersRef
          .where("username", "==", this.thisitem.seller)
          .get();
        console.log(user);

        //let query = await roomsRef.where("users", '', [this.$store.state.user.uid).where("users", "array-contains", user.docs[0].id).get()
        /*await roomsRef.where("users", "==", ["5RlZazMyPgdoHgGfjTud", "B1Xk6qliE2ceNJN6HsoCk2MQO2K2"]).get()
 .then((querySnapshot) => {
    querySnapshot.forEach((doc) => {
      console.log(doc.data())
        console.log(doc.id, ' => ', doc.data());
    });
});*/
        if (user.docs[0]) {
          let query = await roomsRef
            .where("users", "==", [user.docs[0].id, this.$store.state.user.uid])
            .get();
          let otherquery = await roomsRef
            .where("users", "==", [this.$store.state.user.uid, user.docs[0].id])
            .get();
          console.log(query.docs[0]);
          /*

await roomsRef.where("users", "array-contains", this.$store.state.user.uid).get()
   .then((querySnapshot) => {
    querySnapshot.forEach((doc) => {
      console.log(doc.data(users))
        console.log(doc.id, ' => ', doc.data());
    });
});*/

          if (query.docs[0] || otherquery.docs[0]) {
            alert("Seller already added or seller not found");
          } else {
            //await usersRef.doc(id).update({ _id: id });
            await roomsRef.add({
              users: [user.docs[0].id, this.$store.state.user.uid],
              lastUpdated: new Date(),
            });
          }
          this.$router.push("/messages");
        } else {
          alert("Seller not found");
        }
      } else {
        alert("Confirm sign in first with email");
      }
    },
  },
};
</script>


