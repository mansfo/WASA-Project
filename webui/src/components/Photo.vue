<script>
import CommentModal from "./CommentModal.vue"
export default{
    props: ["uname", "uid", "ppict", "phId", "uplDate", "img", "likes", "likecnt", "comments", "commcnt", "likeStatus"],
    data() {
        return {
            imgURL: "",
            liked: false,
            propic: "",
            commentsList: [],
            showModal: false,
        };
    },

    computed: {
        sameUser() {
            return this.uid === localStorage.getItem("userid");
        },
    },
    methods: {
        reload() {
            this.$emit("reload");
        },
        loadPhoto() {
            this.liked = this.likeStatus;
            let len = this.uid.length + "/tmp/media/photos/".length;
            this.imgURL = __API_URL__ + "/users/" + this.uid + "/photos/" + this.img.slice(len + 1);
        },
        async delPhoto() {
            if (confirm("Are you sure you want to delete this photo?")) {
                try {
                    await this.$axios.delete("/users/" + localStorage.getItem("userid") + "/photos/" + this.phId);
                    this.$emit("delete", this.phId);
                }
                catch (e) {
                    alert(e.toString());
                }
            }
        },
        getLikes() {
            this.$emit("showLikers", this.uid, this.likes, this.phId);
        },
        async toggleLike() {
            try {
                if (this.liked) {
                    await this.$axios.delete("/users/" + this.uid + "/photos/" + this.phId + "/likes/" + localStorage.getItem("userid"));
                }
                else {
                    await this.$axios.put("/users/" + this.uid + "/photos/" + this.phId + "/likes/" + localStorage.getItem("userid"));
                }
                this.liked = !this.liked;
                this.reload();
            }
            catch (e) {
                alert(e.toString());
            }
        },
        async getComments() {
            this.showModal = !this.showModal
            this.reload()
        },
        appendComment(newComment){
            this.commentsList.push(newComment)
            this.reload()
        },
        removeComment(commId){
            this.commentsList = this.commentsList.filter(item => item.comment_id.comment_id !== commId)
            this.reload()
        },
        toProfile() {
            this.$emit("gotoprofile", this.uid);
        },
        toProfileFromComm(identifier){
            if (this.$route.path.includes("stream")){
                this.$emit("gotoprofile", identifier)
            }
            else{
                this.$emit("gotowriter", identifier)
            }

        }
    },
    async mounted() {
        await this.loadPhoto();
        this.propic = this.ppict !== "../assets/images/noPicture.jpg" ? __API_URL__ + "/users/" + this.uid + "/profile/profile_picture/" + this.ppict.slice(this.uid.length + "/tmp/media/profile_pictures/".length + 1) : this.ppict;
        if (this.comments != null){
            this.commentsList = this.comments
        }
    },
    async updated() {
        await this.loadPhoto();
        this.propic = this.ppict !== "../assets/images/noPicture.jpg" ? __API_URL__ + "/users/" + this.uid + "/profile/profile_picture/" + this.ppict.slice(this.uid.length + "/tmp/media/profile_pictures/".length + 1) : this.ppict;
    },
    components: { CommentModal }
}
</script>

<template>
    <div class="container">
        <CommentModal
        :show="showModal"
        :commList="comments"
        :photo="phId"
        :phOwner="uid"
        @closemodal="getComments"
        @delcomment="removeComment"
        @addcomment="appendComment"
        @gotoprof="toProfileFromComm"
        />
        <div class="card">
            <div @click="toProfile">
                <div class="mini-pro-pic-box">
                    <div v-if="propic !== '../assets/images/noPicture.jpg'">
                        <img class="mini-pro-pict" :src="propic" >
                    </div>
                    <div v-else>
                        <img class="mini-pro-pict" src="../assets/images/noPicture.jpg">
                    </div>
                </div>
                <div class="uid">
                    <h3>{{uid}}</h3>
                </div>
                <div v-if="sameUser">
                    <img style="float: right; margin-right: 10px;" class="mini-pict" src="../assets/images/trash.jpg" @click="delPhoto">
                </div>
            </div>
            <hr style="margin-top: 10px;">
            <div style="overflow: hidden; display: flex; flex-wrap: wrap;">
                <img class="main-img" :src="imgURL" >
            </div>
            <hr style="margin-bottom: 5px;">
            <div class="likes-box">
                <div v-if="liked">
                    <img class="like-pict" src="../assets/images/unlike.jpeg" @click="toggleLike" >
                </div>
                <div v-else>
                    <div v-if="!sameUser">
                        <img class="like-pict" src="../assets/images/like.jpg" @click="toggleLike" >
                    </div>
                    <div v-else>
                        <img class="like-pict" src="../assets/images/like.jpg" >
                    </div>
                </div>
                <div v-if="likecnt == 1">
                    <div class="like-infos">
                        <h6 style="margin-top: 5px; text-align: left; font-size: 0.55cm; margin-left: 165px;" @click="getLikes">{{likecnt}} user likes this photo</h6>
                    </div>
                </div>
                <div v-else-if="likecnt == 0">
                    <div class="like-infos">
                        <h6 style="margin-top: 5px; text-align: left; font-size: 0.55cm; margin-left: 148px;">Nobody likes this photo :(</h6>
                    </div>
                </div>
                <div v-else>
                    <div class="like-infos">
                        <h6 style="margin-top: 5px; text-align: left; font-size: 0.55cm; margin-left: 165px;" @click="getLikes">{{likecnt}} users like this photo</h6>
                    </div>
                </div>
            </div>
            <hr style="margin-top: 0px; margin-bottom: 5px;">
            <div class="comments-box" @click="getComments">
                <img class="like-pict" src="../assets/images/comment.png">
                <div v-if="commcnt !== 1">
                    <div class="comments-infos">
                        <h6 style="font-size: 0.55cm;">There are {{commcnt}} comments</h6>
                    </div>
                </div>
                <div v-else>
                    <div class="comments-infos">
                        <h6 style="font-size: 0.55cm; margin-left: 15px;">There is {{commcnt}} comment</h6>
                    </div>
                </div>
            </div>
            <hr style="margin-top: 0px; margin-bottom: 5px;">
            <div class="date">
                <h6 style="font-family: 'Times New Roman', Times, serif; text-align: center;">Posted on {{uplDate.slice(0, 10)+" at "+uplDate.slice(11, 16)}}</h6>
            </div>
        </div>
    </div>
</template>

<style>
.card {
  width: 640px;
  height: 750px;
  margin-left: 340px;
  margin-top: 30px;
  margin-bottom: 5px;
  border-color: darkblue;
  border-style: groove;
}
.mini-pro-pic-box{
    overflow: hidden;
    float: left;
    height: 50px;
    width: 50px;
    margin-left: 20px;
    border-radius: 50%;
    border-color: darkblue;
    border-style: groove;
    margin-top: 10px;
}
.mini-pro-pict{
    min-height: 50px;
    max-height: 50px;
    margin-left: 0;
    position: relative;
    left: 50%;
    transform: translateX(-50%);
}
.mini-pict{
    float: left;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    border-color: darkblue;
    border-style: groove;
    margin-top: 10px;
}
.like-pict{
    float: left;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    border-color: darkblue;
    border-style: groove;
    margin-top: 2px;
    margin-left: 25px; 
}
.main-img{
    position: relative;
    left: 50%;
    transform: translateX(-50%);
    max-height: 500px;
    min-height: 500px;
}
.uid{
    float: left;
    margin-left: 35px;
    margin-top: 18px;
}

.like-infos{
    float: left;
    font-family: 'Times New Roman', Times, serif;
}

.comments-infos{
    text-align: center;
    font-family: 'Times New Roman', Times, serif;
    margin-left: 165px;
    float: left;
    margin-top: 5px;
}
</style>

<style scoped>
hr {
    height: 2px;
    border-width: 0px;
    background-color: darkblue;
    color: darkblue;
}
</style>