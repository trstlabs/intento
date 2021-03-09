
<template>
  <div>
    <v-app>
      <v-main  :class="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-3'">
    
          
         
        
          

    
           

            
              <v-sheet class="ma-6" rounded="lg" elevation="4"> 
                
                <div >
                 
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