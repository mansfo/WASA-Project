<script>
import Photo from "../components/Photo.vue";
import DropDownMenu from "../components/DropDownMenu.vue";

export default {
    data() {
        return {
            username: "",
            newName: "",
            followers: [],
            following: [],
            photos: [],
            prof_picture: "",
            photo_count: 0,
            followers_count: 0,
            following_count: 0,
            banned: false,
            followed: false,
            optMenuActive: false,
        };
    },
    watch: {
        currentUser(uid, olduid){
            if (uid !== olduid){
                this.getProfile()
            }
        },
        currentProPic(newPic, oldPic){
            if (newPic !== oldPic){
                this.getProfile()
            }
        }
    },
    computed: {
        sameUser() {
            return this.$route.params.uid === localStorage.getItem("userid");
        },
        currentUser() {
            return this.$route.params.uid
        },
        currentProPic() {
            return this.prof_picture
        },
    },
    methods: {
        async getProfile() {
            if (this.$route.params.uid === undefined) {
                return;
            }
            try {
                let response = await this.$axios.get("/users/" + this.$route.params.uid + "/profile");
                this.username = response.data.user.username.username;
                this.prof_picture = response.data.profile_picture.uri_image != "../assets/images/noPicture.jpg"? __API_URL__+"/users/"+this.$route.params.uid+"/profile/profile_picture/"+response.data.profile_picture.uri_image.slice(this.$route.params.uid.length + "/tmp/media/profile_pictures/".length + 1) : response.data.profile_picture.uri_image;
                this.followed = response.data.followers != null ? response.data.followers.find(obj => obj.user.user_id.user_id === localStorage.getItem("userid")) : false;
                this.photo_count = response.data.photo_counter;
                this.followers_count = response.data.followers_counter;
                this.following_count = response.data.following_counter;
                this.photos = response.data.photos != null ? response.data.photos : [];
                this.followers = response.data.followers != null ? response.data.followers : [];
                this.following = response.data.following != null ? response.data.following : [];
                this.banned = response.status == 206
            }
            catch (e) {
                alert(e.toString())
            }
        },
        options() {
            this.optMenuActive = !this.optMenuActive;
        },
        async postPhoto() {
            let inputFile = document.getElementById("fileUploader");
            const file = inputFile.files[0];
            const reader = new FileReader();
            reader.readAsArrayBuffer(file);
            reader.onload = async () => {
                await this.$axios.post("/users/" + localStorage.getItem("userid") + "/photos", reader.result, {
                    headers: {
                        "Content-Type": file.type
                    },
                });
                this.getProfile()
            };
        },
        deletePhoto(photoIdentifier) {
            this.photos = this.photos.filter(item => item.photo_id.photo_id !== photoIdentifier)
            this.getProfile()
        },
        getLikes(phOwner, likers, imId){
            this.$emit("getlikes", phOwner, "LIKES", likers, imId)
        },
        async follow() {
            try {
                if (this.followed) {
                    await this.$axios.delete("/users/" + localStorage.getItem("userid") + "/following/" + this.$route.params.uid);
                }
                else {
                    await this.$axios.put("/users/" + localStorage.getItem("userid") + "/following/" + this.$route.params.uid);
                }
                this.followed = !this.followed;
                this.getProfile()
            }
            catch (e) {
                talert(e.toString());
            }
        },
        getFollowers(){
            this.$emit("getfollowers", this.$route.params.uid, "FOLLOWERS", this.followers, "")
        },
        getFollowing(){
            this.$emit("getfollowing", this.$route.params.uid, "FOLLOWING", this.following, "")
        },
        async banUser() {
            try {
                if (this.banned) {
                    await this.$axios.delete("/users/" + localStorage.getItem("userid") + "/banned/" + this.$route.params.uid);
                }
                else {
                    await this.$axios.put("/users/" + localStorage.getItem("userid") + "/banned/" + this.$route.params.uid);
                }
                this.banned = !this.banned;
                this.getProfile()
            }
            catch (e) {
                alert(e.toString());
            }
        },
        getBanned(){
            this.$emit("getbanned", "banned", "BAN", [], "")
        },
        async updateNickname(name) {
            this.newName = name;
            try {
                let response = await this.$axios.put("/users/"+this.$route.params.uid+"/profile/username", {
                    username: this.newName
                });
                this.username = response.data.username;
                localStorage.removeItem('username')
                localStorage.setItem('username', this.username)
                this.getProfile()
            }
            catch (e) {
                if (e.toString() === "AxiosError: Request failed with status code 400"){
                    alert("The username must be from 3 to 16 characters long!")
                } else if (e.toString() === "AxiosError: Request failed with status code 403"){
                    alert("This username is already taken by someone else!")
                } else {
                    alert(e.toString())
                }
            }
        },
        viewWriter(owner){
            this.$emit("toProfileFromComm", "/users/"+owner+"/profile")
        }
    },
    async mounted() {
        await this.getProfile();
    },
    
    components: { DropDownMenu, Photo }
}
</script>

<template>
    <div>
        <div class="usr-info">
            <div class="p-pict-box">
                <div v-if="prof_picture !== '../assets/images/noPicture.jpg'">
                    <img class="p-pict" :src="prof_picture" >
                </div>
                <div v-else>
                    <img class="p-pict" src="../assets/images/noPicture.jpg">
                </div>
            </div>
            <div class="names">
                <h3>{{username}}</h3>
                <h6>@{{this.$route.params.uid}}</h6>
            </div>
            <div v-if="!sameUser && !banned">
                <div class="follow-btn-box">
                    <div v-if="followed">
                        <button class="menu-btn" style="margin-top: 25px;" @click="follow"><b>Unfollow</b></button>
                    </div>
                    <div v-else>
                        <button class="menu-btn" style="margin-top: 25px;" @click="follow"><b>Follow</b></button>
                    </div>
                </div>
            </div>
            <div v-else-if="sameUser">
                <div class="follow-btn-box">
                    <input id="fileUploader" style="display: none;" type="file" @change="postPhoto" ref="fileInput" accept=".jpg, .png, .jpeg">
                    <button class="menu-btn" style="margin-top: 25px; margin-left: 10px;" @click="$refs.fileInput.click()"><b>Post something!</b></button>
                </div>
            </div>
            <img class="imp-img" src="../assets/images/impostazioni.png"  @click="options">
            <div style="float: right; overflow: hidden; z-index: 100; margin-right: 120px;">
                <DropDownMenu v-if="optMenuActive"
                @updateUsername="updateNickname"
                :newname="newName"
                @ban="banUser"
                @viewBans="getBanned"
                @reloadP="getProfile"
                @close="options"
                :activeState="optMenuActive"
                :banState="banned" />
            </div>
            <div v-if="!banned">
                <div @click="getFollowing" class="container">
                    <h5 class="h5-prof" style="margin-right: 25px;">Following: {{following_count}}</h5>
                </div>
                <div @click="getFollowers" class="container">
                    <h5 class="h5-prof">Followers: {{followers_count}}</h5>
                </div>
                <h5 class="h5-prof">Photos: {{photo_count}}</h5>
            </div>
        </div>
        <div class="usr-photos">
            <div v-if="banned">
                <h2>You have banned @{{this.$route.params.uid}}, remove the ban to see his/her photos!</h2>
            </div>
            <div v-else-if="photo_count == 0 && !sameUser">
                <h2>{{username}} has not posted anything yet :(</h2>
            </div>
            <div v-else-if="photo_count == 0 && sameUser">
                <h2>Why don't you post something?</h2>
            </div>
            <div v-else>
                <Photo
                v-for="(photo,index) in photos"
                :key="index"
                :uname="photo.user.user.username.username"
                :uid="photo.user.user.user_id.user_id"
                :ppict="photo.user.path_to_image"
                :phId="photo.photo_id.photo_id"
                :uplDate="photo.date"
                :img="photo.image"
                :likes="photo.likes"
                :likecnt="photo.like_counter"
                :comments="photo.comments"
                :commcnt="photo.comment_counter"
                :likeStatus="photo.like_status"
                @reload="getProfile"
                @showLikers="getLikes"
                @delete="deletePhoto"
                @gotowriter="viewWriter"
                />
            </div>
        </div>
    </div>
</template>

<style>
.usr-info{
    background-color: white;
    padding-bottom: 100px;
    z-index: -1;
}
.p-pict-box{
    overflow: hidden; 
    float: left;
    height: 90px;
    width: 90px;
    margin-left: 100px;
    border-radius: 50%;
    border-color: darkblue;
    border-style: groove;
    margin-top: 5px;
}
.p-pict{
    max-height: 90px;
    min-height: 90px;
    position: relative;
    left: 50%;
    transform: translateX(-50%);
}
.names{
    float: left;
    padding-top: 20px;
    margin-left: 50px;
}
.follow-btn-box{
    float: left;
    margin-left: 50px;
}
h2{
    font-family: 'Times New Roman', Times, serif; 
    color: darkgray;
    text-align: center;
}
h3, h4{
    font-family: 'Times New Roman', Times, serif;
}
h4{
    font-style: italic;
}
.h5-prof{
    float: right;
    padding-left: 50px;
    font-family: 'Times New Roman', Times, serif;
    padding-top: 35px;
    margin-right: 50px;
}
.usr-photos{
    margin-top: 20px;
    z-index: -1;
}
.imp-img{
    height: 50px;
    float: right;
    margin-right: 100px;
    margin-top: 25px;
}
</style>