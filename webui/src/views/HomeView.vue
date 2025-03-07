<script>
import Photo from "../components/Photo.vue";
export default {
    data: function () {
        return {
            photos: [],
        };
    },
    methods: {
        async homepage() {
            try {
                let response = await this.$axios.get("/users/" + localStorage.getItem("userid") + "/stream");
                if (response.data != null) {
                    this.photos = response.data;
                }
            }
            catch (e) {
                alert(e.toString());
            }
        },
        getLikes(phOwner, likers, imId){
            this.$emit("getlikes", phOwner, "LIKES", likers, imId)
        },
        viewProfile(owner){
            this.$emit("toProfileFromPhoto", "/users/"+owner+"/profile")
        }
    },
    async mounted() {
        await this.homepage();
    },

    components: { Photo }
}

</script>

<template>
    <div class="container fluid">
		<div v-if="photos.length === 0" class="row">
		    <h1 style="text-align: center; margin-top: 40px;">No content... Have you tried following someone?</h1>
		</div>
        <div v-else>
            <Photo
            v-for="(photo, index) in photos"
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
            @showLikers="getLikes"
            @reload="homepage"
            @gotoprofile="viewProfile"
            />
        </div>
	</div>
</template>

<style>
h1 {
    font-family: 'Times New Roman', Times, serif;
    color: lightgray;
}
</style>