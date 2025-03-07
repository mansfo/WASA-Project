<script>
export default{
    data(){
        return{
            text: "",
        }
    },
    methods:{
        logout(){
            localStorage.removeItem("userid")
            this.$emit("logoutNavbar", false)
        },
        homepage(){
            this.$emit("toPageNavbar", "/users/"+localStorage.getItem("userid")+"/stream")
        },
        search(){
            this.$emit("searchNavbar", this.text, "SEARCH", [], "")
            this.text = ""
        },
        profile(){
            this.$router.push("/users/"+localStorage.getItem("userid")+"/profile")
        },
    }
}

</script>

<template>
    <nav class="menu">
        <div class="site-name-home">
            <div class="home-btn">
                <button style="margin-top:0px; font-size: 1cm; background-color: darkblue; color: white; border-color: darkblue;" @click="homepage" type="button"><b>WASAPhoto</b></button>
            </div>
            <div class="search-section">
                <form class="search-form">
                    <input v-model="text" type="search" placeholder="Search user" style="padding: 5px; padding-right: 10px; font-size:0.5cm; font-family: 'Times New Roman', Times, serif;">
                    <button type="submit" @click.prevent="search" class="menu-btn" style="margin-left: 5px;">Search</button>
                </form>
            </div>
            <div class="buttons">
                <button @click="profile" class="menu-btn" type="button"><b>My Profile</b></button>
                <button @click="logout" class="menu-btn" type="button"><b>Logout</b></button>
            </div>
        </div>
    </nav>
</template>

<style>
.menu{
    background-color: darkblue;
    margin-top: 0px;
    padding-bottom: 90px;
}
.site-name-home{
    margin-left: 50px;
    color: white;
    font-family: 'Times New Roman', Times, serif;
}
.search-section{
    margin-left: 30px;
    float: left;
    margin-top: 20px;
}
.search-form{
    margin-left: 50px;
}
.home-btn{
    float: left;
    margin-top: 10px;
}
.buttons{
    float: right;
    margin-top: 20px;
    margin-right: 50px;
}
.menu-btn{
    padding-top: 5px;
    padding-bottom: 5px;
    padding-left: 20px;
    padding-right: 20px;
    background-color: darkblue;
    border-color: darkblue;
    color: white;
    font-family: 'Times New Roman', Times, serif;
    font-size: 0.6cm;
    margin-left: 20px;
    margin-top: 0px;
}
</style>