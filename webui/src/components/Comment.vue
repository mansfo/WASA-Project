<script>

export default{
    props:['commId', 'photo', 'owner', 'authorId', 'authorProPic', 'message', 'date'],
    data(){
        return{
            propic:"",
        }
    },
    computed:{
        mycomment(){
            return this.authorId === localStorage.getItem("userid")
        },
    },
    methods:{
        async deleteComm(){
            await this.$axios.delete("/users/"+this.owner+"/photos/"+this.photo+"/comments/"+this.commId)
            this.$emit("delcomm", this.commId)
        },
        viewprofile(){
            this.$emit("viewProf", this.authorId)
        },
    },

    mounted(){
        this.propic = this.authorProPic !== "../assets/images/noPicture.jpg" ? __API_URL__ + "/users/" + this.authorId + "/profile/profile_picture/" + this.authorProPic.slice(this.authorId.length + "/tmp/media/profile_pictures/".length + 1) : this.authorProPic;
    },
    updated(){
       this.propic = this.authorProPic !== "../assets/images/noPicture.jpg" ? __API_URL__ + "/users/" + this.authorId + "/profile/profile_picture/" + this.authorProPic.slice(this.authorId.length + "/tmp/media/profile_pictures/".length + 1) : this.authorProPic; 
    }
}
</script>

<template>
    <div class="container">
        <div class="row">
            <div style="float: left;">
                <div class="comm-infos">
                    <div class="propic-box">
                        <div v-if="propic !== '../assets/images/noPicture.jpg'">
                            <img class="propic-writer" :src="propic" >
                        </div>
                        <div v-else>
                            <img class="propic-writer" src="../assets/images/noPicture.jpg">
                        </div>
                    </div>
                    <div class="name" @click="viewprofile">
                        <h4>{{authorId}}</h4>
                    </div>
                    <div style="float: right;">
                        <h6 style="color: darkgray; padding-top: 4px;">{{date}}</h6>
                        <div v-if="mycomment">
                            <img class="del-btn" src="../assets/images/trash.jpg" @click="deleteComm">
                        </div>
                    </div>
                </div>
                <div style="text-align: left; width: 600px; margin-top: 10px;">
                    <div class="col-12" style="word-wrap: break-word;">{{message}}</div>
                </div>
            </div>
        </div>
        <hr>
    </div>
</template>

<style scoped>
h4{
    font-size: 0.6cm;
    font-family: 'Times New Roman', Times, serif;
    font-style: italic;
    text-align: left;
}
.propic-box{
    border-radius: 50%;
    width: 30px;
    height: 30px;
    border-color: darkgray;
    border-style: groove;
    border-width: 0.05cm;
    overflow: hidden;
    float: left;
}
.propic-writer{
    min-height: 30px;
    max-height: 30px;
    position: relative;
    left: 50%;
    transform: translateX(-50%);
}
.name{
    float: left;
    margin-left: 15px;
}
.comm-box{
    padding-bottom: 10px;
    padding: 10px;
    margin-top: 10px;
}
.comm-infos{
    padding-top: 0px;
    padding-bottom: 30px;
}
.comm-body{
    padding-bottom: 30px;
    padding-top: 10px;
}
.paragraph{
    margin-bottom: 0px;
    word-wrap: break-word;
    max-width: 650px;
    min-width: 650px;
    float: left;
}
.comm-text{
    text-align: left;
    font-family: 'Times New Roman', Times, serif;
    font-size: 0.5cm;
}
.del-btn{
    border-radius: 50%;
    width: 30px;
    height: 30px;
    border-color: darkblue;
    border-style: groove;
    border-width: 0.05cm;
}
</style>