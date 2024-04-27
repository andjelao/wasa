<script setup>
import { axios } from 'axios';
</script>

<script>
export default {
	data: function () {
		return {
			errormsg: null,
			username: '',
            usermessage: '',
		}
	},
	methods: {
		validateuser(username) {
			const pattern = /^(?=.{3,16}$)(^.*?$)$/;
			return pattern.test(username);

		},

	

	async login() {
		if (this.validateuser(this.username)) {
			try {
				let response = await this.$axios.post("/session", {
					"username": this.username
				});

				this.errormsg = null;

                localStorage.setItem("token",response.data.Identifier);

                if (response.status === 200) {
						this.usermessage = `Welcome back, ${this.username}`;
					} else if (response.status === 201) {
						this.usermessage = `Welcome, ${this.username}`;
					}

                setTimeout(() => {
						// Navigate to profile page
						this.$router.push("/"+this.username+"/profile");
					}, 1000);
                
                
                
			} catch (err) {
				this.errormsg = err.message;
			}

		}else{
			alert("The length of the username must be between 3 and 16 characters");
		}



	},
}
}


</script>

<template >
	<main>

		<header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPhoto</a>
			
		</header>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div style="height: 300px; display: flex; justify-content: center; align-items: center; ">

			<div class="col-md-2">
				<form @submit.prevent="login()">
					<div style="padding:5px;">
						<input type="text" v-model="username" placeholder="Username" style="border-color: #650C96;" />
					</div>
					<div class="btn-toolbar mb-2 mb-md-0">
						<div class="btn-group me-2" style="padding:5px;">
							<button type="submit" class="btn btn-sm btn-outline-secondary"
								style=" background-color:#008bff; color:antiquewhite; border-color:#650C96;">
								Login / Sign Up
							</button>
						</div>
					</div>

				</form>
			</div>
		</div>
        <p v-if="usermessage" class="text-center display-4">{{ usermessage }}</p>
		<RouterView></RouterView>
	</main>
</template>
<style></style>