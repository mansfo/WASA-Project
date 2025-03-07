<script>
export default{
    props: ['username', 'userid', 'profPict'],
    
    methods: {
        selectedUser(){
            this.$emit("toProfile", this.userid)
        }
    }
}
</script>

<template>
    <div @click="selectedUser" class="user-info" style="padding: 10px;">
        <div class="p-pict-box">
            <div v-if="profPict !== '../assets/images/noPicture.jpg'">
                    <img class="p-pict" :src="profPict" >
                </div>
                <div v-else>
                    <img class="p-pict" src="../assets/images/noPicture.jpg">
                </div>
        </div>
        <div class="names">
            <h3>{{username}}</h3>
            <h6>@{{userid}}</h6>
        </div>
    </div>
</template>

<style>
.names{
    font-family: 'Times New Roman', Times, serif;
}
</style>