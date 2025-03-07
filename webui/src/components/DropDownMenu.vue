<script>
export default{
    props: ['activeState', 'banState'],
    data(){
        return{
            newname: "",
            bannedUsers: [],
        }
    },
    computed:{
        isMyProfile(){
            return this.$route.params.uid === localStorage.getItem("userid")
        }
    },
    methods: {
        getBans(){
            this.$emit("viewBans")
        },
        reloadPage(){
            this.$emit("reloadP")
        },
        async updateProfPict() {
            let inputFile = document.getElementById("pPictUploader");
            const file = inputFile.files[0];
            const reader = new FileReader();
            reader.readAsArrayBuffer(file);
            reader.onload = async () => {
                await this.$axios.put("/users/" + localStorage.getItem("userid") + "/profile/profile_picture", reader.result, {
                    headers: {
                        "Content-Type": file.type
                    },
                });
                this.reloadPage()
            };
        },
        updateNickname(){
            this.$emit("updateUsername", this.newname)
            this.newname = ""
        },
        toggleBan(){
            this.$emit("ban")
        },
        closeMenu(){
            this.$emit("close")
        },
    },
}
</script>

<template>
    <div class="back-mask">
        <div class="menu-opts">
            <div v-if="activeState">
                <ul>
                <div v-if="isMyProfile">
                    <hr>
                    <h1 style="text-align: center; font-family: 'Times New Roman', Times, serif; color: black;"><i>Options</i></h1>
                    <hr>
                    <div style="text-align: left;">
                        <form style="text-align: center;">
                            <div>
                                <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">Update your Username</h4></li>
                                <input style="padding: 5px; margin: 5px; font-size:0.5cm; font-family: 'Times New Roman', Times, serif;" type="text" v-model="newname" placeholder="Your new Username">
                                <button class="opts-btn" @click.prevent="updateNickname">Update</button>
                            </div>
                        </form>
                    </div>
                    <hr>
                    <div style="text-align: center;">
                    <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">Choose a new profile picture</h4></li>
                        <input id="pPictUploader" style="display: none;" type="file" @change="updateProfPict" ref="pPictInput" accept=".jpg, .png, .jpeg">
                        <button class="opts-btn" style="margin-top: 10px;" @click="$refs.pPictInput.click()">Change Picture</button>
                    </div>
                    <hr>
                    <div style="text-align: center;">
                    <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">See who you have banned</h4></li>
                        <button class="opts-btn" style="margin: 10px;" @click="getBans">Banned Users</button>
                    </div>
                </div>
                <div v-else>
                    <hr>
                    <div v-if="banState">
                        <div style="text-align: center;">
                            <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">Unban this user</h4></li>
                            <button class="opts-btn close-btn" @click="toggleBan">Remove Ban</button>
                        </div>
                    </div>
                    <div v-else>
                        <div style="text-align: center;">
                            <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">Ban this user</h4></li>
                            <button class="opts-btn close-btn" @click="toggleBan">Ban</button>
                        </div>
                    </div>
                </div>
                <hr>
                <div style="text-align: center;">
                    <li><h4 style="font-family: 'Times New Roman', Times, serif; margin-top: 10px;">Close this menu</h4></li>
                    <button class="opts-btn" @click="closeMenu">Close</button>
                </div>
                <hr>
                </ul>
            </div>
        </div>
    </div>
</template>

<style>
.back-mask{
    position: fixed;
    z-index: 9998;
    top: 0;
    right: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    transition: opacity 0.3s ease;
}
.menu-opts{
    background-color: white;
    width: 350px;
    height: 100%;
    position: fixed;
    right: 0;
}
.close-btn{
    text-align: center;
    margin-top: 50px;
}
.opts-btn{
    text-align: center; 
    margin-top: 10px; 
    padding-left: 60px; 
    padding-right: 60px;
    padding-top: 5px;
    padding-bottom: 5px;
    background-color: darkblue;
    border-color: darkblue;
    color: white;
    font-family: 'Times New Roman', Times, serif;
    font-size: 0.6cm;
}
</style>

<style scoped>
hr{
    margin-top: 25px;
    margin-bottom: 25px;
}
</style>