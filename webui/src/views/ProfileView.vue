<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: localStorage.getItem("token"),
            banned: false,
            myprofile: true,
            newUsername: "",
            showChangeUsernameInput: false,
		}
	},
    methods: {


        async toggleChangeUsernameInput() {
            this.showChangeUsernameInput = !this.showChangeUsernameInput;
            this.newUsername = '';
        },
        async validateUsername(username) {
            const pattern = /^(?=.{3,16}$)(^.*?$)$/;
			return pattern.test(username);
        },
        async changeUsername() {

            if (this.validateUsername(this.newUsername)) {
                try {
                    console.log("newusernameinput: "+this.newUsername)

                    // Retrieve authentication token from localStorage
                    const token = localStorage.getItem("token");
                    // Set authorization header
                    const headers = { Authorization: `Bearer ${token}` };


                    let response = await this.$axios.put("/users/" + this.username, {
                        newUsername: this.newUsername}, { headers });

                    localStorage.setItem("token",response.data.Identifier);
                    this.username = response.data.Identifier

                    this.showChangeUsernameInput = false;
                    this.$router.push("/"+this.username+"/profile");
                } catch (error) {
                    if(error.response && error.response.status === 409){
                        alert("This username is already taken, please choose a different one!")
                    }else if (error.response && error.response.status === 400){
                        console.error("Bad Request: ", error)
                    }else{
                        console.log("Error while updating username", error)
                }
          

        }

      } else {
        alert("The length of the new username must be between min 3 and max 16 characters")

      }

    },
    },       
    
}
</script>

<template >
	<main>

		<header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPhoto</a>
			
		</header>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <div class="container-fluid" style="margin-left: 400px; margin-top: 30px; height: max-content;">
            <div class="row">

                <div class="col col-lg-3" style="font-size: large; ">
                    <div class="profile-info">
                        <b style="margin-bottom: 10px">@{{ username }}</b>

                        <div v-if="banned">
                        You have been banned by this user!
                        </div>
                        <div v-else-if="myprofile" style="margin-bottom: 20px;margin-top: 20px;">
                            <button class="btn btn-primary btn-md" type="button" @click="toggleChangeUsernameInput()">
                                Change username
                            </button>
                            <br><br>
                            <div v-if="showChangeUsernameInput">
                                <input v-model="newUsername" type="text" placeholder="Type new username..." />
                                <button type="button" class="btn btn-primary btn-sm me-2" style="margin-left: 10px;" @click="changeUsername()">Save</button>
                                <button type="button" class=" btn btn-danger btn-sm" @click="toggleChangeUsernameInput()">Close</button>
                            </div>
            </div>
        </div>
        </div>
        </div>
        </div>

    </main>
</template>
<style></style>


