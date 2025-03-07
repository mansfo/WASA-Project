<script>
export default {
	data: function() {
		return {
			username: "",
			identifier: "",
			validName: false,
		}
	},
	methods: {
		async doLogin() {
			try {
				let response = await this.$axios.post('/session', {
				    username: this.username
				});
				this.identifier = response.data.user_id.user_id;
				this.identifier.trim()
				
			    localStorage.setItem('userid', this.identifier)
			    localStorage.setItem('username', this.username)
				this.$router.push('/users/'+this.identifier+'/stream')
				this.$emit("updateLog", true)
				this.validName = true
				
			} catch (e) {
				alert(e.toString());
			}
		},
	},
	async mounted(){
		this.$emit("updateLog", false)
		if (this.validName){
		    this.$router.push('/users/'+this.identifier+'/stream')
		} else {
			this.$router.push('/session')
		}
	}
}
</script>

<template>
	<div class="login-page">
		<div class="login-title">
			<h1 style="color: darkblue; font-size: 2cm;"><b>WASAPhoto</b></h1>
		    <h5 style="color: black; font-style: italic;">Keep in touch with your friends</h5>
		    <h5 style="color: black; font-style: italic;">thanks to WASAPhoto!</h5>
		</div>
		<div class="log-form">
			<div id="login-square">
				<form id="login-form" method="post" @submit.prevent="doLogin">
					<div class="text-field">
						<h4 style="color: black; font-family: 'Times New Roman', Times, serif; text-align: left;"><i>Your Username</i></h4>
						<input style="padding-left: 5px;" type="text" v-model="username" placeholder="Username" /><br>
						<button id="logbtn" type="button" @click = "doLogin"><h4><b>Login</b></h4></button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style>
.login-page{
	margin-left: 0px;
	margin-top: 0px;
	margin-right: 0px;
	margin-bottom: 0px;
}
.login-title, .log-form{
	text-align: center;
	padding: 30px;
	margin-top: 10%;
}
h1, h5{
	font-family: 'Times New Roman', Times, serif;
}
h5{
	font-size: 0.8cm;
}
.login-title{
	float: left;
	margin-left: 20%;
}
.log-form{
	float: right;
	margin-right: 20%;
	margin-top: 11%;
}
#logbtn{
	background-color: darkblue; 
	padding-left: 70px; 
	padding-right: 70px;
	margin-top: 10px;
	margin-left: 5px;
	margin-right: 5px;
	font-family: 'Times New Roman', Times, serif; 
	color: white;
	border-color: blue;
}
#login-square{
	visibility: visible;
	background-color: white;
	padding: 30px;
	padding-left: 80px;
	padding-right: 80px;
	border-radius: 10%;
}
</style>