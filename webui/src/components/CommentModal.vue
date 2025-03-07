<script>
import Comment from "./Comment.vue"
export default{
    data() {
        return {
            commtext: "",
            commId: "",
        };
    },
    props: ["show", "commList", "photo", "phOwner"],
    methods: {
        closeModal() {
            this.$emit("closemodal");
        },
        async postComment() {
            try{
                let response = await this.$axios.post("/users/"+this.phOwner+"/photos/"+this.photo+"/comments", {
                    comment: this.commtext,
                },{
                    headers:{
                        'Content-Type': 'application/json'
                    }
                })
                this.$emit("addcomment", response.data)
                this.commtext= ""
            }catch(e){
                alert(e.toString())
            }
        },
        deleteComment(commIdentifier){
            this.$emit("delcomment", commIdentifier)
        },
        profile(writer){
            this.$emit("gotoprof", writer)
        },
    },
    components: { Comment }
}
</script>

<template>
    <div v-if="show" class="modal-mask">
        <div class="modal-container">
            <div class="modal-header">
                <div class="header">
                    <h3 style="text-align: center; margin-left: 50px;"><i>Comments</i></h3>
                </div>
                <button class="modal-default-button" @click="closeModal"><b>Close</b></button>
            </div>
            <hr style="margin-top: 10px; margin-bottom: 10px;">
            <div class="modal-body">
                <div v-if="commList === null">
                    <h1 style="margin-top: 20px; text-align: center;">There are no comments :(</h1>
                    <hr>
                </div>
                <div v-else>
                    <Comment
                    v-for="(comment, index) in commList"
                    :key="index"
                    :commId="comment.comment_id.comment_id"
                    :photo="photo"
                    :owner="phOwner"
                    :authorId="comment.user.user.user_id.user_id"
                    :authorProPic="comment.user.path_to_image"
                    :message="comment.comment"
                    :date="comment.date.slice(0,10)"
                    @delcomm="deleteComment"
                    @viewProf="profile"
                    />
                </div>
            </div>
            <div class="modal-footer">
                <img class="mini-pict" style="margin-bottom: 10px;" src="../assets/images/comment.png">
                <textarea style="float: left;" placeholder="Add a comment!" id="form-comment" rows="2" cols="60" maxlength="300" v-model="commtext"></textarea>
                <button class="send-btn" @click.prevent="postComment" :disabled="commtext.length < 1 || commtext.length > 300">Publish!</button>
            </div>
        </div>
    </div>
</template>

<style>
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  transition: opacity 0.3s ease;
}

.header{
    justify-content: center;
    width: 800px;
}

.modal-container {
  width: 800px;
  position: absolute;
  margin-top: 0px;
  margin-left: 560px;
  padding: 20px 30px;
  background-color:white;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-header {
  margin-top: 0;
  color: darkblue;
  text-align: center;
}

.modal-body {
  padding-bottom: 5px;
  overflow-y: auto;
  max-height: 600px;
}

.modal-default-button, .send-btn {
  float: right;
  background-color: darkblue;
  font-family: 'Times New Roman', Times, serif;
  color: white;
  padding-left: 20px;
  padding-bottom: 5px;
  padding-right: 20px;
  padding-top: 5px;
}
</style>