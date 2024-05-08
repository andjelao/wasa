<script setup>
import LoadingSpinner from '../components/LoadingSpinner.vue'
import Stream_Photo from '../components/Stream.vue';
</script>


<script>
export default {

    components: {
        LoadingSpinner,
        Stream_Photo
  },
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: localStorage.getItem("token"),
            banned: false,
            myprofile: true,
            newUsername: "",
            showChangeUsernameInput: false,
            isFollower: false,
            profile_username: "",
            isBanned: false,
            followersNumber: 0,
            followedNumber: 0,
            postsNumber: 0,
            followers: [],
            followed: [],
            bannedList: [],
            searchQuery: "",
            caption: "",


            
            posts: [],
            photos: [],
		}
	},

    watch: {
    '$route.params.username': function(newVal, oldVal) {

        // Handler function
        // Perform actions based on the changes to the route parameter 'username'
        if (newVal !== oldVal){
                this.errormsg=null
                this.photos=[]
                console.log(this.photos)
                this.refresh();
                console.log("new user event")
				return;
            }
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
                        this.errormsg= "Bad Request";
                    }else{
                        this.errormsg= "Error while updating username";
                    }   

                }

            } else {
                alert("The length of the new username must be between min 3 and max 16 characters")

            }
        },   

        async myProfile() {
                    this.profile_username = this.$route.params.username;
                    if (this.profile_username != this.username){
                        this.myprofile = false
                    } else {
                        this.myprofile = true
                    }
        },

        async postPhoto(){
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };

            this.loading = true;
            this.errormsg = null;

            const image = document.getElementById("fileInput").files[0];
            document.getElementById("fileInput").value = "";

            try {
                console.log(image)
                let fd = new FormData();
                fd.append("photo", image);
                console.log("appended photo")
                fd.append("caption", this.caption);
                let response = await this.$axios.post("/photos/", fd,{ headers })
                this.postsNumber+=1
                this.photos.unshift(response.data.photo);

            } catch (error) {
                this.errormsg= "Error while uploading the photo"
            }
            this.loading = false;
            this.caption=""
           
            document.querySelector('#newPostModal .btn-close').click();
            
        },

        async refresh() {
            this.banned= false
            this.username= localStorage.getItem("token")
            this.myProfile()
            this.searchQuery= ""
            if (localStorage.getItem("token")== null){
                this.$router.push("/");
                return
            }

            if (!this.myprofile){
                
                    // Retrieve authentication token from localStorage
                    const token = localStorage.getItem("token");
                    // Set authorization header
                    const headers = { Authorization: `Bearer ${token}` };
                    try {
                        let response = await this.$axios.get("/users/" +this.profile_username+"/profile",{ headers })
                        this.photos = response.data.userPhotos
                        // update followers
                        if (response.data.followersCount == null) {
                            this.followersNumber = 0;
                        } else {
                            this.followersNumber = response.data.followersCount
                        }

                        // update followers list
                        if (response.data.followersList == []) {
                            this.followers = [];
                        } else {
                            this.followers = response.data.followersList
                        }

                        
                        // update followed
                        if (response.data.followingCount == null) {
                            this.followedNumber = 0;
                        } else {
                            this.followedNumber = response.data.followingCount
                        }

                        // update followed list
                        if (response.data.followingList == []) {
                            this.followed = [];
                        } else {
                            this.followed = response.data.followingList
                        }



                        // update number of posts
                        if (response.data.photosCount== null) {
                            this.postsNumber = 0;
                        } else {
                            this.postsNumber = response.data.photosCount
                        }

                        //update banned list
                        try {
                            let bannedlistResponse = await this.$axios.get("/users/"+this.profile_username+"/banned/",{ headers })
                            this.bannedList = bannedlistResponse.data
                        }catch (error){
                            if (error.response.status === 404) {
                                this.bannedList = [];
                            }
                        }



                        // check if i follow this user
                        try {
                            let followResponse = await this.$axios.post("/users/" +this.username+"/followed/", {
                                followedUsername: this.profile_username},{ headers })
                            this.isFollower = false;
                            await this.$axios.delete("/users/" +this.username+"/followed/"+this.profile_username,{ headers })
                        }catch (error){
                            if (error.response.status === 409) {
                                // Set the 'isFollower' variable to true
                                this.isFollower = true;
                            }
                        }

                        // check if i banned this user
                        try {
                            let bannedResponse = await this.$axios.delete("/users/"+this.username+"/banned/"+this.profile_username,{ headers })
                            if (bannedResponse.status===204) {
                                this.isBanned=true
                                this.Ban()
                            }
                            
                        
                        }catch (error){
                            if (error.response.status === 404) {
                                this.isBanned = false;
                            }
                        }
                        
                     } catch (error){
                        if (error.response.status === 403) {
                            // Set the 'banned' variable to true
                            this.banned = true;
                        } else if (error.response.status === 404) {
                            // Set error message to "User not found"
                            this.errormsg = "User not found";
                        }
                    }

                    
            }
            else{
                // Retrieve authentication token from localStorage
                const token = localStorage.getItem("token");
                // Set authorization header
                const headers = { Authorization: `Bearer ${token}` };
                let response = await this.$axios.get("/users/" +this.profile_username+"/profile",{ headers })

                this.photos = response.data.userPhotos


                // update followers
                if (response.data.followersCount == null) {
                    this.followersNumber = 0;
                } else {
                    this.followersNumber = response.data.followersCount
                }

                // update followers list
                if (response.data.followersList == []) {
                    this.followers = [];
                } else {
                    this.followers = response.data.followersList
                }

                    
                // update followed
                if (response.data.followingCount == null) {
                    this.followedNumber = 0;
                } else {
                    this.followedNumber = response.data.followingCount
                }

                // update followed list
                if (response.data.followingList == []) {
                    this.followed = [];
                } else {
                    this.followed = response.data.followingList
                }



                // update number of posts
                if (response.data.photosCount== null) {
                    this.postsNumber = 0;
                } else {
                    this.postsNumber = response.data.photosCount
                }

                //update banned list
                try {
                    let bannedlistResponse = await this.$axios.get("/users/"+this.profile_username+"/banned/",{ headers })
                    this.bannedList = bannedlistResponse.data
                }catch (error){
                    if (error.response.status === 404) {
                        this.bannedList = [];
                    }
                }
            }
        },

        async Unfollow() {
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            await this.$axios.delete("/users/" +this.username+"/followed/"+this.profile_username,{ headers })
            this.isFollower = false
            this.followersNumber-=1

            try {
            let response= await this.$axios.get("/users/" +this.profile_username+"/followers/",{ headers })
            this.followers= response.data

            }
            catch(error){ 
                if (error.response.status === 404) {
                    this.followers =  []
                }
               }
            
        },

        async Follow() {
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            await this.$axios.post("/users/" +this.username+"/followed/", {
                        followedUsername: this.profile_username},{ headers })
            this.isFollower = true
            this.followersNumber+=1
            let response= await this.$axios.get("/users/" +this.profile_username+"/followers/",{ headers })
            this.followers= response.data
            
        },

        async Unban() {
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            await this.$axios.delete("/users/" +this.username+"/banned/"+this.profile_username,{ headers })
            this.isBanned = false
        },

        async Ban() {
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            await this.$axios.post("/users/" +this.username+"/banned/", {
                        BannedUsername: this.profile_username},{ headers })
            this.isBanned = true
        },
         async delPost(post) {
            console.log("post: " + post)
            for (let i = 0; i < this.photos.length; i++) {
                if (this.photos[i].photoId === post.photoId) {
                    console.log("found ")
                    this.photos.splice(i, 1);
                    i--;
                }

            };
            this.postsNumber -= 1;
            },

        async searchUser(){
            this.$router.push("/"+this.searchQuery+"/profile");
            //this.refresh()
        },

        
    }, 
    mounted() {
        this.refresh();
    },  
    }  
    
</script>

<template >
	<main>

		<header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPhoto</a>
            <router-link :to="'/'+username+'/profile'" class="nav-link">  
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#user" />
                </svg>
                My profile
            </router-link>
            <router-link :to="'/' + username + '/photostream'" class="nav-link">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#image" />
                </svg>
                My Photostream
            </router-link>
			<form class="d-flex" @submit.prevent="searchUser">
                <input class="form-control me-2" type="search" placeholder="Search users" aria-label="Search"  v-model="searchQuery">
                <button class="btn btn-outline-light border-0" type="submit" :disabled="searchQuery.trim() === ''">Search</button>
            </form>
            
		</header>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <div  v-if="errormsg === null" class="container-fluid" style="margin-left: 400px; margin-top: 30px; height: max-content;">
            <div class="row">

                <div class="col col-lg-3" style="font-size: large; ">
                    <div class="profile-info">
                        <b style="margin-bottom: 10px">@{{ profile_username }}</b>

                        <div v-if="banned">
                            You have been banned by this user!
                        </div>

                        <div v-else-if="myprofile" style="margin-bottom: 20px;margin-top: 20px;">
                        
                            <button class="btn btn-primary btn-md ms-4" type="button" @click="toggleChangeUsernameInput()">
                                Change username
                            </button>

                            <button class="btn btn-primary btn-md ms-4" type="button" data-bs-toggle="modal" data-bs-target="#newPostModal">
                                Add new post
                            </button>
                        
                            <br><br>
                            <div v-if="showChangeUsernameInput">
                                <input v-model="newUsername" type="text" placeholder="Type new username..." />
                                <button type="button" class="btn btn-primary btn-sm me-2" style="margin-left: 10px;" @click="changeUsername()">Save</button>
                                <button type="button" class=" btn btn-danger btn-sm" @click="toggleChangeUsernameInput()">Close</button>
                            </div>
                        </div>
        
                        <div v-else-if="!myprofile" style="margin-bottom: 20px; margin-top: 20px;">
                            <div class="row">
                                <div class="col-md-6 mb-2">            
                                            
                                    <div v-if="isFollower">
                                        <button class="btn btn-outline-primary btn-lg" type="button" @click="Unfollow()">
                                            <i class="bi-person-dash-fill"></i>
                                                Unfollow
                                            </button>
                                    </div>
                                                    
                                    <div v-else>
                                        <button class="btn btn-primary btn-lg" type="button" @click="Follow()">
                                            <i class="bi-person-plus-fill"></i>
                                                Follow
                                        </button>
                                    </div>
                                    <br>
                                    <div class="col-md-6 mb-2">
                                        <div v-if="isBanned">
                                            <button class="btn btn-success btn-lg" type="button" @click="Unban()">
                                                <i class="bi-person-check-fill"></i>
                                                Unban
                                            </button>
                                        </div>
                                        <div v-else>
                                            <button class=" btn btn-danger btn-lg" type="button" @click="Ban()">
                                                <i class="bi-person-x-fill"></i>
                                                Ban
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                                
                    </div>
                </div>

                <div v-if="!banned" class="row" style="margin-bottom: 20px;">
                    <div class="col col-sm-1" style="font-size: medium;">
                        
                        <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#followersModal">
                        <strong>Followers</strong>
                        </button>
                    </div>
                    <div class="col col-sm-1" style="font-size: medium;">
                        {{ followersNumber }}
                    </div>
                    

                    <div class="col col-sm-1" style="font-size: medium;">
                        <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#followedModal">
                        <strong>Following</strong>
                        </button>
                    </div>
                    <div class="col col-sm-1" style="font-size: medium;">
                        {{ followedNumber }}
                    </div>
                    <div class="col col-sm-1" style="font-size: medium;">
                        <strong>Posts</strong>
                    </div>
                    <div class="col col-sm-1" style="font-size: medium;">
                        {{ postsNumber }}
                    </div>
                    <div v-if="myprofile" class="col col-sm-1" style="font-size: medium;">
                        <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#bannedUsersModal">
                            <strong>Banned list</strong>
                        </button>
                    </div>
                </div>

                

                <div class="modal fade" id="followersModal" tabindex="-1" aria-labelledby="followersModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="followersModalLabel">Followers</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                            <ul>
                            <!-- Iterate over followers and display their names -->
                            <li v-for="follower in followers" :key="follower.followerUsername"><strong>@{{ follower.followerUsername}}</strong></li>
                            </ul>
                        </div>
                        </div>
                    </div>
                    </div>


                <div class="modal fade" id="followedModal" tabindex="-1" aria-labelledby="followedModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="followedModalLabel">Following</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                            <ul>
                            <!-- Iterate over followers and display their names -->
                            <li v-for="followed in followed" :key="followed.followedUsername"><strong>@{{ followed.followedUsername}}</strong></li>
                            </ul>
                        </div>
                        </div>
                    </div>
                    </div>
                
                <!-- Banned Users Modal -->
                <div class="modal fade" id="bannedUsersModal" tabindex="-1" aria-labelledby="bannedUsersModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                            <div class="modal-header">
                                <h5 class="modal-title" id="bannedUsersModalLabel">Banned Users</h5>
                                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                            </div>
                            <div class="modal-body">
                                <ul>
                                <!-- Iterate over followers and display their names -->
                                <li v-for="user in bannedList" :key="user.bannedUsername"><strong>@{{ user.bannedUsername}}</strong></li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>


                <!-- new post Modal -->
                <div class="modal fade" id="newPostModal" tabindex="-1" aria-labelledby="newPostModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                            <div class="modal-header">
                                <h5 class="modal-title" id="newPostModalLabel">Add new post</h5>
                                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                            </div>
                            <div class="modal-body">
                                <form @submit.prevent="postPhoto" class="form-container"
                            enctype="multipart/form-data">
                            
                            
                            <input class="form-control" type="file" id="fileInput" accept="image/jpeg, image/png"
                                style="width:fit-content;">
                            <br>  
                             <div class="mb-3">
                                <label for="captionInput" class="form-label">Caption:</label>
                                <input type="text" class="form-control" id="captionInput" v-model="caption">
                            </div>
                            <br>
                            <button type="submit" class="btn btn-success btn-sm" id="submitBut">Upload</button>
                            <br>
                        </form>
                                
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="!banned" class="col">
                <Stream_Photo :posts="photos"  @delete-post="delPost"></Stream_Photo>
                <div>
                <div v-if="loading">
                    <LoadingSpinner></LoadingSpinner>
                </div>
                </div>
            </div>

                
            </div>
        </div>  
        <RouterView></RouterView>      
    </main>
</template>
<style></style>


