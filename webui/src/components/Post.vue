<script>
export default {

    props: {
        post: {
            type: Object
        }
     },

    data: function () {
        return {
             my_post: false,
             username: '',
             photo_id: null,
             caption: "",
             datapost: null,
             photo: null,
             isLiked: false,
             num_likes: 0,
             likes: null,
             comments: [],
             commentCount: 0,
             byteArray: null,
            token : localStorage.getItem("token"),
            newComment: '',
             
            }
    },

    emits: ["delete-post"],

    methods: {


        toggleLike() {
             if (this.isLiked) {
                this.Unlike();
            } else {
                this.Like();
            }
        },

        async initialize() {
            
            this.caption = this.post.caption
            this.username = this.post.author;
            this.photo_id = this.post.photoId
            
            if (this.username == localStorage.getItem("token")) {
                this.my_post = true;
            } else {
                this.my_post = false;
            }

            
            


            const imageData = this.post.photoEncoded;
            const byteCharacters = atob(imageData);
            const byteNumbers = new Array(byteCharacters.length);
            for (let i = 0; i < byteCharacters.length; i++) {
                byteNumbers[i] = byteCharacters.charCodeAt(i);
            }
            this.byteArray = new Uint8Array(byteNumbers);
            // Now you have the image data as a byte array (byteArray)
            // Create an image element
            
            // is the post liked
            const token = localStorage.getItem("token");
            this.likes = this.post.likes
            for (let i = 0; i < this.likes.length; i++) {
                if (this.likes[i].username === token) {
                   this.isLiked = true;
                    break;
                }
            }



            this.datapost = this.post.uploadDateTime
            
            this.num_likes= this.post.likesCount
            console.log(this.num_likes)
            this.comments = this.post.comments
            this.commentCount= this.post.commentsCount
        },
        async deletePhoto() {
            // Implement delete post functionality
            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            try {
                await this.$axios.delete("/photos/" +this.photo_id,{ headers })
                this.$emit("delete-post", this.post);
            }
            catch(error){    }   

        },

        async Like() {
            // Implement like post functionality

            // Retrieve authentication token from localStorage
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            try {
                await this.$axios.post("/photos/" +this.photo_id+"/likes/",{},{ headers })
                this.isLiked=true
                this.num_likes+=1
                let response= await this.$axios.get("/photos/" +this.photo_id+"/likes/",{ headers })
                this.likes= response.data
            }
            catch(error){ 
                
               }
            
                    
        },


        async Unlike() {
            // Implement unlike post functionality
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            try {
                await this.$axios.delete("/photos/" +this.photo_id+"/likes/"+token,{ headers})
                this.isLiked= false
                this.num_likes-=1
                let response= await this.$axios.get("/photos/" +this.photo_id+"/likes/",{ headers })
                this.likes= response.data
            }
            catch(error){
                if (error.response.status === 404) {
                    this.likes =  []
                }
                } 
                  
        },

        async deleteComment(commentId) {
            const token = localStorage.getItem("token");
            // Set authorization header
            const headers = { Authorization: `Bearer ${token}` };
            await this.$axios.delete("/photos/"+this.post.photoId+"/comments/"+commentId,{ headers})
            for (let i = 0; i < this.comments.length; i++) {
                if (this.comments[i].commentId === commentId) {
                this.comments.splice(i, 1);
                this.commentCount-=1
                i--;
                }
            
            /*try{
                let response=  await this.$axios.get("/photos/"+this.post.photoId+"/comments/")
                this.comments= response.data
            }
            catch(error){
                if (error.response.status === 404) {
                    this.comments =  []
                }
                } */
            }
         },
        
        async postComment() {
            if (this.newComment.trim() !== '') {
                try {
                    const token = localStorage.getItem("token");
                    // Set authorization header
                    const headers = { Authorization: `Bearer ${token}` };
                    let response= await this.$axios.post("/photos/" +this.photo_id+"/comments/",
                    { "text":this.newComment }, { headers})

                    this.comments.push(response.data)
                    this.newComment=""
                    this.commentCount+=1


                }catch(error){

                }
                this.newComment=""
                document.querySelector('#newcommentModal' + this.photo_id + ' .btn-close').click();
            }
        }
    },

     mounted() {
        try{
        this.initialize();
        
       /* // Set the image source to the decoded image data
        const blob = new Blob([this.byteArray], { type: 'image/jpeg' });
        const imgUrl = URL.createObjectURL(blob);
        const img = document.createElement("img");
        img.src = imgUrl;

        // Append the image element to a container in your HTML
        document.getElementById("imageContainer").appendChild(img);
        }*/}catch(error){
            console.log(error)
                }
                 
            
    },

    computed: {
        imageUrl() {
            if (this.byteArray) {
                const blob = new Blob([this.byteArray], { type: 'image/jpeg' });
                return URL.createObjectURL(blob);
            }
            return null;
        }
    },


    

}
</script>


<template>
    <div class="instagram-post">
        <!-- Username in the top left corner -->
        <div class="info-section">
        <div class="username">@{{username }}</div>

        <!-- Trash bin button in the top right corner -->
        <button  v-if="my_post" class="delete-button" @click="deletePhoto">
        <svg class="feather">
            <use href="/feather-sprite-v4.29.0.svg#trash-2" />
        </svg>
        </button>
        </div>
        <!-- Photo -->
        <div  id="imageContainer" class="photo" >
        <img v-if="imageUrl" :src="imageUrl" class="photo" />
        
        </div>


        <!-- Caption below the photo -->
        <div class="caption" >{{ caption }}</div>

        <!-- Heart icon button in the bottom left corner -->
        
        <!--<div class="icons-container">-->
         <div class="likes-container">
        <button class="heart-button" @click="toggleLike">
        <svg class="feather">
            <use v-if="isLiked" href="/feather-sprite-v4.29.0.svg#heart" fill="red" />
            <use v-else href="/feather-sprite-v4.29.0.svg#heart" />
        </svg>
        
        </button>
        <span class="like-count" data-bs-toggle="modal" :data-bs-target="'#likesModal' + post.photoId">{{ num_likes }} likes</span> <!-- Number of likes -->
        </div>

        <!-- Comment icon button in the bottom right corner -->
        <div class="comments-container">
        <button class="comment-button" data-bs-toggle="modal" :data-bs-target="'#newcommentModal' + post.photoId">
        <svg class="feather">
            <use href="/feather-sprite-v4.29.0.svg#message-circle" />
        </svg>
         
        </button>
            <span class="comment-count" data-bs-toggle="modal" :data-bs-target="'#commentsModal' + post.photoId">{{ commentCount }} comments</span> <!-- Number of comments -->
        <!--</div>-->
        </div>
        <!-- Upload date and time of the photo -->
        <div class="upload-date">{{datapost }}</div>


        <!-- Likes modal -->

       <div class="modal fade" :id="'likesModal' + post.photoId" tabindex="-1" aria-labelledby="likesModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="likesModalLabel">Likes</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                            <ul>
                            <!-- Iterate over followers and display their names -->
                            <li v-for="like in likes" :key="like.username"><strong>@{{ like.username}}</strong></li>
                            </ul>
                        </div>
                        </div>
                    </div>
                    </div>

        <!-- Comments modal -->
        <div class="modal fade" :id="'commentsModal' + post.photoId" tabindex="-1" aria-labelledby="commentsModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="commentsModalLabel">Comments</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                            <ul>
                            <!-- Iterate over followers and display their names -->
                            <li v-for="comment in comments" :key="comment.commentId"><strong>@{{ comment.author}}</strong> / {{comment.date}} : {{ comment.text}}
                            <button v-if="comment.author === token" @click="deleteComment(comment.commentId)" class="btn btn-danger btn-sm">Delete</button>
                            </li>
                            </ul>

                            
                        </div>
                        </div>
                    </div>
                    </div>
        <!-- New comments modal -->
        <div class="modal fade" :id="'newcommentModal' + post.photoId" tabindex="-1" aria-labelledby="newcommentModalLabel" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered">
                        <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="newcommentModalLabel">Post a new comment</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                        <!-- New comment form -->
                            <form @submit.prevent="postComment">
                                <div class="mb-3">
                                    <label for="commentText" class="form-label">Comment:</label>
                                    <textarea class="form-control" id="commentText" v-model="newComment"></textarea>
                                </div>
                                <button type="submit" class="btn btn-primary" :disabled="newComment.trim() === ''">Post</button>
                                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                            </form>
                            </div>
                        </div>
                    </div>
                    </div>



  </div>







</template> 



<style scoped>
    .instagram-post {
  position: relative;
  display: inline-block;
  max-width: 450px;
  border: 1px solid #ccc;
  border-radius: 5px;
  overflow: hidden;
  padding-bottom: 40px; /* Add padding to make space for the icons */
}

.username {
  
  font-weight: bold;
}

.delete-button {
  
  background-color: transparent;
  border: none;
  cursor: pointer;
   
}

.photo {
  width: 100%;
  height: 400px;
  overflow: hidden;
   margin-top: 30px;
}

.icons-container {
  position: absolute;
  bottom: 10px; /* Adjust the position */
  padding: 5px;
  display: flex; /* Use flexbox */
  align-items: center; /* Align items vertically */
  
}
.likes-container {
    position: absolute;
    bottom: 10px;
    left: 10px; /* Adjusted to bottom left */
    display: flex;
    align-items: center;
}

.comments-container {
    position: absolute;
    bottom: 10px;
    right: 10px; /* Adjusted to bottom right */
    display: flex;
    align-items: center;
}

.heart-button, .comment-button{
  margin-right: 5px;
  
  background-color: transparent;
  border: none;
  cursor: pointer;
}


.caption {
  padding: 10px;
  margin-top: 5px; /* Add margin to separate from the photo */
  font-size: 14px;
}

.upload-date {
  margin-top: 10px;
  font-size: 12px;
}

.like-count, .comment-count{
  cursor: pointer;
  font-size: 12px;
}

.info-section {
        background-color: white;
        position: absolute;
        top: 0;
        left: 0;
         width: 100%; /* Adjusted width for delete button */
        padding: 10px;
        box-sizing: border-box;
         display: flex;
         justify-content: space-between;
}


</style>