<script setup>
import Post from "./Post.vue"

</script>

<script>

export default {

    props: {
        posts: {
            type: Object // Photoposts to display
        }
    },

    watch: {
        posts: function (new_posts) {
            this.posts_ = new_posts;
        }
    },

    data: function () {
        return {
            posts_: this.posts
        }
    },

    emits: ["delete-post"],
    

    methods: {

        async initialize() {

        },
        async DelPost(post) {

            this.$emit("delete-post", post);
        },


        

    },
    mounted() {

        this.initialize();

    }
}
</script>

<template>

    <div v-for="post in posts" :key="post.photoId" class="col-md-4 mx-2" style=" margin-top: 20px;">
        <Post :post="post" @delete-post="DelPost"></Post>
    </div>

</template>

<style>

</style>