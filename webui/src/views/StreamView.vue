<script setup>
import {RouterView } from 'vue-router'

import Stream_Photo from '../components/Stream.vue';
</script>

<script>
export default {

    components: {
        Stream_Photo
  },
  data: function() {
		return { 
            errormsg: null,
			loading: false,
			username: localStorage.getItem("token"),
            photos: [],
            exists: true,
        }
         },
    methods: {
        async initialize() {
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };

            try{
                let response = await this.$axios.get("/users/"+this.username+"/photo-stream",{ headers })
                this.photos=response.data
                this.posts= response.data
            }catch (error){
                            if (error.response.status === 404) {
                                // Set the 'isFollower' variable to true
                                this.exists = false;
                            }
                        }


        },
        async delPost(post) {
			
			for (let i = 0; i < this.posts.length; i++) {
				if (this.posts[i].photoId === post.photoId) {
					this.posts.splice(i, 1);
					i--;
				}

			};
			this.n_posts -= 1;
		},

    },
    mounted() {
		this.initialize();
		

	},
}
</script>

<template>
<main>
        <header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
                    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPhoto</a>
                     <div class="col-md-6 text-center">
                    <router-link :to="'/'+username+'/profile'" class="nav-link">  
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#user" />
                        </svg>
                        My profile
                    </router-link>
                     </div>
                </header>
                <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
                <div style="margin-left: 450px; margin-top: 30px; font-size: large;">
			<strong> My Stream:</strong>
		</div>
        <div v-if="exists" class="col" style="margin-left: 750px; margin-top: 30px;">
			
			<Stream_Photo :posts="photos" @delete-post="delPost()"></Stream_Photo>
			<div>
				<div v-if="loading">
					<LoadingSpinner></LoadingSpinner>
				</div>
			</div>
		</div>
		<RouterView :key="$route.fullPath"></RouterView>
</main>
</template>

<style></style>