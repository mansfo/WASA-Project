<script>
import User from "../components/User.vue"
export default{
    props: ['users', 'page'],
    data() {
        return {
            usrid: "",
            apiURL: "",
        };
    },
    methods: {
        goToProfile(id) {
            this.usrid = id;
            this.$router.push("/users/" + this.usrid + "/profile");
            this.usrid = "";
        }
    },
    mounted(){
        this.apiURL = __API_URL__
    },

    components: { User }
}
</script>

<template>
    <div class="container">
        <div v-if="users.length !== 0">
            <div class="row" style="background-color: white;">
                <User
                v-for="(user, index) in users"
                :key="index"
                :username="user.user.username.username"
                :userid="user.user.user_id.user_id"
                :profPict="user.path_to_image !== '../assets/images/noPicture.jpg' ? apiURL+'/users/'+user.user.user_id.user_id+'/profile/profile_picture/'+user.path_to_image.slice(user.user.user_id.user_id.length + '/tmp/media//profile_pictures/'.length) : '../assets/images/noPicture.jpg'"
                @toProfile="goToProfile"
                />
            </div>
        </div>
        <div v-else-if="page==='SEARCH'">
            <h2 style="text-align: center; margin-top: 20px;">There are no users matching the name you have written :(</h2>
        </div>
        <div v-else-if="page==='BAN'">
            <h2 style="text-align: center; margin-top: 20px;">You haven't banned anyone yet</h2>
        </div>
        <div v-else-if="page==='FOLLOWERS'">
            <h2 style="text-align: center; margin-top: 20px;">No users are following @{{this.$route.params.uid}} :(</h2>
        </div>
        <div v-else-if="page==='FOLLOWING'">
            <h2 style="text-align: center; margin-top: 20px;">@{{this.$route.params.uid}} is not following anyone at this moment</h2>
        </div>
        <div v-else-if="page==='LIKES'">
            <h2 style="text-align: center; margin-top: 20px;">No one likes this photo :(</h2>
        </div>
    </div>
</template>