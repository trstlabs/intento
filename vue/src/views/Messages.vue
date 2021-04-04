
<template>
  <div>
    <v-app>
      <v-main  :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
    
          
         
        
          

    
           

            
              <v-sheet class="ma-6" rounded="lg" elevation="4"> 
                
                <div >
                  <v-alert v-if="!this.$store.state.user && this.$store.state.account.address "
  type="warning" dismissible
>Confirm the verification link sent to the Email linked to your Google account</v-alert>
                  <message-box  />
                <!--  <message-boxtest :currentUserId="this.$store.state.account.address"/>-->
               
                </div>
              </v-sheet>
         
            
       
  
      </v-main>
    </v-app>
  </div>
</template>

<script>

export default {
 


  
}
</script>