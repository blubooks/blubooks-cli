<template>
        <div>hallo</div>

</template>
<script lang="ts" setup>

import  { Previewer } from "pagedjs";

import { onMounted } from 'vue';

  onMounted(()=>{
    let paged = new Previewer(null);
  
    let flow = paged.preview(null,null,null).then((flow: any) => {
  
      console.log("Rendered", flow.total, "pages.");
    })

    console.log(flow)
  });



</script>

<style lang="css">
@page {
  size: A4;
 
}
</style>