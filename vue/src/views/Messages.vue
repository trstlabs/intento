
<template>
  <div>
    <v-app>
   <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
    
          
         
        
           <v-container class="pb-6 pt-0 px-0" >

    
            <v-col cols="12" class="mx-auto pa-0">

            
              <v-sheet class="mt-12" rounded="lg" elevation="4"> 
                
                <div >
                  <!--<v-alert v-if="!this.$store.state.user && this.$store.state.account.address "
  type="warning" dismissible
>Confirm the verification link sent to the Email linked to your Google account</v-alert>-->
                  <message-box  />
                <!--  <message-boxtest :currentUserId="this.$store.state.account.address"/>-->
               
                </div>
              </v-sheet>
            </v-col>
           </v-container>
       
  
      </v-main>
    </v-app>
  </div>
</template>

<script>

export default {
 


  
}
</script>