<template>
  <div class="bl-doc-header">
    <div>header</div>
  </div>
  <div class="bl-doc-footer">
    <div>footer</div>
  </div>
  <div class="bl-doc-frontpage">hallo</div>
  <div class="bl-doc-page-break-before">hallo</div>
  <h1></h1>

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

<style lang="css" scope>

@media screen, pagedjs-ignore {
    body {
        background-color: #cccccc;
    }
    .pagedjs_page {
        background-color: #ffffff;
        border: 1px solid #000000;
        margin: 40px auto;
    }
  }

 
</style>
<style lang="css">


  

@page {
  
  @top-center {
    content: element(header);
  }
  
  @bottom-center {
    content: element(footer);
  }
  
}

.bl-doc-header {
  display: block;
  position: running(header);
}

.bl-doc-footer {
  display: block;
  position: running(footer);
}

.bl-doc-page-break-before {
  break-before: page;
}

 
</style>