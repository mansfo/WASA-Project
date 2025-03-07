<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
		    isLogged: false,
			searchedName: "",
			searchRes: [],
			searchPagetype: "",
			photoId: "",
		}
	},
	methods:{
		logout(log){
			this.isLogged = log
			this.$router.replace("/session")
		},
		updateLogged(newlog){
			this.isLogged = newlog
		},
		newView(view){
			this.$router.push(view)
		},
		async getUsers(name, pagetype, usrarray, phId){
			this.searchedName = name
			this.searchPagetype = pagetype
			try{
				if (this.searchPagetype === "SEARCH"){
					let response = await this.$axios.get("/users/"+localStorage.getItem("userid")+"/search/"+this.searchedName)
					this.searchRes = response.data != null? response.data : []
					this.$router.push("/users/"+localStorage.getItem("userid")+"/search/"+this.searchedName)
				} else if (this.searchPagetype === "BAN"){
					let response = await this.$axios.get("/users/"+localStorage.getItem("userid")+"/banned")
					this.searchRes = response.data != null? response.data : []
					this.$router.push("/users/"+localStorage.getItem("userid")+"/banned")
				} else if (this.searchPagetype === "FOLLOWERS"){
					this.searchRes = usrarray
					this.$router.push("/users/"+this.searchedName+"/followers")
				} else if (this.searchPagetype === "FOLLOWING"){
					this.searchRes = usrarray
					this.$router.push("/users/"+this.searchedName+"/following")
				} else if (this.searchPagetype === "LIKES"){
					this.searchRes = usrarray
					this.photoId = phId
					this.$router.push("/users/"+this.searchedName+"/photos/"+this.photoId+"/likes")
				} else {
					alert("Something went wrong loading UserListView")
				}
					
			}catch(e){
				alert(e.toString())
			}
		}
	},
	created(){
		this.$emit("logoutNavbar", false)
	},
	mounted(){
		if (localStorage.getItem("userid")){
			this.isLogged = true
		}
		else{
			this.$router.replace("/session")
		}
	},
}
</script>

<template>
    <div>
		<div>
			<div>
				<main>
					<Navbar v-if="isLogged"
					@logoutNavbar="logout"
					@toPageNavbar="newView"
					@searchNavbar="getUsers"/>
					<RouterView 
					@updateLog="updateLogged"
					@searchNavbar="getUsers"
					:searchedUser="searchedName"
					:users="searchRes"
					@getfollowers="getUsers"
					@getfollowing="getUsers"
					@getbanned="getUsers"
					@getlikes="getUsers"
					:page="searchPagetype"
					@toProfileFromPhoto="newView"
					@toProfileFromComm="newView"
					/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
